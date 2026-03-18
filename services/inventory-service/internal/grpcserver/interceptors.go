// Package grpcserver предоставляет gRPC-сервер для inventory-service.
// Включает interceptor для логирования запросов.
package grpcserver

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
)

// LoggingInterceptor логирует входящие gRPC запросы
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v\n%s", r, debug.Stack())
		}
	}()

	resp, err := handler(ctx, req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("gRPC error on %s: %v (%s)", info.FullMethod, err, duration)
	} else {
		log.Printf("gRPC success on %s (%s)", info.FullMethod, duration)
	}

	return resp, err
}

// RecoveryInterceptor обрабатывает паники в gRPC запросах
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic in gRPC handler %s: %v\n%s", info.FullMethod, r, debug.Stack())
			}
		}()

		return handler(ctx, req)
	}
}
