// Package main запускает cart-service — gRPC-сервис для управления корзиной покупок.
// Хранит данные корзины в Redis с настраиваемым TTL.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kinos/cart-service/config"
	"kinos/cart-service/internal/catalog"
	"kinos/cart-service/internal/grpcserver"
	"kinos/cart-service/internal/inventory"
	"kinos/cart-service/internal/repository"
	"kinos/cart-service/internal/service"
	pb "kinos/proto/cart"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.NewConfig()
	if cfg.RedisAddr == "" {
		log.Fatal("REDIS_ADDRESS env var must be set")
	}

	// Подключение к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	log.Println("Connected to Redis")

	// gRPC-клиенты к внешним сервисам
	catalogClient, err := catalog.NewCatalogClient(cfg.CatalogGRPCAddr)
	if err != nil {
		log.Fatalf("Failed to connect to catalog service: %v", err)
	}
	defer catalogClient.Close()

	inventoryClient, err := inventory.NewInventoryClient(cfg.InventoryGRPCAddr)
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer inventoryClient.Close()

	// Репозиторий и сервис
	cartRepo := repository.NewCartRepository(rdb, cfg.CartTTL)
	cartSvc := service.NewCartService(cartRepo, catalogClient, inventoryClient)

	// gRPC-сервер
	cartServer := grpcserver.NewCartServer(cartSvc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCartServiceServer(grpcServer, cartServer)
	reflection.Register(grpcServer)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Cart gRPC server started on :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Cart Service...")
	grpcServer.GracefulStop()
	log.Println("Cart Service exited")
}
