package service

import (
	"context"
	"testing"
	"time"

	"kinos/cart-service/internal/models"
	pbCatalog "kinos/proto/catalog"
	pbInventory "kinos/proto/inventory"
)

type mockCartRepo struct {
	cart *models.Cart
}

func (m *mockCartRepo) AddItem(ctx context.Context, userID, productID uint64, productName string, price float64, quantity uint32) error {
	for _, item := range m.cart.Items {
		if item.ProductID == productID {
			item.Quantity += quantity
			return nil
		}
	}
	m.cart.Items = append(m.cart.Items, &models.CartItem{
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		Price:       price,
		AddedAt:     time.Now(),
	})
	return nil
}

func (m *mockCartRepo) UpdateItem(ctx context.Context, userID, productID uint64, quantity uint32) error {
	return nil
}

func (m *mockCartRepo) GetCart(ctx context.Context, userID uint64) (*models.Cart, error) {
	return m.cart, nil
}

func (m *mockCartRepo) SaveCart(ctx context.Context, cart *models.Cart) error {
	m.cart = cart
	return nil
}

func (m *mockCartRepo) RemoveItem(ctx context.Context, userID, productID uint64) error {
	return nil
}

func (m *mockCartRepo) ClearCart(ctx context.Context, userID uint64) error {
	return nil
}

func (m *mockCartRepo) GetItemsCount(ctx context.Context, userID uint64) (int, error) {
	return len(m.cart.Items), nil
}

type mockCatalogClient struct{}

func (m *mockCatalogClient) GetProductByID(ctx context.Context, id uint64) (*pbCatalog.Product, error) {
	return &pbCatalog.Product{Id: id, Name: "Item", Price: 100}, nil
}

func (m *mockCatalogClient) Close() error {
	return nil
}

type mockInventoryClient struct {
	available int32
}

func (m *mockInventoryClient) GetInventory(ctx context.Context, productID uint64) (*pbInventory.Inventory, error) {
	return &pbInventory.Inventory{ProductId: productID, AvailableQuantity: m.available}, nil
}

func (m *mockInventoryClient) GetListInventory(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) (*pbInventory.ListInventoryResponse, error) {
	return nil, nil
}

func (m *mockInventoryClient) CreateInventory(ctx context.Context, productID uint64, quantity int32, location string) (*pbInventory.Inventory, error) {
	return nil, nil
}

func (m *mockInventoryClient) UpdateInventory(ctx context.Context, id uint64, quantity int32, location string) (*pbInventory.Inventory, error) {
	return nil, nil
}

func (m *mockInventoryClient) DeleteInventory(ctx context.Context, id uint64) error {
	return nil
}

func (m *mockInventoryClient) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) (*pbInventory.ReserveStockResponse, error) {
	return nil, nil
}

func (m *mockInventoryClient) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (*pbInventory.ReleaseReservationResponse, error) {
	return nil, nil
}

func (m *mockInventoryClient) GetListWarehouse(ctx context.Context, limit, offset int32) (*pbInventory.ListWarehouseResponse, error) {
	return nil, nil
}

func (m *mockInventoryClient) CreateWarehouse(ctx context.Context, name, city, street, building, building2 string) (*pbInventory.Warehouse, error) {
	return nil, nil
}

func (m *mockInventoryClient) UpdateWarehouse(ctx context.Context, id uint64, name, city, street, building, building2 string) (*pbInventory.Warehouse, error) {
	return nil, nil
}

func (m *mockInventoryClient) DeleteWarehouse(ctx context.Context, id uint64) error {
	return nil
}

func TestCartService_AddItemChecksTotalQuantity(t *testing.T) {
	repo := &mockCartRepo{
		cart: &models.Cart{
			UserID: 1,
			Items: []*models.CartItem{
				{ProductID: 10, ProductName: "Item", Quantity: 3, Price: 100},
			},
		},
	}

	svc := NewCartService(repo, &mockCatalogClient{}, &mockInventoryClient{available: 5})

	_, err := svc.AddItem(context.Background(), 1, 10, 3)
	if err != ErrInsufficientStock {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
}
