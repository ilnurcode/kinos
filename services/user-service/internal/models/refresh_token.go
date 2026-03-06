// Package models предоставляет модели данных для user-service.
// Включает модели User и RefreshToken.
package models

import "time"

type RefreshToken struct {
	TokenHash string
	UserID    uint64
	ExpiresAt time.Time
	CreatedAt time.Time
}
