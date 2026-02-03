package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/catalog/internal/domain"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = p.CreatedAt
	_, err := r.pool.Exec(ctx,
		`INSERT INTO products (id, category_id, name, slug, description, price, brand, seller_id, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		p.ID, p.CategoryID, p.Name, p.Slug, nullStr(p.Description), p.Price, nullStr(p.Brand), p.SellerID, p.Status, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var p domain.Product
	var desc, brand *string
	err := r.pool.QueryRow(ctx,
		`SELECT id, category_id, name, slug, description, price, brand, seller_id, status, created_at, updated_at FROM products WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.CategoryID, &p.Name, &p.Slug, &desc, &p.Price, &brand, &p.SellerID, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	if desc != nil {
		p.Description = *desc
	}
	if brand != nil {
		p.Brand = *brand
	}
	return &p, nil
}

func (r *ProductRepository) List(ctx context.Context, categoryID *uuid.UUID, sellerID string, status string, limit, offset int) ([]*domain.Product, error) {
	query := `SELECT id, category_id, name, slug, description, price, brand, seller_id, status, created_at, updated_at FROM products WHERE 1=1`
	args := []interface{}{}
	idx := 1
	if categoryID != nil {
		query += ` AND category_id = $` + fmt.Sprint(idx)
		args = append(args, *categoryID)
		idx++
	}
	if sellerID != "" {
		query += ` AND seller_id = $` + fmt.Sprint(idx)
		args = append(args, sellerID)
		idx++
	}
	if status != "" {
		query += ` AND status = $` + fmt.Sprint(idx)
		args = append(args, status)
		idx++
	}
	query += ` ORDER BY updated_at DESC LIMIT $` + fmt.Sprint(idx) + ` OFFSET $` + fmt.Sprint(idx+1)
	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Product
	for rows.Next() {
		var p domain.Product
		var desc, brand *string
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Slug, &desc, &p.Price, &brand, &p.SellerID, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if desc != nil {
			p.Description = *desc
		}
		if brand != nil {
			p.Brand = *brand
		}
		list = append(list, &p)
	}
	return list, rows.Err()
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	p.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE products SET category_id = $2, name = $3, slug = $4, description = $5, price = $6, brand = $7, status = $8, updated_at = $9 WHERE id = $1`,
		p.ID, p.CategoryID, p.Name, p.Slug, nullStr(p.Description), p.Price, nullStr(p.Brand), p.Status, p.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
