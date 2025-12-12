package services

import (
	"errors"

	"github.com/test-repo-golang-support/internal/auth"
	"github.com/test-repo-golang-support/models"
)

// AuthUserService handles authentication-related user operations
type AuthUserService struct {
	authenticator *auth.Authenticator
}

// NewAuthUserService creates a new AuthUserService instance
func NewAuthUserService(authenticator *auth.Authenticator) *AuthUserService {
	return &AuthUserService{
		authenticator: authenticator,
	}
}

// LoginUser logs in a user and returns a token
// BUG: Calls GenerateTokn() instead of GenerateToken() (method name typo)
func (s *AuthUserService) LoginUser(email, password string) (string, error) {
	// Get user by email
	user, err := auth.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	// BUG: Method name typo - should be GenerateToken()
	token, err := s.authenticator.GenerateTokn(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserInfo retrieves user information
// BUG: Accesses user.EmailAddress when field is actually Email
// BUG: Calls non-existent function GetUserPermissions()
func (s *AuthUserService) GetUserInfo(userID string) (map[string]interface{}, error) {
	// This would normally fetch from database
	user := &models.User{
		ID:    userID,
		Email: "user@example.com",
	}

	// BUG: Field name is Email, not EmailAddress
	email := user.EmailAddress

	// BUG: Function GetUserPermissions() doesn't exist
	permissions, err := GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          user.ID,
		"email":       email,
		"permissions": permissions,
	}, nil
}

// ValidateUserToken validates a user token
func (s *AuthUserService) ValidateUserToken(token string) (bool, error) {
	return s.authenticator.ValidateToken(token)
}

