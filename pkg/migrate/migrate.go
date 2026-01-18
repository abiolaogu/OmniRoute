// Package migrate provides database migration utilities for OmniRoute services.
package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration represents a single migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// Migrator handles database migrations
type Migrator struct {
	db             *sql.DB
	migrationsPath string
	tableName      string
}

// NewMigrator creates a new migrator
func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
		tableName:      "schema_migrations",
	}
}

// Init creates the migrations tracking table
func (m *Migrator) Init(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`, m.tableName)

	_, err := m.db.ExecContext(ctx, query)
	return err
}

// Up applies all pending migrations
func (m *Migrator) Up(ctx context.Context) error {
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	for _, migration := range migrations {
		if _, ok := applied[migration.Version]; ok {
			continue // Already applied
		}

		if migration.UpSQL == "" {
			continue // No up migration
		}

		fmt.Printf("Applying migration %s: %s\n", migration.Version, migration.Name)

		tx, err := m.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}

		if _, err := tx.ExecContext(ctx, migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", migration.Version, err)
		}

		if _, err := tx.ExecContext(ctx,
			fmt.Sprintf("INSERT INTO %s (version, name) VALUES ($1, $2)", m.tableName),
			migration.Version, migration.Name,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", migration.Version, err)
		}

		fmt.Printf("✓ Applied %s\n", migration.Version)
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down(ctx context.Context) error {
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	// Find the last applied migration
	var lastApplied *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if _, ok := applied[migrations[i].Version]; ok {
			lastApplied = migrations[i]
			break
		}
	}

	if lastApplied == nil {
		fmt.Println("No migrations to roll back")
		return nil
	}

	if lastApplied.DownSQL == "" {
		return fmt.Errorf("migration %s has no down SQL", lastApplied.Version)
	}

	fmt.Printf("Rolling back migration %s: %s\n", lastApplied.Version, lastApplied.Name)

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, lastApplied.DownSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("rollback migration %s: %w", lastApplied.Version, err)
	}

	if _, err := tx.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE version = $1", m.tableName),
		lastApplied.Version,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("remove migration record %s: %w", lastApplied.Version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit rollback %s: %w", lastApplied.Version, err)
	}

	fmt.Printf("✓ Rolled back %s\n", lastApplied.Version)
	return nil
}

// Status prints the migration status
func (m *Migrator) Status(ctx context.Context) error {
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")
	for _, migration := range migrations {
		status := "[ ] Pending"
		if appliedAt, ok := applied[migration.Version]; ok {
			status = fmt.Sprintf("[✓] Applied at %s", appliedAt.Format(time.RFC3339))
		}
		fmt.Printf("%s %s: %s\n", status, migration.Version, migration.Name)
	}

	return nil
}

func (m *Migrator) loadMigrations() ([]*Migration, error) {
	files, err := os.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	migrationMap := make(map[string]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		// Parse filename: 001_initial_schema.up.sql or 001_initial_schema.down.sql
		parts := strings.Split(strings.TrimSuffix(name, ".sql"), ".")
		if len(parts) != 2 {
			continue
		}

		direction := parts[1]
		versionName := parts[0]
		versionParts := strings.SplitN(versionName, "_", 2)
		if len(versionParts) != 2 {
			continue
		}

		version := versionParts[0]
		migrationName := versionParts[1]

		content, err := os.ReadFile(filepath.Join(m.migrationsPath, name))
		if err != nil {
			return nil, fmt.Errorf("read migration file %s: %w", name, err)
		}

		if migrationMap[version] == nil {
			migrationMap[version] = &Migration{
				Version: version,
				Name:    migrationName,
			}
		}

		if direction == "up" {
			migrationMap[version].UpSQL = string(content)
		} else if direction == "down" {
			migrationMap[version].DownSQL = string(content)
		}
	}

	var migrations []*Migration
	for _, m := range migrationMap {
		migrations = append(migrations, m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (m *Migrator) getAppliedMigrations(ctx context.Context) (map[string]*time.Time, error) {
	rows, err := m.db.QueryContext(ctx,
		fmt.Sprintf("SELECT version, applied_at FROM %s ORDER BY version", m.tableName),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]*time.Time)
	for rows.Next() {
		var version string
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, err
		}
		applied[version] = &appliedAt
	}

	return applied, rows.Err()
}
