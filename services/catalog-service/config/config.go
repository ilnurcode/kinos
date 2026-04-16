// Package config предоставляет конфигурацию для catalog-service.
// Включает загрузку переменных окружения (DB_URL).
package config

import (
	"os"
)

type Config struct {
	DBURL string
}

func NewConfig() *Config {
	return &Config{
		DBURL: os.Getenv("DB_URL"),
	}
}
