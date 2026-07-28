package repositories

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"sort"
)

//go:embed schema/*.sql
var schemaFS embed.FS

func RunMigrations(db *sql.DB) error {
	if err := createMigrationsTable(db); err != nil {
		return err
	}
	applied, err := loadApplied(db)
	if err != nil {
		return err
	}
	files, err := schemaFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if applied[file] {
			continue
		}
		if err := applyMigration(db, file); err != nil {
			return err
		}
	}
	slog.Info("schema migrations up to date")
	return nil
}

func applyMigration(db *sql.DB, file string) error {
	content, err := schemaFS.ReadFile("schema/" + file)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", file, err)
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for %s: %w", file, err)
	}
	if _, err := tx.Exec(string(content)); err != nil {
		tx.Rollback()
		return fmt.Errorf("migration failed for %s: %w", file, err)
	}
	if _, err := tx.Exec("INSERT INTO schema_migrations (filename) VALUES (?)", file); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration %s: %w", file, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", file, err)
	}
	slog.Info("applied migration", "file", file)
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename   TEXT PRIMARY KEY,
			applied_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return nil
}

func loadApplied(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT filename FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to load applied migrations: %w", err)
	}
	defer rows.Close()
	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		applied[name] = true
	}
	return applied, rows.Err()
}

func schemaFiles() ([]string, error) {
	entries, err := schemaFS.ReadDir("schema")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded schema directory: %w", err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}
