package db

import (
	"database/sql"
	"testing"

	"codedock/internal/testdb"
)

func openPGTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := testdb.Open(t)
	if err := RunMigrations(db); err != nil {
		t.Fatalf("migrate isolated test database: %v", err)
	}
	return db
}

func TestPostgresPing(t *testing.T) {
	db := openPGTestDB(t)
	if err := db.Ping(); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}
}
