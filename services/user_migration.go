package services

import (
	"context"
	"errors"

	"github.com/test-repo-golang-support/models"
)

// UserMigrationService handles migration from old User to UserRefactored
// This file demonstrates knowledge graph traversal:
// - References models.User (old) and models.UserRefactored (new)
// - Should detect that other files still use old User.Email
type UserMigrationService struct {
	oldUsers map[string]*models.User
	newUsers map[string]*models.UserRefactored
}

// NewUserMigrationService creates a migration service
func NewUserMigrationService() *UserMigrationService {
	return &UserMigrationService{
		oldUsers: make(map[string]*models.User),
		newUsers: make(map[string]*models.UserRefactored),
	}
}

// MigrateUser migrates an old user to the new format
// BUG: This method accesses user.Email which was renamed to EmailAddress
// Knowledge graph should detect this relationship and flag the issue
func (s *UserMigrationService) MigrateUser(ctx context.Context, oldUser *models.User) (*models.UserRefactored, error) {
	if oldUser == nil {
		return nil, errors.New("user is nil")
	}

	// BUG: Accessing oldUser.Email - field was renamed to EmailAddress in UserRefactored
	// Knowledge graph should detect this cross-file relationship issue
	newUser := &models.UserRefactored{
		BaseEntity: models.BaseEntity{
			ID:        oldUser.ID,
			CreatedAt: oldUser.CreatedAt,
			UpdatedAt: oldUser.UpdatedAt,
		},
		Timestamps: models.Timestamps{
			DeletedAt: oldUser.DeletedAt,
		},
		FirstName:    oldUser.FirstName,
		LastName:      oldUser.LastName,
		EmailAddress:  oldUser.Email, // BUG: Should detect that Email field exists but EmailAddress doesn't on old User
		Role:          oldUser.Role,
		Active:        oldUser.Active,
	}

	s.newUsers[oldUser.ID] = newUser
	return newUser, nil
}

// FindUserByEmailAddress finds a user by email address
// BUG: This method calls FindByEmail which still uses old User.Email field
// Knowledge graph should detect that FindByEmail in service.go uses user.Email
// but we're now working with EmailAddress
func (s *UserMigrationService) FindUserByEmailAddress(ctx context.Context, email string) (*models.UserRefactored, error) {
	// BUG: FindByEmail uses old User.Email, but we need EmailAddress
	// Knowledge graph should detect this relationship mismatch
	oldUser, err := s.FindByEmail(ctx, email) // This method uses user.Email internally
	if err != nil {
		return nil, err
	}

	// Migrate to new format
	return s.MigrateUser(ctx, oldUser)
}

// FindByEmail is a helper that would use the old UserService
// This demonstrates that the knowledge graph should trace through
// service relationships to find all affected code
func (s *UserMigrationService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range s.oldUsers {
		if user.Email == email { // BUG: This field access should be flagged
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// UpdateUserEmail updates user email using new method signature
// BUG: Calls UpdateEmailAddress with wrong number of arguments
// Old UpdateEmail() took 1 arg, new UpdateEmailAddress() takes 2 args
// Knowledge graph should detect this method signature change
func (s *UserMigrationService) UpdateUserEmail(ctx context.Context, userID, email string) error {
	user, exists := s.newUsers[userID]
	if !exists {
		return errors.New("user not found")
	}

	// BUG: UpdateEmailAddress requires 2 args (email, verified) but only passing 1
	// Knowledge graph should detect method signature mismatch
	user.UpdateEmailAddress(email) // Should be: user.UpdateEmailAddress(email, false)
	return nil
}

