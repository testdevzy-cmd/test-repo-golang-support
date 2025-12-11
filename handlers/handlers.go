package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/test-repo-golang-support/models"
	"github.com/test-repo-golang-support/services"
)

// Handler wraps the user service and provides HTTP handlers
type Handler struct {
	service *services.UserService
	logger  *log.Logger
}

// NewHandler creates a new Handler instance
func NewHandler(service *services.UserService, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// =====================================
// HTTP Handlers
// =====================================

// GetUsers handles GET /users - returns all users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.service.ReadAll(ctx)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// GetUser handles GET /users/{id} - returns a specific user
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.Read(ctx, id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser handles POST /users - creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate email
	if !services.ValidateEmail(input.Email) {
		h.respondError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Create new user
	userID := services.GenerateUserID()
	user := services.CreateUser(userID, input.FirstName, input.LastName, input.Email)
	user.SetRole(input.Role)

	// Save user
	if err := h.service.Write(ctx, user); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.logger.Printf("Created user: %s (%s)", user.FullName(), user.ID)

	h.respondJSON(w, http.StatusCreated, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser handles PUT /users/{id} - updates an existing user
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	// Get existing user
	user, err := h.service.Read(ctx, id)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update user using pointer receiver methods
	if input.FirstName != "" || input.LastName != "" {
		user.UpdateName(input.FirstName, input.LastName)
	}
	if input.Email != "" {
		user.UpdateEmail(input.Email)
	}
	if input.Role != "" {
		user.SetRole(input.Role)
	}

	// Save updated user
	if err := h.service.Write(ctx, user); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser handles DELETE /users/{id} - deletes a user
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if user exists
	exists, _ := h.service.Exists(ctx, id)
	if !exists {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	h.respondJSON(w, http.StatusOK, models.APIResponse{
		Code:    models.ResponseOK,
		Message: "User deleted successfully",
	})
}

// HealthCheck handles GET /health - returns server health status
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// =====================================
// Helper Methods
// =====================================

// respondJSON sends a JSON response (pointer receiver)
func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("Error encoding response: %v", err)
	}
}

// respondError sends an error response (pointer receiver)
func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, models.APIResponse{
		Code:    models.ResponseError,
		Message: message,
	})
}

// =====================================
// Middleware Functions
// =====================================

// LoggingMiddleware logs incoming requests
func LoggingMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Printf("Started %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
			logger.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("Panic recovered: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// =====================================
// Router Setup
// =====================================

// SetupRoutes configures all routes for the application
func SetupRoutes(h *Handler, logger *log.Logger) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(CORSMiddleware)
	router.Use(LoggingMiddleware(logger))
	router.Use(RecoveryMiddleware(logger))

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", h.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	api.HandleFunc("/users", h.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")

	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Test Server - API v1")
	}).Methods("GET")

	return router
}

