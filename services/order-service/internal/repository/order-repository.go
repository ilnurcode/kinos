package repository

import (
	"context"
	"errors"
	"fmt"

	"kinos/order-service/internal/errs"
	"kinos/order-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepositoryInterface interface {
	CreateOrder(ctx context.Context, order *models.Order) (uint64, error)
	GetOrder(ctx context.Context, orderID uint64) (*models.Order, error)
	GetListOrders(ctx context.Context, limit, offset int32, status string) ([]*models.Order, int32, error)
	GetUserOrders(ctx context.Context, userID uint64, limit, offset int32) ([]*models.Order, int32, error)
	UpdateOrderStatus(ctx context.Context, orderID uint64, status models.OrderStatus) error
	AddOrderItem(ctx context.Context, item *models.OrderItem) error
}

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) (uint64, error) {
	querier := GetQuerier(ctx, r.pool)

	var orderID uint64
	err := querier.QueryRow(ctx, `
		INSERT INTO orders (user_id, total, status, delivery_address, phone, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING order_id
	`, order.UserID, order.Total, order.Status, order.DeliveryAddress, order.Phone, order.Comment).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	order.ID = orderID
	return orderID, nil
}

func (r *OrderRepository) AddOrderItem(ctx context.Context, item *models.OrderItem) error {
	querier := GetQuerier(ctx, r.pool)
	_, err := querier.Exec(ctx, `
		INSERT INTO order_items (order_id, product_id, product_name, quantity, price, subtotal)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, item.OrderID, item.ProductID, item.ProductName, item.Quantity, item.Price, item.Subtotal)
	return err
}

func (r *OrderRepository) GetOrder(ctx context.Context, orderID uint64) (*models.Order, error) {
	querier := GetQuerier(ctx, r.pool)
	order := &models.Order{}

	err := querier.QueryRow(ctx, `
		SELECT order_id, user_id, total, status, delivery_address, phone, comment, created_at, updated_at
		FROM orders WHERE order_id = $1
	`, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.DeliveryAddress,
		&order.Phone,
		&order.Comment,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	items, err := r.getOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (r *OrderRepository) getOrderItems(ctx context.Context, orderID uint64) ([]models.OrderItem, error) {
	querier := GetQuerier(ctx, r.pool)
	rows, err := querier.Query(ctx, `
		SELECT item_id, order_id, product_id, product_name, quantity, price, subtotal
		FROM order_items WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.Quantity, &item.Price, &item.Subtotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderRepository) GetListOrders(ctx context.Context, limit, offset int32, status string) ([]*models.Order, int32, error) {
	querier := GetQuerier(ctx, r.pool)

	var (
		rows pgx.Rows
		err  error
	)
	if status != "" {
		rows, err = querier.Query(ctx, `
			SELECT order_id, user_id, total, status, delivery_address, phone, comment, created_at, updated_at
			FROM orders
			WHERE status = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, status, limit, offset)
	} else {
		rows, err = querier.Query(ctx, `
			SELECT order_id, user_id, total, status, delivery_address, phone, comment, created_at, updated_at
			FROM orders
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.DeliveryAddress, &order.Phone, &order.Comment, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	var total int32
	if status != "" {
		if err := querier.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE status = $1`, status).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := querier.QueryRow(ctx, `SELECT COUNT(*) FROM orders`).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	return orders, total, nil
}

func (r *OrderRepository) GetUserOrders(ctx context.Context, userID uint64, limit, offset int32) ([]*models.Order, int32, error) {
	querier := GetQuerier(ctx, r.pool)
	rows, err := querier.Query(ctx, `
		SELECT order_id, user_id, total, status, delivery_address, phone, comment, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.DeliveryAddress, &order.Phone, &order.Comment, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	var total int32
	if err := querier.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID uint64, status models.OrderStatus) error {
	querier := GetQuerier(ctx, r.pool)
	result, err := querier.Exec(ctx, `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE order_id = $2
	`, status, orderID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}
	return nil
}
