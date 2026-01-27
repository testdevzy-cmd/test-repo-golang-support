package services

import (
	"context"
	"errors"

	"github.com/test-repo-golang-support/internal/core/domain/entities"
	"github.com/test-repo-golang-support/internal/core/domain/repositories"
	"github.com/test-repo-golang-support/internal/core/domain/valueobjects"
)

// TransactionService handles transaction-related business logic
// Multi-layer relationship:
// TransactionService -> TransactionRepository -> Transaction
// TransactionService -> AccountService -> AccountRepository -> Account
type TransactionService struct {
	txRepo      repositories.TransactionRepository
	accountRepo repositories.AccountRepository
}

// NewTransactionService creates a new TransactionService
func NewTransactionService(txRepo repositories.TransactionRepository, accountRepo repositories.AccountRepository) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		accountRepo: accountRepo,
	}
}

// GetTransaction retrieves a transaction by ID
func (s *TransactionService) GetTransaction(ctx context.Context, id string) (*entities.Transaction, error) {
	return s.txRepo.FindByID(ctx, id)
}

// GetAccountTransactions retrieves all transactions for an account
// Knowledge graph path: TransactionService -> TransactionRepository.FindByAccountID -> Transaction.AccountID -> Account.ID
func (s *TransactionService) GetAccountTransactions(ctx context.Context, accountID string) ([]*entities.Transaction, error) {
	// Verify account exists
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if !account.IsActive() {
		return nil, errors.New("account is not active")
	}

	return s.txRepo.FindByAccountID(ctx, accountID)
}

// CreateDeposit creates a deposit transaction
// Multi-layer relationship involving:
// - Transaction entity
// - Account entity (balance update)
// - Money value object
func (s *TransactionService) CreateDeposit(ctx context.Context, accountID string, amount valueobjects.Money) (*entities.Transaction, error) {
	// Validate amount
	if !amount.IsPositive() {
		return nil, errors.New("deposit amount must be positive")
	}

	// Get account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Check if account is active
	// Knowledge graph: TransactionService -> Account.IsActive
	if !account.IsActive() {
		return nil, errors.New("account is not active")
	}

	// Create transaction
	tx := entities.NewTransaction(generateTransactionID(), accountID, entities.TransactionTypeDeposit, amount.Amount)

	// Update account balance
	// BUG: Gets balance as float but Money expects structured value
	// Knowledge graph should detect: Account.Balance (float64) vs Money.Amount
	newBalance := account.GetBalance() + amount.Amount
	account.UpdateBalance(newBalance)

	// Save transaction and account
	if err := s.txRepo.Save(ctx, tx); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	tx.Complete()
	return tx, nil
}

// CreateWithdrawal creates a withdrawal transaction
func (s *TransactionService) CreateWithdrawal(ctx context.Context, accountID string, amount valueobjects.Money) (*entities.Transaction, error) {
	if !amount.IsPositive() {
		return nil, errors.New("withdrawal amount must be positive")
	}

	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// BUG: Comparing float64 (account.GetBalance()) with Money.Amount directly
	// Knowledge graph should detect type relationship issues
	if account.GetBalance() < amount.Amount {
		return nil, errors.New("insufficient balance")
	}

	tx := entities.NewTransaction(generateTransactionID(), accountID, entities.TransactionTypeWithdrawal, amount.Amount)

	newBalance := account.GetBalance() - amount.Amount
	account.UpdateBalance(newBalance)

	if err := s.txRepo.Save(ctx, tx); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	tx.Complete()
	return tx, nil
}

// CreateTransfer creates a transfer between accounts
// Complex multi-layer relationship:
// TransactionService -> 2x Account -> 2x Transaction
func (s *TransactionService) CreateTransfer(ctx context.Context, sourceAccountID, targetAccountID string, amount valueobjects.Money) (*entities.Transaction, error) {
	if sourceAccountID == targetAccountID {
		return nil, errors.New("cannot transfer to same account")
	}

	// Get source account
	sourceAccount, err := s.accountRepo.FindByID(ctx, sourceAccountID)
	if err != nil {
		return nil, errors.New("source account not found")
	}

	// Get target account
	targetAccount, err := s.accountRepo.FindByID(ctx, targetAccountID)
	if err != nil {
		return nil, errors.New("target account not found")
	}

	// Validate both accounts are active
	if !sourceAccount.IsActive() || !targetAccount.IsActive() {
		return nil, errors.New("both accounts must be active")
	}

	// Check sufficient balance
	if sourceAccount.GetBalance() < amount.Amount {
		return nil, errors.New("insufficient balance")
	}

	// Create transfer transaction
	tx := entities.NewTransaction(generateTransactionID(), sourceAccountID, entities.TransactionTypeTransfer, amount.Amount)
	tx.SourceAccountID = sourceAccountID
	tx.TargetAccountID = targetAccountID

	// Update balances
	sourceAccount.UpdateBalance(sourceAccount.GetBalance() - amount.Amount)
	targetAccount.UpdateBalance(targetAccount.GetBalance() + amount.Amount)

	// Save all changes
	if err := s.txRepo.Save(ctx, tx); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Save(ctx, sourceAccount); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Save(ctx, targetAccount); err != nil {
		return nil, err
	}

	tx.Complete()
	return tx, nil
}

// ProcessPendingTransactions processes all pending transactions
// Knowledge graph should track: TransactionService -> TransactionRepository.GetPendingTransactions -> Transaction.IsPending
func (s *TransactionService) ProcessPendingTransactions(ctx context.Context) error {
	pendingTxs, err := s.txRepo.GetPendingTransactions(ctx)
	if err != nil {
		return err
	}

	for _, tx := range pendingTxs {
		// Process each transaction
		// Knowledge graph: calls Transaction.Complete or Transaction.Fail
		if tx.IsPending() {
			tx.Complete()
			if err := s.txRepo.Save(ctx, tx); err != nil {
				tx.Fail()
				s.txRepo.Save(ctx, tx)
			}
		}
	}

	return nil
}

func generateTransactionID() string {
	return "tx_" + "12345" // Simplified for demo
}

