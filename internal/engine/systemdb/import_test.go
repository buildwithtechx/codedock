package systemdb

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"codedock.run/codedock/internal/repositories"
)

func openImportTarget(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("CODEDOCK_TEST_PG_URL")
	if dsn == "" {
		t.Skip("CODEDOCK_TEST_PG_URL is not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `DROP SCHEMA public CASCADE`); err != nil {
		t.Fatalf("drop public schema: %v", err)
	}
	if _, err := db.ExecContext(ctx, `CREATE SCHEMA public`); err != nil {
		t.Fatalf("create public schema: %v", err)
	}
	if err := repositories.RunMigrationsDialect(db, repositories.DriverPostgres); err != nil {
		t.Fatalf("migrate target: %v", err)
	}
	return db
}

func buildLegacyFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "codedock.db")
	litedb, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer litedb.Close()
	if err := repositories.RunMigrationsDialect(litedb, repositories.DriverSQLite); err != nil {
		t.Fatalf("migrate fixture: %v", err)
	}
	login := time.Date(2026, 3, 14, 9, 26, 53, 0, time.UTC)
	if _, err := litedb.Exec(`INSERT INTO organizations (id, name, slug, is_team) VALUES ('org-1', 'Acme', 'acme', 1)`); err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	if _, err := litedb.Exec(`INSERT INTO users (id, email, name, password_hash, role, email_verified, is_active, last_login) VALUES ('user-1', 'ada@example.com', 'Ada', 'hash', 'owner', 1, 1, ?)`, login); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if _, err := litedb.Exec(`INSERT INTO users (id, email, name, password_hash, is_active) VALUES ('user-2', 'inactive@example.com', 'Inactive', 'hash', 0)`); err != nil {
		t.Fatalf("seed second user: %v", err)
	}
	if _, err := litedb.Exec(`INSERT INTO organization_members (id, organization_id, user_id, email, role) VALUES ('mem-1', 'org-1', 'user-1', 'ada@example.com', 'owner')`); err != nil {
		t.Fatalf("seed membership: %v", err)
	}
	return path
}

func TestImportFromSQLite(t *testing.T) {
	pgdb := openImportTarget(t)
	ctx := context.Background()
	path := buildLegacyFixture(t)
	if err := ImportFromSQLite(ctx, pgdb, path); err != nil {
		t.Fatalf("import: %v", err)
	}
	var users, members int
	if err := pgdb.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&users); err != nil {
		t.Fatalf("count users: %v", err)
	}
	if err := pgdb.QueryRowContext(ctx, `SELECT COUNT(*) FROM organization_members`).Scan(&members); err != nil {
		t.Fatalf("count members: %v", err)
	}
	if users != 2 || members != 1 {
		t.Fatalf("expected 2 users and 1 member, got %d and %d", users, members)
	}
	var active bool
	var login time.Time
	if err := pgdb.QueryRowContext(ctx, `SELECT is_active, last_login FROM users WHERE id = 'user-1'`).Scan(&active, &login); err != nil {
		t.Fatalf("read imported user: %v", err)
	}
	if !active {
		t.Fatal("expected imported user to be active")
	}
	if login.Year() != 2026 || login.Month() != 3 || login.Day() != 14 {
		t.Fatalf("expected imported login timestamp, got %v", login)
	}
	var inactive bool
	if err := pgdb.QueryRowContext(ctx, `SELECT is_active FROM users WHERE id = 'user-2'`).Scan(&inactive); err != nil {
		t.Fatalf("read inactive user: %v", err)
	}
	if inactive {
		t.Fatal("expected imported user to stay inactive")
	}
	var team bool
	if err := pgdb.QueryRowContext(ctx, `SELECT is_team FROM organizations WHERE id = 'org-1'`).Scan(&team); err != nil {
		t.Fatalf("read imported organization: %v", err)
	}
	if !team {
		t.Fatal("expected imported organization to stay a team")
	}
}

func TestImportRefusesNonEmptyTarget(t *testing.T) {
	pgdb := openImportTarget(t)
	ctx := context.Background()
	path := buildLegacyFixture(t)
	if err := ImportFromSQLite(ctx, pgdb, path); err != nil {
		t.Fatalf("first import: %v", err)
	}
	if err := ImportFromSQLite(ctx, pgdb, path); err == nil {
		t.Fatal("expected second import to be refused")
	}
}
