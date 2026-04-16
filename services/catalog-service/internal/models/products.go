// Package models предоставляет модели данных для catalog-service.
// Включает модели Category, Product, Manufacturer и фильтры.
package models

type Product struct {
	ID              uint64
	Name            string
	ManufacturersID uint64
	CategoryID      uint64
	Price           float64
}

type ProductFilter struct {
	NameContains    *string
	ManufacturersID *uint64
	CategoryID      *uint64
	PriceMin        *float64
	PriceMax        *float64
}
