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

	// Seed some initial data
	seedData(userService)

	// Initialize handlers
	handler := handlers.NewHandler(userService, logger)

	// Setup routes
	router := handlers.SetupRoutes(handler, logger)

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

// seedData adds some initial test users
func seedData(svc *services.UserService) {
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
		_ = svc.Write(ctx, user)
	}
}

// init function runs before main
func init() {
	log.Println("Initializing Go Test Server...")
}

