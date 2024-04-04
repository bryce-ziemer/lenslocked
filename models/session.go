package models

import (
	"bryce-ziemer/github.com/lenslocked/rand"
	"database/sql"
	"fmt"
)

const (
	// The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating new session.  When looking up the session
	// this will be left empty as we only styore the hash of a session token in our DB
	// and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating each session token.
	// If this value is not set or is less than the MinBytesPerToken const it will be ignored and MinBytesPerToken will be used.
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)

	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	// TODO: Hash the session token
	session := Session{
		UserID: userID,
		Token:  token,
		// TODO set the token hash

	}
	// TODO store the session in our DB

	return &session, nil

}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO Implement SessionService.User
	return nil, nil
}
