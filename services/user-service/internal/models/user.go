// Package models предоставляет модели данных для user-service.
// Включает модели User и RefreshToken.
package models

type User struct {
	Id       uint64
	Username string
	Email    string
	Phone    string
	Password string
	Role     string
}
