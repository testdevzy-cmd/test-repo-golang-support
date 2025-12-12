package models

import (
	"time"
)

// =====================================
// REFACTORING: User model changes
// These changes break relationships with other files
// =====================================

// UserRefactored represents a refactored user model
// BUG PATTERN 1: Field renamed from Email to EmailAddress
// This breaks: services/service.go (line 119), handlers/handlers.go (lines 93, 141)
// Knowledge graph should detect these cross-file references
type UserRefactored struct {
	BaseEntity
	Timestamps
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"` // Changed from Email
	Role         string `json:"role"`
	Active       bool   `json:"active"`
}

// UpdateEmailAddress updates the user's email address
// BUG PATTERN 2: Method signature changed - now requires additional parameter
// Old: UpdateEmail(email string)
// New: UpdateEmailAddress(email string, verified bool)
// This breaks: handlers/handlers.go (line 141) which calls UpdateEmail()
func (u *UserRefactored) UpdateEmailAddress(email string, verified bool) {
	u.EmailAddress = email
	u.UpdatedAt = time.Now()
	// verified parameter added but not used - breaks existing callers
}

// =====================================
// REFACTORING: Type alias changes
// =====================================

// UserIdentifier is a new type definition replacing UserID alias
// BUG PATTERN 3: Changed from type alias (UserID = string) to type definition
// This breaks: services/service.go, handlers/org_handlers.go, models/models.go
// Knowledge graph should detect all usages of UserID
type UserIdentifier string

// ConvertUserID converts old UserID to new UserIdentifier
// BUG: Uses UserID type alias which is defined in models.go
// Knowledge graph should detect this type relationship
func ConvertUserID(oldID string) UserIdentifier {
	return UserIdentifier(oldID)
}

// =====================================
// REFACTORING: Interface method signature change
// =====================================

// UserValidator interface for user validation
// BUG PATTERN 4: Validate method signature changed
// Old: Validate() error
// New: Validate(ctx context.Context) (bool, error)
// This breaks: models/models.go User.Validate() implementation
type UserValidator interface {
	Validate(ctx interface{}) (bool, error) // Changed signature - requires context
}

// ValidateUserRefactored validates a refactored user
// This should implement UserValidator but signature doesn't match
func (u *UserRefactored) ValidateUserRefactored() error {
	if u.ID == "" {
		return nil // Missing error return
	}
	if u.EmailAddress == "" {
		return nil // Missing error return
	}
	return nil
}

// =====================================
// REFACTORING: Method removed/changed
// =====================================

// GetEmail is a new method replacing direct Email field access
// BUG PATTERN 5: Old code accesses user.Email directly
// New code should use user.GetEmail() but old code still uses .Email
// This breaks: services/service.go FindByEmail, models/models.go String()
func (u UserRefactored) GetEmail() string {
	return u.EmailAddress
}

// SetEmail is a new method for setting email
// Old code uses user.Email = value, new code should use SetEmail()
func (u *UserRefactored) SetEmail(email string) {
	u.EmailAddress = email
	u.UpdatedAt = time.Now()
}

