package models

import (
	"database/sql"
	"errors"

	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {

		return nil
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`
	_, err = m.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil

}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashed_password []byte

	smtp := "SELECT id, hashed_password FROM users WHERE email=$1"

	err := m.DB.QueryRow(smtp, email).Scan(&id, &hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
	}

	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {

		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil

}

// check user an id in database
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id=$1)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
