package validator

import (
	"testing"
)

func TestValidator_ValidateRegister(t *testing.T) {
	v := &Validator{}

	tests := []struct {
		name    string
		input   RegisterInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: RegisterInput{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Phone:    "+79991234567",
			},
			wantErr: false,
		},
		{
			name: "short username",
			input: RegisterInput{
				Username: "te",
				Email:    "test@example.com",
				Password: "password123",
				Phone:    "+79991234567",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: RegisterInput{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
				Phone:    "+79991234567",
			},
			wantErr: true,
		},
		{
			name: "short password",
			input: RegisterInput{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "short",
				Phone:    "+79991234567",
			},
			wantErr: true,
		},
		{
			name: "invalid phone",
			input: RegisterInput{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Phone:    "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateRegister(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateLogin(t *testing.T) {
	v := &Validator{}

	tests := []struct {
		name    string
		input   LoginInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			input: LoginInput{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "short password",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateLogin(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
