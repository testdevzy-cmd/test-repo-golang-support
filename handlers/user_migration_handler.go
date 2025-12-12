package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/test-repo-golang-support/models"
	"github.com/test-repo-golang-support/services"
)

// UserMigrationHandler handles user migration endpoints
// This demonstrates knowledge graph traversal:
// - Uses models.UserRefactored (new)
// - Should detect that other handlers still use old models.User
type UserMigrationHandler struct {
	migrationService *services.UserMigrationService
	logger           *log.Logger
}

// NewUserMigrationHandler creates a new migration handler
func NewUserMigrationHandler(service *services.UserMigrationService, logger *log.Logger) *UserMigrationHandler {
	return &UserMigrationHandler{
		migrationService: service,
		logger:           logger,
	}
}

// MigrateUser handles POST /migrate/users/{id} - migrates a user
// BUG: This handler uses UserRefactored but the main handlers still use User
// Knowledge graph should detect that both User and UserRefactored are used
// and flag potential inconsistencies
func (h *UserMigrationHandler) MigrateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// BUG: Creating UserRefactored but other parts of codebase use User
	// Knowledge graph should detect this type mismatch across files
	newUser := &models.UserRefactored{
		BaseEntity: models.BaseEntity{
			ID: userID,
		},
		EmailAddress: "new@example.com", // Using new field name
	}

	// BUG: Accessing EmailAddress directly - but if this was passed from
	// old User model, it would have Email field, not EmailAddress
	// Knowledge graph should detect field name mismatch
	h.logger.Printf("Migrated user: %s with email: %s", newUser.ID, newUser.EmailAddress)

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "User migrated successfully",
		Data:    newUser,
	})
}

// CreateRefactoredUser handles POST /users/refactored - creates a refactored user
// BUG: This creates UserRefactored but the main CreateUser handler creates User
// Knowledge graph should detect that two different user types are being created
// and flag potential data inconsistency
func (h *UserMigrationHandler) CreateRefactoredUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"` // JSON field is still "email" but struct field is EmailAddress
		Role      string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// BUG: Using UpdateEmailAddress with wrong signature
	// Should pass 2 args: (email, verified) but only passing email
	// Knowledge graph should detect method signature mismatch
	user := &models.UserRefactored{
		BaseEntity: models.BaseEntity{
			ID: services.GenerateUserID(),
		},
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	user.UpdateEmailAddress(input.Email) // BUG: Missing second parameter (verified bool)

	h.respondJSON(w, http.StatusCreated, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Refactored user created",
		Data:    user,
	})
}

// GetUserEmail handles GET /users/{id}/email - gets user email
// BUG: This method might receive old User type but tries to access EmailAddress
// Knowledge graph should detect type mismatch
func (h *UserMigrationHandler) GetUserEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// BUG: Assuming UserRefactored but might receive old User type
	// Knowledge graph should detect that User.Email and UserRefactored.EmailAddress
	// are related but different field names
	user := &models.UserRefactored{
		BaseEntity: models.BaseEntity{ID: userID},
	}

	// BUG: If this was an old User, it would have .Email, not .EmailAddress
	// Knowledge graph should detect this field access mismatch
	email := user.EmailAddress // Would fail if user was old User type with Email field

	h.respondJSON(w, http.StatusOK, map[string]string{
		"email": email,
	})
}

// Helper methods
func (h *UserMigrationHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *UserMigrationHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, models.APIResponse{
		Code:    models.ResponseError,
		Message: message,
	})
}

