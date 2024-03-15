package usecases

import (
	"database/sql"
	"github.com/pahuss/otus/models"
	"github.com/pahuss/otus/repository"
	"github.com/ramonmacias/go-auth-middleware/auth"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	Auth              *Auth
	UserRepository    *repository.UserRepository
	SessionRepository *repository.SessionRepository
	Session           *auth.Session
	Redis             *redis.Client
}

func (s UserService) Profile(email string) (*models.Profile, error) {
	profile, err := s.UserRepository.UserProfile(email)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func NewUserService(db *sql.DB, rc *redis.Client, sr *repository.SessionRepository) *UserService {
	r := &repository.UserRepository{
		Db: db,
	}
	return &UserService{
		SessionRepository: sr,
		UserRepository:    r,
		Redis:             rc,
		Auth:              &Auth{UserRepository: r, SessionRepository: sr},
	}
}
