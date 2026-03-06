// Package repository предоставляет репозиторий для работы с категориями в базе данных.
// Включает CRUD-операции и поиск категорий.
package repository

import (
	"context"
	"fmt"
	"kinos/catalog-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepositoryInterface interface {
	DeleteCategory(ctx context.Context, categoryID uint64) error
	UpdateCategory(ctx context.Context, category *models.Category) error
	CreateCategory(ctx context.Context, category *models.Category) (uint64, error)
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	GetCategoryByID(ctx context.Context, id uint64) (*models.Category, error)
	GetListCategory(ctx context.Context, limit, offset int32) ([]*models.Category, int32, error)
}

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

func (r *CategoryRepository) GetListCategory(ctx context.Context, limit, offset int32) ([]*models.Category, int32, error) {
	querier := GetQuerier(ctx, r.DB)
	rows, err := querier.Query(ctx, "SELECT category_id, category_name FROM category LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get list category: %w", err)
	}
	defer rows.Close()
	var categories []*models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, 0, fmt.Errorf("failed to scan list category: %w", err)
		}
		categories = append(categories, &category)

	}
	var total int32
	err = querier.QueryRow(ctx, "SELECT COUNT(*) FROM category").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total category: %w", err)
	}
	return categories, total, nil
}

func (r *CategoryRepository) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT * FROM category WHERE category_name=$1", name).Scan(&category.Id, &category.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uint64) (*models.Category, error) {
	var category models.Category
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT * FROM category WHERE category_id=$1", id).Scan(&category.Id, &category.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) (uint64, error) {
	var categoryID uint64
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "INSERT INTO category (category_name) VALUES ($1) RETURNING category_id", category.Name).Scan(&categoryID)
	if err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}
	return categoryID, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "UPDATE category SET category_name=$1 WHERE category_id=$2", category.Name, category.Id)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, categoryID uint64) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "DELETE FROM category WHERE category_id=$1", categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
