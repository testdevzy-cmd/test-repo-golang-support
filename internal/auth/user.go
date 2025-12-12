package auth

import (
	"errors"

	"github.com/test-repo-golang-support/pkg/utils"
)

// User represents an authenticated user
type User struct {
	ID           string
	Email        string
	PasswordHash string
}

// GetUserByEmail retrieves a user by email
// BUG: Calls ValidateToken with wrong number of arguments (2 instead of 1)
// BUG: Imported with alias 'auth' but calling methods without prefix
func GetUserByEmail(email string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	
	// BUG: Imported as 'auth' but not using the prefix
	// BUG: ValidateToken expects 1 argument but passing 2
	authenticator := NewAuthenticator("secret", 3600)
	valid, err := authenticator.ValidateToken("token", email)
	if err != nil {
		return nil, err
	}
	
	if !valid {
		return nil, errors.New("authentication failed")
	}
	
	// This would normally query a database
	return &User{
		ID:    "user_123",
		Email: email,
	}, nil
}

// AuthenticateUser authenticates a user with email and password
func AuthenticateUser(email, password string) (*User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	
	// BUG: Calls HashPasswords() (plural) which doesn't exist - should be HashPassword()
	hashedPassword := utils.HashPasswords(password)
	
	// Password validation would happen here
	if user.PasswordHash != hashedPassword {
		return nil, errors.New("invalid password")
	}
	
	return user, nil
}

