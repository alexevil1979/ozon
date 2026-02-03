package db

import (
	"context"
	"embed"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed all:migrations
var migrationsFS embed.FS

const migrationsDir = "migrations"

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	names, err := listMigrations()
	if err != nil {
		return err
	}
	for _, name := range names {
		applied, err := isApplied(ctx, pool, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		body, err := migrationsFS.ReadFile(migrationsDir + "/" + name)
		if err != nil {
			return err
		}
		sql := strings.TrimSpace(string(body))
		if sql == "" {
			_ = markApplied(ctx, pool, name)
			continue
		}
		if _, err := pool.Exec(ctx, sql); err != nil {
			return err
		}
		if err := markApplied(ctx, pool, name); err != nil {
			return err
		}
	}
	return nil
}

func listMigrations() ([]string, error) {
	entries, err := migrationsFS.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func isApplied(ctx context.Context, pool *pgxpool.Pool, name string) (bool, error) {
	_, _ = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`)
	var n int
	err := pool.QueryRow(ctx, `SELECT 1 FROM schema_migrations WHERE name = $1`, name).Scan(&n)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func markApplied(ctx context.Context, pool *pgxpool.Pool, name string) error {
	_, err := pool.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, name)
	return err
}
