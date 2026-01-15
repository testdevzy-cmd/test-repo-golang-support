package entities

import (
	"time"
)

// AccountStatus represents account status
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusSuspended AccountStatus = "suspended"
	AccountStatusClosed   AccountStatus = "closed"
)

// AccountType represents account type
type AccountType string

const (
	AccountTypePersonal   AccountType = "personal"
	AccountTypeBusiness   AccountType = "business"
	AccountTypeEnterprise AccountType = "enterprise"
)

// Account represents a user account in the domain layer
// This entity is at the deepest level and will be referenced by multiple layers
type Account struct {
	ID          string
	OwnerID     string
	Email       string        // Field name that will be referenced by upper layers
	AccountType AccountType
	Status      AccountStatus
	Balance     float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewAccount creates a new Account entity
func NewAccount(id, ownerID, email string, accountType AccountType) *Account {
	now := time.Now()
	return &Account{
		ID:          id,
		OwnerID:     ownerID,
		Email:       email,
		AccountType: accountType,
		Status:      AccountStatusActive,
		Balance:     0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsActive checks if account is active (value receiver)
func (a Account) IsActive() bool {
	return a.Status == AccountStatusActive
}

// GetBalance returns the account balance (value receiver)
func (a Account) GetBalance() float64 {
	return a.Balance
}

// UpdateBalance updates the account balance (pointer receiver)
func (a *Account) UpdateBalance(amount float64) {
	a.Balance = amount
	a.UpdatedAt = time.Now()
}

// Suspend suspends the account (pointer receiver)
func (a *Account) Suspend() {
	a.Status = AccountStatusSuspended
	a.UpdatedAt = time.Now()
}

// Activate activates the account (pointer receiver)
func (a *Account) Activate() {
	a.Status = AccountStatusActive
	a.UpdatedAt = time.Now()
}

