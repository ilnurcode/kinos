// Package repository предоставляет репозиторий для работы с refresh-токенами в базе данных.
// Включает сохранение, поиск и удаление refresh-токенов.
package repository

import (
	"context"
	"fmt"
	"time"

	"kinos/user-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshInterface interface {
	Save(ctx context.Context, hash string, userID uint64, exp time.Time) error
	Find(ctx context.Context, hash string) (*models.RefreshToken, error)
	Delete(ctx context.Context, hash string) error
}

func NewRefreshRepository(db *pgxpool.Pool) *RefreshRepository {
	return &RefreshRepository{
		DB: db,
	}
}

type RefreshRepository struct {
	DB *pgxpool.Pool
}

func (r *RefreshRepository) Save(ctx context.Context, hash string, userID uint64, exp time.Time) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, `Insert into refresh_tokens (token_hash, user_id, expires_at, created_at) values ($1, $2, $3, $4)`, hash, userID, exp, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

func (r *RefreshRepository) Find(ctx context.Context, hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, `Select token_hash, user_id, expires_at, created_at from refresh_tokens where token_hash = $1`, hash).Scan(&token.TokenHash, &token.UserID, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}
	return &token, nil
}

func (r *RefreshRepository) Delete(ctx context.Context, hash string) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "DELETE from refresh_tokens where token_hash = $1", hash)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}
