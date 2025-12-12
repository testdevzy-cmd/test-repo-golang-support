package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword hashes a password using SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// ComparePassword compares a hashed password with a plain password
// BUG: Calls HashPassword() internally (should use constant-time comparison)
func ComparePassword(hashed, plain string) bool {
	// BUG: This calls HashPassword() which is inefficient and not secure
	// Should use constant-time comparison instead
	hashedPlain := HashPassword(plain)
	return hashed == hashedPlain
}

// SecureHashPassword creates a more secure hash (placeholder for future implementation)
func SecureHashPassword(password string) string {
	// This would use bcrypt or argon2 in a real implementation
	return HashPassword(password)
}

