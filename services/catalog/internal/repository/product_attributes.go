package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/catalog/internal/domain"
)

func (r *ProductRepository) GetAttributeValues(ctx context.Context, productID uuid.UUID) ([]domain.ProductAttributeValue, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT product_id, attribute_id, value FROM product_attribute_values WHERE product_id = $1`,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.ProductAttributeValue
	for rows.Next() {
		var v domain.ProductAttributeValue
		if err := rows.Scan(&v.ProductID, &v.AttributeID, &v.Value); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *ProductRepository) ReplaceAttributeValues(ctx context.Context, productID uuid.UUID, values []domain.ProductAttributeValue) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM product_attribute_values WHERE product_id = $1`, productID)
	if err != nil {
		return err
	}
	for _, v := range values {
		v.ProductID = productID
		_, err := r.pool.Exec(ctx,
			`INSERT INTO product_attribute_values (product_id, attribute_id, value) VALUES ($1, $2, $3)`,
			v.ProductID, v.AttributeID, v.Value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
