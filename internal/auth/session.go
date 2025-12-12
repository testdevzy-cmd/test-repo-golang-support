package auth

import (
	"time"
)

// Session represents a user session
// Similar concepts to User: ID, expiry, validation
// This tests embedding-based semantic search
type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	IsValid   bool
	Token     string
}

// NewSession creates a new session
func NewSession(userID, token string, expiry time.Duration) *Session {
	return &Session{
		ID:        generateSessionID(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiry),
		IsValid:   true,
		Token:     token,
	}
}

// IsExpired checks if the session has expired (value receiver)
func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsActive checks if the session is active (value receiver)
func (s Session) IsActive() bool {
	return s.IsValid && !s.IsExpired()
}

// Validate validates the session (value receiver)
func (s Session) Validate() bool {
	return s.IsActive()
}

// Invalidate invalidates the session (pointer receiver)
func (s *Session) Invalidate() {
	s.IsValid = false
}

// Extend extends the session expiry (pointer receiver)
func (s *Session) Extend(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
}

// GetUserID returns the user ID associated with the session (value receiver)
func (s Session) GetUserID() string {
	return s.UserID
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return "session_" + time.Now().Format("20060102150405")
}

