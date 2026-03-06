// Package service предоставляет бизнес-логику для управления производителями.
// Включает создание, обновление, удаление и получение производителей.
package service

import (
	"context"
	"fmt"
	"kinos/catalog-service/internal/models"
	"kinos/catalog-service/internal/repository"
	"kinos/catalog-service/internal/validator"
)

type ManufacturersServiceInterface interface {
	GetListManufacturers(ctx context.Context, limit, offset int32)
	GetManufacturer(ctx context.Context, name string) (*models.Manufacturer, error)
	DeleteManufacturer(ctx context.Context, id uint64) error
	UpdateManufacturer(ctx context.Context, id uint64, name string) (*models.Manufacturer, error)
	CreateManufacturer(ctx context.Context, name string) (*models.Manufacturer, error)
}

type ManufacturersService struct {
	repo      repository.ManufacturersRepositoryInterface
	validator validator.ValidatorInterface
	txManager *repository.TxManager
}

func NewManufacturersService(repo repository.ManufacturersRepositoryInterface, validator validator.ValidatorInterface, txManager *repository.TxManager) *ManufacturersService {
	return &ManufacturersService{repo: repo, validator: validator, txManager: txManager}
}

func (s *ManufacturersService) CreateManufacturer(ctx context.Context, name string) (*models.Manufacturer, error) {
	if err := s.validator.ValidateManufactures(validator.ManufacturersInput{Name: name}); err != nil {
		return nil, fmt.Errorf("failed validate manufacturers: %v", err)
	}
	var manufacturer *models.Manufacturer
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		manufacturerID, err := s.repo.CreateManufacturers(txCtx, &models.Manufacturer{Name: name})
		if err != nil {
			return fmt.Errorf("failed create manufacturers: %v", err)
		}
		manufacturer = &models.Manufacturer{Id: manufacturerID, Name: name}
		return nil
	})
	return manufacturer, err
}

func (s *ManufacturersService) UpdateManufacturer(ctx context.Context, id uint64, name string) (*models.Manufacturer, error) {
	if err := s.validator.ValidateManufactures(validator.ManufacturersInput{Name: name}); err != nil {
		return nil, fmt.Errorf("failed validate manufacturers: %v", err)
	}
	var manufacturer *models.Manufacturer
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		manufacturer = &models.Manufacturer{Id: id, Name: name}
		if err := s.repo.UpdateManufacturers(txCtx, manufacturer); err != nil {
			return fmt.Errorf("failed update manufacturers: %v", err)
		}
		return nil
	})
	return manufacturer, err
}

func (s *ManufacturersService) DeleteManufacturer(ctx context.Context, id uint64) error {
	return s.repo.DeleteManufacturers(ctx, id)
}

func (s *ManufacturersService) GetManufacturer(ctx context.Context, name string) (*models.Manufacturer, error) {
	return s.repo.GetManufacturerByName(ctx, name)
}

func (s *ManufacturersService) GetListManufacturers(ctx context.Context, limit, offset int32) ([]*models.Manufacturer, int32, error) {
	return s.repo.GetListManufacturers(ctx, limit, offset)
}
