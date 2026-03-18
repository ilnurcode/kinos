// Package main предоставляет gRPC-сервис для управления запасами товаров.
// Обрабатывает учет товаров, резервирование и управление складскими запасами.
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

	"kinos/inventory-service/config"
	"kinos/inventory-service/internal/grpcserver"
	"kinos/inventory-service/internal/repository"
	"kinos/inventory-service/internal/service"
	"kinos/inventory-service/internal/validator"

	pb "kinos/proto/inventory"

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
	if cfg.DBURL == "" {
		log.Fatal("DBURL env variable must be set")
	}

	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(cfg.DBURL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	inventoryRepo := repository.NewInventoryRepository(pool)
	warehouseRepo := repository.NewWarehouseRepository(pool)
	val := &validator.Validator{}
	txManager := repository.NewTxManager(pool)

	inventorySvc := service.NewInventoryService(inventoryRepo, val, txManager)
	warehouseSvc := service.NewWarehouseService(warehouseRepo, txManager)

	list, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcserver.LoggingInterceptor))
	inventoryServer := grpcserver.NewInventoryServer(inventorySvc, warehouseSvc)
	pb.RegisterInventoryServiceServer(grpcServer, inventoryServer)
	reflection.Register(grpcServer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("gRPC server started on :8083")
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down gRPC server...")
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
