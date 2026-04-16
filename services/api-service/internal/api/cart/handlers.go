// Package cart предоставляет HTTP-обработчики для управления корзиной.
package cart

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CartClient CartClient
}

func NewHandler(cc CartClient) *Handler {
	return &Handler{CartClient: cc}
}

func (h *Handler) GetCart(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	cart, err := h.CartClient.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения корзины"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cart": cart, "total": cart.Total})
}

func (h *Handler) AddItem(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	var body struct {
		ProductID uint64 `json:"product_id"`
		Quantity  uint32 `json:"quantity"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	if body.ProductID == 0 || body.Quantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id и quantity обязательны"})
		return
	}
	cart, err := h.CartClient.AddItem(c.Request.Context(), userID, body.ProductID, body.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в корзину"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cart": cart, "total": cart.Total})
}

func (h *Handler) RemoveItem(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	productID, _ := strconv.ParseUint(c.Param("product_id"), 10, 64)
	cart, err := h.CartClient.RemoveItem(c.Request.Context(), userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления из корзины"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cart": cart, "total": cart.Total})
}

func (h *Handler) UpdateItem(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	productID, _ := strconv.ParseUint(c.Param("product_id"), 10, 64)
	var body struct {
		Quantity uint32 `json:"quantity"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	cart, err := h.CartClient.UpdateItem(c.Request.Context(), userID, productID, body.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления количества"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cart": cart, "total": cart.Total})
}

func (h *Handler) ClearCart(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	err = h.CartClient.ClearCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка очистки корзины"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Корзина очищена"})
}

func (h *Handler) GetItemsCount(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	count, err := h.CartClient.GetItemsCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения количества"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

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
