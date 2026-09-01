package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Environment  = "ENV"
	Port         = "PORT"
	DBURL        = "DB_URL"
	SecretKey    = "SECRET_KEY"
	MetricsToken = "METRICS_TOKEN"
)

type App struct {
	Environment    string
	Port           string
	DatabaseURL    string
	SecretKey      string
	MetricsToken   string
	AllowedOrigins []string
	AccessTokenTTL time.Duration
}

func Load() (App, error) {
	cfg := App{
		Environment:    valueOrDefault(Environment, "development"),
		Port:           valueOrDefault(Port, "8000"),
		DatabaseURL:    strings.TrimSpace(os.Getenv(DBURL)),
		SecretKey:      os.Getenv(SecretKey),
		MetricsToken:   strings.TrimSpace(os.Getenv(MetricsToken)),
		AllowedOrigins: splitAndTrim(valueOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		AccessTokenTTL: 15 * time.Minute,
	}

	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return App{}, fmt.Errorf("PORT must be numeric: %w", err)
	}
	if cfg.DatabaseURL == "" {
		return App{}, fmt.Errorf("DB_URL is required")
	}
	if len(cfg.SecretKey) < 32 {
		return App{}, fmt.Errorf("SECRET_KEY must contain at least 32 characters")
	}
	if raw := strings.TrimSpace(os.Getenv("ACCESS_TOKEN_TTL")); raw != "" {
		ttl, err := time.ParseDuration(raw)
		if err != nil || ttl <= 0 {
			return App{}, fmt.Errorf("ACCESS_TOKEN_TTL must be a positive duration")
		}
		cfg.AccessTokenTTL = ttl
	}

	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}
