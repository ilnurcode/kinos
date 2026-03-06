// Package grpcserver предоставляет gRPC-сервер для catalog-service.
// Включает interceptor для логирования запросов.
package grpcserver

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}
