// Package validator предоставляет валидацию входных данных для регистрации и входа.
// Использует библиотеку go-playground/validator для проверки полей.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidatorInterface interface {
	ValidateRegister(input RegisterInput) error
	ValidateLogin(input LoginInput) error
}

var validate = validator.New()

// E.164 формат: +[код страны][номер] (например, +79991234567)
// Допускается от 8 до 15 цифр после +
var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

type RegisterInput struct {
	Username string `validate:"required,min=3,max=30"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=40"`
	Phone    string
}

type LoginInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=40"`
}

type Validator struct{}

func (v *Validator) ValidateRegister(input RegisterInput) error {
	// Проверяем обязательные поля
	if err := validate.Struct(input); err != nil {
		return err
	}

	// Кастомная валидация телефона
	if err := v.validatePhone(input.Phone); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("телефон обязателен")
	}

	// Удаляем пробелы, дефисы и скобки
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Проверяем формат E.164
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("неверный формат телефона. Пример: +79991234567")
	}

	return nil
}

func (v *Validator) ValidateLogin(input LoginInput) error {
	return validate.Struct(input)
}
