package models

import (
	"bryce-ziemer/github.com/lenslocked/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating each password reset token.
	// If this value is not set or is less than the MinBytesPerToken const it will be ignored and MinBytesPerToken will be used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

// use when attempting to create a new password reset token
func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// verify we have a valifd email address for a user, and get that usersID
	email = strings.ToLower(email)
	var userID int
	row := service.DB.QueryRow(`
	SELECT id FROM users WHERE email = $1;`, email)
	err := row.Scan(&userID)
	if err != nil {
		// TODO: consider returning a specific error when the user does not exist
		return nil, fmt.Errorf("Create: %w", err)
	}

	// Build the PasswordReset
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)

	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: service.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	// Insert the PasswordRest into the DB
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

func (service *PasswordResetService) Consume(toke string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordEmailService.Consume()")

}

func (service *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
