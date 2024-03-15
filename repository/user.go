package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pahuss/otus/models"
)

var ErrNoUser = errors.New("user: user not fond")
var ErrProfileExist = errors.New("profile exist")

type UserRepository struct{ Db *sql.DB }

//func NewUserRepository(db *sql.DB) UserRepository {
//	return UserRepository{Db: db}
//}

func (r UserRepository) InsertProfile(regForm *models.UserRegistration) (int64, error) {
	var id int64
	query := "INSERT INTO `user` (`email`, `first_name`, `last_name`, `password`, `city_id`) VALUES (?, ?, ?, ?, ?)"
	insertResult, err := r.Db.ExecContext(context.Background(), query, regForm.Email, regForm.FirstName, regForm.LastName, regForm.Password, nil)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		ok := errors.As(err, &mysqlErr)

		if !ok {
			return id, err
		}

		if mysqlErr.Number == 1062 {
			return id, ErrProfileExist
		}
	}
	id, err = insertResult.LastInsertId()
	return id, err
}

func (r UserRepository) UserProfile(email string) (models.Profile, error) {
	var user models.Profile
	row := r.Db.QueryRow("SELECT id, first_name, last_name, email, age FROM user WHERE email = ?", email)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrNoUser
		}
		return user, fmt.Errorf("UserProfile %d: %v", email, err)
	}
	return user, nil
}

func (r UserRepository) PasswordByEmail(email string) (models.Credentials, error) {
	c := models.Credentials{}
	row := r.Db.QueryRow("SELECT id, password FROM user WHERE email = ?", email)
	if err := row.Scan(&c.ID, &c.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c, ErrNoUser
		}
		return c, fmt.Errorf("PasswordByEmail %s: %v", email, err)
	}
	return c, nil
}
