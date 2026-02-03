package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/user/internal/domain"
)

type ProfileRepository struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	var p domain.Profile
	err := r.pool.QueryRow(ctx,
		`SELECT user_id, display_name, avatar_url, created_at, updated_at FROM user_profiles WHERE user_id = $1`,
		userID,
	).Scan(&p.UserID, &p.DisplayName, &p.AvatarURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProfileRepository) Create(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_profiles (user_id, created_at, updated_at) VALUES ($1, now(), now()) ON CONFLICT (user_id) DO NOTHING`,
		userID,
	)
	return err
}

func (r *ProfileRepository) Update(ctx context.Context, userID string, displayName *string, avatarURL *string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE user_profiles SET display_name = COALESCE($2, display_name), avatar_url = COALESCE($3, avatar_url), updated_at = now() WHERE user_id = $1`,
		userID, displayName, avatarURL,
	)
	return err
}
