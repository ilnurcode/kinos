// Package config предоставляет конфигурацию для inventory-service.
// Включает загрузку переменных окружения (DB_URL).
package config

import (
	"os"
)

type Config struct {
	DBURL    string
	GRPCPort string
}

func NewConfig() *Config {
	return &Config{
		DBURL:    os.Getenv("DB_URL"),
		GRPCPort: getEnv("GRPC_PORT", "8083"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
