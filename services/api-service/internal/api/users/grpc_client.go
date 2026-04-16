// Package users предоставляет gRPC-клиент для взаимодействия с user-service.
// Используется для вызова методов аутентификации, управления пользователями и валидации токенов.
package users

import (
	"context"
	pb "kinos/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type UserClientInterface interface {
	Register(ctx context.Context, username, email, password, phone string) (*pb.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*pb.AuthResponse, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) (*pb.RevokeResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*pb.AuthResponse, error)
	GetProfile(ctx context.Context, token string) (*pb.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, token, username, email, phone string) (*pb.UpdateProfileResponse, error)
	UpdateRole(ctx context.Context, token, role string, userID uint64) (*pb.UpdateRoleResponse, error)
	DeleteUser(ctx context.Context, token string, userID uint64) (*pb.DeleteUserResponse, error)
	ValidateAccess(ctx context.Context, token string) (*pb.ValidateAccessResponse, error)
	GetUsers(ctx context.Context, token string, limit, offset int32) (*pb.GetUsersResponse, error)
}
type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewUserServiceClient(conn)
	return &UserClient{conn: conn, client: c}, nil
}

func (uc *UserClient) Close() error {
	if uc == nil || uc.conn == nil {
		return nil
	}
	return uc.conn.Close()
}

func (uc *UserClient) Register(ctx context.Context, username, email, password, phone string) (*pb.AuthResponse, error) {
	req := &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
		Phone:    phone,
	}
	return uc.client.Register(ctx, req)
}

func (uc *UserClient) Login(ctx context.Context, email, password string) (*pb.AuthResponse, error) {
	req := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}
	return uc.client.Login(ctx, req)
}

func (uc *UserClient) Refresh(ctx context.Context, refreshToken string) (*pb.AuthResponse, error) {
	req := &pb.RefreshRequest{
		RefreshToken: refreshToken,
	}
	return uc.client.Refresh(ctx, req)
}

func (uc *UserClient) RevokeRefreshToken(ctx context.Context, refreshToken string) (*pb.RevokeResponse, error) {
	req := &pb.RevokeRequest{
		RefreshToken: refreshToken,
	}
	return uc.client.Revoke(ctx, req)
}

func (uc *UserClient) GetProfile(ctx context.Context, token string) (*pb.UserProfileResponse, error) {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pb.GetProfileRequest{}
	return uc.client.GetProfile(ctx, req)
}

func (uc *UserClient) UpdateProfile(ctx context.Context, token, username, email, phone string) (*pb.UpdateProfileResponse, error) {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pb.UpdateProfileRequest{
		Username: username,
		Email:    email,
		Phone:    phone,
	}
	resp, err := uc.client.UpdateProfile(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *UserClient) UpdateRole(ctx context.Context, token, role string, userID uint64) (*pb.UpdateRoleResponse, error) {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pb.UpdateRoleRequest{
		UserId: userID,
		Role:   role,
	}
	resp, err := uc.client.UpdateRole(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *UserClient) DeleteUser(ctx context.Context, token string, userID uint64) (*pb.DeleteUserResponse, error) {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pb.DeleteUserRequest{UserId: userID}
	resp, err := uc.client.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *UserClient) ValidateAccess(ctx context.Context, token string) (*pb.ValidateAccessResponse, error) {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pb.ValidateAccessRequest{AccessToken: token}
	return uc.client.ValidateAccess(ctx, req)
}

func (uc *UserClient) GetUsers(ctx context.Context, token string, limit, offset int32) (*pb.GetUsersResponse, error) {
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	req := &pb.GetUsersRequest{
		Limit:  limit,
		Offset: offset,
	}
	return uc.client.GetUsers(ctx, req)
}
