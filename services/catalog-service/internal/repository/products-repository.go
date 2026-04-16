// Package repository предоставляет репозиторий для работы с товарами в базе данных.
// Включает CRUD-операции, поиск и фильтрацию товаров.
package repository

import (
	"context"
	"fmt"
	"kinos/catalog-service/internal/models"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductsRepositoryInterface interface {
	GetProductByID(ctx context.Context, id uint64) (*models.Product, error)
	GetProductByName(ctx context.Context, name string) (*models.Product, error)
	GetListProduct(ctx context.Context, filter models.ProductFilter, limit, offset int32) ([]*models.Product, int32, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, productID uint64) error
	CreateProduct(ctx context.Context, product *models.Product) (uint64, error)
}

type ProductsRepository struct {
	DB *pgxpool.Pool
}

func NewProductsRepository(db *pgxpool.Pool) *ProductsRepository {
	return &ProductsRepository{
		DB: db,
	}
}

func (r *ProductsRepository) GetProductByID(ctx context.Context, id uint64) (*models.Product, error) {
	var product models.Product
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT * FROM products WHERE product_id=$1", id).Scan(&product.ID, &product.Name, &product.ManufacturersID, &product.CategoryID, &product.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}
	return &product, nil
}

func (r *ProductsRepository) GetProductByName(ctx context.Context, name string) (*models.Product, error) {
	var products models.Product
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT * FROM products WHERE product_name=$1", name).Scan(&products.ID, &products.Name, &products.ManufacturersID, &products.CategoryID, &products.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by name: %w", err)
	}
	return &products, nil
}

func (r *ProductsRepository) GetListProduct(ctx context.Context, filter models.ProductFilter, limit, offset int32) ([]*models.Product, int32, error) {
	var products []*models.Product
	querier := GetQuerier(ctx, r.DB)
	where, args := buildFilterQuery(filter)
	query := fmt.Sprintf(`SELECT  product_id, product_name, manufacturers_id, category_id, price FROM products %s LIMIT $%d OFFSET $%d`, where, len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	rows, err := querier.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.ManufacturersID, &product.CategoryID, &product.Price); err != nil {
			return nil, 0, fmt.Errorf("failed to scan products: %w", err)
		}
		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}
	var total int32
	// Для COUNT запроса нужно пересоздать аргументы (без limit/offset)
	_, countArgs := buildFilterQuery(filter)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", where)
	err = querier.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total products: %w", err)
	}

	return products, total, nil

}

func buildFilterQuery(filter models.ProductFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	var where string
	argIdx := 1
	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.ManufacturersID != nil {
		conditions = append(conditions, fmt.Sprintf("manufacturers_id = $%d", argIdx))
		args = append(args, *filter.ManufacturersID)
		argIdx++
	}
	if filter.PriceMin != nil {
		conditions = append(conditions, fmt.Sprintf("price >=$%d", argIdx))
		args = append(args, *filter.PriceMin)
		argIdx++
	}
	if filter.PriceMax != nil {
		conditions = append(conditions, fmt.Sprintf("price <=$%d", argIdx))
		args = append(args, *filter.PriceMax)
		argIdx++
	}
	if filter.NameContains != nil {
		conditions = append(conditions, fmt.Sprintf("product_name ILIKE $%d", argIdx))
		args = append(args, *filter.NameContains)
		argIdx++
	}
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}
	return where, args
}

func (r *ProductsRepository) CreateProduct(ctx context.Context, product *models.Product) (uint64, error) {
	var productID uint64
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "INSERT INTO products (product_name, manufacturers_id, category_id, price) VALUES ($1, $2, $3, $4) RETURNING product_id", product.Name, product.ManufacturersID, product.CategoryID, product.Price).Scan(&productID)
	if err != nil {
		return 0, fmt.Errorf("failed to update product: %w", err)
	}
	return productID, nil
}

func (r *ProductsRepository) DeleteProduct(ctx context.Context, productID uint64) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "DELETE FROM products WHERE product_id=$1", productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (r *ProductsRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "UPDATE products SET product_name=$1, manufacturers_id=$2, category_id=$3, price=$4 WHERE product_id=$5", product.Name, product.ManufacturersID, product.CategoryID, product.Price, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}
