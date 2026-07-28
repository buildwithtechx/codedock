package repositories

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"
)

//go:embed schema/*.sql
var schemaFS embed.FS

func RunMigrations(db *sql.DB) error {
	sub, err := fs.Sub(schemaFS, "schema")
	if err != nil {
		return fmt.Errorf("failed to sub schema fs: %w", err)
	}
	return runMigrations(db, sub)
}

func runMigrations(db *sql.DB, fsys fs.FS) error {
	if err := createMigrationsTable(db); err != nil {
		return err
	}
	applied, err := loadApplied(db)
	if err != nil {
		return err
	}

	files, err := schemaFiles(fsys)
	if err != nil {
		return err
	}
	for _, file := range files {
		content, err := fs.ReadFile(fsys, file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}
		sum := fileChecksum(content)
		if storedSum, done := applied[file]; done {
			if storedSum != sum {
				return fmt.Errorf("migration %s was modified after being applied (checksum mismatch)", file)
			}
			continue
		}
		if err := applyMigration(db, file, content, sum); err != nil {
			return err
		}
	}
	slog.Info("schema migrations up to date")
	return nil
}

func applyMigration(db *sql.DB, file string, content []byte, sum string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for %s: %w", file, err)
	}
	rollback := func(cause error) error {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%w (rollback: %v)", cause, rbErr)
		}
		return cause
	}
	if _, err := tx.Exec(string(content)); err != nil {
		return rollback(fmt.Errorf("migration failed for %s: %w", file, err))
	}
	if _, err := tx.Exec("INSERT INTO schema_migrations (filename, checksum) VALUES (?, ?)", file, sum); err != nil {
		return rollback(fmt.Errorf("failed to record migration %s: %w", file, err))
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
			checksum   TEXT NOT NULL DEFAULT '',
			applied_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return nil
}

func loadApplied(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT filename, checksum FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to load applied migrations: %w", err)
	}
	defer rows.Close()
	applied := make(map[string]string)
	for rows.Next() {
		var name, sum string
		if err := rows.Scan(&name, &sum); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		applied[name] = sum
	}
	return applied, rows.Err()
}

func schemaFiles(fsys fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read schema directory: %w", err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}

func fileChecksum(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}
