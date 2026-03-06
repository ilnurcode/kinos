// Package grpcmiddleware предоставляет gRPC middleware для аутентификации и авторизации.
// Включает перехватчик для проверки JWT-токенов и добавления user_id/role в контекст.
package grpcmiddleware

import (
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"context"

	"kinos/user-service/internal/service"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"
const RoleKey ctxKey = "role"

func AuthUnaryInterceptor(tokenSvc *service.TokenService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/user.UserService/Register" || info.FullMethod == "/user.UserService/Login" || info.FullMethod == "/user.UserService/Refresh" || info.FullMethod == "/user.UserService/Revoke" {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization required")
		}
		parts := strings.SplitN(authHeaders[0], " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid auth header")
		}
		tokenStr := parts[1]
		claims, err := tokenSvc.ParseAccessTokenClaims(ctx, tokenStr)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		return handler(ctx, req)
	}
}
