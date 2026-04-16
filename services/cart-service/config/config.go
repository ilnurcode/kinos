package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCPort          string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	CatalogGRPCAddr   string
	InventoryGRPCAddr string
	CartTTL           time.Duration
}

func NewConfig() *Config {
	redisDB := 0
	if val := os.Getenv("REDIS_DB"); val != "" {
		db, err := strconv.Atoi(val)
		if err != nil {
			redisDB = 0
		} else {
			redisDB = db
		}
	}

	cartTTL := 24 * time.Hour
	if val := os.Getenv("CART_TTL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cartTTL = d
		}
	}

	return &Config{
		GRPCPort:          getEnv("CART_GRPC_PORT", "8084"),
		RedisAddr:         os.Getenv("REDIS_ADDRESS"),
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		RedisDB:           redisDB,
		CatalogGRPCAddr:   getEnv("CATALOG_GRPC_ADDR", "catalog-service:8082"),
		InventoryGRPCAddr: getEnv("INVENTORY_GRPC_ADDR", "inventory-service:8083"),
		CartTTL:           cartTTL,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
