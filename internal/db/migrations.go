package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"codedock/internal/platform/database"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func RunMigrations(db *sql.DB) error {
	sub, err := fs.Sub(MigrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to sub migrations fs: %w", err)
	}
	return database.RunMigrationsFromFS(db, sub)
}

func RunMigrationsFromFS(db *sql.DB, fsys fs.FS) error {
	return database.RunMigrationsFromFS(db, fsys)
}

func SplitStatements(content string) []string {
	return database.SplitStatements(content)
}

func schemaFiles(fsys fs.FS) ([]string, error) {
	return database.SchemaFiles(fsys)
}

func fileChecksum(content []byte) string {
	return database.FileChecksum(content)
}
