package app

import (
	"context"
	"github.com/ramonmacias/go-auth-middleware/auth"
)

type App struct {
	CurrentSession auth.Session
	Ctx            context.Context
}

func NewApp(ctx context.Context) *App {
	return &App{
		Ctx: ctx,
	}
}
