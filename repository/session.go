package repository

import (
	"context"
	"github.com/ramonmacias/go-auth-middleware/auth"
	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	Redis *redis.Client
}

func (sr SessionRepository) Add(s auth.Session) error {
	ctx := context.Background()
	return sr.Redis.Set(ctx, s.UserID.String(), s.Email, 0).Err()
}

func (sr SessionRepository) Get(key string) (auth.Session, error) {
	ctx := context.Background()
	email, err := sr.Redis.Get(ctx, key).Result()

	s := auth.Session{}

	if err != nil {
		return s, err
	}

	s.Email = email
	return s, nil
}
