package config

import (
	"time"
	
    "github.com/caarlos0/env/v6"
)

type Config struct {
    DBURL      string        `env:"DB_URL,required"`
    JWTSecret  string        `env:"JWT_SECRET,required"`
    ServerPort string        `env:"SERVER_PORT" envDefault:"8080"`
    LogLevel   string        `env:"LOG_LEVEL" envDefault:"info"`
	Timeout    time.Duration `env:"TIMEOUT" envDefault:"5s"`
}

func Load() (*Config, error) {
    var cfg Config
    if err := env.Parse(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
