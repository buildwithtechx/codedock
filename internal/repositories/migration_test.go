package repositories

import (
	"database/sql"
	"strings"
	"testing"
	"testing/fstest"

	_ "modernc.org/sqlite"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestMigrationFreshDatabase(t *testing.T) {
	db := openTestDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE users (id TEXT PRIMARY KEY);")},
		"002_add_email.sql":    {Data: []byte("ALTER TABLE users ADD COLUMN email TEXT;")},
	}
	if err := runMigrations(db, fsys); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("failed to scan migration count: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 applied migrations, got %d", count)
	}
}

func TestMigrationIdempotent(t *testing.T) {
	db := openTestDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE users (id TEXT PRIMARY KEY);")},
	}
	if err := runMigrations(db, fsys); err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	if err := runMigrations(db, fsys); err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("failed to scan migration count: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 record after two runs, got %d", count)
	}
}

func TestMigrationOrdering(t *testing.T) {
	db := openTestDB(t)
	fsys := fstest.MapFS{
		"002_add_email.sql":    {Data: []byte("ALTER TABLE users ADD COLUMN email TEXT;")},
		"001_create_users.sql": {Data: []byte("CREATE TABLE users (id TEXT PRIMARY KEY);")},
	}
	if err := runMigrations(db, fsys); err != nil {
		t.Fatalf("migrations applied out of order or failed: %v", err)
	}
}

func TestMigrationFailedRollback(t *testing.T) {
	db := openTestDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE users (id TEXT PRIMARY KEY);")},
		"002_bad.sql":          {Data: []byte("NOT VALID SQL !!!;")},
	}
	if err := runMigrations(db, fsys); err == nil {
		t.Fatal("expected migration to fail on bad SQL")
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("failed to scan migration count: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected only 001 to be recorded, got %d", count)
	}
	var exists int
	if err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&exists); err != nil {
		t.Fatalf("failed to check users table existence: %v", err)
	}
	if exists != 1 {
		t.Fatal("expected users table from 001 to still exist after 002 rollback")
	}
}

func TestMigrationChecksumMismatch(t *testing.T) {
	db := openTestDB(t)
	original := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE users (id TEXT PRIMARY KEY);")},
	}
	if err := runMigrations(db, original); err != nil {
		t.Fatalf("initial migration failed: %v", err)
	}
	modified := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE users (id TEXT PRIMARY KEY, injected TEXT);")},
	}
	err := runMigrations(db, modified)
	if err == nil {
		t.Fatal("expected error on modified migration")
	}
	if !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch error, got: %v", err)
	}
}

func TestMigrationSQLOnlyFiles(t *testing.T) {
	db := openTestDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE users (id TEXT PRIMARY KEY);")},
		"README.md":            {Data: []byte("not sql")},
		".gitkeep":             {Data: []byte("")},
	}
	if err := runMigrations(db, fsys); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("failed to scan migration count: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected only .sql files to be applied, got %d migrations", count)
	}
}
