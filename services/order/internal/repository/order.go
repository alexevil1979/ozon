package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/order/internal/domain"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, o *domain.Order, items []domain.OrderItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO orders (id, buyer_id, status, total, payment_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		o.ID, o.BuyerID, o.Status, o.Total, nullStr(o.PaymentID), o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return err
	}
	for _, it := range items {
		it.OrderID = o.ID
		if it.ID == uuid.Nil {
			it.ID = uuid.New()
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (id, order_id, product_id, name, price, quantity) VALUES ($1, $2, $3, $4, $5, $6)`,
			it.ID, it.OrderID, it.ProductID, it.Name, it.Price, it.Quantity,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var o domain.Order
	var paymentID *string
	err := r.pool.QueryRow(ctx,
		`SELECT id, buyer_id, status, total, payment_id, created_at, updated_at FROM orders WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.BuyerID, &o.Status, &o.Total, &paymentID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	if paymentID != nil {
		o.PaymentID = *paymentID
	}
	return &o, nil
}

func (r *OrderRepository) ListByBuyerID(ctx context.Context, buyerID string, limit, offset int) ([]*domain.Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, buyer_id, status, total, payment_id, created_at, updated_at FROM orders WHERE buyer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		buyerID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Order
	for rows.Next() {
		var o domain.Order
		var paymentID *string
		if err := rows.Scan(&o.ID, &o.BuyerID, &o.Status, &o.Total, &paymentID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		if paymentID != nil {
			o.PaymentID = *paymentID
		}
		list = append(list, &o)
	}
	return list, rows.Err()
}

func (r *OrderRepository) GetItems(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, order_id, product_id, name, price, quantity FROM order_items WHERE order_id = $1 ORDER BY name`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.OrderItem
	for rows.Next() {
		var it domain.OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Name, &it.Price, &it.Quantity); err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	return list, rows.Err()
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, paymentID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE orders SET status = $2, payment_id = $3, updated_at = $4 WHERE id = $1`,
		id, status, nullStr(paymentID), time.Now().UTC(),
	)
	return err
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
