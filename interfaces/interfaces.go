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

// =====================================
// New Interfaces for Incremental Testing
// =====================================

// Cacheable interface for cache operations
type Cacheable interface {
	CacheKey() string
	CacheTTL() int
}

// Invalidatable interface for cache invalidation
type Invalidatable interface {
	InvalidateCache() error
	InvalidateRelated() error
}

// CacheManager combines caching interfaces
// Demonstrates interface embedding
type CacheManager interface {
	Cacheable     // Embedded interface
	Invalidatable // Embedded interface
	WarmCache(ctx context.Context) error
}

// Notifier interface for notification operations
type Notifier interface {
	Notify(ctx context.Context, message string) error
	NotifyAsync(ctx context.Context, message string) error
}

// EmailNotifier extends Notifier for email
type EmailNotifier interface {
	Notifier // Embedded interface
	SendEmail(ctx context.Context, to, subject, body string) error
}

// PushNotifier extends Notifier for push notifications
type PushNotifier interface {
	Notifier // Embedded interface
	SendPush(ctx context.Context, deviceID, title, body string) error
}

// MultiChannelNotifier combines multiple notification channels
type MultiChannelNotifier interface {
	EmailNotifier // Embedded interface
	PushNotifier  // Embedded interface
}

// Searchable interface for search operations
type Searchable interface {
	Search(ctx context.Context, query string) ([]interface{}, error)
	SearchWithFilters(ctx context.Context, query string, filters map[string]interface{}) ([]interface{}, error)
}

// Indexable interface for search indexing
type Indexable interface {
	Index(ctx context.Context, id string, data interface{}) error
	Reindex(ctx context.Context) error
	DeleteIndex(ctx context.Context, id string) error
}

// SearchEngine combines search interfaces
type SearchEngine interface {
	Searchable // Embedded interface
	Indexable  // Embedded interface
}

// =====================================
// Organization-specific Interfaces
// =====================================

// OrgReader interface for organization read operations
type OrgReader interface {
	ReadOrg(ctx context.Context, id string) (*models.Organization, error)
	ReadAllOrgs(ctx context.Context) (models.OrgList, error)
	ReadOrgsByOwner(ctx context.Context, ownerID string) (models.OrgList, error)
}

// OrgWriter interface for organization write operations
type OrgWriter interface {
	WriteOrg(ctx context.Context, org *models.Organization) error
	DeleteOrg(ctx context.Context, id string) error
}

// OrgReadWriter combines OrgReader and OrgWriter
type OrgReadWriter interface {
	OrgReader // Embedded interface
	OrgWriter // Embedded interface
}

// OrgRepository is a comprehensive interface for organization data access
type OrgRepository interface {
	OrgReadWriter // Embedded composite interface
	CountOrgs(ctx context.Context) (int, error)
	OrgExists(ctx context.Context, id string) (bool, error)
}

// MembershipManager interface for membership operations
type MembershipManager interface {
	AddMember(ctx context.Context, membership *models.Membership) error
	RemoveMember(ctx context.Context, userID, orgID string) error
	GetMembers(ctx context.Context, orgID string) ([]*models.Membership, error)
	GetMembership(ctx context.Context, userID, orgID string) (*models.Membership, error)
	UpdateMemberRole(ctx context.Context, userID, orgID string, role models.MemberRole) error
}

// OrgService combines organization repository with membership management
type OrgService interface {
	OrgRepository     // Embedded interface
	MembershipManager // Embedded interface
}

