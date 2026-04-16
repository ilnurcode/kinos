// Package grpcserver предоставляет gRPC-сервер для cart-service.
package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"kinos/cart-service/internal/service"
	pb "kinos/proto/cart"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CartServer struct {
	pb.UnimplementedCartServiceServer
	cartService service.CartServiceInterface
}

func NewCartServer(cartService service.CartServiceInterface) *CartServer {
	return &CartServer{
		cartService: cartService,
	}
}

// AddItem добавляет товар в корзину
func (s *CartServer) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.Cart, error) {
	log.Printf("AddItem request: user_id=%d, product_id=%d, quantity=%d", req.UserId, req.ProductId, req.Quantity)

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id обязателен")
	}
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id обязателен")
	}
	if req.Quantity == 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity должен быть больше 0")
	}

	cart, err := s.cartService.AddItem(ctx, req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "товар не найден")
		}
		if errors.Is(err, service.ErrInsufficientStock) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		log.Printf("AddItem error: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка добавления в корзину: %v", err)
	}

	return cartToProto(cart), nil
}

// GetCart получает корзину пользователя
func (s *CartServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.Cart, error) {
	log.Printf("GetCart request: user_id=%d", req.UserId)

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id обязателен")
	}

	cart, err := s.cartService.GetCart(ctx, req.UserId)
	if err != nil {
		log.Printf("GetCart error: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка получения корзины: %v", err)
	}

	return cartToProto(cart), nil
}

// RemoveItem удаляет товар из корзины
func (s *CartServer) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.Cart, error) {
	log.Printf("RemoveItem request: user_id=%d, product_id=%d", req.UserId, req.ProductId)

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id обязателен")
	}
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id обязателен")
	}

	cart, err := s.cartService.RemoveItem(ctx, req.UserId, req.ProductId)
	if err != nil {
		log.Printf("RemoveItem error: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка удаления товара: %v", err)
	}

	return cartToProto(cart), nil
}

// UpdateItem обновляет количество товара в корзине
func (s *CartServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.Cart, error) {
	log.Printf("UpdateItem request: user_id=%d, product_id=%d, quantity=%d", req.UserId, req.ProductId, req.Quantity)

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id обязателен")
	}
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id обязателен")
	}

	cart, err := s.cartService.UpdateItem(ctx, req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientStock) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		log.Printf("UpdateItem error: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка обновления товара: %v", err)
	}

	return cartToProto(cart), nil
}

// ClearCart очищает корзину
func (s *CartServer) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*emptypb.Empty, error) {
	log.Printf("ClearCart request: user_id=%d", req.UserId)

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id обязателен")
	}

	err := s.cartService.ClearCart(ctx, req.UserId)
	if err != nil {
		log.Printf("ClearCart error: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка очистки корзины: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetItemsCount возвращает количество товаров в корзине
func (s *CartServer) GetItemsCount(ctx context.Context, req *pb.GetItemsCountRequest) (*pb.ItemsCountResponse, error) {
	log.Printf("GetItemsCount request: user_id=%d", req.UserId)

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id обязателен")
	}

	count, err := s.cartService.GetItemsCount(ctx, req.UserId)
	if err != nil {
		log.Printf("GetItemsCount error: %v", err)
		return nil, status.Errorf(codes.Internal, "ошибка получения количества: %v", err)
	}

	return &pb.ItemsCountResponse{Count: int32(count)}, nil
}

// cartToProto преобразует модель корзины в proto сообщение
func cartToProto(cart *service.Cart) *pb.Cart {
	if cart == nil {
		return &pb.Cart{}
	}

	items := make([]*pb.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, &pb.CartItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			AddedAt:     item.AddedAt.Format(time.RFC3339),
		})
	}

	return &pb.Cart{
		UserId:    cart.UserID,
		Items:     items,
		Total:     cart.Total,
		UpdatedAt: cart.UpdatedAt.Format(time.RFC3339),
	}
}

// Start запускает gRPC сервер
func (s *CartServer) Start(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCartServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	log.Printf("Cart gRPC server started on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
