// Package grpcserver предоставляет gRPC-сервер для user-service.
// Реализует методы UserService: Register, Login, Refresh, GetProfile, UpdateRole и другие.
package grpcserver

import (
	"context"
	"log"

	pb "kinos/proto/user"
	"kinos/user-service/internal/grpcmiddleware"
	"kinos/user-service/internal/repository"
	"kinos/user-service/internal/service"
	"kinos/user-service/internal/validator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	authSvc   *service.AuthService
	userSvc   repository.UserInterface
	tokenSvc  *service.TokenService
	validator validator.ValidatorInterface
}

func NewUserServer(authSvc *service.AuthService, userSvc repository.UserInterface, validator validator.ValidatorInterface, tokenSvc *service.TokenService) *UserServer {
	return &UserServer{
		authSvc:   authSvc,
		userSvc:   userSvc,
		tokenSvc:  tokenSvc,
		validator: validator,
	}
}

func (s *UserServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	limit := req.Limit
	offset := req.Offset
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	users, total, err := s.userSvc.GetAllUsers(ctx, limit, offset)
	if err != nil {
		log.Printf("failed get users: %v", err)
	}
	var result []*pb.UserItem
	for _, user := range users {
		result = append(result, &pb.UserItem{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
		})
	}
	return &pb.GetUsersResponse{
		Users: result,
		Total: total,
	}, nil
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	log.Printf("Register request: email=%s, username=%s", req.Email, req.Username)
	input := validator.RegisterInput{Username: req.Username, Email: req.Email, Password: req.Password, Phone: req.Phone}
	err := s.validator.ValidateRegister(input)
	if err != nil {
		log.Printf("Validation error: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}
	access, refresh, exp, err := s.authSvc.Register(ctx, req.Username, req.Email, req.Password, req.Phone)
	if err != nil {
		log.Printf("Register service error: %v", err)
		// Если это уже gRPC статус с понятным сообщением, возвращаем его
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "ошибка при регистрации: %v", err)
	}
	log.Printf("User registered successfully: email=%s, user_id=%v", req.Email, access)
	return &pb.AuthResponse{
		AccessToken:      access,
		RefreshToken:     refresh,
		RefreshExpiresAt: exp.Unix(),
	}, nil
}

func (s *UserServer) ValidateAccess(ctx context.Context, req *pb.ValidateAccessRequest) (*pb.ValidateAccessResponse, error) {
	claims, err := s.tokenSvc.ParseAccessTokenClaims(ctx, req.AccessToken)
	if err != nil {
		log.Printf("Validate access token error: %v", err)
		return &pb.ValidateAccessResponse{Valid: false}, nil
	}
	return &pb.ValidateAccessResponse{UserId: claims.UserID, Role: claims.Role, Valid: true}, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	log.Printf("Login request: email=%s", req.Email)
	input := validator.LoginInput{Email: req.Email, Password: req.Password}
	err := s.validator.ValidateLogin(input)
	if err != nil {
		log.Printf("Validation error: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}
	access, refresh, exp, err := s.authSvc.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("Login service error: %v", err)
		// Если это уже gRPC статус с понятным сообщением, возвращаем его
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "ошибка при входе: %v", err)
	}
	log.Printf("User logged in successfully: email=%s", req.Email)
	return &pb.AuthResponse{
		AccessToken:      access,
		RefreshToken:     refresh,
		RefreshExpiresAt: exp.Unix(),
	}, nil
}

func (s *UserServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error) {
	access, refresh, exp, err := s.authSvc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		log.Printf("Refresh service error: %v", err)
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "ошибка обновления токена: %v", err)
	}
	return &pb.AuthResponse{
		AccessToken:      access,
		RefreshToken:     refresh,
		RefreshExpiresAt: exp.Unix(),
	}, nil
}

func (s *UserServer) Revoke(ctx context.Context, req *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	if err := s.authSvc.RevokeRefresh(ctx, req.RefreshToken); err != nil {
		log.Printf("Revoke service error: %v", err)
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "ошибка отзыва токена: %v", err)
	}
	return &pb.RevokeResponse{Success: true}, nil
}

func (s *UserServer) GetProfile(ctx context.Context, _ *pb.GetProfileRequest) (*pb.UserProfileResponse, error) {
	userIDRaw := ctx.Value(grpcmiddleware.UserIDKey)
	if userIDRaw == nil {
		return nil, status.Error(codes.Unauthenticated, "требуется аутентификация")
	}
	userID := userIDRaw.(uint64)
	user, err := s.userSvc.FindUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "пользователь не найден")
	}
	return &pb.UserProfileResponse{
		UserId:   user.Id,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
	}, nil
}

func (s *UserServer) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	userIDRaw := ctx.Value(grpcmiddleware.UserIDKey)
	if userIDRaw == nil {
		return nil, status.Error(codes.Unauthenticated, "требуется аутентификация")
	}
	currentUserID := userIDRaw.(uint64)

	roleVal := ctx.Value(grpcmiddleware.RoleKey)
	if roleVal == nil {
		return nil, status.Error(codes.Unauthenticated, "требуется аутентификация")
	}
	if roleVal.(string) != "admin" {
		return nil, status.Error(codes.PermissionDenied, "доступ запрещен: требуется роль администратора")
	}

	// Запрет изменения роли самому себе
	if currentUserID == req.UserId {
		return nil, status.Error(codes.PermissionDenied, "нельзя изменить собственную роль")
	}

	_, err := s.userSvc.FindUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "пользователь не найден")
	}
	if req.Role != "admin" && req.Role != "user" {
		return nil, status.Error(codes.InvalidArgument, "недопустимая роль: разрешены только 'admin' или 'user'")
	}
	err = s.userSvc.UpdateRole(ctx, req.UserId, req.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка обновления роли: %v", err)
	}
	return &pb.UpdateRoleResponse{Success: true}, nil
}

func (s *UserServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	userIDRaw := ctx.Value(grpcmiddleware.UserIDKey)
	if userIDRaw == nil {
		return nil, status.Error(codes.Unauthenticated, "требуется аутентификация")
	}
	userID := userIDRaw.(uint64)
	_, err := s.userSvc.FindUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "пользователь не найден")
	}
	err = s.authSvc.UpdateProfile(ctx, userID, req.Username, req.Email, req.Phone)
	if err != nil {
		// Если это уже gRPC статус с понятным сообщением, возвращаем его
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "ошибка обновления профиля: %v", err)
	}
	return &pb.UpdateProfileResponse{Success: true}, nil
}
