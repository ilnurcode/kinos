// Package validator предоставляет валидацию входных данных для регистрации и входа.
// Использует библиотеку go-playground/validator для проверки полей.
package validator

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorInterface interface {
	ValidateRegister(input RegisterInput) error
	ValidateLogin(input LoginInput) error
}

var validate = validator.New()

type RegisterInput struct {
	Username string `validate:"required,min=3,max=30"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=40"`
	Phone    string `validate:"required,e164"`
}

type LoginInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=40"`
}

type Validator struct{}

func (v *Validator) ValidateRegister(input RegisterInput) error {
	return validate.Struct(input)
}

func (v *Validator) ValidateLogin(input LoginInput) error {
	return validate.Struct(input)
}
