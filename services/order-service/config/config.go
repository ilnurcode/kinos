package config

import "os"

type Config struct {
	GRPCPort          string
	DBURL             string
	InventoryGRPCAddr string
}

func NewConfig() *Config {
	return &Config{
		GRPCPort:          getEnv("GRPC_PORT", "8085"),
		DBURL:             getEnv("DB_URL", ""),
		InventoryGRPCAddr: getEnv("INVENTORY_GRPC_ADDR", "inventory-service:8083"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
