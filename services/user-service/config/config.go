// Package config предоставляет конфигурацию для user-service.
// Включает загрузку переменных окружения (DB_URL, SECRET_KEY).
package config

import (
	"os"
)

type Config struct {
	DBURL     string
	SecretKey string
}

func NewConfig() *Config {
	return &Config{
		DBURL:     os.Getenv("DB_URL"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}
}
