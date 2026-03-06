// Package service предоставляет бизнес-логику для аутентификации и управления пользователями.
// Включает регистрацию, вход, обновление токенов и управление ролями.
package service

import (
	"context"
	"time"

	"kinos/user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, username, email, password, phone string) (string, string, time.Time, error)
	Login(ctx context.Context, email, password string) (string, string, time.Time, error)
	Refresh(ctx context.Context, oldRefreshToken string) (string, string, time.Time, error)
	RevokeRefresh(ctx context.Context, refreshToken string) error
	UpdateRole(ctx context.Context, userID uint64, newRole string) error
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", time.Time{}, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	var access, refresh string
	var exp time.Time
	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		userID, err := s.UserRepo.CreateUser(txCtx, username, email, string(hashedPassword), phone)
		if err != nil {
			return status.Error(codes.Internal, "failed to create user")
		}
		access, refresh, exp, err = s.TokenService.CreateTokens(txCtx, userID, "user")
		if err != nil {
			return status.Error(codes.Internal, "failed to create tokens")
		}
		return nil
	})
	if err != nil {
		return "", "", time.Time{}, status.Error(codes.Internal, err.Error())
	}
	return access, refresh, exp, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, time.Time, error) {
	user, err := s.UserRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", time.Time{}, status.Errorf(codes.NotFound, "User not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", time.Time{}, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	access, refresh, exp, err := s.TokenService.CreateTokens(ctx, user.Id, user.Role)
	if err != nil {
		return "", "", time.Time{}, status.Error(codes.Internal, "failed to create tokens")
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
		return status.Error(codes.Internal, "failed to update role")
	}
	return nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uint64, username, email, phone string) error {
	return s.UserRepo.UpdateProfile(ctx, userID, username, email, phone)
}
