package service

import (
	"context"
	"errors"
	"mini-ledger/internal/model"
	"mini-ledger/internal/store"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidEntry    = errors.New("invalid entry")
)

type LedgerService struct {
	ledgerStore  *store.LedgerStore
	accountStore *store.AccountStore
}

func NewLedgerService(ls *store.LedgerStore, as *store.AccountStore) *LedgerService {
	return &LedgerService{ledgerStore: ls, accountStore: as}
}

func (s *LedgerService) CreateEntry(
	ctx context.Context,
	accountID string,
	direction model.EntryDirection,
	amount string,
	category *string,
	description *string,
	occurredAt time.Time,
) (model.LedgerEntry, error) {

	accountID = strings.TrimSpace(accountID)
	if accountID == "" {
		return model.LedgerEntry{}, ErrInvalidEntry
	}

	// Validate direction
	if direction != model.Debit && direction != model.Credit {
		return model.LedgerEntry{}, ErrInvalidEntry
	}

	// Validate amount string (basic; later we’ll use decimal)
	amount = strings.TrimSpace(amount)
	if amount == "" || amount == "0" || amount == "0.00" {
		return model.LedgerEntry{}, ErrInvalidEntry
	}

	// Validate occurred_at
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	// Ensure account exists (clean 404 instead of FK 500)
	exists, err := s.ledgerStore.AccountExists(ctx, accountID)
	if err != nil {
		return model.LedgerEntry{}, err
	}
	if !exists {
		return model.LedgerEntry{}, ErrAccountNotFound
	}

	e := model.LedgerEntry{
		ID:          uuid.NewString(),
		AccountID:   accountID,
		Direction:   direction,
		Amount:      amount,
		Category:    category,
		Description: description,
		OccurredAt:  occurredAt.UTC(),
	}

	return s.ledgerStore.Create(ctx, e)
}

func (s *LedgerService) ListEntries(ctx context.Context, accountID string, limit int) ([]model.LedgerEntry, error) {
	return s.ledgerStore.ListByAccount(ctx, accountID, limit)
}

func (s *LedgerService) GetBalance(ctx context.Context, accountID string) (string, error) {
	// optional: validate account exists for clearer error
	exists, err := s.ledgerStore.AccountExists(ctx, accountID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrAccountNotFound
	}
	return s.ledgerStore.GetBalance(ctx, accountID)
}