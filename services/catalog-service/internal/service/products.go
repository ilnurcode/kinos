// Package service предоставляет бизнес-логику для управления товарами.
// Включает создание, обновление, удаление и получение товаров с фильтрацией.
package service

import (
	"context"
	"fmt"

	"kinos/catalog-service/internal/errs"
	"kinos/catalog-service/internal/models"
	"kinos/catalog-service/internal/repository"
	"kinos/catalog-service/internal/validator"
)

// Алиасы ошибок из пакета errs
var (
	ErrProductNotFound = errs.ErrProductNotFound
	ErrProductExists   = errs.ErrProductExists
)

type ProductServiceInterface interface {
	CreateProduct(ctx context.Context, name string, manufacturerID, categoryID uint64, price float64) (*models.Product, error)
	UpdateProduct(ctx context.Context, productID uint64, name string, manufacturerID, categoryID uint64, price float64) (*models.Product, error)
	DeleteProduct(ctx context.Context, productID uint64) error
	GetProduct(ctx context.Context, name string) (*models.Product, error)
	GetProductByID(ctx context.Context, productID uint64) (*models.Product, error)
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
		return nil, fmt.Errorf("ошибка валидации товара: %w", err)
	}
	var product *models.Product
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		// Проверяем существование производителя
		if _, err := s.manRep.GetManufacturerByID(txCtx, manufacturerID); err != nil {
			return errs.ErrManufacturerNotFound
		}
		// Проверяем существование категории
		if _, err := s.catRep.GetCategoryByID(txCtx, categoryID); err != nil {
			return errs.ErrCategoryNotFound
		}
		// Проверяем существование товара с таким именем
		existing, _ := s.prodRep.GetProductByName(txCtx, name)
		if existing != nil {
			return errs.ErrProductExists
		}
		productID, err := s.prodRep.CreateProduct(txCtx, &models.Product{
			Name:            name,
			ManufacturersID: manufacturerID,
			CategoryID:      categoryID,
			Price:           price,
		})
		if err != nil {
			return fmt.Errorf("ошибка создания товара: %w", err)
		}
		product = &models.Product{
			ID:   productID,
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
		return nil, fmt.Errorf("ошибка валидации товара: %w", err)
	}

	var product *models.Product
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		// Проверяем существование производителя
		if _, err := s.manRep.GetManufacturerByID(txCtx, manufacturerID); err != nil {
			return errs.ErrManufacturerNotFound
		}
		// Проверяем существование категории
		if _, err := s.catRep.GetCategoryByID(txCtx, categoryID); err != nil {
			return errs.ErrCategoryNotFound
		}
		// Проверяем существование товара по ID
		if _, err := s.prodRep.GetProductByID(txCtx, productID); err != nil {
			return errs.ErrProductNotFound
		}
		// Проверяем, не используется ли имя другим товаром
		productByName, _ := s.prodRep.GetProductByName(txCtx, name)
		if productByName != nil && productByName.ID != productID {
			return errs.ErrProductExists
		}
		product = &models.Product{
			ID:              productID,
			Name:            name,
			ManufacturersID: manufacturerID,
			CategoryID:      categoryID,
			Price:           price,
		}
		if err := s.prodRep.UpdateProduct(txCtx, product); err != nil {
			return fmt.Errorf("ошибка обновления товара: %w", err)
		}
		return nil
	})
	return product, err
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID uint64) error {
	// Проверяем существование товара
	_, err := s.prodRep.GetProductByID(ctx, productID)
	if err != nil {
		return errs.ErrProductNotFound
	}
	return s.prodRep.DeleteProduct(ctx, productID)
}

func (s *ProductService) GetProduct(ctx context.Context, name string) (*models.Product, error) {
	product, err := s.prodRep.GetProductByName(ctx, name)
	if err != nil {
		return nil, errs.ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, productID uint64) (*models.Product, error) {
	product, err := s.prodRep.GetProductByID(ctx, productID)
	if err != nil {
		return nil, errs.ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) GetListProduct(ctx context.Context, filter models.ProductFilter, limit, offset int32) ([]*models.Product, int32, error) {
	products, total, err := s.prodRep.GetListProduct(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения списка товаров: %w", err)
	}
	return products, total, nil
}
