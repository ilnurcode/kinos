// Package main предоставляет gRPC-сервис для управления пользователями и аутентификацией.
// Обрабатывает регистрацию, вход, обновление токенов и управление профилями пользователей.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "kinos/proto/user"
	"kinos/user-service/config"
	"kinos/user-service/internal/grpcmiddleware"
	"kinos/user-service/internal/grpcserver"
	"kinos/user-service/internal/repository"
	"kinos/user-service/internal/service"
	"kinos/user-service/internal/validator"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	_ = godotenv.Load()
	cfg := config.NewConfig()
	if cfg.DBURL == "" || cfg.SecretKey == "" {
		log.Fatal("DB_URL and SECRET_KEY env vars must be set")
	}
	if len(cfg.SecretKey) < 32 {
		log.Fatal("SECRET_KEY must be at least 32 bytes long")
	}

	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(cfg.DBURL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(pool)
	refreshRepo := repository.NewRefreshRepository(pool)
	txManager := repository.NewTxManager(pool)

	tokenSvc := service.NewTokenService(refreshRepo, userRepo, txManager, cfg.SecretKey, time.Minute*15, time.Hour*24*30)

	authSvc := service.NewAuthService(userRepo, tokenSvc, txManager)
	val := &validator.Validator{}
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.AuthUnaryInterceptor(tokenSvc)))
	userSrv := grpcserver.NewUserServer(authSvc, userRepo, val, tokenSvc)
	pb.RegisterUserServiceServer(grpcServer, userSrv)
	reflection.Register(grpcServer)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Println("grpc server started on :8081")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-quit
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server exited")
}

func runMigrations(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open DB for migrations: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("error while migrating database: %v", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error while running migrations: %v", err)
	}
	return nil
}
