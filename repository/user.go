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

func (r UserRepository) InsertProfile(user *models.User) (int64, error) {
	var id int64
	id, err := r.insertUser(user)

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
	return id, err
}

func (r UserRepository) Search(user models.Profile) ([]models.Profile, error) {
	var profiles []models.Profile
	firstNameQuery := user.FirstName + "%"
	lastNameQuery := user.LastName + "%"
	rows, err := r.Db.QueryContext(context.Background(), "SELECT id, first_name, last_name, email, age FROM user WHERE first_name LIKE ? and last_name LIKE ? ORDER BY id", firstNameQuery, lastNameQuery)
	if err != nil {
		return profiles, err
	}
	defer rows.Close()
	for rows.Next() {
		var profile models.Profile
		err = rows.Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Email, &profile.Age)
		if err != nil {
			return profiles, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (r UserRepository) UserProfile(email string) (models.Profile, error) {
	var profile models.Profile
	user, err := r.getUser(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profile, ErrNoUser
		}
		return profile, fmt.Errorf("UserProfile %d: %v", email, err)
	}
	profile = models.UserToProfile(user)
	return profile, nil
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

func (r UserRepository) insertUser(user *models.User) (int64, error) {
	var id int64
	query := "INSERT INTO `user` (`email`, `first_name`, `last_name`, `password`, `age`, `city_id`) VALUES (?, ?, ?, ?, ?, ?)"
	insertResult, err := r.Db.ExecContext(context.Background(), query, user.Email, user.FirstName, user.LastName, user.Password, user.Age, nil)

	if err != nil {
		return id, err
	}

	id, err = insertResult.LastInsertId()
	return id, err
}

func (r UserRepository) getUser(email string) (models.User, error) {
	var user models.User
	row := r.Db.QueryRow("SELECT id, first_name, last_name, email, age FROM user WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Age)
	return user, err
}
