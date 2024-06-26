package models

import (
	"bryce-ziemer/github.com/lenslocked/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
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
	// Could move all this to a tokenManager object and have a new() method to generate new tokens (woudld return (token, tokenHash, err))
	// essentially moving all the logic to create the token into somethign else (in this case a manager)
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)

	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	//Option 1
	//1. Query for a user's session
	//2. If found, update the user's session
	//3. If not found, create a new sessions for the user

	// Option 2
	//1. Try to update the user's session
	//2. If err, create a new session
	row := ss.DB.QueryRow(`
	INSERT INTO sessions (user_id, token_hash)
	VALUES ($1, $2) ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2
	RETURNING id;`, session.UserID, session.TokenHash)

	err = row.Scan(&session.ID) // fills out session.UserID

	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}
	return &session, nil

}

func (ss *SessionService) User(token string) (*User, error) {
	//1. Hash the session token
	tokenHash := ss.hash(token)

	//2. Query for the session with that hash
	var user User
	row := ss.DB.QueryRow(`
	SELECT users.id, users.email, users.password_hash
	FROM sessions
		JOIN users on users.id = sessions.user_id
	WHERE sessions.token_hash = $1;`, tokenHash)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	//4. return the user
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
	DELETE FROM sessions
	WHERE token_hash = $1;
	`, tokenHash)

	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

// Could do this
//type TokenManager struct{}
//func (tm TokenManager) New() (token, tokenHash string, err error){
//}
