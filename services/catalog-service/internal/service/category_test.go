package service

import (
	"context"
	"testing"

	"kinos/catalog-service/internal/models"
	"kinos/catalog-service/internal/repository"
	"kinos/catalog-service/internal/validator"
)

// Mock-репозиторий для тестов
type mockCategoryRepo struct {
	categories map[uint64]*models.Category
	nextID     uint64
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{
		categories: make(map[uint64]*models.Category),
		nextID:     1,
	}
}

func (m *mockCategoryRepo) CreateCategory(ctx context.Context, category *models.Category) (uint64, error) {
	category.Id = m.nextID
	m.categories[m.nextID] = category
	m.nextID++
	return category.Id, nil
}

func (m *mockCategoryRepo) UpdateCategory(ctx context.Context, category *models.Category) error {
	if _, ok := m.categories[category.Id]; !ok {
		return repository.ErrNotFound
	}
	m.categories[category.Id] = category
	return nil
}

func (m *mockCategoryRepo) DeleteCategory(ctx context.Context, id uint64) error {
	if _, ok := m.categories[id]; !ok {
		return repository.ErrNotFound
	}
	delete(m.categories, id)
	return nil
}

func (m *mockCategoryRepo) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	for _, c := range m.categories {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockCategoryRepo) GetCategoryByID(ctx context.Context, id uint64) (*models.Category, error) {
	if c, ok := m.categories[id]; ok {
		return c, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockCategoryRepo) GetListCategory(ctx context.Context, limit, offset int32) ([]*models.Category, int32, error) {
	var result []*models.Category
	for _, c := range m.categories {
		result = append(result, c)
	}
	return result, int32(len(result)), nil
}

// Mock-менеджер транзакций
type mockTxManager struct{}

func (m *mockTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCategoryService_CreateCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	val := &validator.Validator{}
	txManager := &mockTxManager{}

	svc := NewCategoryService(repo, val, txManager)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid category",
			input:   "Электроника",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := svc.CreateCategory(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && category == nil {
				t.Error("CreateCategory() returned nil category")
			}
			if !tt.wantErr && category.Name != tt.input {
				t.Errorf("CreateCategory() name = %v, want %v", category.Name, tt.input)
			}
		})
	}
}

func TestCategoryService_GetListCategory(t *testing.T) {
	repo := newMockCategoryRepo()
	val := &validator.Validator{}
	txManager := &mockTxManager{}

	svc := NewCategoryService(repo, val, txManager)

	// Создадим несколько категорий
	_, _ = svc.CreateCategory(context.Background(), "Категория 1")
	_, _ = svc.CreateCategory(context.Background(), "Категория 2")

	categories, total, err := svc.GetListCategory(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("GetListCategory() error = %v", err)
	}
	if total != 2 {
		t.Errorf("GetListCategory() total = %v, want 2", total)
	}
	if len(categories) != 2 {
		t.Errorf("GetListCategory() len = %v, want 2", len(categories))
	}
}
