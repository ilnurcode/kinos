// Package service предоставляет бизнес-логику для inventory-service.
package service

import (
	"context"
	"kinos/inventory-service/internal/model"
	"kinos/inventory-service/internal/repository"
	"kinos/inventory-service/internal/validator"
)

type InventoryService struct {
	inventoryRepo repository.InventoryInterface
	validator     *validator.Validator
	txManager     repository.TxManagerInterface
}

func NewInventoryService(inventoryRepo repository.InventoryInterface, validator *validator.Validator, txManager repository.TxManagerInterface) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		validator:     validator,
		txManager:     txManager,
	}
}

func (s *InventoryService) CreateInventory(ctx context.Context, productID uint64, quantity int32, location string) (*model.Inventory, error) {
	input := validator.CreateInventoryInput{
		ProductID: productID,
		Quantity:  quantity,
		Location:  location,
	}
	if err := s.validator.ValidateCreateInventory(input); err != nil {
		return nil, err
	}

	var inventory *model.Inventory
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error
		inventory, err = s.inventoryRepo.Create(txCtx, productID, quantity, location)
		return err
	})
	return inventory, err
}

func (s *InventoryService) UpdateInventory(ctx context.Context, id uint64, quantity int32, location string) (*model.Inventory, error) {
	input := validator.UpdateInventoryInput{
		ID:       id,
		Quantity: quantity,
		Location: location,
	}
	if err := s.validator.ValidateUpdateInventory(input); err != nil {
		return nil, err
	}

	var inventory *model.Inventory
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error
		inventory, err = s.inventoryRepo.Update(txCtx, id, quantity, location)
		return err
	})
	return inventory, err
}

func (s *InventoryService) GetInventoryByProductID(ctx context.Context, productID uint64) (*model.Inventory, error) {
	return s.inventoryRepo.GetByProductID(ctx, productID)
}

func (s *InventoryService) GetListInventory(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) ([]*model.Inventory, int32, error) {
	return s.inventoryRepo.GetList(ctx, limit, offset, productID, location, minQuantity)
}

func (s *InventoryService) DeleteInventory(ctx context.Context, id uint64) error {
	return s.inventoryRepo.Delete(ctx, id)
}

func (s *InventoryService) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error {
	if quantity <= 0 {
		return validator.ErrInvalidQuantity
	}

	return s.txManager.Do(ctx, func(txCtx context.Context) error {
		return s.inventoryRepo.ReserveStock(txCtx, productID, quantity, reservationID)
	})
}

func (s *InventoryService) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error) {
	var released int32
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error
		released, err = s.inventoryRepo.ReleaseReservation(txCtx, productID, reservationID)
		return err
	})
	return released, err
}
