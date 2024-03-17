package usecases

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pahuss/otus/models"
	"github.com/pahuss/otus/repository"
	"github.com/ramonmacias/go-auth-middleware/auth"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
)

type Auth struct {
	Redis             *redis.Client
	UserRepository    *repository.UserRepository
	SessionRepository *repository.SessionRepository
}

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInternal = errors.New("internal error")
var ErrProfileExist = errors.New("profile exist")
var ErrValidate = errors.New("validate error")

func (a Auth) Register(registration *models.UserRegistration) (*models.Profile, error) {
	passwordHash, err := encodePassword(registration.Password)
	if err != nil {
		return nil, ErrProfileExist
	}
	registration.Password = passwordHash
	id, err := a.UserRepository.InsertProfile(registration)
	log.Println("Register insert", id, err)
	if err != nil {
		if errors.Is(err, repository.ErrProfileExist) {
			err = ErrProfileExist
		}
		return nil, err
	}

	return &models.Profile{
		ID:        id,
		LastName:  registration.LastName,
		FirstName: registration.FirstName,
		Email:     registration.Email,
		Age:       0,
		Hobbies:   "",
		City:      "",
	}, nil
}

func (a Auth) Login(c *models.CredentialsForm) (auth.Session, error) {
	s := auth.Session{}
	passwordHash, err := encodePassword(c.Password)
	if err != nil {
		return s, err
	}

	creds, err := a.UserRepository.PasswordByEmail(c.Email)
	creds.Email = c.Email

	if err != nil {
		if errors.Is(err, repository.ErrNoUser) {
			return s, ErrInvalidCredentials
		}
		return s, err
	}

	if passwordHash != creds.Password {
		return s, ErrInvalidCredentials
	}

	UUID, err := uuid.NewGen().NewV4()
	if err != nil {
		return s, err
	}

	s = auth.Session{
		Email:  creds.Email,
		UserID: UUID,
	}

	err = a.SessionRepository.Add(s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func encodePassword(password string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, password)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
