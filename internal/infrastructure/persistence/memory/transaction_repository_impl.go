package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/test-repo-golang-support/internal/core/domain/entities"
	"github.com/test-repo-golang-support/internal/core/domain/repositories"
)

// TransactionRepositoryImpl implements TransactionRepository using in-memory storage
// Knowledge graph: TransactionRepositoryImpl implements repositories.TransactionRepository
type TransactionRepositoryImpl struct {
	transactions map[string]*entities.Transaction
	mu           sync.RWMutex
}

// Ensure interface compliance
var _ repositories.TransactionRepository = (*TransactionRepositoryImpl)(nil)

// NewTransactionRepository creates a new in-memory transaction repository
func NewTransactionRepository() *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		transactions: make(map[string]*entities.Transaction),
	}
}

// FindByID finds a transaction by ID
func (r *TransactionRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, exists := r.transactions[id]
	if !exists {
		return nil, errors.New("transaction not found")
	}
	return tx, nil
}

// FindByAccountID finds transactions by account ID
// Knowledge graph should track:
// TransactionRepositoryImpl.FindByAccountID -> Transaction.AccountID relationship
func (r *TransactionRepositoryImpl) FindByAccountID(ctx context.Context, accountID string) ([]*entities.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Transaction
	for _, tx := range r.transactions {
		// BUG: Accessing tx.AccountID - if field is renamed, breaks here
		if tx.AccountID == accountID {
			result = append(result, tx)
		}
	}
	return result, nil
}

// Save saves a transaction
func (r *TransactionRepositoryImpl) Save(ctx context.Context, tx *entities.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tx.ID == "" {
		return errors.New("transaction ID is required")
	}
	r.transactions[tx.ID] = tx
	return nil
}

// GetPendingTransactions gets all pending transactions
// Knowledge graph should track:
// TransactionRepositoryImpl.GetPendingTransactions -> Transaction.IsPending() method call
func (r *TransactionRepositoryImpl) GetPendingTransactions(ctx context.Context) ([]*entities.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Transaction
	for _, tx := range r.transactions {
		// Calls entity method - knowledge graph should track this
		if tx.IsPending() {
			result = append(result, tx)
		}
	}
	return result, nil
}

