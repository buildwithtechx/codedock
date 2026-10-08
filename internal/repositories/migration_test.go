package repositories

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrationFreshDatabase(t *testing.T) {
	db := openPGBareDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE mig_users (id TEXT PRIMARY KEY);")},
		"002_add_email.sql":    {Data: []byte("ALTER TABLE mig_users ADD COLUMN email TEXT;")},
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
	db := openPGBareDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE mig_users (id TEXT PRIMARY KEY);")},
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
	db := openPGBareDB(t)
	fsys := fstest.MapFS{
		"002_add_email.sql":    {Data: []byte("ALTER TABLE mig_users ADD COLUMN email TEXT;")},
		"001_create_users.sql": {Data: []byte("CREATE TABLE mig_users (id TEXT PRIMARY KEY);")},
	}
	if err := runMigrations(db, fsys); err != nil {
		t.Fatalf("migrations applied out of order or failed: %v", err)
	}
}

func TestMigrationFailedRollback(t *testing.T) {
	db := openPGBareDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE mig_users (id TEXT PRIMARY KEY);")},
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
	if err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'mig_users'").Scan(&exists); err != nil {
		t.Fatalf("failed to check users table existence: %v", err)
	}
	if exists != 1 {
		t.Fatal("expected users table from 001 to still exist after 002 rollback")
	}
}

func TestMigrationChecksumMismatch(t *testing.T) {
	db := openPGBareDB(t)
	original := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE mig_users (id TEXT PRIMARY KEY);")},
	}
	if err := runMigrations(db, original); err != nil {
		t.Fatalf("initial migration failed: %v", err)
	}
	modified := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE mig_users (id TEXT PRIMARY KEY, injected TEXT);")},
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
	db := openPGBareDB(t)
	fsys := fstest.MapFS{
		"001_create_users.sql": {Data: []byte("CREATE TABLE mig_users (id TEXT PRIMARY KEY);")},
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
func TestSplitStatements(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    []string
	}{
		{name: "simple", content: "SELECT 1;SELECT 2;", want: []string{"SELECT 1", "SELECT 2"}},
		{name: "quoted identifier", content: "SELECT \"we;ird\";", want: []string{"SELECT \"we;ird\""}},
		{name: "line comment", content: "-- trailing; comment\nSELECT 1;", want: []string{"-- trailing; comment\nSELECT 1"}},
		{name: "dollar quote", content: "SELECT $$a;b$$;SELECT 2;", want: []string{"SELECT $$a;b$$", "SELECT 2"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := splitStatements(tc.content)
			if len(got) != len(tc.want) {
				t.Fatalf("expected %d statements, got %d (%q)", len(tc.want), len(got), got)
			}
			for i := range tc.want {
				if strings.TrimSpace(got[i]) != tc.want[i] {
					t.Fatalf("statement %d: expected %q, got %q", i, tc.want[i], got[i])
				}
			}
		})
	}
}
