package auth

// TokenValidator interface for token validation
type TokenValidator interface {
	// Validate validates a token and returns true if valid
	Validate(token string) bool
}

// JWTValidator implements token validation using JWT
// BUG: Method signature doesn't match interface - returns (bool, error) instead of just bool
type JWTValidator struct {
	secretKey string
}

// Validate validates a JWT token
// BUG: Interface expects bool return, but this returns (bool, error)
func (j *JWTValidator) Validate(token string) (bool, error) {
	if token == "" {
		return false, nil
	}
	
	// JWT validation logic would go here
	return true, nil
}

// NewValidator creates a new JWTValidator instance
func NewValidator(secretKey string) TokenValidator {
	return &JWTValidator{
		secretKey: secretKey,
	}
}

// SetSecretKey sets the secret key for validation
func (j *JWTValidator) SetSecretKey(secretKey string) {
	j.secretKey = secretKey
}

