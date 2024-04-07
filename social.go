package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pahuss/otus/api"
	app2 "github.com/pahuss/otus/app"
	"github.com/pahuss/otus/midlleware"
	"github.com/pahuss/otus/repository"
	"github.com/pahuss/otus/usecases"
	"github.com/ramonmacias/go-auth-middleware/auth"
	"github.com/ramonmacias/go-auth-middleware/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Load environment error")
	}

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Starting server")

	app := app2.NewApp(context.Background())
	app.InitDb(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBHOST"), "tcp")
	app.InitRedis(os.Getenv("REDISHOST"), "", 0)

	sessionRepository := &repository.SessionRepository{
		Redis: app.Redis,
	}

	us := usecases.NewUserService(app.Db, app.Redis, sessionRepository)

	apiHandler := api.ApiHandler{
		Db:          app.Db,
		UserService: us,
		App:         app,
	}
	authProvider := midlleware.SimpleAuthProvider{
		App:               app,
		SessionRepository: sessionRepository,
	}

	router := mux.NewRouter()

	router.HandleFunc("/login", apiHandler.Login).Methods("POST")
	router.HandleFunc("/user/register", apiHandler.Register).Methods("POST")

	// пока что оставляю не закрытым автроризацией
	router.HandleFunc("/user/search", apiHandler.Search).Methods("GET")

	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.HandleFunc("/user/get/{id}", apiHandler.Profile).Methods("GET")
	apiRouter.HandleFunc("/user/profile", apiHandler.Profile).Methods("GET")

	apiRouter.Use(middleware.AuthAPI(authProvider, func(userSession *auth.Session) error {
		fmt.Sprintln(userSession)
		return nil
	}))
	log.Fatal(http.ListenAndServe(":8080", router))
}
