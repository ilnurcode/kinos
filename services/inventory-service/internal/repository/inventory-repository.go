// Package repository предоставляет репозиторий для работы с запасами товаров.
package repository

import (
	"context"
	"strconv"

	"kinos/inventory-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryInterface interface {
	Create(ctx context.Context, productID uint64, quantity int32, location string) (*model.Inventory, error)
	Update(ctx context.Context, id uint64, quantity int32, location string) (*model.Inventory, error)
	GetByProductID(ctx context.Context, productID uint64) (*model.Inventory, error)
	GetList(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) ([]*model.Inventory, int32, error)
	Delete(ctx context.Context, id uint64) error
	ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error
	ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error)
}

type InventoryRepository struct {
	pool *pgxpool.Pool
}

func NewInventoryRepository(pool *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{
		pool: pool,
	}
}

func (r *InventoryRepository) Create(ctx context.Context, productID uint64, quantity int32, location string) (*model.Inventory, error) {
	query := `
		INSERT INTO inventory (product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at)
		VALUES ($1, $2, 0, $2, $3, NOW())
		RETURNING id, product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at
	`
	var inv model.Inventory
	err := r.pool.QueryRow(ctx, query, productID, quantity, location).Scan(
		&inv.Id,
		&inv.ProductId,
		&inv.Quantity,
		&inv.ReservedQuantity,
		&inv.AvailableQuantity,
		&inv.WarehouseLocation,
		&inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) Update(ctx context.Context, id uint64, quantity int32, location string) (*model.Inventory, error) {
	query := `
		UPDATE inventory
		SET quantity = $2,
		    available_quantity = $2 - reserved_quantity,
		    warehouse_location = $3,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at
	`
	var inv model.Inventory
	err := r.pool.QueryRow(ctx, query, id, quantity, location).Scan(
		&inv.Id,
		&inv.ProductId,
		&inv.Quantity,
		&inv.ReservedQuantity,
		&inv.AvailableQuantity,
		&inv.WarehouseLocation,
		&inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) GetByProductID(ctx context.Context, productID uint64) (*model.Inventory, error) {
	query := `
		SELECT id, product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at
		FROM inventory
		WHERE product_id = $1
	`
	var inv model.Inventory
	err := r.pool.QueryRow(ctx, query, productID).Scan(
		&inv.Id,
		&inv.ProductId,
		&inv.Quantity,
		&inv.ReservedQuantity,
		&inv.AvailableQuantity,
		&inv.WarehouseLocation,
		&inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) GetList(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) ([]*model.Inventory, int32, error) {
	countQuery := `SELECT COUNT(*) FROM inventory WHERE 1=1`
	query := `
		SELECT id, product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at
		FROM inventory
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if productID > 0 {
		query += " AND product_id = $" + strconv.Itoa(argIndex)
		countQuery += " AND product_id = $" + strconv.Itoa(argIndex)
		args = append(args, productID)
		argIndex++
	}

	if location != "" {
		query += " AND warehouse_location = $" + strconv.Itoa(argIndex)
		countQuery += " AND warehouse_location = $" + strconv.Itoa(argIndex)
		args = append(args, location)
		argIndex++
	}

	if minQuantity > 0 {
		query += " AND available_quantity >= $" + strconv.Itoa(argIndex)
		countQuery += " AND available_quantity >= $" + strconv.Itoa(argIndex)
		args = append(args, minQuantity)
		argIndex++
	}

	query += " ORDER BY product_id LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	countQuery += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)

	args = append(args, limit, offset)

	var total int32
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var inventories []*model.Inventory
	for rows.Next() {
		var inv model.Inventory
		err := rows.Scan(
			&inv.Id,
			&inv.ProductId,
			&inv.Quantity,
			&inv.ReservedQuantity,
			&inv.AvailableQuantity,
			&inv.WarehouseLocation,
			&inv.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, &inv)
	}

	return inventories, total, nil
}

func (r *InventoryRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM inventory WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *InventoryRepository) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error {
	query := `
		UPDATE inventory
		SET reserved_quantity = reserved_quantity + $1,
		    available_quantity = available_quantity - $1,
		    updated_at = NOW()
		WHERE product_id = $2 AND available_quantity >= $1
	`
	result, err := r.pool.Exec(ctx, query, quantity, productID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (r *InventoryRepository) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error) {
	query := `
		UPDATE inventory
		SET reserved_quantity = GREATEST(0, reserved_quantity - 1),
		    available_quantity = available_quantity + 1,
		    updated_at = NOW()
		WHERE product_id = $1 AND reserved_quantity > 0
		RETURNING 1
	`
	var dummy int
	err := r.pool.QueryRow(ctx, query, productID).Scan(&dummy)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

var ErrInsufficientStock = &InsufficientStockError{}

type InsufficientStockError struct{}

func (e *InsufficientStockError) Error() string {
	return "недостаточно товара на складе"
}
