// Package models предоставляет модели данных для catalog-service.
// Включает модели Category, Product, Manufacturer и фильтры.
package models

type Manufacturer struct {
	ID   uint64
	Name string
}
