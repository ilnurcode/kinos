package validator

import (
	"testing"
)

func TestValidator_ValidateCategory(t *testing.T) {
	v := &Validator{}

	tests := []struct {
		name    string
		input   CategoryInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: CategoryInput{
				Name: "Электроника",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: CategoryInput{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "name too long",
			input: CategoryInput{
				Name: "Очень длинное название категории которое превышает максимальную длину в сто символов и должно вызвать ошибку валидации",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCategory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateManufacturer(t *testing.T) {
	v := &Validator{}

	tests := []struct {
		name    string
		input   ManufacturersInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: ManufacturersInput{
				Name: "Samsung",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: ManufacturersInput{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "name too long",
			input: ManufacturersInput{
				Name: "Очень длинное название производителя которое превышает максимальную длину в сто символов и должно вызвать ошибку валидации потому что название слишком длинное",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateManufactures(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateManufactures() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateProduct(t *testing.T) {
	v := &Validator{}

	tests := []struct {
		name    string
		input   ProductInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: ProductInput{
				Name:            "Смартфон",
				ManufacturersID: 1,
				CategoryID:      1,
				Price:           999.99,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: ProductInput{
				Name:            "",
				ManufacturersID: 1,
				CategoryID:      1,
				Price:           999.99,
			},
			wantErr: true,
		},
		{
			name: "zero manufacturer ID",
			input: ProductInput{
				Name:            "Смартфон",
				ManufacturersID: 0,
				CategoryID:      1,
				Price:           999.99,
			},
			wantErr: true,
		},
		{
			name: "zero category ID",
			input: ProductInput{
				Name:            "Смартфон",
				ManufacturersID: 1,
				CategoryID:      0,
				Price:           999.99,
			},
			wantErr: true,
		},
		{
			name: "zero price",
			input: ProductInput{
				Name:            "Смартфон",
				ManufacturersID: 1,
				CategoryID:      1,
				Price:           0,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			input: ProductInput{
				Name:            "Смартфон",
				ManufacturersID: 1,
				CategoryID:      1,
				Price:           -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateProduct(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
