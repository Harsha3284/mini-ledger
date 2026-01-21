package model

import "time"

type EntryDirection string

const (
	Debit  EntryDirection = "debit"
	Credit EntryDirection = "credit"
)

type LedgerEntry struct {
	ID          string         `json:"id"`
	AccountID   string         `json:"account_id"`
	Direction   EntryDirection `json:"direction"`
	Amount      string         `json:"amount"`
	Category    *string        `json:"category,omitempty"`
	Description *string        `json:"description,omitempty"`
	OccurredAt  time.Time      `json:"occurred_at"`
	CreatedAt   time.Time      `json:"created_at"`
}
