// Package service предоставляет бизнес-логику для управления складами.
package service

import (
	"context"
	"errors"

	"kinos/inventory-service/internal/model"
	"kinos/inventory-service/internal/repository"
)

type WarehouseService struct {
	warehouseRepo repository.WarehouseInterface
	txManager     repository.TxManagerInterface
}

func NewWarehouseService(warehouseRepo repository.WarehouseInterface, txManager repository.TxManagerInterface) *WarehouseService {
	return &WarehouseService{
		warehouseRepo: warehouseRepo,
		txManager:     txManager,
	}
}

func (s *WarehouseService) CreateWarehouse(ctx context.Context, name, city, street, building, building2 string) (*model.Warehouse, error) {
	if name == "" || city == "" || street == "" {
		return nil, errors.New("название, город и улица обязательны")
	}

	var warehouse *model.Warehouse
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error
		warehouse, err = s.warehouseRepo.Create(txCtx, name, city, street, building, building2)
		return err
	})
	return warehouse, err
}

func (s *WarehouseService) GetListWarehouse(ctx context.Context, limit, offset int32) ([]*model.Warehouse, int32, error) {
	return s.warehouseRepo.GetList(ctx, limit, offset)
}

func (s *WarehouseService) DeleteWarehouse(ctx context.Context, id uint64) error {
	return s.warehouseRepo.Delete(ctx, id)
}
