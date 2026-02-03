// Package db — PostgreSQL connection via pgx. Prepared statements for security.
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a pgx connection pool. Use context for timeout.
func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	cfg.ConnConfig.DefaultQueryExecMode = pgxpool.QueryExecModeCacheStatement
	return pgxpool.NewWithConfig(ctx, cfg)
}
