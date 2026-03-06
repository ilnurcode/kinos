// Package repository предоставляет репозиторий для работы с производителями в базе данных.
// Включает CRUD-операции и поиск производителей.
package repository

import (
	"context"
	"fmt"
	"kinos/catalog-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ManufacturersRepositoryInterface interface {
	DeleteManufacturers(ctx context.Context, manufacturerID uint64) error
	UpdateManufacturers(ctx context.Context, manufacturer *models.Manufacturer) error
	GetManufacturerByID(ctx context.Context, id uint64) (*models.Manufacturer, error)
	CreateManufacturers(ctx context.Context, manufacturer *models.Manufacturer) (uint64, error)
	GetManufacturerByName(ctx context.Context, name string) (*models.Manufacturer, error)
	GetListManufacturers(ctx context.Context, limit, offset int32) ([]*models.Manufacturer, int32, error)
}
type ManufacturersRepository struct {
	DB *pgxpool.Pool
}

func NewManufacturersRepository(db *pgxpool.Pool) *ManufacturersRepository {
	return &ManufacturersRepository{
		DB: db,
	}
}

func (r *ManufacturersRepository) GetManufacturerByName(ctx context.Context, name string) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT manufacturers_id, manufacturers_name FROM manufacturers WHERE manufacturers_name=$1", name).Scan(&manufacturer.Id, &manufacturer.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get manufacturer by name: %w", err)
	}
	return &manufacturer, nil
}

func (r *ManufacturersRepository) GetManufacturerByID(ctx context.Context, id uint64) (*models.Manufacturer, error) {
	var manufacturer models.Manufacturer
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT manufacturers_id, manufacturers_name FROM manufacturers WHERE manufacturers_id=$1", id).Scan(&manufacturer.Id, &manufacturer.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get manufacturer by name: %w", err)
	}
	return &manufacturer, nil
}

func (r *ManufacturersRepository) GetListManufacturers(ctx context.Context, limit, offset int32) ([]*models.Manufacturer, int32, error) {
	querier := GetQuerier(ctx, r.DB)
	rows, err := querier.Query(ctx, "SELECT manufacturers_id, manufacturers_name FROM manufacturers LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get list manufacturers: %w", err)
	}
	defer rows.Close()
	var manufacturers []*models.Manufacturer
	for rows.Next() {
		var manufacturer models.Manufacturer
		if err := rows.Scan(&manufacturer.Id, &manufacturer.Name); err != nil {
			return nil, 0, fmt.Errorf("failed to scan list category: %w", err)
		}
		manufacturers = append(manufacturers, &manufacturer)

	}
	var total int32
	err = querier.QueryRow(ctx, "SELECT COUNT(*) FROM manufacturers").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total category: %w", err)
	}
	return manufacturers, total, nil
}

func (r *ManufacturersRepository) CreateManufacturers(ctx context.Context, manufacturer *models.Manufacturer) (uint64, error) {
	var manufacturerID uint64
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "INSERT INTO manufacturers (manufacturers_name) VALUES ($1) RETURNING manufacturers_id", manufacturer.Name).Scan(&manufacturerID)
	if err != nil {
		return 0, fmt.Errorf("failed to create manufacturers: %w", err)
	}
	return manufacturerID, nil
}

func (r *ManufacturersRepository) UpdateManufacturers(ctx context.Context, manufacturer *models.Manufacturer) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "UPDATE manufacturers SET manufacturers_name=$1 WHERE manufacturers_id = $2", manufacturer.Name, manufacturer.Id)
	if err != nil {
		return fmt.Errorf("failed to update manufacturers: %w", err)
	}
	return nil
}

func (r *ManufacturersRepository) DeleteManufacturers(ctx context.Context, manufacturerID uint64) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "DELETE FROM manufacturers WHERE manufacturers_id=$1", manufacturerID)
	if err != nil {
		return fmt.Errorf("failed to delete manufacturers: %w", err)
	}
	return nil
}
