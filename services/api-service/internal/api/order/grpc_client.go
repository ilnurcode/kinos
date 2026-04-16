// Package order предоставляет gRPC-клиент для взаимодействия с order-service.
package order

import (
	"context"

	pbCart "kinos/proto/cart"
	pb "kinos/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrderClient интерфейс для работы с заказами
type OrderClient interface {
	CreateOrder(ctx context.Context, userID uint64, items []*OrderItem, address, phone, comment string) (*pb.Order, error)
	GetOrder(ctx context.Context, orderID uint64) (*pb.Order, error)
	GetListOrders(ctx context.Context, limit, offset int32, status string) (*pb.ListOrdersResponse, error)
	GetUserOrders(ctx context.Context, userID uint64, limit, offset int32) (*pb.ListOrdersResponse, error)
	CancelOrder(ctx context.Context, orderID uint64) (*pb.Order, error)
	Close() error
}

// CartClient интерфейс для работы с корзиной
type CartClient interface {
	GetCart(ctx context.Context, userID uint64) (*pbCart.Cart, error)
	ClearCart(ctx context.Context, userID uint64) error
}

// OrderItem элемент заказа
type OrderItem struct {
	ProductID   uint64
	ProductName string
	Quantity    uint32
	Price       float64
	Subtotal    float64
}

// OrderClientImpl реализация OrderClient
type OrderClientImpl struct {
	client pb.OrderServiceClient
	conn   *grpc.ClientConn
}

// NewOrderClient создаёт новый gRPC клиент к Order Service
func NewOrderClient(address string) (*OrderClientImpl, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &OrderClientImpl{
		client: pb.NewOrderServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close закрывает соединение
func (c *OrderClientImpl) Close() error {
	return c.conn.Close()
}

func (c *OrderClientImpl) CreateOrder(ctx context.Context, userID uint64, items []*OrderItem, address, phone, comment string) (*pb.Order, error) {
	pbItems := make([]*pb.OrderItem, 0, len(items))
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		})
	}

	req := &pb.CreateOrderRequest{
		UserId:          userID,
		Items:           pbItems,
		DeliveryAddress: address,
		Phone:           phone,
		Comment:         comment,
	}
	return c.client.CreateOrder(ctx, req)
}

func (c *OrderClientImpl) GetOrder(ctx context.Context, orderID uint64) (*pb.Order, error) {
	req := &pb.GetOrderRequest{OrderId: orderID}
	return c.client.GetOrder(ctx, req)
}

func (c *OrderClientImpl) GetListOrders(ctx context.Context, limit, offset int32, status string) (*pb.ListOrdersResponse, error) {
	req := &pb.GetListOrdersRequest{
		Limit:  limit,
		Offset: offset,
		Status: status,
	}
	return c.client.GetListOrders(ctx, req)
}

func (c *OrderClientImpl) GetUserOrders(ctx context.Context, userID uint64, limit, offset int32) (*pb.ListOrdersResponse, error) {
	req := &pb.GetUserOrdersRequest{
		UserId: userID,
		Limit:  limit,
		Offset: offset,
	}
	return c.client.GetUserOrders(ctx, req)
}

func (c *OrderClientImpl) CancelOrder(ctx context.Context, orderID uint64) (*pb.Order, error) {
	req := &pb.CancelOrderRequest{OrderId: orderID}
	return c.client.CancelOrder(ctx, req)
}
