package store

import (
	"context"
	"mini-ledger/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerStore struct {
	db *pgxpool.Pool
}

func NewLedgerStore(db *pgxpool.Pool) *LedgerStore {
	return &LedgerStore{db: db}
}

func (s *LedgerStore) Create(ctx context.Context, e model.LedgerEntry) (model.LedgerEntry, error) {
	const q = `
		INSERT INTO ledger_entries (id, account_id, direction, amount, category, description, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at;
	`
	err := s.db.QueryRow(ctx, q,
		e.ID, e.AccountID, string(e.Direction), e.Amount, e.Category, e.Description, e.OccurredAt,
	).Scan(&e.CreatedAt)
	if err != nil {
		return model.LedgerEntry{}, err
	}
	return e, nil
}

func (s *LedgerStore) ListByAccount(ctx context.Context, accountID string, limit int) ([]model.LedgerEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	const q = `
		SELECT id, account_id, direction, amount, category, description, occurred_at, created_at
		FROM ledger_entries
		WHERE account_id = $1
		ORDER BY occurred_at DESC
		LIMIT $2;
	`

	rows, err := s.db.Query(ctx, q, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.LedgerEntry, 0, limit)
	for rows.Next() {
		var e model.LedgerEntry
		var dir string

		if err := rows.Scan(&e.ID, &e.AccountID, &dir, &e.Amount, &e.Category, &e.Description, &e.OccurredAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		e.Direction = model.EntryDirection(dir)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *LedgerStore) GetBalance(ctx context.Context, accountID string) (string, error) {
	// Balance = credits - debits (simple rule)
	const q = `
		SELECT
		  COALESCE(SUM(CASE WHEN direction='credit' THEN amount ELSE 0 END), 0) -
		  COALESCE(SUM(CASE WHEN direction='debit'  THEN amount ELSE 0 END), 0) AS balance
		FROM ledger_entries
		WHERE account_id = $1;
	`
	var balance string
	if err := s.db.QueryRow(ctx, q, accountID).Scan(&balance); err != nil {
		return "", err
	}
	return balance, nil
}

func (s *LedgerStore) AccountExists(ctx context.Context, accountID string) (bool, error) {
	const q = `SELECT 1 FROM accounts WHERE id=$1;`
	var one int
	err := s.db.QueryRow(ctx, q, accountID).Scan(&one)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}