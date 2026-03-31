// Package service предоставляет бизнес-логику для управления категориями.
// Включает создание, обновление, удаление и получение категорий.
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
	ErrCategoryNotFound = errs.ErrCategoryNotFound
	ErrCategoryExists   = errs.ErrCategoryExists
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
		return nil, fmt.Errorf("ошибка валидации категории: %v", err)
	}
	var category *models.Category
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		existing, _ := s.repo.GetCategoryByName(txCtx, name)
		if existing != nil {
			return errs.ErrCategoryExists
		}
		categoryID, err := s.repo.CreateCategory(txCtx, &models.Category{Name: name})
		if err != nil {
			return fmt.Errorf("ошибка создания категории: %v", err)
		}
		category = &models.Category{Id: categoryID, Name: name}
		return nil
	})
	return category, err
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint64, name string) (*models.Category, error) {
	if err := s.validator.ValidateCategory(validator.CategoryInput{Name: name}); err != nil {
		return nil, fmt.Errorf("ошибка валидации категории: %v", err)
	}
	var category *models.Category
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		// Проверяем существует ли категория с таким именем (кроме текущей)
		existing, _ := s.repo.GetCategoryByName(txCtx, name)
		if existing != nil && existing.Id != id {
			return errs.ErrCategoryExists
		}
		// Проверяем существует ли категория с таким id
		_, err := s.repo.GetCategoryByID(txCtx, id)
		if err != nil {
			return errs.ErrCategoryNotFound
		}
		category = &models.Category{Id: id, Name: name}
		if err := s.repo.UpdateCategory(txCtx, category); err != nil {
			return fmt.Errorf("ошибка обновления категории: %v", err)
		}
		return nil
	})
	return category, err
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint64) error {
	// Проверяем существование категории
	_, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return errs.ErrCategoryNotFound
	}
	return s.repo.DeleteCategory(ctx, id)
}

func (s *CategoryService) GetCategory(ctx context.Context, name string) (*models.Category, error) {
	category, err := s.repo.GetCategoryByName(ctx, name)
	if err != nil {
		return nil, errs.ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) GetListCategory(ctx context.Context, limit, offset int32) ([]*models.Category, int32, error) {
	return s.repo.GetListCategory(ctx, limit, offset)
}
