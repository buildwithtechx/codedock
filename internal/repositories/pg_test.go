package repositories

import (
	"database/sql"
	"os"
	"testing"
)

func openPostgresTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("CODEDOCK_TEST_PG_URL")
	if dsn == "" {
		t.Skip("CODEDOCK_TEST_PG_URL is not set")
	}
	db, err := sql.Open(DriverPostgres, dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.Ping(); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}
	return db
}
