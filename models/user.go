package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	var user User
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user.Email = email
	user.PasswordHash = string(passwordHash)

	row := us.DB.QueryRow(`INSERT INTO users(email, password_hash) VALUES($1, $2) RETURNING id;`, email, string(passwordHash))

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, err
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	var user User
	user.Email = email
	row := us.DB.QueryRow(`
    SELECT id, password_hash 
    FROM users WHERE email=$1;`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenicate: %w", err)
	}
	fmt.Println("success!!")
	return &user, nil
}

func (us *UserService) UpdatedPassword(userID int, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	_, err = us.DB.Exec(`UPDATE users SET password_hash = $2 where users.id = $1`, userID, string(passwordHash))
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}
	return nil
}
