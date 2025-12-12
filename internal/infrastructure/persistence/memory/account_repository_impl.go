package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/test-repo-golang-support/internal/core/domain/entities"
	"github.com/test-repo-golang-support/internal/core/domain/repositories"
)

// AccountRepositoryImpl implements AccountRepository using in-memory storage
// Knowledge graph should track interface implementation:
// AccountRepositoryImpl implements repositories.AccountRepository
type AccountRepositoryImpl struct {
	accounts map[string]*entities.Account
	mu       sync.RWMutex
}

// Ensure interface compliance
var _ repositories.AccountRepository = (*AccountRepositoryImpl)(nil)

// NewAccountRepository creates a new in-memory account repository
func NewAccountRepository() *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		accounts: make(map[string]*entities.Account),
	}
}

// FindByID finds an account by ID
// Knowledge graph: AccountRepositoryImpl.FindByID implements AccountRepository.FindByID
func (r *AccountRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	if !exists {
		return nil, errors.New("account not found")
	}
	return account, nil
}

// FindByOwnerID finds accounts by owner ID
func (r *AccountRepositoryImpl) FindByOwnerID(ctx context.Context, ownerID string) ([]*entities.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Account
	for _, account := range r.accounts {
		if account.OwnerID == ownerID {
			result = append(result, account)
		}
	}
	return result, nil
}

// FindByEmail finds an account by email
// BUG PATTERN: This accesses account.Email which might be renamed
// Knowledge graph should track:
// AccountRepositoryImpl.FindByEmail -> entities.Account.Email
// If Account.Email is renamed, knowledge graph should flag this
func (r *AccountRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, account := range r.accounts {
		// BUG: Direct field access to Email - if renamed, breaks here
		if account.Email == email {
			return account, nil
		}
	}
	return nil, errors.New("account not found")
}

// Save saves an account
func (r *AccountRepositoryImpl) Save(ctx context.Context, account *entities.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if account.ID == "" {
		return errors.New("account ID is required")
	}
	r.accounts[account.ID] = account
	return nil
}

// Delete deletes an account
func (r *AccountRepositoryImpl) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[id]; !exists {
		return errors.New("account not found")
	}
	delete(r.accounts, id)
	return nil
}

// GetBalance gets account balance
// Knowledge graph should track: AccountRepositoryImpl.GetBalance -> Account.GetBalance()
func (r *AccountRepositoryImpl) GetBalance(ctx context.Context, id string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	if !exists {
		return 0, errors.New("account not found")
	}
	// Calls entity method - knowledge graph should track this relationship
	return account.GetBalance(), nil
}

