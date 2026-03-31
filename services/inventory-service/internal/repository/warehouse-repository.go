// Package repository предоставляет репозиторий для работы со складами.
package repository

import (
	"context"

	"kinos/inventory-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WarehouseInterface interface {
	Create(ctx context.Context, name, city, street, building, building2 string) (*model.Warehouse, error)
	Update(ctx context.Context, id uint64, name, city, street, building, building2 string) (*model.Warehouse, error)
	GetList(ctx context.Context, limit, offset int32) ([]*model.Warehouse, int32, error)
	Delete(ctx context.Context, id uint64) error
}

type WarehouseRepository struct {
	pool *pgxpool.Pool
}

func NewWarehouseRepository(pool *pgxpool.Pool) *WarehouseRepository {
	return &WarehouseRepository{
		pool: pool,
	}
}

func (r *WarehouseRepository) Create(ctx context.Context, name, city, street, building, building2 string) (*model.Warehouse, error) {
	query := `
		INSERT INTO warehouses (name, city, street, building, building2, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, name, city, street, building, building2, created_at
	`
	var w model.Warehouse
	err := r.pool.QueryRow(ctx, query, name, city, street, building, building2).Scan(
		&w.Id,
		&w.Name,
		&w.City,
		&w.Street,
		&w.Building,
		&w.Building2,
		&w.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WarehouseRepository) Update(ctx context.Context, id uint64, name, city, street, building, building2 string) (*model.Warehouse, error) {
	query := `
		UPDATE warehouses
		SET name = $2, city = $3, street = $4, building = $5, building2 = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, city, street, building, building2, updated_at
	`
	var w model.Warehouse
	err := r.pool.QueryRow(ctx, query, id, name, city, street, building, building2).Scan(
		&w.Id,
		&w.Name,
		&w.City,
		&w.Street,
		&w.Building,
		&w.Building2,
		&w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WarehouseRepository) GetList(ctx context.Context, limit, offset int32) ([]*model.Warehouse, int32, error) {
	query := `
		SELECT id, name, city, street, building, building2, created_at
		FROM warehouses
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var warehouses []*model.Warehouse
	for rows.Next() {
		var w model.Warehouse
		err := rows.Scan(
			&w.Id,
			&w.Name,
			&w.City,
			&w.Street,
			&w.Building,
			&w.Building2,
			&w.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		warehouses = append(warehouses, &w)
	}

	// Получаем общее количество
	var total int32
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM warehouses`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func (r *WarehouseRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM warehouses WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
