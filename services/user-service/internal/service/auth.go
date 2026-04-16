// Package service предоставляет бизнес-логику для аутентификации и управления пользователями.
// Включает регистрацию, вход, обновление токенов и управление ролями.
package service

import (
	"context"
	"errors"
	"time"

	usrErrs "kinos/user-service/internal/errs"
	"kinos/user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// bcryptCost — стоимость хеширования bcrypt. 12 — рекомендуется для production.
const bcryptCost = 12

type Auth interface {
	Register(ctx context.Context, username, email, password, phone string) (string, string, time.Time, error)
	Login(ctx context.Context, email, password string) (string, string, time.Time, error)
	Refresh(ctx context.Context, oldRefreshToken string) (string, string, time.Time, error)
	RevokeRefresh(ctx context.Context, refreshToken string) error
	UpdateRole(ctx context.Context, userID uint64, newRole string) error
	DeleteUser(ctx context.Context, userID uint64) error
}
type AuthService struct {
	UserRepo     repository.UserInterface
	TokenService *TokenService
	txManager    repository.TransactionManager
}

func NewAuthService(UserRepo repository.UserInterface, TokenService *TokenService, txManager repository.TransactionManager) *AuthService {
	return &AuthService{
		UserRepo:     UserRepo,
		TokenService: TokenService,
		txManager:    txManager,
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password, phone string) (string, string, time.Time, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", "", time.Time{}, status.Errorf(codes.Internal, "ошибка при хешировании пароля: %v", err)
	}
	var access, refresh string
	var exp time.Time
	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		userID, err := s.UserRepo.CreateUser(txCtx, username, email, string(hashedPassword), phone)
		if err != nil {
			if errors.Is(err, usrErrs.ErrUserExists) {
				return status.Error(codes.AlreadyExists, "пользователь с таким email уже существует")
			}
			return status.Errorf(codes.Internal, "ошибка при создании пользователя: %v", err)
		}
		access, refresh, exp, err = s.TokenService.CreateTokens(txCtx, userID, "user")
		if err != nil {
			return status.Errorf(codes.Internal, "ошибка при создании токенов: %v", err)
		}
		return nil
	})
	if err != nil {
		// Если это уже gRPC статус, возвращаем как есть
		if _, ok := status.FromError(err); ok {
			return "", "", time.Time{}, err
		}
		return "", "", time.Time{}, status.Errorf(codes.Internal, "ошибка при регистрации: %v", err)
	}
	return access, refresh, exp, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, time.Time, error) {
	user, err := s.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", time.Time{}, status.Error(codes.Unauthenticated, "неверный email или пароль")
		}
		return "", "", time.Time{}, status.Errorf(codes.Internal, "ошибка поиска пользователя: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", time.Time{}, status.Error(codes.Unauthenticated, "неверный email или пароль")
	}
	access, refresh, exp, err := s.TokenService.CreateTokens(ctx, user.Id, user.Role)
	if err != nil {
		return "", "", time.Time{}, status.Errorf(codes.Internal, "ошибка при создании токенов: %v", err)
	}
	return access, refresh, exp, nil
}

func (s *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (string, string, time.Time, error) {
	return s.TokenService.RotateRefresh(ctx, oldRefreshToken)
}

func (s *AuthService) RevokeRefresh(ctx context.Context, refreshToken string) error {
	return s.TokenService.RevokeRefresh(ctx, refreshToken)
}

func (s *AuthService) UpdateRole(ctx context.Context, userID uint64, newRole string) error {
	err := s.UserRepo.UpdateRole(ctx, userID, newRole)
	if err != nil {
		return status.Errorf(codes.Internal, "ошибка при обновлении роли: %v", err)
	}
	return nil
}

func (s *AuthService) DeleteUser(ctx context.Context, userID uint64) error {
	err := s.UserRepo.DeleteUser(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return status.Error(codes.NotFound, "пользователь не найден")
		}
		return status.Errorf(codes.Internal, "ошибка при удалении пользователя: %v", err)
	}
	return nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uint64, username, email, phone string) error {
	err := s.UserRepo.UpdateProfile(ctx, userID, username, email, phone)
	if err != nil {
		if errors.Is(err, usrErrs.ErrEmailExists) {
			return status.Error(codes.AlreadyExists, "email уже используется")
		}
		if errors.Is(err, usrErrs.ErrNotFound) {
			return status.Error(codes.NotFound, "пользователь не найден")
		}
		return status.Errorf(codes.Internal, "ошибка при обновлении профиля: %v", err)
	}
	return nil
}
