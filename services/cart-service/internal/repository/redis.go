package repository

import (
	"kinos/cart-service/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword, DB: cfg.RedisDB})
	return &RedisClient{client: client}
}

func (r *RedisClient) Close() error{
	return r.client.Close()
}