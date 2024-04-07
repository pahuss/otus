package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pahuss/otus/app"
	"github.com/pahuss/otus/models"
	"github.com/pahuss/otus/usecases"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var ErrNoUser = errors.New("user: user not fond")

type ApiHandler struct {
	Db          *sql.DB
	UserService *usecases.UserService
	App         *app.App
}

func (api ApiHandler) Register(w http.ResponseWriter, r *http.Request) {
	var regForm models.UserRegistration
	err := decodeJSONBody(w, r, &regForm)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	p, err := api.UserService.Auth.Register(&regForm)

	if errors.Is(err, usecases.ErrInternal) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errors.Is(err, usecases.ErrProfileExist) {
		http.Error(w, "Profile exist", http.StatusUnprocessableEntity)
		return
	}

	if errors.Is(err, usecases.ErrValidate) {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(p)
}

func (api ApiHandler) Profile(w http.ResponseWriter, r *http.Request) {
	if api.App.CurrentSession.Email == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	p, err := api.UserService.Profile(api.App.CurrentSession.Email)
	if err != nil {
		if errors.Is(err, ErrNoUser) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(p)
}

func (api ApiHandler) Search(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("firstName")
	lastName := r.URL.Query().Get("lastName")

	if firstName == "" && lastName == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	searchData := models.Profile{
		FirstName: firstName,
		LastName:  lastName,
	}

	p, err := api.UserService.Search(searchData)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(p)
}

func (api ApiHandler) UserProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	u, err := userByID(id, api.Db)
	if err != nil {
		if errors.Is(err, ErrNoUser) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (api ApiHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login")
	var c models.CredentialsForm
	err := decodeJSONBody(w, r, &c)

	if c.Email == "" || c.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := api.UserService.Auth.Login(&c)
	log.Println(err)
	if err != nil {
		http.Error(w, "Bad password", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "api/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func userByID(uid any, db *sql.DB) (models.User, error) {
	id, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		panic(err)
	}
	// An album to hold data from the returned row.
	var user models.User

	row := db.QueryRow("SELECT id, first_name, last_name, email, age FROM user WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrNoUser
		}
		return user, fmt.Errorf("userByID %d: %v", id, err)
	}
	return user, nil
}

func asJson(value any) ([]byte, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
