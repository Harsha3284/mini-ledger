package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("db: parse config: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("db: create pool: %w", err)
	}

	// verify connectivity
	ctxPing, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return pool, nil
}
