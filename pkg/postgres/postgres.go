package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Params holds connection settings. If DSN is set, the other fields are ignored.
type Params struct {
	DSN                  string
	MinConns             int32
	MaxConns             int32
	MaxConnLifetime      time.Duration
}

// NewPool creates a pgx connection pool with sensible defaults.
func NewPool(ctx context.Context, p Params) (*pgxpool.Pool, error) {
	if p.DSN == "" {
		return nil, fmt.Errorf("dsn required")
	}

	cfg, err := pgxpool.ParseConfig(p.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	// Apply pool tuning.
	if p.MinConns > 0 {
		cfg.MinConns = p.MinConns
	}
	if p.MaxConns > 0 {
		cfg.MaxConns = p.MaxConns
	}
	if p.MaxConnLifetime > 0 {
		cfg.MaxConnLifetime = p.MaxConnLifetime
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	// Verify connectivity up front.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
