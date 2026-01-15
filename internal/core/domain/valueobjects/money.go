package valueobjects

import (
	"errors"
	"fmt"
)

// Currency represents a currency code
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

// Money represents a monetary value with currency
// This is a value object that will be used across multiple layers
type Money struct {
	Amount   float64
	Currency Currency
}

// NewMoney creates a new Money value object
func NewMoney(amount float64, currency Currency) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

// Add adds two Money values (value receiver)
// BUG PATTERN: Returns Money but some callers might expect *Money
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("currency mismatch")
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Subtract subtracts two Money values (value receiver)
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("currency mismatch")
	}
	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}, nil
}

// IsPositive checks if the amount is positive (value receiver)
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// IsZero checks if the amount is zero (value receiver)
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// String returns string representation (value receiver)
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount, m.Currency)
}

// Multiply multiplies the amount (value receiver)
func (m Money) Multiply(factor float64) Money {
	return Money{
		Amount:   m.Amount * factor,
		Currency: m.Currency,
	}
}

