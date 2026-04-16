package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kinos/cart-service/internal/catalog"
	"kinos/cart-service/internal/inventory"
	"kinos/cart-service/internal/repository"
)

// Ошибки сервиса
var (
	ErrProductNotFound   = errors.New("товар не найден")
	ErrInsufficientStock = errors.New("недостаточно товара на складе")
	ErrInvalidQuantity   = errors.New("недопустимое количество")
	ErrCartNotFound      = errors.New("корзина не найдена")
)

// CartItem элемент корзины
type CartItem struct {
	ProductID   uint64
	ProductName string
	Quantity    uint32
	Price       float64
	AddedAt     time.Time
}

// Cart корзина
type Cart struct {
	UserID    uint64
	Items     []CartItem
	Total     float64
	UpdatedAt time.Time
}

type CartServiceInterface interface {
	AddItem(ctx context.Context, userID, productID uint64, quantity uint32) (*Cart, error)
	GetCart(ctx context.Context, userID uint64) (*Cart, error)
	RemoveItem(ctx context.Context, userID, productID uint64) (*Cart, error)
	UpdateItem(ctx context.Context, userID, productID uint64, quantity uint32) (*Cart, error)
	ClearCart(ctx context.Context, userID uint64) error
	GetItemsCount(ctx context.Context, userID uint64) (int, error)
}

type CartService struct {
	repo      repository.CartRepositoryInterface
	catalog   catalog.CatalogClient
	inventory inventory.InventoryClientInterface
}

func NewCartService(repo repository.CartRepositoryInterface, catalog catalog.CatalogClient, inventory inventory.InventoryClientInterface) *CartService {
	return &CartService{
		repo:      repo,
		catalog:   catalog,
		inventory: inventory,
	}
}

func (c *CartService) AddItem(ctx context.Context, userID, productID uint64, quantity uint32) (*Cart, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	// Получаем информацию о товаре из Catalog
	product, err := c.catalog.GetProductByID(ctx, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	currentCart, err := c.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения корзины: %w", err)
	}

	totalQuantity := quantity
	for _, item := range currentCart.Items {
		if item.ProductID == productID {
			totalQuantity += item.Quantity
			break
		}
	}

	// Проверяем наличие в Inventory
	inv, err := c.inventory.GetInventory(ctx, productID)
	if err != nil {
		return nil, ErrInsufficientStock
	}

	if inv.AvailableQuantity < int32(totalQuantity) {
		return nil, ErrInsufficientStock
	}

	// Добавляем в корзину
	err = c.repo.AddItem(ctx, userID, productID, product.Name, product.Price, quantity)
	if err != nil {
		return nil, fmt.Errorf("ошибка добавления в корзину: %w", err)
	}

	return c.getCart(ctx, userID)
}

func (c *CartService) GetCart(ctx context.Context, userID uint64) (*Cart, error) {
	return c.getCart(ctx, userID)
}

func (c *CartService) RemoveItem(ctx context.Context, userID, productID uint64) (*Cart, error) {
	err := c.repo.RemoveItem(ctx, userID, productID)
	if err != nil {
		return nil, err
	}
	return c.getCart(ctx, userID)
}

func (c *CartService) UpdateItem(ctx context.Context, userID, productID uint64, quantity uint32) (*Cart, error) {
	if quantity <= 0 {
		return c.RemoveItem(ctx, userID, productID)
	}

	// Проверяем наличие в Inventory
	inv, err := c.inventory.GetInventory(ctx, productID)
	if err != nil {
		return nil, ErrInsufficientStock
	}

	if inv.AvailableQuantity < int32(quantity) {
		return nil, ErrInsufficientStock
	}

	err = c.repo.UpdateItem(ctx, userID, productID, quantity)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления товара: %w", err)
	}

	return c.getCart(ctx, userID)
}

func (c *CartService) ClearCart(ctx context.Context, userID uint64) error {
	return c.repo.ClearCart(ctx, userID)
}

func (c *CartService) GetItemsCount(ctx context.Context, userID uint64) (int, error) {
	return c.repo.GetItemsCount(ctx, userID)
}

// getCart вспомогательная функция для получения корзины
func (c *CartService) getCart(ctx context.Context, userID uint64) (*Cart, error) {
	repoCart, err := c.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Конвертируем в нашу модель
	cart := &Cart{
		UserID:    repoCart.UserID,
		Items:     make([]CartItem, 0, len(repoCart.Items)),
		Total:     repoCart.Total,
		UpdatedAt: repoCart.UpdatedAt,
	}

	for _, item := range repoCart.Items {
		cart.Items = append(cart.Items, CartItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    uint32(item.Quantity),
			Price:       item.Price,
			AddedAt:     item.AddedAt,
		})
	}

	return cart, nil
}
