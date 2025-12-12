package auth

import (
	"errors"
	"time"
)

// DefaultExpiry is the default token expiry time (24 hours)
const DefaultExpiry = 24 * time.Hour

// Authenticator handles authentication operations
type Authenticator struct {
	secretKey   string
	tokenExpiry time.Duration
	userStore   interface{}
}

// NewAuthenticator creates a new Authenticator instance
func NewAuthenticator(secretKey string, tokenExpiry time.Duration) *Authenticator {
	return &Authenticator{
		secretKey:   secretKey,
		tokenExpiry: tokenExpiry,
	}
}

// ValidateToken validates a token
// BUG: Calls non-existent method Validate() on token string
func (a *Authenticator) ValidateToken(token string) (bool, error) {
	// BUG: token is a string, doesn't have Validate() method
	if !token.Validate() {
		return false, errors.New("invalid token")
	}
	return true, nil
}

// GenerateToken generates a new token for a user
// BUG: References a.secret instead of a.secretKey (typo)
func (a *Authenticator) GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID is required")
	}
	
	// BUG: Should be a.secretKey, not a.secret
	token := a.secret + ":" + userID
	return token, nil
}

// SetUserStore sets the user store
func (a *Authenticator) SetUserStore(store interface{}) {
	a.userStore = store
}

