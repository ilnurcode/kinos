// Package service предоставляет бизнес-логику для управления категориями.
// Включает создание, обновление, удаление и получение категорий.
package service

import (
	"context"
	"fmt"
	"kinos/catalog-service/internal/models"
	"kinos/catalog-service/internal/repository"
	"kinos/catalog-service/internal/validator"
)

type CategoryServiceInterface interface {
	CreateCategory(ctx context.Context, name string) (*models.Category, error)
	UpdateCategory(ctx context.Context, id uint64, name string) (*models.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
	GetCategory(ctx context.Context, name string) (*models.Category, error)
	GetListCategory(ctx context.Context, limit, offset int32) ([]*models.Category, int32, error)
}

type txManagerInterface interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type CategoryService struct {
	repo      repository.CategoryRepositoryInterface
	validator validator.ValidatorInterface
	txManager txManagerInterface
}

func NewCategoryService(rep repository.CategoryRepositoryInterface, validator validator.ValidatorInterface, txManager txManagerInterface) *CategoryService {
	return &CategoryService{repo: rep, validator: validator, txManager: txManager}
}

func (s *CategoryService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	if err := s.validator.ValidateCategory(validator.CategoryInput{Name: name}); err != nil {
		return nil, fmt.Errorf("validate category error: %v", err)
	}
	var category *models.Category
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		categoryID, err := s.repo.CreateCategory(txCtx, &models.Category{Name: name})
		if err != nil {
			return fmt.Errorf("failed create category: %v", err)
		}
		category = &models.Category{Id: categoryID, Name: name}
		return nil
	})
	return category, err
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint64, name string) (*models.Category, error) {
	if err := s.validator.ValidateCategory(validator.CategoryInput{Name: name}); err != nil {
		return nil, fmt.Errorf("validate category error: %v", err)
	}
	var category *models.Category
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		category = &models.Category{Id: id, Name: name}
		if err := s.repo.UpdateCategory(txCtx, category); err != nil {
			return fmt.Errorf("failed update category: %v", err)
		}
		return nil
	})
	return category, err
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *CategoryService) GetCategory(ctx context.Context, name string) (*models.Category, error) {
	return s.repo.GetCategoryByName(ctx, name)
}

func (s *CategoryService) GetListCategory(ctx context.Context, limit, offset int32) ([]*models.Category, int32, error) {
	return s.repo.GetListCategory(ctx, limit, offset)
}
