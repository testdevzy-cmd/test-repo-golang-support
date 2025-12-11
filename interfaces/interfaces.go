package interfaces

import (
	"context"

	"github.com/test-repo-golang-support/models"
)

// Reader interface for read operations
type Reader interface {
	Read(ctx context.Context, id string) (*models.User, error)
	ReadAll(ctx context.Context) (models.UserList, error)
}

// Writer interface for write operations
type Writer interface {
	Write(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

// ReadWriter combines Reader and Writer interfaces
// Demonstrates interface embedding
type ReadWriter interface {
	Reader // Embedded interface
	Writer // Embedded interface
}

// Validator interface for validation operations
type Validator interface {
	Validate() error
}

// Serializer interface for serialization
type Serializer interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

// Entity combines multiple interfaces
// Demonstrates multiple interface embedding
type Entity interface {
	Validator  // Embedded interface
	Serializer // Embedded interface
}

// Repository is a comprehensive interface for data access
// Embeds ReadWriter and adds additional methods
type Repository interface {
	ReadWriter // Embedded composite interface
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
}

// Logger interface for logging operations
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// AuditLogger extends Logger with audit capabilities
type AuditLogger interface {
	Logger // Embedded interface
	Audit(action string, userID string, details map[string]interface{})
}

// Service interface combines repository access with logging
type Service interface {
	Repository  // Embedded interface
	AuditLogger // Embedded interface
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// EventEmitter interface for event-driven architecture
type EventEmitter interface {
	Emit(event string, data interface{}) error
	Subscribe(event string, handler func(data interface{})) error
}

// EventDrivenService combines Service with event capabilities
type EventDrivenService interface {
	Service      // Embedded composite interface
	EventEmitter // Embedded interface
}

