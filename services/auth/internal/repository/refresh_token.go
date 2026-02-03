package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshTokenRow for DB. We store token_hash only (never plain refresh token).
type RefreshTokenRow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// RefreshTokenRepository persists refresh token hashes.
type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

// NewRefreshTokenRepository returns a refresh token repository.
func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

// Create inserts a refresh token record.
func (r *RefreshTokenRepository) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	return err
}

// GetByTokenHash returns row if token hash exists and not expired.
func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshTokenRow, error) {
	var row RefreshTokenRow
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1 AND expires_at > now()`,
		tokenHash,
	).Scan(&row.ID, &row.UserID, &row.TokenHash, &row.ExpiresAt, &row.CreatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

// DeleteByID removes a refresh token (e.g. on refresh rotation).
func (r *RefreshTokenRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE id = $1`, id)
	return err
}

// DeleteExpired removes expired tokens (can be run by cron).
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	res, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE expires_at <= now()`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
