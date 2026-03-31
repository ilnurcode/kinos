// Package errs предоставляет типизированные ошибки для inventory-service.
// Включает ошибки для запасов и складов.
package errs

import "errors"

// Ошибки запасов (inventory)
var (
	ErrInventoryNotFound = errors.New("запасы не найдены")
	ErrInventoryExists   = errors.New("запасы для этого товара уже существуют")
	ErrInvalidQuantity   = errors.New("недопустимое количество")
	ErrInvalidProductID  = errors.New("недопустимый ID товара")
	ErrInsufficientStock = errors.New("недостаточно товара на складе")
	ErrInvalidLocation   = errors.New("недопустимое местоположение склада")
)

// Ошибки складов (warehouse)
var (
	ErrWarehouseNotFound = errors.New("склад не найден")
	ErrWarehouseExists   = errors.New("склад с таким названием уже существует")
	ErrWarehouseRequired = errors.New("название, город и улица обязательны")
)

// Ошибки валидации
var (
	ErrNameRequired   = errors.New("название обязательно")
	ErrCityRequired   = errors.New("город обязателен")
	ErrStreetRequired = errors.New("улица обязательна")
)
