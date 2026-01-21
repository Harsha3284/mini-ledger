package model

import "time"

type AccountType string
type Currency string

const (
	AccountChecking AccountType = "checking"
	AccountSavings  AccountType = "savings"
	AccountCredit   AccountType = "credit"
)

const (
	USD Currency = "USD"
)

type Account struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      AccountType `json:"type"`
	Currency  Currency    `json:"currency"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
