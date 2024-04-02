package models

import (
	"bryce-ziemer/github.com/lenslocked/rand"
	"database/sql"
	"fmt"
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
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	token, err := rand.SessionToken()

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
