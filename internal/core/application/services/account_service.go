package services

import (
	"context"
	"errors"

	"github.com/test-repo-golang-support/internal/core/domain/entities"
	"github.com/test-repo-golang-support/internal/core/domain/repositories"
	"github.com/test-repo-golang-support/internal/core/domain/valueobjects"
)

// AccountService handles account-related business logic
// This is at the application layer, referencing domain layer
// Knowledge graph should track: AccountService -> AccountRepository -> Account
type AccountService struct {
	repo repositories.AccountRepository
}

// NewAccountService creates a new AccountService
func NewAccountService(repo repositories.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

// GetAccount retrieves an account by ID
// Knowledge graph path: AccountService.GetAccount -> AccountRepository.FindByID -> Account
func (s *AccountService) GetAccount(ctx context.Context, id string) (*entities.Account, error) {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountByEmail retrieves an account by email
// BUG: Calls repo.FindByEmail which uses Account.Email field
// If Account.Email is renamed, this chain breaks
// Knowledge graph should detect: AccountService -> AccountRepository.FindByEmail -> Account.Email
func (s *AccountService) GetAccountByEmail(ctx context.Context, email string) (*entities.Account, error) {
	return s.repo.FindByEmail(ctx, email)
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(ctx context.Context, ownerID, email string, accountType entities.AccountType) (*entities.Account, error) {
	// Check if account with email already exists
	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("account with email already exists")
	}

	account := entities.NewAccount(generateAccountID(), ownerID, email, accountType)
	if err := s.repo.Save(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountBalance gets the account balance as Money value object
// Knowledge graph path: AccountService -> AccountRepository.GetBalance -> Account.GetBalance()
// Also references: valueobjects.Money
func (s *AccountService) GetAccountBalance(ctx context.Context, id string) (valueobjects.Money, error) {
	balance, err := s.repo.GetBalance(ctx, id)
	if err != nil {
		return valueobjects.Money{}, err
	}
	return valueobjects.NewMoney(balance, valueobjects.CurrencyUSD), nil
}

// UpdateAccountBalance updates the account balance
// BUG: Calls account.UpdateBalance but might be renamed to SetBalance
// Knowledge graph should track: AccountService -> Account.UpdateBalance
func (s *AccountService) UpdateAccountBalance(ctx context.Context, id string, amount float64) error {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	account.UpdateBalance(amount)
	return s.repo.Save(ctx, account)
}

// SuspendAccount suspends an account
// Knowledge graph should track: AccountService.SuspendAccount -> Account.Suspend
func (s *AccountService) SuspendAccount(ctx context.Context, id string) error {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	account.Suspend()
	return s.repo.Save(ctx, account)
}

func generateAccountID() string {
	return "acc_" + "12345" // Simplified for demo
}

