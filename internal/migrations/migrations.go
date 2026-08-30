package migrations

import (
	"context"
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

//go:embed sql/*.up.sql sql/*.down.sql
var migrationFiles embed.FS

type migration struct {
	version  int64
	name     string
	sql      string
	checksum string
	down     string
}

// Up applies each migration exactly once, in version order.
func Up(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			checksum TEXT NOT NULL DEFAULT '',
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}
	// Keep existing installations compatible with the checksum column introduced later.
	if err := db.WithContext(ctx).Exec(`ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS checksum TEXT NOT NULL DEFAULT ''`).Error; err != nil {
		return fmt.Errorf("add migration checksum column: %w", err)
	}

	var applied []struct {
		Version  int64
		Name     string
		Checksum string
	}
	if err := db.WithContext(ctx).Table("schema_migrations").Select("version, name, checksum").Scan(&applied).Error; err != nil {
		return fmt.Errorf("read applied migrations: %w", err)
	}
	appliedSet := make(map[int64]struct {
		name, checksum string
	}, len(applied))
	for _, record := range applied {
		appliedSet[record.Version] = struct{ name, checksum string }{record.Name, record.Checksum}
	}

	migrations, err := load()
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		if record, ok := appliedSet[migration.version]; ok {
			if record.name != migration.name || record.checksum != migration.checksum {
				return fmt.Errorf("migration %d was modified after it was applied", migration.version)
			}
			continue
		}
		if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(migration.sql).Error; err != nil {
				return fmt.Errorf("execute migration %d (%s): %w", migration.version, migration.name, err)
			}
			return tx.Exec("INSERT INTO schema_migrations (version, name, checksum) VALUES (?, ?, ?)", migration.version, migration.name, migration.checksum).Error
		}); err != nil {
			return err
		}
	}
	return nil
}

// Down rolls back the most recently applied migration.
func Down(ctx context.Context, db *gorm.DB) error {
	var applied struct {
		Version int64
		Name    string
	}
	if err := db.WithContext(ctx).Table("schema_migrations").Order("version DESC").Take(&applied).Error; err != nil {
		return fmt.Errorf("find latest migration: %w", err)
	}
	items, err := load()
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.version != applied.Version || item.name != applied.Name {
			continue
		}
		if strings.TrimSpace(item.down) == "" {
			return fmt.Errorf("migration %d (%s) has no rollback", item.version, item.name)
		}
		if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(item.down).Error; err != nil {
				return fmt.Errorf("rollback migration %d (%s): %w", item.version, item.name, err)
			}
			return tx.Exec("DELETE FROM schema_migrations WHERE version = ?", item.version).Error
		}); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("migration %d (%s) is not present in source", applied.Version, applied.Name)
}

func load() ([]migration, error) {
	entries, err := fs.Glob(migrationFiles, "sql/*.up.sql")
	if err != nil {
		return nil, fmt.Errorf("discover migrations: %w", err)
	}
	result := make([]migration, 0, len(entries))
	for _, entry := range entries {
		base := filepath.Base(entry)
		parts := strings.SplitN(base, "_", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid migration filename %q", base)
		}
		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil || version < 1 {
			return nil, fmt.Errorf("invalid migration version in %q", base)
		}
		contents, err := migrationFiles.ReadFile(entry)
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", base, err)
		}
		name := strings.TrimSuffix(parts[1], ".up.sql")
		down, err := migrationFiles.ReadFile(strings.TrimSuffix(entry, ".up.sql") + ".down.sql")
		if err != nil {
			return nil, fmt.Errorf("read rollback for %q: %w", base, err)
		}
		result = append(result, migration{version: version, name: name, sql: string(contents), checksum: fmt.Sprintf("%x", sha256.Sum256(contents)), down: string(down)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].version < result[j].version })
	for i := 1; i < len(result); i++ {
		if result[i-1].version == result[i].version {
			return nil, fmt.Errorf("duplicate migration version %d", result[i].version)
		}
	}
	return result, nil
}
