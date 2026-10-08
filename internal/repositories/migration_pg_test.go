package repositories

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type pgColumnExpectation struct {
	table        string
	column       string
	dataType     string
	defaultValue string
}

func TestRunMigrationsPostgres(t *testing.T) {
	db := openPostgresTestDB(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `DROP SCHEMA public CASCADE`); err != nil {
		t.Fatalf("drop public schema: %v", err)
	}
	if _, err := db.ExecContext(ctx, `CREATE SCHEMA public`); err != nil {
		t.Fatalf("create public schema: %v", err)
	}
	if err := RunMigrationsDialect(db, DriverPostgres); err != nil {
		t.Fatalf("postgres migrations failed: %v", err)
	}
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatalf("count applied migrations: %v", err)
	}
	if count != 16 {
		t.Fatalf("expected 16 applied migrations, got %d", count)
	}
	for _, table := range []string{"users", "projects", "deployments", "managed_servers", "migration_runs", "traffic_buckets", "schema_migrations"} {
		var exists bool
		query := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)`
		if err := db.QueryRowContext(ctx, query, table).Scan(&exists); err != nil {
			t.Fatalf("check table %s: %v", table, err)
		}
		if !exists {
			t.Fatalf("expected table %s to exist", table)
		}
	}
	expectations := []pgColumnExpectation{
		{table: "deployments", column: "trigger", dataType: "text"},
		{table: "users", column: "is_active", dataType: "boolean", defaultValue: "true"},
		{table: "attention_issues", column: "first_seen", dataType: "timestamp with time zone"},
		{table: "schema_migrations", column: "applied_at", dataType: "timestamp with time zone", defaultValue: "CURRENT_TIMESTAMP"},
	}
	for _, want := range expectations {
		var dataType, columnDefault string
		query := `SELECT data_type, COALESCE(column_default, '') FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1 AND column_name = $2`
		if err := db.QueryRowContext(ctx, query, want.table, want.column).Scan(&dataType, &columnDefault); err != nil {
			t.Fatalf("describe %s.%s: %v", want.table, want.column, err)
		}
		if dataType != want.dataType {
			t.Fatalf("expected %s.%s to be %s, got %s", want.table, want.column, want.dataType, dataType)
		}
		if want.defaultValue != "" && columnDefault != want.defaultValue && columnDefault != "now()" {
			t.Fatalf("expected %s.%s default %s, got %s", want.table, want.column, want.defaultValue, columnDefault)
		}
	}
	if err := RunMigrationsDialect(db, DriverPostgres); err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatalf("recount applied migrations: %v", err)
	}
	if count != 16 {
		t.Fatalf("expected 16 records after two runs, got %d", count)
	}
}

func TestRunMigrationsUnknownDriver(t *testing.T) {
	db := openTestDB(t)
	if err := RunMigrationsDialect(db, "bogus"); err == nil {
		t.Fatal("expected error for unknown driver")
	}
}
