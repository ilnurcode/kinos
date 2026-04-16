// Package catalog предоставляет HTTP-обработчики для управления каталогом товаров через API.
// Включает handlers для работы с категориями, производителями и товарами.
package catalog

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "kinos/proto/catalog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	client *CatalogClient
}

func NewHandler(c *CatalogClient) *Handler {
	return &Handler{
		client: c,
	}
}

// CreateCategory godoc
// @Summary Создать категорию
// @Description Создать новую категорию товаров (только для администраторов)
// @Tags catalog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name body string true "Название категории"
// @Success 201 {object} object "Созданная категория"
// @Failure 400 {object} object{error=string} "Некорректные данные"
// @Router /api/admin/catalog/categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.CreateCategory(ctx, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// UpdateCategory godoc
// @Summary Обновить категорию
// @Description Обновить существующую категорию товаров (только для администраторов)
// @Tags catalog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID категории"
// @Param name body string true "Название категории"
// @Success 200 {object} object "Обновленная категория"
// @Failure 400 {object} object{error=string} "Некорректные данные"
// @Router /api/admin/catalog/categories/:id [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.UpdateCategory(ctx, id, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	if _, err := h.client.DeleteCategory(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetCategory(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.GetCategory(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetListCategory(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset <= 0 {
		offset = 0
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.GetListCategory(ctx, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateManufacturer(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.CreateManufacturer(ctx, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateManufacturer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.UpdateManufacturer(ctx, id, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteManufacturer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	if _, err := h.client.DeleteManufacturer(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetManufacturers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.GetManufacturer(ctx, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetManufacturersList(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset <= 0 {
		offset = 0
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.GetListManufacturer(ctx, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var req struct {
		Name           string  `json:"name" binding:"required,min=1,max=100"`
		ManufacturerID uint64  `json:"manufacturer_id" binding:"required,gt=0"`
		CategoryID     uint64  `json:"category_id" binding:"required,gt=0"`
		Price          float64 `json:"price" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.CreateProduct(ctx, req.Name, req.ManufacturerID, req.CategoryID, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		Name           string  `json:"name" binding:"required,min=1,max=100"`
		ManufacturerID uint64  `json:"manufacturer_id" binding:"required,gt=0"`
		CategoryID     uint64  `json:"category_id" binding:"required,gt=0"`
		Price          float64 `json:"price" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.UpdateProduct(ctx, id, req.Name, req.ManufacturerID, req.CategoryID, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	if _, err := h.client.DeleteProduct(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetProduct(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.GetProduct(ctx, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetProductList(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset <= 0 {
		offset = 0
	}
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	manufacturerID, _ := strconv.ParseUint(c.Query("manufacturer_id"), 10, 64)
	priceMin, _ := strconv.ParseFloat(c.Query("price_min"), 64)
	priceMax, _ := strconv.ParseFloat(c.Query("price_max"), 64)
	nameContains := c.Query("name_contains")
	req := &pb.GetListProductRequest{
		Limit:          int32(limit),
		Offset:         int32(offset),
		CategoryId:     categoryID,
		ManufacturerId: manufacturerID,
		PriceMin:       priceMin,
		PriceMax:       priceMax,
		NameContains:   nameContains,
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.client.GetListProduct(ctx, req.Limit, req.Offset, req.CategoryId, req.ManufacturerId, req.PriceMax, req.PriceMin, req.NameContains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func requestContext(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 10*time.Second)
}
