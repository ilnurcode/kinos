// Package validator предоставляет валидацию входных данных для каталога.
// Использует библиотеку go-playground/validator для проверки категорий, производителей и товаров.
package validator

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorInterface interface {
	ValidateManufactures(input ManufacturersInput) error
	ValidateCategory(input CategoryInput) error
	ValidateProduct(input ProductInput) error
}

var validate = validator.New()

type ManufacturersInput struct {
	Name string `validate:"required,min=1,max=100"`
}
type CategoryInput struct {
	Name string `validate:"required,min=1,max=100"`
}
type ProductInput struct {
	Name            string  `validate:"required,min=1,max=100"`
	ManufacturersID uint64  `validate:"required,gt=0"`
	CategoryID      uint64  `validate:"required,gt=0"`
	Price           float64 `validate:"required,gt=0"`
}

type Validator struct{}

func (v *Validator) ValidateManufactures(input ManufacturersInput) error {
	return validate.Struct(input)
}
func (v *Validator) ValidateCategory(input CategoryInput) error {
	return validate.Struct(input)
}
func (v *Validator) ValidateProduct(input ProductInput) error {
	return validate.Struct(input)
}
