// Package errs предоставляет типизированные ошибки для каталога сервисов.
// Включает ошибки для категорий, производителей и товаров.
package errs

import "errors"

// Ошибки категорий
var (
	ErrCategoryNotFound = errors.New("категория не найдена")
	ErrCategoryExists   = errors.New("категория с таким названием уже существует")
)

// Ошибки производителей
var (
	ErrManufacturerNotFound = errors.New("производитель не найден")
	ErrManufacturerExists   = errors.New("производитель с таким названием уже существует")
)

// Ошибки товаров
var (
	ErrProductNotFound = errors.New("товар не найден")
	ErrProductExists   = errors.New("товар с таким названием уже существует")
)

// Ошибки валидации
var (
	ErrNameRequired     = errors.New("название обязательно")
	ErrPriceInvalid     = errors.New("неверная цена")
	ErrCategoryRequired = errors.New("категория обязательна")
)
