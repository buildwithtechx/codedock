package repositories

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresDialectProbe(t *testing.T) {
	db := openPostgresTestDB(t)
	ctx := context.Background()
	if DriverName(db) != DriverPostgres {
		t.Fatalf("expected postgres driver, got %q", DriverName(db))
	}
	if _, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS dialect_probe`); err != nil {
		t.Fatalf("drop probe table: %v", err)
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE dialect_probe (id TEXT PRIMARY KEY, name TEXT NOT NULL, detail TEXT NOT NULL DEFAULT '')`); err != nil {
		t.Fatalf("create probe table: %v", err)
	}
	t.Cleanup(func() { db.ExecContext(context.Background(), `DROP TABLE IF EXISTS dialect_probe`) })
	insert := Rebind(DriverPostgres, `INSERT INTO dialect_probe (id, name, detail) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET name = excluded.name, detail = excluded.detail`)
	if _, err := db.ExecContext(ctx, insert, "row-1", "Pilot", "line one"); err != nil {
		t.Fatalf("rebound upsert: %v", err)
	}
	if _, err := db.ExecContext(ctx, insert, "row-1", "Pilot", "line two"); err != nil {
		t.Fatalf("rebound upsert conflict path: %v", err)
	}
	truncate := `UPDATE dialect_probe SET detail = right(detail || $1, 8192) WHERE id = $2`
	if _, err := db.ExecContext(ctx, truncate, "\nline three", "row-1"); err != nil {
		t.Fatalf("tail truncate: %v", err)
	}
	var detail string
	err := db.QueryRowContext(ctx, Rebind(DriverPostgres, `SELECT detail FROM dialect_probe WHERE name ILIKE ?`), "pil%").Scan(&detail)
	if err != nil {
		t.Fatalf("case-insensitive match: %v", err)
	}
	if detail == "" {
		t.Fatal("expected stored detail")
	}
}
