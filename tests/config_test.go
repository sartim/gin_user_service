package tests

import (
	"testing"

	"gin-shop-api/internal/config"
)

func TestConfigRequiresStrongSecret(t *testing.T) {
	t.Setenv(config.DBURL, "postgres://user:password@localhost/database")
	t.Setenv(config.SecretKey, "short")

	if _, err := config.Load(); err == nil {
		t.Fatal("expected weak secret to be rejected")
	}
}

func TestConfigDefaults(t *testing.T) {
	t.Setenv(config.DBURL, "postgres://user:password@localhost/database")
	t.Setenv(config.SecretKey, "test-secret-with-at-least-32-characters")
	t.Setenv(config.Port, "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Port != "8000" {
		t.Fatalf("expected port 8000, got %s", cfg.Port)
	}
}
