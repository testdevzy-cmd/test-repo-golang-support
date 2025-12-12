package commands

import (
	"context"

	"github.com/test-repo-golang-support/internal/core/application/services"
	"github.com/test-repo-golang-support/internal/core/domain/entities"
	"github.com/test-repo-golang-support/internal/core/domain/valueobjects"
)

// CreateAccountCommand represents a command to create an account
// Command pattern - knowledge graph should track command -> service -> repository -> entity chain
type CreateAccountCommand struct {
	OwnerID     string
	Email       string
	AccountType entities.AccountType
}

// CreateAccountHandler handles CreateAccountCommand
// Knowledge graph path: CreateAccountHandler -> AccountService -> AccountRepository -> Account
type CreateAccountHandler struct {
	accountService *services.AccountService
}

// NewCreateAccountHandler creates a new handler
func NewCreateAccountHandler(accountService *services.AccountService) *CreateAccountHandler {
	return &CreateAccountHandler{
		accountService: accountService,
	}
}

// Handle executes the create account command
// BUG: Uses Email field which might be renamed in entity layer
// Knowledge graph should track 4-level deep relationship:
// CreateAccountHandler -> AccountService.CreateAccount -> AccountRepository.FindByEmail -> Account.Email
func (h *CreateAccountHandler) Handle(ctx context.Context, cmd CreateAccountCommand) (*entities.Account, error) {
	return h.accountService.CreateAccount(ctx, cmd.OwnerID, cmd.Email, cmd.AccountType)
}

// DepositCommand represents a command to deposit money
type DepositCommand struct {
	AccountID string
	Amount    float64
	Currency  valueobjects.Currency
}

// DepositHandler handles deposit commands
// Multi-layer relationship:
// DepositHandler -> TransactionService -> TransactionRepository + AccountRepository
type DepositHandler struct {
	txService *services.TransactionService
}

// NewDepositHandler creates a new deposit handler
func NewDepositHandler(txService *services.TransactionService) *DepositHandler {
	return &DepositHandler{
		txService: txService,
	}
}

// Handle executes the deposit command
// Knowledge graph should track value object usage:
// DepositHandler -> valueobjects.NewMoney -> TransactionService.CreateDeposit
func (h *DepositHandler) Handle(ctx context.Context, cmd DepositCommand) (*entities.Transaction, error) {
	money := valueobjects.NewMoney(cmd.Amount, cmd.Currency)
	return h.txService.CreateDeposit(ctx, cmd.AccountID, money)
}

// TransferCommand represents a command to transfer money
type TransferCommand struct {
	SourceAccountID string
	TargetAccountID string
	Amount          float64
	Currency        valueobjects.Currency
}

// TransferHandler handles transfer commands
// Most complex relationship chain:
// TransferHandler -> TransactionService -> (TransactionRepo + AccountRepo) -> (Transaction + 2x Account)
type TransferHandler struct {
	txService *services.TransactionService
}

// NewTransferHandler creates a new transfer handler
func NewTransferHandler(txService *services.TransactionService) *TransferHandler {
	return &TransferHandler{
		txService: txService,
	}
}

// Handle executes the transfer command
func (h *TransferHandler) Handle(ctx context.Context, cmd TransferCommand) (*entities.Transaction, error) {
	money := valueobjects.NewMoney(cmd.Amount, cmd.Currency)
	return h.txService.CreateTransfer(ctx, cmd.SourceAccountID, cmd.TargetAccountID, money)
}

// SuspendAccountCommand represents a command to suspend an account
type SuspendAccountCommand struct {
	AccountID string
	Reason    string
}

// SuspendAccountHandler handles suspend account commands
type SuspendAccountHandler struct {
	accountService *services.AccountService
}

// NewSuspendAccountHandler creates a new handler
func NewSuspendAccountHandler(accountService *services.AccountService) *SuspendAccountHandler {
	return &SuspendAccountHandler{
		accountService: accountService,
	}
}

// Handle executes the suspend account command
// Knowledge graph path: SuspendAccountHandler -> AccountService.SuspendAccount -> Account.Suspend
func (h *SuspendAccountHandler) Handle(ctx context.Context, cmd SuspendAccountCommand) error {
	return h.accountService.SuspendAccount(ctx, cmd.AccountID)
}

