// Package service предоставляет бизнес-логику для управления товарами.
// Включает создание, обновление, удаление и получение товаров с фильтрацией.
package service

import (
	"context"
	"fmt"
	"kinos/catalog-service/internal/models"
	"kinos/catalog-service/internal/repository"
	"kinos/catalog-service/internal/validator"
)

type ProductServiceInterface interface {
	CreateProduct(ctx context.Context, name string, manufacturerID, categoryID uint64, price float64) (*models.Product, error)
	UpdateProduct(ctx context.Context, productID uint64, name string, manufacturerID, categoryID uint64, price float64) (*models.Product, error)
	DeleteProduct(ctx context.Context, productID uint64) error
	GetProduct(ctx context.Context, name string) (*models.Product, error)
	GetListProduct(ctx context.Context, filter models.ProductFilter, limit, offset int32) ([]*models.Product, int32, error)
}
type ProductService struct {
	prodRep   repository.ProductsRepository
	manRep    repository.ManufacturersRepository
	catRep    repository.CategoryRepository
	val       validator.ValidatorInterface
	txManager *repository.TxManager
}

func NewProductService(prodRep *repository.ProductsRepository, manRep *repository.ManufacturersRepository, catRep *repository.CategoryRepository, validator validator.ValidatorInterface, txManager *repository.TxManager) *ProductService {
	return &ProductService{
		prodRep:   *prodRep,
		manRep:    *manRep,
		catRep:    *catRep,
		val:       validator,
		txManager: txManager,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string, manufacturerID, categoryID uint64, price float64) (*models.Product, error) {
	if err := s.val.ValidateProduct(validator.ProductInput{
		Name:            name,
		ManufacturersID: manufacturerID,
		CategoryID:      categoryID,
		Price:           price,
	}); err != nil {
		return nil, fmt.Errorf("failed validate product: %w", err)
	}
	var product *models.Product
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		if _, err := s.manRep.GetManufacturerByID(txCtx, manufacturerID); err != nil {
			return fmt.Errorf("failed get Manufacturer by ID: %w", err)
		}
		if _, err := s.catRep.GetCategoryByID(txCtx, categoryID); err != nil {
			return fmt.Errorf("failed get Category by ID: %w", err)
		}
		productID, err := s.prodRep.CreateProduct(txCtx, &models.Product{
			Name:            name,
			ManufacturersId: manufacturerID,
			CategoryId:      categoryID,
			Price:           price,
		})
		if err != nil {
			return fmt.Errorf("failed create product: %w", err)
		}
		product = &models.Product{
			Id:   productID,
			Name: name,
		}
		return nil

	})
	return product, err
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID uint64, name string, manufacturerID, categoryID uint64, price float64) (*models.Product, error) {
	if err := s.val.ValidateProduct(validator.ProductInput{
		Name:            name,
		ManufacturersID: manufacturerID,
		CategoryID:      categoryID,
		Price:           price,
	}); err != nil {
		return nil, fmt.Errorf("failed validate product: %w", err)
	}

	var product *models.Product
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		if _, err := s.manRep.GetManufacturerByID(txCtx, manufacturerID); err != nil {
			return fmt.Errorf("failed get Manufacturer by ID: %w", err)
		}
		if _, err := s.catRep.GetCategoryByID(txCtx, categoryID); err != nil {
			return fmt.Errorf("failed get Category by ID: %w", err)
		}
		product = &models.Product{
			Id:              productID,
			Name:            name,
			ManufacturersId: manufacturerID,
			CategoryId:      categoryID,
			Price:           price,
		}
		if err := s.prodRep.UpdateProduct(txCtx, product); err != nil {
			return fmt.Errorf("failed update product: %w", err)
		}
		product = &models.Product{
			Id: productID,
		}
		return nil
	})
	return product, err
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID uint64) error {
	return s.prodRep.DeleteProduct(ctx, productID)
}

func (s *ProductService) GetProduct(ctx context.Context, name string) (*models.Product, error) {
	product, err := s.prodRep.GetProductByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed get product by name: %w", err)
	}
	return product, nil
}

func (s *ProductService) GetListProduct(ctx context.Context, filter models.ProductFilter, limit, offset int32) ([]*models.Product, int32, error) {
	products, total, err := s.prodRep.GetListProduct(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get products: %w", err)
	}
	return products, total, nil
}
