package service

import (
	"context"
	"strconv"
	"testing"
	"time"

	"kinos/inventory-service/internal/errs"
	"kinos/inventory-service/internal/model"
	"kinos/inventory-service/internal/validator"
)

type mockInventoryRepo struct {
	inventories  map[uint64]*model.Inventory
	reservations map[string]int32
	nextID       uint64
}

func newMockInventoryRepo() *mockInventoryRepo {
	return &mockInventoryRepo{
		inventories:  make(map[uint64]*model.Inventory),
		reservations: make(map[string]int32),
		nextID:       1,
	}
}

func reservationKey(reservationID string, productID uint64) string {
	return reservationID + ":" + strconv.FormatUint(productID, 10)
}

func (m *mockInventoryRepo) Create(ctx context.Context, productID uint64, quantity int32, location string) (*model.Inventory, error) {
	inv := &model.Inventory{
		Id:                m.nextID,
		ProductId:         productID,
		Quantity:          quantity,
		ReservedQuantity:  0,
		AvailableQuantity: quantity,
		WarehouseLocation: location,
		UpdatedAt:         time.Now(),
	}
	m.nextID++
	m.inventories[inv.Id] = inv
	return inv, nil
}

func (m *mockInventoryRepo) GetByProductID(ctx context.Context, productID uint64) (*model.Inventory, error) {
	for _, inv := range m.inventories {
		if inv.ProductId == productID {
			return inv, nil
		}
	}
	return nil, errs.ErrInventoryNotFound
}

func (m *mockInventoryRepo) Update(ctx context.Context, id uint64, quantity int32, location string) (*model.Inventory, error) {
	inv, ok := m.inventories[id]
	if !ok {
		return nil, errs.ErrInventoryNotFound
	}
	inv.Quantity = quantity
	inv.AvailableQuantity = quantity - inv.ReservedQuantity
	inv.WarehouseLocation = location
	inv.UpdatedAt = time.Now()
	return inv, nil
}

func (m *mockInventoryRepo) Delete(ctx context.Context, id uint64) error {
	delete(m.inventories, id)
	return nil
}

func (m *mockInventoryRepo) GetList(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) ([]*model.Inventory, int32, error) {
	result := make([]*model.Inventory, 0)
	for _, inv := range m.inventories {
		if productID > 0 && inv.ProductId != productID {
			continue
		}
		if location != "" && inv.WarehouseLocation != location {
			continue
		}
		if minQuantity > 0 && inv.AvailableQuantity < minQuantity {
			continue
		}
		result = append(result, inv)
	}
	return result, int32(len(result)), nil
}

func (m *mockInventoryRepo) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error {
	key := reservationKey(reservationID, productID)
	if _, exists := m.reservations[key]; exists {
		return nil
	}

	for _, inv := range m.inventories {
		if inv.ProductId == productID {
			if inv.AvailableQuantity < quantity {
				return errs.ErrInsufficientStock
			}
			inv.ReservedQuantity += quantity
			inv.AvailableQuantity -= quantity
			m.reservations[key] = quantity
			return nil
		}
	}
	return errs.ErrInventoryNotFound
}

func (m *mockInventoryRepo) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error) {
	key := reservationKey(reservationID, productID)
	quantity, exists := m.reservations[key]
	if !exists {
		return 0, nil
	}
	delete(m.reservations, key)

	for _, inv := range m.inventories {
		if inv.ProductId == productID {
			inv.ReservedQuantity -= quantity
			inv.AvailableQuantity += quantity
			return quantity, nil
		}
	}
	return 0, errs.ErrInventoryNotFound
}

type mockTxManager struct{}

func (m *mockTxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestInventoryService_ReserveAndReleaseStock(t *testing.T) {
	repo := newMockInventoryRepo()
	val := &validator.Validator{}
	txMgr := &mockTxManager{}
	svc := NewInventoryService(repo, val, txMgr)

	_, _ = svc.CreateInventory(context.Background(), 100, 50, "warehouse-A")

	if err := svc.ReserveStock(context.Background(), 100, 10, "res-1"); err != nil {
		t.Fatalf("unexpected reserve error: %v", err)
	}

	if err := svc.ReserveStock(context.Background(), 100, 10, "res-1"); err != nil {
		t.Fatalf("expected idempotent reserve, got %v", err)
	}

	inv, _ := svc.GetInventoryByProductID(context.Background(), 100)
	if inv.ReservedQuantity != 10 {
		t.Fatalf("expected reserved 10, got %d", inv.ReservedQuantity)
	}
	if inv.AvailableQuantity != 40 {
		t.Fatalf("expected available 40, got %d", inv.AvailableQuantity)
	}

	released, err := svc.ReleaseReservation(context.Background(), 100, "res-1")
	if err != nil {
		t.Fatalf("unexpected release error: %v", err)
	}
	if released != 10 {
		t.Fatalf("expected released 10, got %d", released)
	}

	inv, _ = svc.GetInventoryByProductID(context.Background(), 100)
	if inv.ReservedQuantity != 0 {
		t.Fatalf("expected reserved 0, got %d", inv.ReservedQuantity)
	}
	if inv.AvailableQuantity != 50 {
		t.Fatalf("expected available 50, got %d", inv.AvailableQuantity)
	}
}

func TestInventoryService_ValidatesReservationID(t *testing.T) {
	repo := newMockInventoryRepo()
	val := &validator.Validator{}
	txMgr := &mockTxManager{}
	svc := NewInventoryService(repo, val, txMgr)

	if err := svc.ReserveStock(context.Background(), 100, 1, ""); err != errs.ErrInvalidReservationID {
		t.Fatalf("expected ErrInvalidReservationID, got %v", err)
	}

	if _, err := svc.ReleaseReservation(context.Background(), 100, ""); err != errs.ErrInvalidReservationID {
		t.Fatalf("expected ErrInvalidReservationID, got %v", err)
	}
}
