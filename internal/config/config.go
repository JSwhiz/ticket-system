package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config структура для хранения конфигурации
type Config struct {
	DBURL      string
	JWTSecret  string
	ServerPort string
	LogLevel   string
}

// Load загружает конфигурацию из .env
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		DBURL:      os.Getenv("DB_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("SERVER_PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}, nil
}