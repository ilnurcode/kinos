// Package order предоставляет HTTP-обработчики для управления заказами.
package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	OrderClient OrderClient
	CartClient  CartClient
}

func NewHandler(oc OrderClient, cc CartClient) *Handler {
	return &Handler{
		OrderClient: oc,
		CartClient:  cc,
	}
}

// CreateOrder godoc
// @Summary Создать новый заказ
// @Description Создать новый заказ из текущей корзины пользователя
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrderRequest true "Информация о заказе"
// @Success 200 {object} object "Созданный заказ"
// @Failure 400 {object} object{error=string} "Некорректные данные"
// @Failure 401 {object} object{error=string} "Требуется аутентификация"
// @Router /api/orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}

	// Получаем корзину пользователя
	cart, err := h.CartClient.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения корзины"})
		return
	}

	if cart == nil || len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Корзина пуста"})
		return
	}

	var body struct {
		DeliveryAddress string `json:"delivery_address"`
		Phone           string `json:"phone"`
		Comment         string `json:"comment"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	if body.DeliveryAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Адрес доставки обязателен"})
		return
	}

	if body.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Телефон обязателен"})
		return
	}

	// Конвертируем элементы корзины в элементы заказа
	orderItems := make([]*OrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		orderItems = append(orderItems, &OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Price * float64(item.Quantity),
		})
	}

	// Создаем заказ
	order, err := h.OrderClient.CreateOrder(c.Request.Context(), userID, orderItems, body.DeliveryAddress, body.Phone, body.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания заказа"})
		return
	}

	// Очищаем корзину
	_ = h.CartClient.ClearCart(c.Request.Context(), userID)

	c.JSON(http.StatusCreated, gin.H{"order": order})
}

// CreateOrderRequest запрос на создание заказа
type CreateOrderRequest struct {
	DeliveryAddress string `json:"delivery_address"`
	Phone           string `json:"phone"`
	Comment         string `json:"comment"`
}

// GetOrder godoc
// @Summary Получить заказ по ID
// @Description Получить информацию о заказе по его ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID заказа"
// @Success 200 {object} object "Заказ"
// @Failure 404 {object} object{error=string} "Заказ не найден"
// @Router /api/orders/:id [get]
func (h *Handler) GetOrder(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заказа"})
		return
	}

	order, err := h.OrderClient.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}

	// Проверяем, что заказ принадлежит пользователю
	if order.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// GetUserOrders godoc
// @Summary Получить заказы пользователя
// @Description Получить список заказов текущего пользователя
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int32 false "Лимит" default(20)
// @Param offset query int32 false "Смещение" default(0)
// @Success 200 {object} object "Список заказов"
// @Router /api/orders/my [get]
func (h *Handler) GetUserOrders(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	orders, err := h.OrderClient.GetUserOrders(c.Request.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заказов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders.Orders, "total": orders.Total})
}

// CancelOrder godoc
// @Summary Отменить заказ
// @Description Отменить заказ (только в статусе pending)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID заказа"
// @Success 200 {object} object "Отмененный заказ"
// @Failure 400 {object} object{error=string} "Нельзя отменить заказ"
// @Router /api/orders/:id/cancel [post]
func (h *Handler) CancelOrder(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID заказа"})
		return
	}

	// Получаем заказ
	order, err := h.OrderClient.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}

	// Проверяем, что заказ принадлежит пользователю
	if order.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	// Отменяем заказ
	updatedOrder, err := h.OrderClient.CancelOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя отменить заказ в текущем статусе"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": updatedOrder})
}

// getUserIDFromContext получает ID пользователя из контекста
func getUserIDFromContext(c *gin.Context) (uint64, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user_id not found")
	}

	id, ok := userID.(uint64)
	if !ok {
		return 0, errors.New("invalid user_id type")
	}

	return id, nil
}
