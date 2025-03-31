package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPasswordResetTokenExpired = errors.New("models: password reset token has expired")
	ErrPasswordResetTokenInvalid = errors.New("models: invalid token")
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created bc only TokenHash is stored in DB.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

const (
	// DefaultResetDuration is the default time that a PasswordReset is
	// valid for.
	DefaultResetDuration = 1 * time.Hour
)

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int

	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

// Create a new reset token which can be used by the email service
func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	var userID int
	row := service.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, email)
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmailNotFound
		}
		return nil, fmt.Errorf("create: %w", err)
	}

	tm := TokenManager{
		BytesPerToken: service.BytesPerToken,
	}
	token, tokenHash, err := tm.New()

	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}

	row = service.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

// We are going to consume a token and return the user associated with it, or return an error if the token wasn't valid for any reason.
func (service *PasswordResetService) Consume(token string) (*User, error) {
	tm := TokenManager{
		BytesPerToken: service.BytesPerToken,
	}
	tokenHash := tm.Hash(token)

	var user User
	var pwReset PasswordReset

	row := service.DB.QueryRow(`
		SELECT password_resets.id,
			password_resets.expires_at,
			users.id,
			users.email, 
			users.password_hash
		FROM password_resets  
			JOIN users ON users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1;
	`, tokenHash)
	err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPasswordResetTokenInvalid
		}
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, ErrPasswordResetTokenExpired
	}
	err = service.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (service *PasswordResetService) delete(id int) error {
	_, err := service.DB.Exec(`
	DELETE FROM password_resets
	WHERE password_resets.id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
