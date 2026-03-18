// Package model предоставляет модели данных для inventory-service.
package model

import "time"

type Inventory struct {
	Id                uint64    `json:"id"`
	ProductId         uint64    `json:"product_id"`
	Quantity          int32     `json:"quantity"`
	ReservedQuantity  int32     `json:"reserved_quantity"`
	AvailableQuantity int32     `json:"available_quantity"`
	WarehouseLocation string    `json:"warehouse_location"`
	UpdatedAt         time.Time `json:"updated_at"`
}
