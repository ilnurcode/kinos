package service

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"kinos/order-service/internal/errs"
	"kinos/order-service/internal/models"
)

type mockOrderRepo struct {
	orders     map[uint64]*models.Order
	nextID     uint64
	createErr  error
	getErr     error
	listErr    error
	updateErr  error
	itemAddErr error
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{
		orders: make(map[uint64]*models.Order),
		nextID: 1,
	}
}

func (m *mockOrderRepo) CreateOrder(ctx context.Context, order *models.Order) (uint64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}
	id := m.nextID
	m.nextID++
	copyOrder := *order
	copyOrder.ID = id
	copyOrder.CreatedAt = time.Now()
	copyOrder.UpdatedAt = time.Now()
	m.orders[id] = &copyOrder
	return id, nil
}

func (m *mockOrderRepo) AddOrderItem(ctx context.Context, item *models.OrderItem) error {
	if m.itemAddErr != nil {
		return m.itemAddErr
	}
	order, ok := m.orders[item.OrderID]
	if !ok {
		return errs.ErrOrderNotFound
	}
	order.Items = append(order.Items, *item)
	return nil
}

func (m *mockOrderRepo) GetOrder(ctx context.Context, id uint64) (*models.Order, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	order, ok := m.orders[id]
	if !ok {
		return nil, errs.ErrOrderNotFound
	}
	copyOrder := *order
	copyOrder.Items = append([]models.OrderItem(nil), order.Items...)
	return &copyOrder, nil
}

func (m *mockOrderRepo) GetListOrders(ctx context.Context, limit, offset int32, status string) ([]*models.Order, int32, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var result []*models.Order
	for _, order := range m.orders {
		if status == "" || string(order.Status) == status {
			copyOrder := *order
			copyOrder.Items = append([]models.OrderItem(nil), order.Items...)
			result = append(result, &copyOrder)
		}
	}
	return result, int32(len(result)), nil
}

func (m *mockOrderRepo) GetUserOrders(ctx context.Context, userID uint64, limit, offset int32) ([]*models.Order, int32, error) {
	var result []*models.Order
	for _, order := range m.orders {
		if order.UserID == userID {
			copyOrder := *order
			copyOrder.Items = append([]models.OrderItem(nil), order.Items...)
			result = append(result, &copyOrder)
		}
	}
	return result, int32(len(result)), nil
}

func (m *mockOrderRepo) UpdateOrderStatus(ctx context.Context, id uint64, status models.OrderStatus) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	order, ok := m.orders[id]
	if !ok {
		return errs.ErrOrderNotFound
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	return nil
}

type mockInventoryClient struct {
	reserveErr  error
	releaseErr  error
	reserveIDs  []string
	releaseIDs  []string
	reserveHits int
	releaseHits int
}

func (m *mockInventoryClient) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error {
	m.reserveHits++
	m.reserveIDs = append(m.reserveIDs, reservationID+":"+strconv.FormatUint(productID, 10))
	return m.reserveErr
}

func (m *mockInventoryClient) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error) {
	m.releaseHits++
	m.releaseIDs = append(m.releaseIDs, reservationID+":"+strconv.FormatUint(productID, 10))
	if m.releaseErr != nil {
		return 0, m.releaseErr
	}
	return 1, nil
}

func (m *mockInventoryClient) Close() error {
	return nil
}

type mockTxManager struct{}

func (m *mockTxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestOrderService_CreateOrderReservesInventory(t *testing.T) {
	repo := newMockOrderRepo()
	inventory := &mockInventoryClient{}
	svc := NewOrderService(repo, &mockTxManager{}, inventory)

	order, err := svc.CreateOrder(context.Background(), 1, []models.OrderItem{
		{ProductID: 10, ProductName: "Item", Quantity: 2, Price: 50, Subtotal: 100},
	}, "addr", "phone", "comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.ID == 0 {
		t.Fatal("expected order ID to be set")
	}
	if order.Status != models.StatusPending {
		t.Fatalf("expected pending status, got %s", order.Status)
	}
	if inventory.reserveHits != 1 {
		t.Fatalf("expected 1 reserve call, got %d", inventory.reserveHits)
	}
	if inventory.reserveIDs[0] != strconv.FormatUint(order.ID, 10)+":10" {
		t.Fatalf("unexpected reservation key %s", inventory.reserveIDs[0])
	}
}

func TestOrderService_CreateOrderReleasesOnReserveFailure(t *testing.T) {
	repo := newMockOrderRepo()
	inventory := &mockInventoryClient{reserveErr: errs.ErrInsufficientStock}
	svc := NewOrderService(repo, &mockTxManager{}, inventory)

	_, err := svc.CreateOrder(context.Background(), 1, []models.OrderItem{
		{ProductID: 10, ProductName: "Item", Quantity: 2, Price: 50, Subtotal: 100},
	}, "addr", "phone", "comment")
	if !errors.Is(err, errs.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
	if inventory.releaseHits != 0 {
		t.Fatalf("expected no release because nothing was reserved, got %d", inventory.releaseHits)
	}
}

func TestOrderService_CancelOrderCompensatesOnStatusFailure(t *testing.T) {
	repo := newMockOrderRepo()
	repo.orders[1] = &models.Order{
		ID:     1,
		UserID: 1,
		Status: models.StatusPending,
		Items: []models.OrderItem{
			{OrderID: 1, ProductID: 10, Quantity: 1},
		},
	}
	repo.updateErr = errors.New("update failed")

	inventory := &mockInventoryClient{}
	svc := NewOrderService(repo, &mockTxManager{}, inventory)

	_, err := svc.CancelOrder(context.Background(), 1)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
	if inventory.releaseHits != 1 {
		t.Fatalf("expected release to be called once, got %d", inventory.releaseHits)
	}
	if inventory.reserveHits != 1 {
		t.Fatalf("expected compensation reserve to be called once, got %d", inventory.reserveHits)
	}
}
