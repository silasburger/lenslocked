package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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

var (
	// A common pattern is to add the package as a prefix to the error for
	// context.
	ErrEmailTaken        = errors.New("models: email address is already in use")
	ErrEmailNonexistent  = errors.New("models: no account exists with that email address")
	ErrPasswordIncorrect = errors.New("models: password entered for email is incorrect")
)

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
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmailNonexistent
		}
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrPasswordIncorrect
		}
		return nil, fmt.Errorf("authenicate: %w", err)
	}
	fmt.Println("success!!")
	return &user, nil
}

func (us *UserService) UpdatePassword(userID int, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	_, err = us.DB.Exec(`
		UPDATE users 
		SET password_hash = $2 
		WHERE users.id = $1`, userID, string(passwordHash))
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (us *UserService) UpdateEmail(userID int, email string) error {
	_, err := us.DB.Exec(`
	UPDATE users 
	SET email = $2 
	WHERE users.id = $1`, userID, email)
	if err != nil {
		return fmt.Errorf("update email: %w", err)
	}
	return nil
}
