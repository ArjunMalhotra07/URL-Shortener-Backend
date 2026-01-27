package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Params holds connection settings.
type Params struct {
	DSN             string
	MaxConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
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
	if p.MaxConns > 0 {
		cfg.MaxConns = p.MaxConns
	}
	if p.MaxConnLifetime > 0 {
		cfg.MaxConnLifetime = p.MaxConnLifetime
	}
	if p.MaxConnIdleTime > 0 {
		cfg.MaxConnIdleTime = p.MaxConnIdleTime
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
