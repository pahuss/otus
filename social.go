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
		Net:    "tcp",
		Addr:   "mysql_db:3306",
		DBName: os.Getenv("DBNAME"),
	}

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Starting server")

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	log.Println(cfg.FormatDSN())
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

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
