package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/catalog/internal/domain"
)

type AttributeRepository struct {
	pool *pgxpool.Pool
}

func NewAttributeRepository(pool *pgxpool.Pool) *AttributeRepository {
	return &AttributeRepository{pool: pool}
}

func (r *AttributeRepository) Create(ctx context.Context, a *domain.Attribute) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO attributes (id, name, slug, value_type, created_at) VALUES ($1, $2, $3, $4, $5)`,
		a.ID, a.Name, a.Slug, a.ValueType, a.CreatedAt,
	)
	return err
}

func (r *AttributeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
	var a domain.Attribute
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, value_type, created_at FROM attributes WHERE id = $1`,
		id,
	).Scan(&a.ID, &a.Name, &a.Slug, &a.ValueType, &a.CreatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AttributeRepository) GetBySlug(ctx context.Context, slug string) (*domain.Attribute, error) {
	var a domain.Attribute
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, value_type, created_at FROM attributes WHERE slug = $1`,
		slug,
	).Scan(&a.ID, &a.Name, &a.Slug, &a.ValueType, &a.CreatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AttributeRepository) List(ctx context.Context, limit, offset int) ([]*domain.Attribute, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, slug, value_type, created_at FROM attributes ORDER BY name LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Attribute
	for rows.Next() {
		var a domain.Attribute
		if err := rows.Scan(&a.ID, &a.Name, &a.Slug, &a.ValueType, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	return list, rows.Err()
}

func (r *AttributeRepository) Update(ctx context.Context, a *domain.Attribute) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE attributes SET name = $2, slug = $3, value_type = $4 WHERE id = $1`,
		a.ID, a.Name, a.Slug, a.ValueType,
	)
	return err
}

func (r *AttributeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM attributes WHERE id = $1`, id)
	return err
}
