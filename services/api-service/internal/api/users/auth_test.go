package users

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "kinos/proto/user"

	"github.com/gin-gonic/gin"
)

// Mock UserClient для тестов
type mockUserClient struct{}

func (m *mockUserClient) Register(ctx context.Context, username, email, password, phone string) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
	}, nil
}

func (m *mockUserClient) Login(ctx context.Context, email, password string) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
	}, nil
}

func (m *mockUserClient) Refresh(ctx context.Context, refreshToken string) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}, nil
}

func (m *mockUserClient) GetProfile(ctx context.Context, token string) (*pb.UserProfileResponse, error) {
	return &pb.UserProfileResponse{
		UserId:   1,
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "+79991234567",
	}, nil
}

func (m *mockUserClient) ValidateAccess(ctx context.Context, token string) (*pb.ValidateAccessResponse, error) {
	return &pb.ValidateAccessResponse{
		UserId: 1,
		Role:   "user",
		Valid:  true,
	}, nil
}

func (m *mockUserClient) GetUsers(ctx context.Context, token string, limit, offset int32) (*pb.GetUsersResponse, error) {
	return &pb.GetUsersResponse{
		Users: []*pb.UserItem{
			{
				Id:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Role:     "user",
			},
		},
		Total: 1,
	}, nil
}

func (m *mockUserClient) RevokeRefreshToken(ctx context.Context, refreshToken string) (*pb.RevokeResponse, error) {
	return &pb.RevokeResponse{Success: true}, nil
}

func (m *mockUserClient) UpdateProfile(ctx context.Context, token, username, email, phone string) (*pb.UpdateProfileResponse, error) {
	return &pb.UpdateProfileResponse{Success: true}, nil
}

func (m *mockUserClient) UpdateRole(ctx context.Context, token, role string, userID uint64) (*pb.UpdateRoleResponse, error) {
	return &pb.UpdateRoleResponse{Success: true}, nil
}

func (m *mockUserClient) DeleteUser(ctx context.Context, token string, userID uint64) (*pb.DeleteUserResponse, error) {
	return &pb.DeleteUserResponse{Success: true}, nil
}

func TestHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       interface{}
		wantStatus int
	}{
		{
			name: "valid request",
			body: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
				"phone":    "+79991234567",
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&mockUserClient{})
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			if w.Code != tt.wantStatus {
				t.Errorf("Register() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
