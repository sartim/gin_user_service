package tests

import (
	"context"
	"os"
	"testing"

	"gin-shop-api/internal/migrations"
	"gin-shop-api/internal/repository"
)

func TestPostgresMigrationsAreIdempotent(t *testing.T) {
	dsn := os.Getenv("INTEGRATION_DB_URL")
	if dsn == "" {
		t.Skip("INTEGRATION_DB_URL is not set")
	}

	db, err := repository.Open(dsn)
	if err != nil {
		t.Fatalf("open integration database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get integration pool: %v", err)
	}
	defer sqlDB.Close()

	ctx := context.Background()
	if err := migrations.Up(ctx, db); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if err := migrations.Up(ctx, db); err != nil {
		t.Fatalf("reapply migrations: %v", err)
	}

	var count int64
	if err := db.Table("schema_migrations").Count(&count).Error; err != nil {
		t.Fatalf("count applied migrations: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one applied migration, got %d", count)
	}

	if err := migrations.Down(ctx, db); err != nil {
		t.Fatalf("rollback migration: %v", err)
	}
}
