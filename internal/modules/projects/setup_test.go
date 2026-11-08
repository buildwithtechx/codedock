package projects

import (
	"database/sql"
	"testing"

	"codedock/internal/testdb"
)

func openPGTestDB(t *testing.T) *sql.DB {
	t.Helper()
	return testdb.Open(t)
}
