// Package models предоставляет модели данных для catalog-service.
// Включает модели Category, Product, Manufacturer и фильтры.
package models

type Product struct {
	Id              uint64
	Name            string
	ManufacturersId uint64
	CategoryId      uint64
	Price           float64
}

type ProductFilter struct {
	NameContains    *string
	ManufacturersId *uint64
	CategoryId      *uint64
	PriceMin        *float64
	PriceMax        *float64
}
