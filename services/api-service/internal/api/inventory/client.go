// Package inventory предоставляет gRPC-клиент для связи с inventory-service.
package inventory

import (
	"context"

	pb "kinos/proto/inventory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryClientInterface interface {
	GetInventory(ctx context.Context, productID uint64) (*pb.Inventory, error)
	GetListInventory(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) (*pb.ListInventoryResponse, error)
	CreateInventory(ctx context.Context, productID uint64, quantity int32, location string) (*pb.Inventory, error)
	UpdateInventory(ctx context.Context, id uint64, quantity int32, location string) (*pb.Inventory, error)
	DeleteInventory(ctx context.Context, id uint64) error
	ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) (*pb.ReserveStockResponse, error)
	ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (*pb.ReleaseReservationResponse, error)

	// Warehouse methods
	GetListWarehouse(ctx context.Context, limit, offset int32) (*pb.ListWarehouseResponse, error)
	CreateWarehouse(ctx context.Context, name, city, street, building, building2 string) (*pb.Warehouse, error)
	UpdateWarehouse(ctx context.Context, id uint64, name, city, street, building, building2 string) (*pb.Warehouse, error)
	DeleteWarehouse(ctx context.Context, id uint64) error
}

type InventoryClient struct {
	client pb.InventoryServiceClient
	conn   *grpc.ClientConn
}

func NewInventoryClient(address string) *InventoryClient {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &InventoryClient{
		client: pb.NewInventoryServiceClient(conn),
		conn:   conn,
	}
}

func (c *InventoryClient) Close() error {
	return c.conn.Close()
}

func (c *InventoryClient) GetInventory(ctx context.Context, productID uint64) (*pb.Inventory, error) {
	return c.client.GetInventory(ctx, &pb.GetInventoryRequest{ProductId: productID})
}

func (c *InventoryClient) GetListInventory(ctx context.Context, limit, offset int32, productID uint64, location string, minQuantity int32) (*pb.ListInventoryResponse, error) {
	return c.client.GetListInventory(ctx, &pb.GetListInventoryRequest{
		Limit:             limit,
		Offset:            offset,
		ProductId:         productID,
		WarehouseLocation: location,
		MinQuantity:       minQuantity,
	})
}

func (c *InventoryClient) CreateInventory(ctx context.Context, productID uint64, quantity int32, location string) (*pb.Inventory, error) {
	return c.client.CreateInventory(ctx, &pb.CreateInventoryRequest{
		ProductId:         productID,
		Quantity:          quantity,
		WarehouseLocation: location,
	})
}

func (c *InventoryClient) UpdateInventory(ctx context.Context, id uint64, quantity int32, location string) (*pb.Inventory, error) {
	return c.client.UpdateInventory(ctx, &pb.UpdateInventoryRequest{
		Id:                id,
		Quantity:          quantity,
		WarehouseLocation: location,
	})
}

func (c *InventoryClient) DeleteInventory(ctx context.Context, id uint64) error {
	_, err := c.client.DeleteInventory(ctx, &pb.DeleteInventoryRequest{Id: id})
	return err
}

func (c *InventoryClient) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) (*pb.ReserveStockResponse, error) {
	return c.client.ReserveStock(ctx, &pb.ReserveStockRequest{
		ProductId:     productID,
		Quantity:      quantity,
		ReservationId: reservationID,
	})
}

func (c *InventoryClient) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (*pb.ReleaseReservationResponse, error) {
	return c.client.ReleaseReservation(ctx, &pb.ReleaseReservationRequest{
		ProductId:     productID,
		ReservationId: reservationID,
	})
}

// Warehouse methods
func (c *InventoryClient) GetListWarehouse(ctx context.Context, limit, offset int32) (*pb.ListWarehouseResponse, error) {
	return c.client.GetListWarehouse(ctx, &pb.GetListWarehouseRequest{
		Limit:  limit,
		Offset: offset,
	})
}

func (c *InventoryClient) CreateWarehouse(ctx context.Context, name, city, street, building, building2 string) (*pb.Warehouse, error) {
	return c.client.CreateWarehouse(ctx, &pb.CreateWarehouseRequest{
		Name:      name,
		City:      city,
		Street:    street,
		Building:  building,
		Building2: building2,
	})
}

func (c *InventoryClient) UpdateWarehouse(ctx context.Context, id uint64, name, city, street, building, building2 string) (*pb.Warehouse, error) {
	return c.client.UpdateWarehouse(ctx, &pb.UpdateWarehouseRequest{
		Id:        id,
		Name:      name,
		City:      city,
		Street:    street,
		Building:  building,
		Building2: building2,
	})
}

func (c *InventoryClient) DeleteWarehouse(ctx context.Context, id uint64) error {
	_, err := c.client.DeleteWarehouse(ctx, &pb.DeleteWarehouseRequest{Id: id})
	return err
}
