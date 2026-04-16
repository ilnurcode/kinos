package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"kinos/cart-service/internal/models"

	"github.com/redis/go-redis/v9"
)

type CartRepositoryInterface interface {
	AddItem(ctx context.Context, userID, productID uint64, productName string, price float64, quantity uint32) error
	UpdateItem(ctx context.Context, userID, productID uint64, quantity uint32) error
	GetCart(ctx context.Context, userID uint64) (*models.Cart, error)
	SaveCart(ctx context.Context, cart *models.Cart) error
	RemoveItem(ctx context.Context, userID, productID uint64) error
	ClearCart(ctx context.Context, userID uint64) error
	GetItemsCount(ctx context.Context, userID uint64) (int, error)
}

type CartRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCartRepository(client *redis.Client, ttl time.Duration) *CartRepository {
	return &CartRepository{
		client: client,
		ttl:    ttl,
	}
}

func (c *CartRepository) AddItem(ctx context.Context, userID, productID uint64, productName string, price float64, quantity uint32) error {
	cart, err := c.GetCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %v", err)
	}
	found := false
	for _, item := range cart.Items {
		if item.ProductID == productID {
			item.Quantity += quantity
			item.Price = price
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, &models.CartItem{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    quantity,
			Price:       price,
			AddedAt:     time.Now(),
		})
	}
	return c.SaveCart(ctx, cart)
}

func (c *CartRepository) GetCart(ctx context.Context, userID uint64) (*models.Cart, error) {
	data, err := c.client.Get(ctx, strconv.FormatUint(userID, 10)).Result()
	if err == redis.Nil {
		return &models.Cart{
			UserID: userID,
			Items:  make([]*models.CartItem, 0),
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}
	var cart *models.Cart
	if err := json.Unmarshal([]byte(data), &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshall cart: %v", err)
	}
	return cart, nil
}

func (c *CartRepository) SaveCart(ctx context.Context, cart *models.Cart) error {
	cart.UpdatedAt = time.Now()
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += float64(item.Quantity) * item.Price
	}
	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %v", err)
	}
	if err := c.client.SetEx(ctx, strconv.FormatUint(cart.UserID, 10), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to save cart: %v", err)
	}
	return nil
}

func (c *CartRepository) RemoveItem(ctx context.Context, userID, productID uint64) error {
	cart, err := c.GetCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %v", err)
	}
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}
	return c.SaveCart(ctx, cart)
}

func (c *CartRepository) ClearCart(ctx context.Context, userID uint64) error {
	return c.client.Del(ctx, strconv.FormatUint(userID, 10)).Err()
}

func (c *CartRepository) GetItemsCount(ctx context.Context, userID uint64) (int, error) {
	cart, err := c.GetCart(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get cart: %v", err)
	}
	return len(cart.Items), nil
}

func (c *CartRepository) UpdateItem(ctx context.Context, userID, productID uint64, quantity uint32) error {
	cart, err := c.GetCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %v", err)
	}
	for _, item := range cart.Items {
		if item.ProductID == productID {
			if quantity <= 0 {
				return c.RemoveItem(ctx, userID, productID)
			}
			item.Quantity = quantity
			break
		}
	}
	return c.SaveCart(ctx, cart)
}
