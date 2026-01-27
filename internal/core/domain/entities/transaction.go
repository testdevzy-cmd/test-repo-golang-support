package entities

import (
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Transaction represents a financial transaction
// References Account entity - knowledge graph should track this relationship
type Transaction struct {
	ID              string
	AccountID       string            // References Account.ID
	Type            TransactionType
	Status          TransactionStatus
	Amount          float64
	Description     string
	SourceAccountID string            // For transfers
	TargetAccountID string            // For transfers
	CreatedAt       time.Time
	ProcessedAt     *time.Time
}

// NewTransaction creates a new Transaction
func NewTransaction(id, accountID string, txType TransactionType, amount float64) *Transaction {
	return &Transaction{
		ID:        id,
		AccountID: accountID,
		Type:      txType,
		Status:    TransactionStatusPending,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
}

// IsPending checks if transaction is pending (value receiver)
func (t Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// IsCompleted checks if transaction is completed (value receiver)
func (t Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// Complete marks transaction as completed (pointer receiver)
func (t *Transaction) Complete() {
	t.Status = TransactionStatusCompleted
	now := time.Now()
	t.ProcessedAt = &now
}

// Fail marks transaction as failed (pointer receiver)
func (t *Transaction) Fail() {
	t.Status = TransactionStatusFailed
	now := time.Now()
	t.ProcessedAt = &now
}

