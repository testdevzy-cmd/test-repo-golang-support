package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/test-repo-golang-support/handlers"
	"github.com/test-repo-golang-support/services"
)

const (
	defaultPort    = "8081"
	defaultTimeout = 15 * time.Second
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize services
	userService := services.NewUserService()
	orgService := services.NewOrganizationService()

	// Seed some initial data
	seedData(userService, orgService)

	// Initialize handlers
	handler := handlers.NewHandler(userService, logger)
	orgHandler := handlers.NewOrgHandler(orgService, logger)

	// Setup routes
	router := handlers.SetupRoutes(handler, logger)

	// Setup organization routes
	api := router.PathPrefix("/api/v1").Subrouter()
	handlers.SetupOrgRoutes(api, orgHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		logger.Printf("Starting server on port %s", port)
		logger.Printf("API endpoints available at http://localhost:%s/api/v1", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-shutdown
	logger.Println("Shutdown signal received, gracefully shutting down...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server stopped gracefully")
}

// seedData adds some initial test users and organizations
func seedData(userSvc *services.UserService, orgSvc *services.OrganizationService) {
	ctx := context.Background()

	// Create test users
	users := []*struct {
		firstName string
		lastName  string
		email     string
		role      string
	}{
		{"John", "Doe", "john.doe@example.com", "admin"},
		{"Jane", "Smith", "jane.smith@example.com", "user"},
		{"Bob", "Wilson", "bob.wilson@example.com", "user"},
	}

	for i, u := range users {
		user := services.CreateUser(
			fmt.Sprintf("user_%d", i+1),
			u.firstName,
			u.lastName,
			u.email,
		)
		user.SetRole(u.role)
		_ = userSvc.Write(ctx, user)
	}

	// Create test organizations
	orgs := []*struct {
		name     string
		industry string
		ownerID  string
	}{
		{"Acme Corp", "Technology", "user_1"},
		{"Global Industries", "Manufacturing", "user_2"},
	}

	for i, o := range orgs {
		org := services.CreateOrganization(
			fmt.Sprintf("org_%d", i+1),
			o.name,
			o.ownerID,
		)
		org.SetIndustry(o.industry)
		_ = orgSvc.WriteOrg(ctx, org)

		// Add owner as member
		membership := services.CreateMembership(o.ownerID, org.ID, "owner")
		_ = orgSvc.AddMember(ctx, membership)
	}
}

// init function runs before main
func init() {
	log.Println("Initializing Go Test Server...")
}

