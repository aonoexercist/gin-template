package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Env             string // development | production
	Port            string
	DatabaseURL     string
	LogLevel        string // debug | info | warn | error
	ShutdownTimeout time.Duration
}

// Load reads configuration from environment variables (12-factor).
func Load() (*Config, error) {
	cfg := &Config{
		Env:         get("APP_ENV", "development"),
		Port:        get("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    get("LOG_LEVEL", "info"),
	}

	timeout, err := time.ParseDuration(get("SHUTDOWN_TIMEOUT", "15s"))
	if err != nil {
		return nil, errors.New("invalid SHUTDOWN_TIMEOUT")
	}
	cfg.ShutdownTimeout = timeout

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	return cfg, nil
}

func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
