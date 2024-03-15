package midlleware

import (
	"github.com/pahuss/otus/app"
	"github.com/pahuss/otus/repository"
	"github.com/ramonmacias/go-auth-middleware/auth"
)

type AuthMiddleware interface{}

type UserAuthMiddleware struct {
	p auth.Provider
}

type SimpleAuthProvider struct {
	App               *app.App
	SessionRepository *repository.SessionRepository
}

func (a SimpleAuthProvider) Sign(s auth.Session) (string, error) {
	return "", nil
}

func (a SimpleAuthProvider) Refresh(token string) (string, error) {
	return token, nil
}

func (a SimpleAuthProvider) Validate(token string) (*auth.Session, error) {
	s, _ := a.SessionRepository.Get(token)
	a.App.CurrentSession = s
	//id, err := uuid.FromString(token)
	//if err != nil {
	//	return &s, err
	//}
	return &s, nil
}

//func (a SimpleAuthProvider) GetSession() *auth.Session {
//	return &a.s
//}
