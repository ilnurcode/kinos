// Package service предоставляет сервис для управления JWT-токенами и refresh-токенами.
// Включает генерацию, валидацию, обновление и отзыв токенов.
package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"kinos/user-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type TokenServiceInterface interface {
	GenerateAccessToken(userID uint64, role string) (string, error)
	GenerateRefreshToken() (string, string, error)
	CreateTokens(ctx context.Context, UserID uint64, role string) (string, string, time.Time, error)
	HashToken(refreshPlain string) string
	ParseAccessTokenClaims(tokenStr string) (*Claims, error)
	RevokeRefresh(ctx context.Context, plain string) error
	RotateRefresh(ctx context.Context, oldPlain string) (string, string, time.Time, error)
}

type TokenService struct {
	refreshRepo repository.RefreshInterface
	userRepo    repository.UserInterface
	txManager   repository.TransactionManager
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewTokenService(refreshRepo repository.RefreshInterface, userRepo repository.UserInterface, txManager repository.TransactionManager, jwtSecret string, accessTTL time.Duration, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		refreshRepo: refreshRepo,
		userRepo:    userRepo,
		txManager:   txManager,
		jwtSecret:   []byte(jwtSecret),
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

type Claims struct {
	UserID uint64 `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *TokenService) GenerateAccessToken(userID uint64, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-service",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) GenerateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	plain := base64.RawURLEncoding.EncodeToString(b)

	hash := s.HashToken(plain)
	return plain, hash, nil
}

func (s *TokenService) CreateTokens(ctx context.Context, UserID uint64, role string) (string, string, time.Time, error) {
	access, err := s.GenerateAccessToken(UserID, role)
	if err != nil {
		return "", "", time.Time{}, err
	}
	refreshPlain, refreshHash, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, err
	}
	exp := time.Now().Add(s.refreshTTL)
	err = s.refreshRepo.Save(ctx, refreshHash, UserID, exp)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return access, refreshPlain, exp, nil
}

func (s *TokenService) HashToken(refreshPlain string) string {
	h := sha256.Sum256([]byte(refreshPlain))
	hash := hex.EncodeToString(h[:])
	return hash
}

func (s *TokenService) RotateRefresh(ctx context.Context, oldPlain string) (string, string, time.Time, error) {
	var resultAccess string
	var resultRefresh string
	var resultExp time.Time
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		oldHash := s.HashToken(oldPlain)
		token, err := s.refreshRepo.Find(ctx, oldHash)
		if err != nil {
			return err
		}
		if time.Now().After(token.ExpiresAt) {
			_ = s.refreshRepo.Delete(ctx, oldHash)
			return fmt.Errorf("refresh token expired")
		}
		if err := s.refreshRepo.Delete(ctx, oldHash); err != nil {
			return err
		}

		user, err := s.userRepo.FindUserByID(ctx, token.UserID)
		if err != nil {
			return err
		}

		newPlain, newHash, err := s.GenerateRefreshToken()
		if err != nil {
			return err
		}
		newExp := time.Now().Add(s.refreshTTL)
		if err := s.refreshRepo.Save(ctx, newHash, token.UserID, newExp); err != nil {
			return err
		}
		access, err := s.GenerateAccessToken(token.UserID, user.Role)
		if err != nil {
			return err
		}
		resultAccess = access
		resultRefresh = newPlain
		resultExp = newExp
		return nil
	})
	if err != nil {
		return "", "", time.Time{}, err
	}
	return resultAccess, resultRefresh, resultExp, nil
}

func (s *TokenService) ParseAccessTokenClaims(ctx context.Context, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (s *TokenService) RevokeRefresh(ctx context.Context, plain string) error {
	hash := s.HashToken(plain)
	return s.refreshRepo.Delete(ctx, hash)
}
