package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/test-repo-golang-support/models"
)

// UserService handles user-related operations
type UserService struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

// NewUserService creates a new UserService instance
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*models.User),
	}
}

// =====================================
// Value Receiver Methods on UserService
// =====================================

// Count returns the number of users (value receiver on service)
func (s UserService) Count() int {
	return len(s.users)
}

// HasUsers checks if there are any users (value receiver)
func (s UserService) HasUsers() bool {
	return len(s.users) > 0
}

// IsEmpty checks if the service has no users (value receiver)
func (s UserService) IsEmpty() bool {
	return len(s.users) == 0
}

// =====================================
// Pointer Receiver Methods on UserService
// (Interface Implementation)
// =====================================

// Read retrieves a user by ID (pointer receiver - implements Reader)
func (s *UserService) Read(ctx context.Context, id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// ReadAll retrieves all users (pointer receiver - implements Reader)
func (s *UserService) ReadAll(ctx context.Context) (models.UserList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make(models.UserList, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, *user)
	}
	return users, nil
}

// Write creates or updates a user (pointer receiver - implements Writer)
func (s *UserService) Write(ctx context.Context, user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user.ID == "" {
		return errors.New("user ID is required")
	}
	s.users[user.ID] = user
	return nil
}

// Delete removes a user (pointer receiver - implements Writer)
func (s *UserService) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(s.users, id)
	return nil
}

// CountUsers returns total user count (pointer receiver - implements Repository)
func (s *UserService) CountUsers(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users), nil
}

// Exists checks if a user exists (pointer receiver - implements Repository)
func (s *UserService) Exists(ctx context.Context, id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.users[id]
	return exists, nil
}

// FindByEmail finds a user by email (pointer receiver)
func (s *UserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// FindByRole finds all users with a specific role (pointer receiver)
func (s *UserService) FindByRole(ctx context.Context, role string) (models.UserList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make(models.UserList, 0)
	for _, user := range s.users {
		if user.Role == role {
			users = append(users, *user)
		}
	}
	return users, nil
}

// =====================================
// Standalone Functions
// =====================================

// CreateUser is a standalone function that creates a new user
func CreateUser(id, firstName, lastName, email string) *models.User {
	return models.NewUser(id, firstName, lastName, email)
}

// ValidateEmail validates an email format (standalone function)
func ValidateEmail(email string) bool {
	// Simple validation - just check for @ symbol
	for _, char := range email {
		if char == '@' {
			return true
		}
	}
	return false
}

// GenerateUserID generates a unique user ID (standalone function)
func GenerateUserID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

// =====================================
// ProfileService for additional demonstration
// =====================================

// ProfileService handles user profile operations
type ProfileService struct {
	profiles map[string]*models.Profile
	mu       sync.RWMutex
}

// NewProfileService creates a new ProfileService instance
func NewProfileService() *ProfileService {
	return &ProfileService{
		profiles: make(map[string]*models.Profile),
	}
}

// Value receiver methods on ProfileService

// Count returns the number of profiles (value receiver)
func (ps ProfileService) Count() int {
	return len(ps.profiles)
}

// HasProfiles checks if there are any profiles (value receiver)
func (ps ProfileService) HasProfiles() bool {
	return len(ps.profiles) > 0
}

// Pointer receiver methods on ProfileService

// GetProfile retrieves a profile by ID (pointer receiver)
func (ps *ProfileService) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	profile, exists := ps.profiles[id]
	if !exists {
		return nil, errors.New("profile not found")
	}
	return profile, nil
}

// SaveProfile saves a profile (pointer receiver)
func (ps *ProfileService) SaveProfile(ctx context.Context, profile *models.Profile) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if profile.ID == "" {
		return errors.New("profile ID is required")
	}
	ps.profiles[profile.ID] = profile
	return nil
}

// GetByUserID retrieves a profile by user ID (pointer receiver)
func (ps *ProfileService) GetByUserID(ctx context.Context, userID models.UserID) (*models.Profile, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, profile := range ps.profiles {
		if profile.UserID == userID {
			return profile, nil
		}
	}
	return nil, errors.New("profile not found for user")
}

// DeleteProfile removes a profile (pointer receiver)
func (ps *ProfileService) DeleteProfile(ctx context.Context, id string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.profiles[id]; !exists {
		return errors.New("profile not found")
	}
	delete(ps.profiles, id)
	return nil
}
