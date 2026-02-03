package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/catalog/internal/domain"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO categories (id, name, slug, parent_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID, c.Name, c.Slug, c.ParentID, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var c domain.Category
	var parentID *uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, parent_id, created_at, updated_at FROM categories WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Name, &c.Slug, &parentID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	c.ParentID = parentID
	return &c, nil
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var c domain.Category
	var parentID *uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, parent_id, created_at, updated_at FROM categories WHERE slug = $1`,
		slug,
	).Scan(&c.ID, &c.Name, &c.Slug, &parentID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	c.ParentID = parentID
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context, parentID *uuid.UUID, limit, offset int) ([]*domain.Category, error) {
	query := `SELECT id, name, slug, parent_id, created_at, updated_at FROM categories`
	args := []interface{}{limit, offset}
	if parentID != nil {
		query += ` WHERE parent_id = $1 ORDER BY name LIMIT $2 OFFSET $3`
		args = append([]interface{}{*parentID}, args...)
	} else {
		query += ` WHERE parent_id IS NULL ORDER BY name LIMIT $1 OFFSET $2`
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Category
	for rows.Next() {
		var c domain.Category
		var pid *uuid.UUID
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &pid, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.ParentID = pid
		list = append(list, &c)
	}
	return list, rows.Err()
}

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	c.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE categories SET name = $2, slug = $3, parent_id = $4, updated_at = $5 WHERE id = $1`,
		c.ID, c.Name, c.Slug, c.ParentID, c.UpdatedAt,
	)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}
