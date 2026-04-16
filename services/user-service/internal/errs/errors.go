// Package errs предоставляет типизированные ошибки для пользовательского сервиса.
// Включает ошибки валидации, дублирования и бизнес-логики.
package errs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Ошибки репозитория
var (
	ErrNotFound           = errors.New("пользователь не найден")
	ErrUserExists         = errors.New("пользователь с таким email уже существует")
	ErrEmailExists        = errors.New("email уже используется")
	ErrInvalidCredentials = errors.New("неверный email или пароль")
)

// Ошибки валидации
var (
	ErrUsernameRequired = errors.New("имя пользователя обязательно")
	ErrUsernameTooShort = errors.New("имя пользователя должно быть не менее 3 символов")
	ErrUsernameTooLong  = errors.New("имя пользователя должно быть не более 30 символов")
	ErrEmailRequired    = errors.New("email обязателен")
	ErrEmailInvalid     = errors.New("неверный формат email")
	ErrPasswordRequired = errors.New("пароль обязателен")
	ErrPasswordTooShort = errors.New("пароль должен быть не менее 8 символов")
	ErrPasswordTooLong  = errors.New("пароль должен быть не более 40 символов")
	ErrPhoneRequired    = errors.New("телефон обязателен")
	ErrPhoneInvalid     = errors.New("неверный формат телефона. Пример: +79991234567")
)

// Ошибки авторизации
var (
	ErrUnauthorized        = errors.New("требуется аутентификация")
	ErrInvalidToken        = errors.New("неверный формат токена")
	ErrTokenExpired        = errors.New("токен истек")
	ErrPermissionDenied    = errors.New("доступ запрещен")
	ErrCannotChangeOwnRole = errors.New("нельзя изменить собственную роль")
)

// Ошибки бизнес-логики
var (
	ErrInternalError         = errors.New("внутренняя ошибка сервера")
	ErrFailedToHashPassword  = errors.New("ошибка при хешировании пароля")
	ErrFailedToCreateUser    = errors.New("ошибка при создании пользователя")
	ErrFailedToCreateTokens  = errors.New("ошибка при создании токенов")
	ErrFailedToUpdateRole    = errors.New("ошибка при обновлении роли")
	ErrFailedToUpdateProfile = errors.New("ошибка при обновлении профиля")
)

// ValidationError описывает ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("ошибка валидации: %s - %s", e.Field, e.Message)
}

// WrapValidationError оборачивает ошибки validator в понятные сообщения
func WrapValidationError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var messages []string
		for _, ve := range validationErrors {
			field := formatFieldName(ve.Field())
			messages = append(messages, formatValidationRule(ve.Tag(), field))
		}
		return errors.New(strings.Join(messages, "; "))
	}

	return err
}

// formatFieldName форматирует имя поля на русский язык
func formatFieldName(field string) string {
	fieldNames := map[string]string{
		"Username": "Имя пользователя",
		"Email":    "Email",
		"Password": "Пароль",
		"Phone":    "Телефон",
	}
	if name, ok := fieldNames[field]; ok {
		return name
	}
	return field
}

// formatValidationRule возвращает понятное сообщение для правила валидации
func formatValidationRule(tag, field string) string {
	messages := map[string]string{
		"required": fmt.Sprintf("%s обязателен", field),
		"email":    fmt.Sprintf("%s должен быть корректным email адресом", field),
		"min":      fmt.Sprintf("%s слишком короткий", field),
		"max":      fmt.Sprintf("%s слишком длинный", field),
		"gt":       fmt.Sprintf("%s должен быть больше", field),
	}
	if msg, ok := messages[tag]; ok {
		return msg
	}
	return fmt.Sprintf("%s не прошел валидацию (%s)", field, tag)
}
