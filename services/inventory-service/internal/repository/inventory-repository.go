package repository

import (
	"context"
	"errors"
	"strconv"

	"kinos/inventory-service/internal/errs"
	"kinos/inventory-service/internal/model"

	"github.com/jackc/pgx/v5"
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
	return &InventoryRepository{pool: pool}
}

func (r *InventoryRepository) Create(ctx context.Context, productID uint64, quantity int32, location string) (*model.Inventory, error) {
	querier := GetQuerier(ctx, r.pool)
	query := `
		INSERT INTO inventory (product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at)
		VALUES ($1, $2, 0, $2, $3, NOW())
		RETURNING id, product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at
	`

	var inv model.Inventory
	err := querier.QueryRow(ctx, query, productID, quantity, location).Scan(
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
	querier := GetQuerier(ctx, r.pool)
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
	err := querier.QueryRow(ctx, query, id, quantity, location).Scan(
		&inv.Id,
		&inv.ProductId,
		&inv.Quantity,
		&inv.ReservedQuantity,
		&inv.AvailableQuantity,
		&inv.WarehouseLocation,
		&inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrInventoryNotFound
		}
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) GetByProductID(ctx context.Context, productID uint64) (*model.Inventory, error) {
	querier := GetQuerier(ctx, r.pool)
	query := `
		SELECT id, product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at
		FROM inventory
		WHERE product_id = $1
	`

	var inv model.Inventory
	err := querier.QueryRow(ctx, query, productID).Scan(
		&inv.Id,
		&inv.ProductId,
		&inv.Quantity,
		&inv.ReservedQuantity,
		&inv.AvailableQuantity,
		&inv.WarehouseLocation,
		&inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrInventoryNotFound
		}
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) GetList(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) ([]*model.Inventory, int32, error) {
	querier := GetQuerier(ctx, r.pool)
	countQuery := `SELECT COUNT(*) FROM inventory WHERE 1=1`
	query := `
		SELECT id, product_id, quantity, reserved_quantity, available_quantity, warehouse_location, updated_at
		FROM inventory
		WHERE 1=1
	`

	filterArgs := []interface{}{}
	argIndex := 1

	if productID > 0 {
		query += " AND product_id = $" + strconv.Itoa(argIndex)
		countQuery += " AND product_id = $" + strconv.Itoa(argIndex)
		filterArgs = append(filterArgs, productID)
		argIndex++
	}

	if location != "" {
		query += " AND warehouse_location = $" + strconv.Itoa(argIndex)
		countQuery += " AND warehouse_location = $" + strconv.Itoa(argIndex)
		filterArgs = append(filterArgs, location)
		argIndex++
	}

	if minQuantity > 0 {
		query += " AND available_quantity >= $" + strconv.Itoa(argIndex)
		countQuery += " AND available_quantity >= $" + strconv.Itoa(argIndex)
		filterArgs = append(filterArgs, minQuantity)
		argIndex++
	}

	query += " ORDER BY product_id LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)

	var total int32
	if err := querier.QueryRow(ctx, countQuery, filterArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args := append(filterArgs, limit, offset)
	rows, err := querier.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var inventories []*model.Inventory
	for rows.Next() {
		var inv model.Inventory
		if err := rows.Scan(
			&inv.Id,
			&inv.ProductId,
			&inv.Quantity,
			&inv.ReservedQuantity,
			&inv.AvailableQuantity,
			&inv.WarehouseLocation,
			&inv.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, &inv)
	}

	return inventories, total, nil
}

func (r *InventoryRepository) Delete(ctx context.Context, id uint64) error {
	querier := GetQuerier(ctx, r.pool)
	_, err := querier.Exec(ctx, `DELETE FROM inventory WHERE id = $1`, id)
	return err
}

func (r *InventoryRepository) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error {
	querier := GetQuerier(ctx, r.pool)

	var existingQuantity int32
	err := querier.QueryRow(ctx, `
		SELECT quantity
		FROM inventory_reservations
		WHERE reservation_id = $1 AND product_id = $2
	`, reservationID, productID).Scan(&existingQuantity)
	if err == nil {
		return nil
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	result, err := querier.Exec(ctx, `
		UPDATE inventory
		SET reserved_quantity = reserved_quantity + $1,
		    available_quantity = available_quantity - $1,
		    updated_at = NOW()
		WHERE product_id = $2 AND available_quantity >= $1
	`, quantity, productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		var exists bool
		if err := querier.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM inventory WHERE product_id = $1)`, productID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return errs.ErrInventoryNotFound
		}
		return errs.ErrInsufficientStock
	}

	_, err = querier.Exec(ctx, `
		INSERT INTO inventory_reservations (reservation_id, product_id, quantity, created_at)
		VALUES ($1, $2, $3, NOW())
	`, reservationID, productID, quantity)
	return err
}

func (r *InventoryRepository) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error) {
	querier := GetQuerier(ctx, r.pool)

	var quantity int32
	err := querier.QueryRow(ctx, `
		DELETE FROM inventory_reservations
		WHERE reservation_id = $1 AND product_id = $2
		RETURNING quantity
	`, reservationID, productID).Scan(&quantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	result, err := querier.Exec(ctx, `
		UPDATE inventory
		SET reserved_quantity = reserved_quantity - $1,
		    available_quantity = available_quantity + $1,
		    updated_at = NOW()
		WHERE product_id = $2
	`, quantity, productID)
	if err != nil {
		return 0, err
	}
	if result.RowsAffected() == 0 {
		return 0, errs.ErrInventoryNotFound
	}

	return quantity, nil
}
