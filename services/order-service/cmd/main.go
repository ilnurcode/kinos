package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kinos/order-service/config"
	"kinos/order-service/internal/grpcserver"
	inventoryclient "kinos/order-service/internal/inventory"
	"kinos/order-service/internal/repository"
	"kinos/order-service/internal/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	godotenv.Load()

	cfg := config.NewConfig()

	// Подключаемся к БД
	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database")

	if err := runMigrations(cfg.DBURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Создаём репозиторий
	orderRepo := repository.NewOrderRepository(pool)
	txManager := repository.NewTxManager(pool)

	inventoryClient, err := inventoryclient.NewClient(cfg.InventoryGRPCAddr)
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer inventoryClient.Close()

	orderSvc := service.NewOrderService(orderRepo, txManager, inventoryClient)

	// Создаём gRPC сервер
	server := grpcserver.NewOrderServer(orderSvc)

	// Запускаем сервер
	log.Printf("Starting Order Service on :%s", cfg.GRPCPort)
	grpcServer, err := server.Start(cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to start Order Service: %v", err)
	}

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Order Service...")
	grpcServer.GracefulStop()
	log.Println("Order Service exited")
}

func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
