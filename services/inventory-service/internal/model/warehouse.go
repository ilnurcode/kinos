// Package model предоставляет модели данных для inventory-service.
package model

import "time"

type Warehouse struct {
	Id        uint64    `json:"id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Street    string    `json:"street"`
	Building  string    `json:"building"`
	Building2 string    `json:"building2"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
