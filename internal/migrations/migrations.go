package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

//go:embed sql/*.sql
var migrationFiles embed.FS

type migration struct {
	version int64
	name    string
	sql     string
}

// Up applies each migration exactly once, in version order.
func Up(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	var applied []int64
	if err := db.WithContext(ctx).Table("schema_migrations").Pluck("version", &applied).Error; err != nil {
		return fmt.Errorf("read applied migrations: %w", err)
	}
	appliedSet := make(map[int64]struct{}, len(applied))
	for _, version := range applied {
		appliedSet[version] = struct{}{}
	}

	migrations, err := load()
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		if _, ok := appliedSet[migration.version]; ok {
			continue
		}
		if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(migration.sql).Error; err != nil {
				return fmt.Errorf("execute migration %d (%s): %w", migration.version, migration.name, err)
			}
			return tx.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)", migration.version, migration.name).Error
		}); err != nil {
			return err
		}
	}
	return nil
}

func load() ([]migration, error) {
	entries, err := fs.Glob(migrationFiles, "sql/*.sql")
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
		result = append(result, migration{version: version, name: strings.TrimSuffix(parts[1], ".sql"), sql: string(contents)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].version < result[j].version })
	for i := 1; i < len(result); i++ {
		if result[i-1].version == result[i].version {
			return nil, fmt.Errorf("duplicate migration version %d", result[i].version)
		}
	}
	return result, nil
}
