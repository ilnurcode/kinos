package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"kinos/order-service/internal/errs"
	"kinos/order-service/internal/models"
	"kinos/order-service/internal/service"
	pb "kinos/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	orderSvc *service.OrderService
}

func NewOrderServer(orderSvc *service.OrderService) *OrderServer {
	return &OrderServer{orderSvc: orderSvc}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cart is empty")
	}

	items := make([]models.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, models.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		})
	}

	order, err := s.orderSvc.CreateOrder(ctx, req.UserId, items, req.DeliveryAddress, req.Phone, req.Comment)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrEmptyCart), errors.Is(err, errs.ErrInvalidUserID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errs.ErrInsufficientStock):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, errs.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
		}
	}

	return orderToProto(order), nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if req.OrderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, err := s.orderSvc.GetOrder(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	return orderToProto(order), nil
}

func (s *OrderServer) GetListOrders(ctx context.Context, req *pb.GetListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, total, err := s.orderSvc.GetListOrders(ctx, req.Limit, req.Offset, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	pbOrders := make([]*pb.Order, 0, len(orders))
	for _, order := range orders {
		pbOrders = append(pbOrders, orderToProto(order))
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
		Total:  total,
	}, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.Order, error) {
	if req.OrderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, err := s.orderSvc.CancelOrder(ctx, req.OrderId)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound), errors.Is(err, errs.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, errs.ErrCannotCancelOrder):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
		}
	}

	return orderToProto(order), nil
}

func (s *OrderServer) GetUserOrders(ctx context.Context, req *pb.GetUserOrdersRequest) (*pb.ListOrdersResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	orders, total, err := s.orderSvc.GetUserOrders(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list user orders: %v", err)
	}

	pbOrders := make([]*pb.Order, 0, len(orders))
	for _, order := range orders {
		pbOrders = append(pbOrders, orderToProto(order))
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
		Total:  total,
	}, nil
}

func orderToProto(order *models.Order) *pb.Order {
	if order == nil {
		return &pb.Order{}
	}

	items := make([]*pb.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &pb.OrderItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		})
	}

	return &pb.Order{
		Id:              order.ID,
		UserId:          order.UserID,
		Items:           items,
		Total:           order.Total,
		Status:          string(order.Status),
		DeliveryAddress: order.DeliveryAddress,
		Phone:           order.Phone,
		Comment:         order.Comment,
		CreatedAt:       order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       order.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *OrderServer) Start(port string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Order gRPC server starting on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			serverErr <- fmt.Errorf("gRPC server error: %w", err)
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return nil, err
	case <-time.After(100 * time.Millisecond):
	}

	return grpcServer, nil
}
