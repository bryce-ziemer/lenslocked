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
	// lower case email so dont get multiple emails mapped to same email (emails are not case sensitive)
	email = strings.ToLower(email)

	// hash our password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	// start filling out user object we wish to be returned
	// note user var is not a pointer
	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	// perform sql insert to insert user
	row := us.DB.QueryRow(`
	INSERT INTO users (email, password_hash)
	VALUES ($1, $2) RETURNING id;`, email, passwordHash)
	err = row.Scan(&user.ID) // fill out id field of user struct instance

	if err != nil {
		return nil, fmt.Errorf("Create User: %w", err)
	}

	return &user, nil // use '&' so user var is a pointer (specified return pointer in function declarationa nd betetr to return pointers since want to use pointers in other methods (such as an update method))
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}

	row := us.DB.QueryRow(`
	SELECT id, password_hash
	FROM users
	WHERE email=$1
	;
	`, email)

	err := row.Scan(&user.ID, &user.PasswordHash)

	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("autenticate: %w", err)
	}
	fmt.Println("Password is correct!")

	return &user, err

}
