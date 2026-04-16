// Package grpcserver предоставляет gRPC-сервер для inventory-service.
// Включает interceptor для логирования запросов.
package grpcserver

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// LoggingInterceptor логирует входящие gRPC запросы
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("gRPC error on %s: %v (%s)", info.FullMethod, err, duration)
	} else {
		log.Printf("gRPC success on %s (%s)", info.FullMethod, duration)
	}

	return resp, err
}
