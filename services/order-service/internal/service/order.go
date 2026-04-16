package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"kinos/order-service/internal/errs"
	inventoryclient "kinos/order-service/internal/inventory"
	"kinos/order-service/internal/models"
	"kinos/order-service/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	orderRepo repository.OrderRepositoryInterface
	txManager repository.TransactionManager
	inventory inventoryclient.Client
}

func NewOrderService(orderRepo repository.OrderRepositoryInterface, txManager repository.TransactionManager, inventory inventoryclient.Client) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		txManager: txManager,
		inventory: inventory,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint64, items []models.OrderItem, address, phone, comment string) (*models.Order, error) {
	if userID == 0 {
		return nil, errs.ErrInvalidUserID
	}
	if len(items) == 0 {
		return nil, errs.ErrEmptyCart
	}

	order := &models.Order{
		UserID:          userID,
		Items:           items,
		Total:           calculateTotal(items),
		Status:          models.StatusPending,
		DeliveryAddress: address,
		Phone:           phone,
		Comment:         comment,
	}

	var (
		createdOrder  *models.Order
		reservationID string
		reservedItems []models.OrderItem
	)

	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		orderID, err := s.orderRepo.CreateOrder(txCtx, order)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		reservationID = strconv.FormatUint(orderID, 10)
		for i := range items {
			items[i].OrderID = orderID
			if err := s.orderRepo.AddOrderItem(txCtx, &items[i]); err != nil {
				return fmt.Errorf("failed to add order item: %w", err)
			}
		}

		for _, item := range items {
			if err := s.inventory.ReserveStock(ctx, item.ProductID, int32(item.Quantity), reservationID); err != nil {
				return mapInventoryError(err)
			}
			reservedItems = append(reservedItems, item)
		}

		createdOrder, err = s.orderRepo.GetOrder(txCtx, orderID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.releaseReservedItems(ctx, reservationID, reservedItems)
		return nil, err
	}

	return createdOrder, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uint64) (*models.Order, error) {
	if orderID == 0 {
		return nil, errs.ErrInvalidOrderID
	}

	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (s *OrderService) GetListOrders(ctx context.Context, limit, offset int32, status string) ([]*models.Order, int32, error) {
	return s.orderRepo.GetListOrders(ctx, limit, offset, status)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uint64, limit, offset int32) ([]*models.Order, int32, error) {
	return s.orderRepo.GetUserOrders(ctx, userID, limit, offset)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uint64) (*models.Order, error) {
	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order for cancellation: %w", err)
	}
	if order.Status != models.StatusPending {
		return nil, errs.ErrCannotCancelOrder
	}

	reservationID := strconv.FormatUint(orderID, 10)
	releasedItems := make([]models.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		releasedQuantity, err := s.inventory.ReleaseReservation(ctx, item.ProductID, reservationID)
		if err != nil {
			s.reReserveItems(ctx, reservationID, releasedItems)
			return nil, mapInventoryError(err)
		}
		if releasedQuantity > 0 {
			releasedItems = append(releasedItems, item)
		}
	}

	if err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		return s.orderRepo.UpdateOrderStatus(txCtx, orderID, models.StatusCancelled)
	}); err != nil {
		s.reReserveItems(ctx, reservationID, releasedItems)
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	return s.GetOrder(ctx, orderID)
}

func calculateTotal(items []models.OrderItem) float64 {
	var total float64
	for i := range items {
		total += items[i].Subtotal
	}
	return total
}

func (s *OrderService) releaseReservedItems(ctx context.Context, reservationID string, items []models.OrderItem) {
	if reservationID == "" {
		return
	}
	for _, item := range items {
		if _, err := s.inventory.ReleaseReservation(ctx, item.ProductID, reservationID); err != nil {
			log.Printf("failed to release reservation for product_id=%d reservation_id=%s: %v", item.ProductID, reservationID, err)
		}
	}
}

func (s *OrderService) reReserveItems(ctx context.Context, reservationID string, items []models.OrderItem) {
	if reservationID == "" {
		return
	}
	for _, item := range items {
		if err := s.inventory.ReserveStock(ctx, item.ProductID, int32(item.Quantity), reservationID); err != nil {
			log.Printf("failed to re-reserve product_id=%d reservation_id=%s: %v", item.ProductID, reservationID, err)
		}
	}
}

func mapInventoryError(err error) error {
	switch status.Code(err) {
	case codes.FailedPrecondition:
		return errs.ErrInsufficientStock
	case codes.NotFound:
		return errs.ErrInventoryNotFound
	default:
		return err
	}
}
