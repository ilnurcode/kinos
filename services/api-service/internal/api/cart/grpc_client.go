// Package cart предоставляет gRPC-клиент для взаимодействия с cart-service.
package cart

import (
	"context"

	pb "kinos/proto/cart"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CartClient интерфейс для работы с корзиной
type CartClient interface {
	GetCart(ctx context.Context, userID uint64) (*pb.Cart, error)
	AddItem(ctx context.Context, userID, productID uint64, quantity uint32) (*pb.Cart, error)
	RemoveItem(ctx context.Context, userID, productID uint64) (*pb.Cart, error)
	UpdateItem(ctx context.Context, userID, productID uint64, quantity uint32) (*pb.Cart, error)
	ClearCart(ctx context.Context, userID uint64) error
	GetItemsCount(ctx context.Context, userID uint64) (int, error)
	Close() error
}

// CartClientImpl реализация CartClient
type CartClientImpl struct {
	client pb.CartServiceClient
	conn   *grpc.ClientConn
}

// NewCartClient создаёт новый gRPC клиент к Cart Service
func NewCartClient(address string) (*CartClientImpl, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &CartClientImpl{
		client: pb.NewCartServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close закрывает соединение
func (c *CartClientImpl) Close() error {
	return c.conn.Close()
}

func (c *CartClientImpl) GetCart(ctx context.Context, userID uint64) (*pb.Cart, error) {
	req := &pb.GetCartRequest{UserId: userID}
	return c.client.GetCart(ctx, req)
}

func (c *CartClientImpl) AddItem(ctx context.Context, userID, productID uint64, quantity uint32) (*pb.Cart, error) {
	req := &pb.AddItemRequest{UserId: userID, ProductId: productID, Quantity: quantity}
	return c.client.AddItem(ctx, req)
}

func (c *CartClientImpl) RemoveItem(ctx context.Context, userID, productID uint64) (*pb.Cart, error) {
	req := &pb.RemoveItemRequest{UserId: userID, ProductId: productID}
	return c.client.RemoveItem(ctx, req)
}

func (c *CartClientImpl) UpdateItem(ctx context.Context, userID, productID uint64, quantity uint32) (*pb.Cart, error) {
	req := &pb.UpdateItemRequest{UserId: userID, ProductId: productID, Quantity: quantity}
	return c.client.UpdateItem(ctx, req)
}

func (c *CartClientImpl) ClearCart(ctx context.Context, userID uint64) error {
	req := &pb.ClearCartRequest{UserId: userID}
	_, err := c.client.ClearCart(ctx, req)
	return err
}

func (c *CartClientImpl) GetItemsCount(ctx context.Context, userID uint64) (int, error) {
	req := &pb.GetItemsCountRequest{UserId: userID}
	resp, err := c.client.GetItemsCount(ctx, req)
	if err != nil {
		return 0, err
	}
	return int(resp.Count), nil
}
