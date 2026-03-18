// Package inventory предоставляет HTTP-обработчики для управления запасами товаров.
package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	InventoryClient InventoryClientInterface
}

func NewHandler(ic InventoryClientInterface) *Handler {
	return &Handler{
		InventoryClient: ic,
	}
}

// GetInventory godoc
// @Summary Получить информацию о запасе товара
// @Description Получить информацию о запасе конкретного товара
// @Tags inventory
// @Accept json
// @Produce json
// @Param product_id query uint64 true "ID товара"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/inventory [get]
func (h *Handler) GetInventory(c *gin.Context) {
	productIDStr := c.Query("product_id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
		return
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	inv, err := h.InventoryClient.GetInventory(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
		return
	}

	c.JSON(http.StatusOK, inv)
}

// GetListInventory godoc
// @Summary Получить список запасов
// @Description Получить список всех запасов с фильтрацией и пагинацией
// @Tags inventory
// @Accept json
// @Produce json
// @Param limit query int32 false "Лимит записей" default(20)
// @Param offset query int32 false "Смещение" default(0)
// @Param product_id query uint64 false "ID товара"
// @Param warehouse_location query string false "Местоположение склада"
// @Param min_quantity query int32 false "Минимальное количество"
// @Success 200 {object} object
// @Router /api/inventory/list [get]
func (h *Handler) GetListInventory(c *gin.Context) {
	limit := int32(20)
	offset := int32(0)
	var productID uint64
	var location string
	var minQuantity int32

	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = int32(l)
		}
	}
	if oStr := c.Query("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o > 0 {
			offset = int32(o)
		}
	}
	if pStr := c.Query("product_id"); pStr != "" {
		if p, err := strconv.ParseUint(pStr, 10, 64); err == nil {
			productID = p
		}
	}
	if lStr := c.Query("warehouse_location"); lStr != "" {
		location = lStr
	}
	if mStr := c.Query("min_quantity"); mStr != "" {
		if m, err := strconv.Atoi(mStr); err == nil && m > 0 {
			minQuantity = int32(m)
		}
	}

	resp, err := h.InventoryClient.GetListInventory(c.Request.Context(), limit, offset, productID, location, minQuantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get inventory list"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateInventory godoc
// @Summary Создать запись о запасе
// @Description Создать новую запись о запасе товара
// @Tags inventory
// @Accept json
// @Produce json
// @Param product_id body uint64 true "ID товара"
// @Param quantity body int32 true "Количество"
// @Param warehouse_location body string true "Местоположение склада"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Router /api/inventory [post]
func (h *Handler) CreateInventory(c *gin.Context) {
	var body struct {
		ProductID         uint64 `json:"product_id"`
		Quantity          int32  `json:"quantity"`
		WarehouseLocation string `json:"warehouse_location"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	inv, err := h.InventoryClient.CreateInventory(c.Request.Context(), body.ProductID, body.Quantity, body.WarehouseLocation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create inventory"})
		return
	}

	c.JSON(http.StatusOK, inv)
}

// UpdateInventory godoc
// @Summary Обновить запись о запасе
// @Description Обновить существующую запись о запасе товара
// @Tags inventory
// @Accept json
// @Produce json
// @Param id path uint64 true "ID записи"
// @Param quantity body int32 true "Количество"
// @Param warehouse_location body string true "Местоположение склада"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/inventory/:id [put]
func (h *Handler) UpdateInventory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Quantity          int32  `json:"quantity"`
		WarehouseLocation string `json:"warehouse_location"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	inv, err := h.InventoryClient.UpdateInventory(c.Request.Context(), id, body.Quantity, body.WarehouseLocation)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
		return
	}

	c.JSON(http.StatusOK, inv)
}

// DeleteInventory godoc
// @Summary Удалить запись о запасе
// @Description Удалить существующую запись о запасе товара
// @Tags inventory
// @Accept json
// @Produce json
// @Param id path uint64 true "ID записи"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /api/inventory/:id [delete]
func (h *Handler) DeleteInventory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.InventoryClient.DeleteInventory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ReserveStock godoc
// @Summary Зарезервировать товар
// @Description Зарезервировать определенное количество товара
// @Tags inventory
// @Accept json
// @Produce json
// @Param product_id body uint64 true "ID товара"
// @Param quantity body int32 true "Количество для резервирования"
// @Param reservation_id body string true "ID резервирования"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 409 {object} object
// @Router /api/inventory/reserve [post]
func (h *Handler) ReserveStock(c *gin.Context) {
	var body struct {
		ProductID     uint64 `json:"product_id"`
		Quantity      int32  `json:"quantity"`
		ReservationID string `json:"reservation_id"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.InventoryClient.ReserveStock(c.Request.Context(), body.ProductID, body.Quantity, body.ReservationID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "insufficient stock"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ReleaseReservation godoc
// @Summary Снять резервирование
// @Description Снять резервирование с товара
// @Tags inventory
// @Accept json
// @Produce json
// @Param product_id body uint64 true "ID товара"
// @Param reservation_id body string true "ID резервирования"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Router /api/inventory/release [post]
func (h *Handler) ReleaseReservation(c *gin.Context) {
	var body struct {
		ProductID     uint64 `json:"product_id"`
		ReservationID string `json:"reservation_id"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.InventoryClient.ReleaseReservation(c.Request.Context(), body.ProductID, body.ReservationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release reservation"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetListWarehouse godoc
// @Summary Получить список складов
// @Description Получить список всех складов
// @Tags inventory
// @Accept json
// @Produce json
// @Param limit query int32 false "Лимит записей" default(20)
// @Param offset query int32 false "Смещение" default(0)
// @Success 200 {object} object
// @Router /api/inventory/warehouses/list [get]
func (h *Handler) GetListWarehouse(c *gin.Context) {
	limit := int32(20)
	offset := int32(0)

	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = int32(l)
		}
	}
	if oStr := c.Query("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o > 0 {
			offset = int32(o)
		}
	}

	resp, err := h.InventoryClient.GetListWarehouse(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get warehouses"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateWarehouse godoc
// @Summary Создать склад
// @Description Создать новый склад
// @Tags inventory
// @Accept json
// @Produce json
// @Param name body string true "Название"
// @Param city body string true "Город"
// @Param street body string true "Улица"
// @Param building body string false "Дом"
// @Success 200 {object} object
// @Router /api/inventory/warehouses [post]
func (h *Handler) CreateWarehouse(c *gin.Context) {
	var body struct {
		Name      string `json:"name"`
		City      string `json:"city"`
		Street    string `json:"street"`
		Building  string `json:"building"`
		Building2 string `json:"building2"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if body.Name == "" || body.City == "" || body.Street == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, city и street обязательны"})
		return
	}

	warehouse, err := h.InventoryClient.CreateWarehouse(c.Request.Context(), body.Name, body.City, body.Street, body.Building, body.Building2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create warehouse"})
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// DeleteWarehouse godoc
// @Summary Удалить склад
// @Description Удалить склад по ID
// @Tags inventory
// @Accept json
// @Produce json
// @Param id path uint64 true "ID склада"
// @Success 200 {object} object
// @Router /api/inventory/warehouses/:id [delete]
func (h *Handler) DeleteWarehouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.InventoryClient.DeleteWarehouse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete warehouse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
