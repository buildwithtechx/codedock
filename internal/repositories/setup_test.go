package repositories

import (
	"database/sql"
	"testing"

	cdkdb "codedock/db"
	"codedock/internal/testdb"
)

func openPGTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := testdb.Open(t)
	if err := cdkdb.RunMigrations(db); err != nil {
		t.Fatalf("migrate isolated test database: %v", err)
	}
	return db
}
