package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Migrator handles database schema migrations using plain SQL files.
// It tracks applied migrations in a schema_migrations table.
type Migrator struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
	dir    string
}

// NewMigrator creates a new Migrator instance.
func NewMigrator(pool *pgxpool.Pool, migrationsDir string, logger *zap.Logger) *Migrator {
	return &Migrator{
		pool:   pool,
		logger: logger,
		dir:    migrationsDir,
	}
}

// MigrateUp applies all pending migrations in order.
func (m *Migrator) MigrateUp(ctx context.Context) error {
	if err := m.ensureMigrationsTable(ctx); err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	files, err := m.getMigrationFiles("up")
	if err != nil {
		return fmt.Errorf("read migration files: %w", err)
	}

	for _, file := range files {
		version := extractVersion(file)
		if _, ok := applied[version]; ok {
			continue
		}

		m.logger.Info("applying migration", zap.String("file", file), zap.String("version", version))

		content, err := os.ReadFile(filepath.Join(m.dir, file))
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", file, err)
		}

		tx, err := m.pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin transaction for migration %s: %w", version, err)
		}

		if _, err := tx.Exec(ctx, string(content)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("execute migration %s: %w", version, err)
		}

		if _, err := tx.Exec(ctx,
			"INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)",
			version, time.Now().UTC(),
		); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", version, err)
		}

		m.logger.Info("migration applied successfully", zap.String("version", version))
	}

	return nil
}

func (m *Migrator) ensureMigrationsTable(ctx context.Context) error {
	_, err := m.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(20) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (m *Migrator) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	rows, err := m.pool.Query(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

func (m *Migrator) getMigrationFiles(direction string) ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, err
	}

	suffix := "." + direction + ".sql"
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}

// extractVersion extracts the version number from a migration filename.
// Expected format: 000001_description.up.sql -> "000001"
func extractVersion(filename string) string {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}
