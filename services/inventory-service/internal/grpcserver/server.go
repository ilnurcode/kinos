// Package grpcserver предоставляет gRPC-сервер для inventory-service.
// Реализует методы InventoryService: CreateInventory, UpdateInventory, GetInventory, ReserveStock и другие.
package grpcserver

import (
	"context"
	"log"
	"time"

	"kinos/inventory-service/internal/service"
	pb "kinos/proto/inventory"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	inventorySvc *service.InventoryService
	warehouseSvc *service.WarehouseService
}

func NewInventoryServer(inventorySvc *service.InventoryService, warehouseSvc *service.WarehouseService) *InventoryServer {
	return &InventoryServer{
		inventorySvc: inventorySvc,
		warehouseSvc: warehouseSvc,
	}
}

func (s *InventoryServer) CreateInventory(ctx context.Context, req *pb.CreateInventoryRequest) (*pb.Inventory, error) {
	log.Printf("CreateInventory request: product_id=%d, quantity=%d, location=%s", req.ProductId, req.Quantity, req.WarehouseLocation)

	inventory, err := s.inventorySvc.CreateInventory(ctx, req.ProductId, req.Quantity, req.WarehouseLocation)
	if err != nil {
		log.Printf("CreateInventory error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create inventory: %v", err)
	}

	return &pb.Inventory{
		Id:                inventory.Id,
		ProductId:         inventory.ProductId,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: inventory.AvailableQuantity,
		WarehouseLocation: inventory.WarehouseLocation,
		UpdatedAt:         inventory.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *InventoryServer) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryRequest) (*pb.Inventory, error) {
	log.Printf("UpdateInventory request: id=%d, quantity=%d, location=%s", req.Id, req.Quantity, req.WarehouseLocation)

	inventory, err := s.inventorySvc.UpdateInventory(ctx, req.Id, req.Quantity, req.WarehouseLocation)
	if err != nil {
		log.Printf("UpdateInventory error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update inventory: %v", err)
	}

	return &pb.Inventory{
		Id:                inventory.Id,
		ProductId:         inventory.ProductId,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: inventory.AvailableQuantity,
		WarehouseLocation: inventory.WarehouseLocation,
		UpdatedAt:         inventory.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *InventoryServer) GetInventory(ctx context.Context, req *pb.GetInventoryRequest) (*pb.Inventory, error) {
	log.Printf("GetInventory request: product_id=%d", req.ProductId)

	inventory, err := s.inventorySvc.GetInventoryByProductID(ctx, req.ProductId)
	if err != nil {
		log.Printf("GetInventory error: %v", err)
		return nil, status.Errorf(codes.NotFound, "inventory not found: %v", err)
	}

	return &pb.Inventory{
		Id:                inventory.Id,
		ProductId:         inventory.ProductId,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: inventory.AvailableQuantity,
		WarehouseLocation: inventory.WarehouseLocation,
		UpdatedAt:         inventory.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *InventoryServer) GetListInventory(ctx context.Context, req *pb.GetListInventoryRequest) (*pb.ListInventoryResponse, error) {
	log.Printf("GetListInventory request: limit=%d, offset=%d", req.Limit, req.Offset)

	limit := req.Limit
	offset := req.Offset
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	inventories, total, err := s.inventorySvc.GetListInventory(ctx, limit, offset, req.ProductId, req.WarehouseLocation, req.MinQuantity)
	if err != nil {
		log.Printf("GetListInventory error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get inventory list: %v", err)
	}

	var result []*pb.Inventory
	for _, inv := range inventories {
		result = append(result, &pb.Inventory{
			Id:                inv.Id,
			ProductId:         inv.ProductId,
			Quantity:          inv.Quantity,
			ReservedQuantity:  inv.ReservedQuantity,
			AvailableQuantity: inv.AvailableQuantity,
			WarehouseLocation: inv.WarehouseLocation,
			UpdatedAt:         inv.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListInventoryResponse{
		Inventory: result,
		Total:     total,
	}, nil
}

func (s *InventoryServer) DeleteInventory(ctx context.Context, req *pb.DeleteInventoryRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteInventory request: id=%d", req.Id)

	err := s.inventorySvc.DeleteInventory(ctx, req.Id)
	if err != nil {
		log.Printf("DeleteInventory error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete inventory: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *InventoryServer) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	log.Printf("ReserveStock request: product_id=%d, quantity=%d, reservation_id=%s", req.ProductId, req.Quantity, req.ReservationId)

	err := s.inventorySvc.ReserveStock(ctx, req.ProductId, req.Quantity, req.ReservationId)
	if err != nil {
		log.Printf("ReserveStock error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to reserve stock: %v", err)
	}

	return &pb.ReserveStockResponse{
		Success:          true,
		ReservationId:    req.ReservationId,
		ReservedQuantity: req.Quantity,
	}, nil
}

func (s *InventoryServer) ReleaseReservation(ctx context.Context, req *pb.ReleaseReservationRequest) (*pb.ReleaseReservationResponse, error) {
	log.Printf("ReleaseReservation request: product_id=%d, reservation_id=%s", req.ProductId, req.ReservationId)

	releasedQuantity, err := s.inventorySvc.ReleaseReservation(ctx, req.ProductId, req.ReservationId)
	if err != nil {
		log.Printf("ReleaseReservation error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to release reservation: %v", err)
	}

	return &pb.ReleaseReservationResponse{
		Success:          true,
		ReleasedQuantity: releasedQuantity,
	}, nil
}

// CreateWarehouse создает новый склад
func (s *InventoryServer) CreateWarehouse(ctx context.Context, req *pb.CreateWarehouseRequest) (*pb.Warehouse, error) {
	log.Printf("CreateWarehouse request: name=%s, city=%s", req.Name, req.City)

	warehouse, err := s.warehouseSvc.CreateWarehouse(ctx, req.Name, req.City, req.Street, req.Building, req.Building2)
	if err != nil {
		log.Printf("CreateWarehouse error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create warehouse: %v", err)
	}

	return &pb.Warehouse{
		Id:        warehouse.Id,
		Name:      warehouse.Name,
		City:      warehouse.City,
		Street:    warehouse.Street,
		Building:  warehouse.Building,
		Building2: warehouse.Building2,
		CreatedAt: warehouse.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *InventoryServer) GetListWarehouse(ctx context.Context, req *pb.GetListWarehouseRequest) (*pb.ListWarehouseResponse, error) {
	log.Printf("GetListWarehouse request: limit=%d, offset=%d", req.Limit, req.Offset)

	limit := req.Limit
	offset := req.Offset
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	warehouses, total, err := s.warehouseSvc.GetListWarehouse(ctx, limit, offset)
	if err != nil {
		log.Printf("GetListWarehouse error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get warehouses: %v", err)
	}

	var result []*pb.Warehouse
	for _, w := range warehouses {
		result = append(result, &pb.Warehouse{
			Id:        w.Id,
			Name:      w.Name,
			City:      w.City,
			Street:    w.Street,
			Building:  w.Building,
			Building2: w.Building2,
			CreatedAt: w.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListWarehouseResponse{
		Warehouses: result,
		Total:      total,
	}, nil
}

func (s *InventoryServer) DeleteWarehouse(ctx context.Context, req *pb.DeleteWarehouseRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteWarehouse request: id=%d", req.Id)

	err := s.warehouseSvc.DeleteWarehouse(ctx, req.Id)
	if err != nil {
		log.Printf("DeleteWarehouse error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete warehouse: %v", err)
	}

	return &emptypb.Empty{}, nil
}
