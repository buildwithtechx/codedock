package repositories

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRebindKeepsSQLitePlaceholders(t *testing.T) {
	query := `SELECT id FROM servers WHERE id = ? AND status = ?`
	if got := Rebind(DriverSQLite, query); got != query {
		t.Fatalf("expected unchanged query, got %q", got)
	}
	if got := Rebind("", query); got != query {
		t.Fatalf("expected unchanged query for empty driver, got %q", got)
	}
}

func TestRebindNumbersPostgresPlaceholders(t *testing.T) {
	query := `UPDATE servers SET status = ?, updated_at = ? WHERE id = ?`
	want := `UPDATE servers SET status = $1, updated_at = $2 WHERE id = $3`
	if got := Rebind(DriverPostgres, query); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDriverNameDefaultsToSQLite(t *testing.T) {
	if got := DriverName(nil); got != "" {
		t.Fatalf("expected empty driver for nil db, got %q", got)
	}
	db, err := sql.Open(DriverSQLite, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if got := DriverName(db); got != DriverSQLite {
		t.Fatalf("expected sqlite driver, got %q", got)
	}
}
