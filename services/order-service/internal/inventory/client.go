package inventory

import (
	"context"

	pb "kinos/proto/inventory"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error
	ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error)
	Close() error
}

type grpcClient struct {
	client pb.InventoryServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &grpcClient{
		client: pb.NewInventoryServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *grpcClient) ReserveStock(ctx context.Context, productID uint64, quantity int32, reservationID string) error {
	_, err := c.client.ReserveStock(ctx, &pb.ReserveStockRequest{
		ProductId:     productID,
		Quantity:      quantity,
		ReservationId: reservationID,
	})
	return err
}

func (c *grpcClient) ReleaseReservation(ctx context.Context, productID uint64, reservationID string) (int32, error) {
	resp, err := c.client.ReleaseReservation(ctx, &pb.ReleaseReservationRequest{
		ProductId:     productID,
		ReservationId: reservationID,
	})
	if err != nil {
		return 0, err
	}
	return resp.ReleasedQuantity, nil
}

func (c *grpcClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
