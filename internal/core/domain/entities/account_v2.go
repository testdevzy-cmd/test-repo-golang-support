package entities

import (
	"time"
)

// =====================================
// BREAKING CHANGE: Account V2
// This file introduces breaking changes at the deepest layer
// Knowledge graph should detect all affected files through multi-level traversal
// =====================================

// AccountV2 is the new version of Account with breaking changes
// BUG PATTERN: Field renames that break 5+ layers above
type AccountV2 struct {
	ID           string
	OwnerID      string
	EmailAddress string        // BREAKING: Renamed from Email
	AcctType     AccountType   // BREAKING: Renamed from AccountType
	AcctStatus   AccountStatus // BREAKING: Renamed from Status
	BalanceAmt   float64       // BREAKING: Renamed from Balance
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewAccountV2 creates a new AccountV2
// BUG: Uses different field names than NewAccount
// Knowledge graph should detect callers still using old field names
func NewAccountV2(id, ownerID, email string, accountType AccountType) *AccountV2 {
	now := time.Now()
	return &AccountV2{
		ID:           id,
		OwnerID:      ownerID,
		EmailAddress: email,      // Changed from Email
		AcctType:     accountType, // Changed from AccountType
		AcctStatus:   AccountStatusActive, // Changed from Status
		BalanceAmt:   0,          // Changed from Balance
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsActive checks if account is active
// BUG: Uses AcctStatus instead of Status
// Knowledge graph should detect: Account.IsActive uses Status, AccountV2.IsActive uses AcctStatus
func (a AccountV2) IsActive() bool {
	return a.AcctStatus == AccountStatusActive // Changed field name
}

// GetBalance returns the account balance
// BUG: Uses BalanceAmt instead of Balance
// Knowledge graph should detect all callers using GetBalance() expecting Balance field
func (a AccountV2) GetBalance() float64 {
	return a.BalanceAmt // Changed field name
}

// UpdateBalance updates the account balance
// BUG PATTERN: Method signature same but internal field different
// Callers won't know the internal field changed
func (a *AccountV2) UpdateBalance(amount float64) {
	a.BalanceAmt = amount // Changed from Balance
	a.UpdatedAt = time.Now()
}

// Suspend suspends the account
func (a *AccountV2) Suspend() {
	a.AcctStatus = AccountStatusSuspended // Changed from Status
	a.UpdatedAt = time.Now()
}

// GetEmail returns the email address
// NEW method - old code accessed Email directly, new code should use GetEmail()
func (a AccountV2) GetEmail() string {
	return a.EmailAddress
}

// =====================================
// MIGRATION HELPER
// This function shows the field mapping between versions
// Knowledge graph should detect these relationships
// =====================================

// MigrateToV2 migrates Account to AccountV2
// This documents the breaking changes:
// - Email -> EmailAddress
// - AccountType field -> AcctType
// - Status -> AcctStatus
// - Balance -> BalanceAmt
func MigrateToV2(old *Account) *AccountV2 {
	return &AccountV2{
		ID:           old.ID,
		OwnerID:      old.OwnerID,
		EmailAddress: old.Email,       // Email -> EmailAddress
		AcctType:     old.AccountType, // AccountType -> AcctType
		AcctStatus:   old.Status,      // Status -> AcctStatus
		BalanceAmt:   old.Balance,     // Balance -> BalanceAmt
		CreatedAt:    old.CreatedAt,
		UpdatedAt:    old.UpdatedAt,
	}
}

// =====================================
// AFFECTED LAYERS (Knowledge Graph should detect):
// =====================================
// Layer 1 (Entity): account.go - Original Account struct
// Layer 2 (Repository): account_repository.go - Uses Account.Email, Account.Balance
// Layer 3 (Service): account_service.go - Calls repository methods, accesses Account fields
// Layer 4 (Command): account_commands.go - Uses AccountService, creates Account entities
// Layer 5 (Handler): account_handler.go - Uses commands, accesses Account.Email directly
// Layer 6 (Infrastructure): account_repository_impl.go - Implements repository, accesses fields
//
// If Account.Email is renamed to EmailAddress:
// - account_repository.go:21 - FindByEmail uses Account.Email
// - account_repository_impl.go:57 - account.Email == email
// - account_handler.go:52 - account.Email
// - account_service.go:36 - calls FindByEmail which uses Email
// - account_commands.go:35 - uses CreateAccount which sets Email

