// Package repository предоставляет репозиторий для работы с пользователями в базе данных.
// Включает CRUD-операции и поиск пользователей.
package repository

import (
	"context"
	"errors"
	"fmt"

	"kinos/user-service/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserInterface interface {
	CreateUser(ctx context.Context, username, email, password, phone string) (uint64, error)
	UpdateProfile(ctx context.Context,
		ID uint64,
		Username,
		Email,
		Phone string) error
	DeleteUser(ctx context.Context, ID uint64) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, ID uint64) (*models.User, error)
	UpdateRole(ctx context.Context, userID uint64, role string) error
	GetAllUsers(ctx context.Context, limit, offset int32) ([]*models.User, int32, error)
}
type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT user_ID, username,email,phone, hashed_password, role FROM users WHERE email = $1", email).Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role)
	if err != nil {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID uint64, role string) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "UPDATE users SET role = $1 WHERE user_ID = $2", role, userID)
	return err
}

func (r *UserRepository) FindUserByID(ctx context.Context, ID uint64) (*models.User, error) {
	var user models.User
	querier := GetQuerier(ctx, r.DB)
	err := querier.QueryRow(ctx, "SELECT user_ID, username,email,phone, hashed_password, role FROM users WHERE user_ID = $1", ID).Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role)
	if err != nil {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, username, email, password, phone string) (uint64, error) {
	querier := GetQuerier(ctx, r.DB)
	var userID uint64
	err := querier.QueryRow(ctx,
		`INSERT INTO users (username, email, phone, hashed_password)
         VALUES ($1, $2, $3, $4)
         RETURNING user_ID`,
		username, email, phone, password,
	).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("user with email %s already exists", email)
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, ID uint64) error {
	querier := GetQuerier(ctx, r.DB)
	_, err := querier.Exec(ctx, "DELETE FROM users WHERE user_ID = $1", ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context,
	ID uint64,
	Username string,
	Email string,
	Phone string) error {
	querier := GetQuerier(ctx, r.DB)
	user, err := r.FindUserByID(ctx, ID)
	if err != nil {
		return fmt.Errorf("failed to find user by ID: %w", err)
	}
	_, err = r.FindUserByEmail(ctx, Email)
	if err == nil && user.Email != Email {
		return fmt.Errorf("user with email %s already exists", Email)
	}
	_, err = querier.Exec(ctx, "Update users set username=$2, email=$3, phone=$4 where user_ID=$1", ID, Username, Email, Phone)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context, limit, offset int32) ([]*models.User, int32, error) {
	var users []*models.User
	querier := GetQuerier(ctx, r.DB)
	rows, err := querier.Query(ctx, "SELECT user_ID, username, email, phone, role FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all users: %w", err)
	}
	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.Role)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan users: %w", err)
		}
		users = append(users, &user)
	}
	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("failed to get all users: %w", rows.Err())
	}
	defer rows.Close()
	var total int32
	err = querier.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}
	return users, total, nil
}
