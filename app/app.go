package app

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/ramonmacias/go-auth-middleware/auth"
	"github.com/redis/go-redis/v9"
	"log"
)

type App struct {
	Db             *sql.DB
	Redis          *redis.Client
	CurrentSession auth.Session
	Ctx            context.Context
}

func NewApp(ctx context.Context) *App {
	return &App{
		Ctx: ctx,
	}
}

func (a *App) InitDb(user, pass, dbname, addr, net string) {
	cfg := mysql.Config{
		User:   user,
		Passwd: pass,
		DBName: dbname,
		Addr:   addr,
		Net:    net,
	}
	var err error
	a.Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) InitRedis(addr, password string, db int) {
	a.Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
