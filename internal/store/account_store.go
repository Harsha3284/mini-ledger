package store

import (
	"context"
	"mini-ledger/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountStore struct {
	db *pgxpool.Pool
}

func NewAccountStore(db *pgxpool.Pool) *AccountStore {
	return &AccountStore{db: db}
}

func (s *AccountStore) Create(ctx context.Context, a model.Account) (model.Account, error) {
	const q = `
		INSERT INTO accounts (id, name, type, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at;
	`
	err := s.db.QueryRow(ctx, q, a.ID, a.Name, string(a.Type), string(a.Currency)).
		Scan(&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return model.Account{}, err
	}
	return a, nil
}

func (s *AccountStore) GetByID(ctx context.Context, id string) (model.Account, error) {
	const q = `
		SELECT id, name, type, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1;
	`
	var a model.Account
	var typ string
	var cur string

	err := s.db.QueryRow(ctx, q, id).Scan(&a.ID, &a.Name, &typ, &cur, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Account{}, ErrNotFound
		}
		return model.Account{}, err
	}

	a.Type = model.AccountType(typ)
	a.Currency = model.Currency(cur)
	return a, nil
}

func (s *AccountStore) List(ctx context.Context, limit int) ([]model.Account, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	const q = `
		SELECT id, name, type, currency, created_at, updated_at
		FROM accounts
		ORDER BY created_at DESC
		LIMIT $1;
	`

	rows, err := s.db.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Account, 0, limit)
	for rows.Next() {
		var a model.Account
		var typ string
		var cur string

		if err := rows.Scan(&a.ID, &a.Name, &typ, &cur, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.Type = model.AccountType(typ)
		a.Currency = model.Currency(cur)
		out = append(out, a)
	}

	return out, rows.Err()
}

func (s *AccountStore) UpdateName(ctx context.Context, id string, name string) (model.Account, error) {
	const q = `
		UPDATE accounts
		SET name = $2
		WHERE id = $1
		RETURNING id, name, type, currency, created_at, updated_at;
	`

	var a model.Account
	var typ, cur string

	err := s.db.QueryRow(ctx, q, id, name).Scan(&a.ID, &a.Name, &typ, &cur, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Account{}, ErrNotFound
		}
		return model.Account{}, err
	}
	a.Type = model.AccountType(typ)
	a.Currency = model.Currency(cur)
	return a, nil
}

func (s *AccountStore) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM accounts WHERE id = $1;`

	ct, err := s.db.Exec(ctx, q, id)
	if err != nil {
		// FK restrict violation = 23503
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return ErrConflict
		}
		return err
	}

	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
