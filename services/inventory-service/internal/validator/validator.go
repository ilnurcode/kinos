// Package validator предоставляет валидацию входных данных для inventory-service.
package validator

import "errors"

var (
	ErrInvalidQuantity  = errors.New("недопустимое количество")
	ErrInvalidProductID = errors.New("недопустимый ID товара")
	ErrInvalidLocation  = errors.New("недопустимое местоположение склада")
)

type CreateInventoryInput struct {
	ProductID uint64
	Quantity  int32
	Location  string
}

type UpdateInventoryInput struct {
	ID       uint64
	Quantity int32
	Location string
}

type Validator struct{}

func (v *Validator) ValidateCreateInventory(input CreateInventoryInput) error {
	if input.ProductID == 0 {
		return ErrInvalidProductID
	}
	if input.Quantity < 0 {
		return ErrInvalidQuantity
	}
	if input.Location == "" {
		return ErrInvalidLocation
	}
	return nil
}

func (v *Validator) ValidateUpdateInventory(input UpdateInventoryInput) error {
	if input.ID == 0 {
		return ErrInvalidProductID
	}
	if input.Quantity < 0 {
		return ErrInvalidQuantity
	}
	if input.Location == "" {
		return ErrInvalidLocation
	}
	return nil
}
