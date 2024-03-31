package models

import (
	"database/sql"
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
	//1. Create the session token - we want to handlke this in the service bc, otherwise, a poor token may be p[assed in and used]
	// TODO: create the session token
	// TODO: Implement SessionService.Create
	return nil, nil

}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO Implement SessionService.User
	return nil, nil
}
