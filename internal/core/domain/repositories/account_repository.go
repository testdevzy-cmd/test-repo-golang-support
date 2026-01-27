package repositories

import (
	"context"

	"github.com/test-repo-golang-support/internal/core/domain/entities"
)

// AccountRepository defines the interface for account data access
// This interface is at the repository layer and will be implemented at infrastructure layer
type AccountRepository interface {
	// FindByID finds an account by ID
	FindByID(ctx context.Context, id string) (*entities.Account, error)
	
	// FindByOwnerID finds accounts by owner ID
	FindByOwnerID(ctx context.Context, ownerID string) ([]*entities.Account, error)
	
	// FindByEmail finds an account by email
	// BUG PATTERN: This method uses Email field which might be renamed
	// Knowledge graph should track this relationship to entities.Account.Email
	FindByEmail(ctx context.Context, email string) (*entities.Account, error)
	
	// Save saves an account
	Save(ctx context.Context, account *entities.Account) error
	
	// Delete deletes an account
	Delete(ctx context.Context, id string) error
	
	// GetBalance gets account balance
	// Knowledge graph should track this calls Account.GetBalance()
	GetBalance(ctx context.Context, id string) (float64, error)
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// FindByID finds a transaction by ID
	FindByID(ctx context.Context, id string) (*entities.Transaction, error)
	
	// FindByAccountID finds transactions by account ID
	// Knowledge graph should track relationship: Transaction.AccountID -> Account.ID
	FindByAccountID(ctx context.Context, accountID string) ([]*entities.Transaction, error)
	
	// Save saves a transaction
	Save(ctx context.Context, tx *entities.Transaction) error
	
	// GetPendingTransactions gets all pending transactions
	// Knowledge graph should track: uses Transaction.IsPending()
	GetPendingTransactions(ctx context.Context) ([]*entities.Transaction, error)
}

// UnitOfWork defines the interface for transaction management
// Combines multiple repositories - knowledge graph should track these relationships
type UnitOfWork interface {
	AccountRepository() AccountRepository
	TransactionRepository() TransactionRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

