package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pahuss/otus/api"
	app2 "github.com/pahuss/otus/app"
	"github.com/pahuss/otus/midlleware"
	"github.com/pahuss/otus/repository"
	"github.com/pahuss/otus/usecases"
	"github.com/ramonmacias/go-auth-middleware/auth"
	"github.com/ramonmacias/go-auth-middleware/middleware"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Load environment error")
	}
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		//Net:    "tcp",
		//Addr:   "127.0.0.1:3306",
		DBName: "social_db",
	}

	client := redis.NewClient(&redis.Options{
		//Addr:     "127.0.0.1:6379",
		//Password: "",
		//DB:       0,
	})

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	//r := createRepository(db)

	//fmt.Println(cfg)
	//mux := http.NewServeMux()

	app := app2.NewApp(context.Background())

	sessionRepository := &repository.SessionRepository{
		Redis: client,
	}

	us := usecases.NewUserService(db, client, sessionRepository)

	apiHandler := api.ApiHandler{
		Db:          db,
		UserService: us,
		App:         app,
	}
	authProvider := midlleware.SimpleAuthProvider{
		App:               app,
		SessionRepository: sessionRepository,
	}
	//s := authProvider.GetSession()
	router := mux.NewRouter()

	router.HandleFunc("/login", apiHandler.Login).Methods("POST")
	router.HandleFunc("/user/register", apiHandler.Register).Methods("POST")

	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.HandleFunc("/user/get/{id}", apiHandler.Profile).Methods("GET")
	apiRouter.HandleFunc("/user/profile", apiHandler.Profile).Methods("GET")

	apiRouter.Use(middleware.AuthAPI(authProvider, func(userSession *auth.Session) error {
		fmt.Sprintln(userSession)
		return nil
	}))
	log.Fatal(http.ListenAndServe(":8080", router))
}
