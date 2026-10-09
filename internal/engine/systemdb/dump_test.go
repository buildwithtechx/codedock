package systemdb

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	_ "github.com/jackc/pgx/v5/stdlib"

	"codedock/internal/repositories"
)

func TestDumpRestoreRoundtrip(t *testing.T) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("docker unavailable: %v", err)
	}
	ctx := context.Background()
	if _, err := docker.Ping(ctx); err != nil {
		t.Skipf("docker daemon unreachable: %v", err)
	}
	dataDir := t.TempDir()
	supervisor := NewSupervisor(docker, "codedock-postgres-dumptest", dataDir, "postgres:16", 5546)
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_ = docker.ContainerStop(cleanupCtx, "codedock-postgres-dumptest", container.StopOptions{})
		_ = docker.ContainerRemove(cleanupCtx, "codedock-postgres-dumptest", container.RemoveOptions{Force: true})
	})
	url, err := supervisor.EnsureRunning(ctx)
	if err != nil {
		t.Fatalf("ensure running: %v", err)
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO organizations (id, name) VALUES ('org-dump', 'Dump')`); err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	password, err := ReadPassword(dataDir)
	if err != nil {
		t.Fatalf("read password: %v", err)
	}
	dump, err := DumpContainer(ctx, "codedock-postgres-dumptest", DBUser, DBName, password)
	if err != nil {
		t.Fatalf("dump: %v", err)
	}
	if !strings.Contains(string(dump), "org-dump") {
		t.Fatal("dump is missing seeded row")
	}
	if _, err := db.ExecContext(ctx, `DROP SCHEMA public CASCADE`); err != nil {
		t.Fatalf("drop schema: %v", err)
	}
	if _, err := db.ExecContext(ctx, `CREATE SCHEMA public`); err != nil {
		t.Fatalf("recreate schema: %v", err)
	}
	if err := RestoreContainer(ctx, "codedock-postgres-dumptest", DBUser, DBName, password, dump); err != nil {
		t.Fatalf("restore: %v", err)
	}
	var name string
	if err := db.QueryRowContext(ctx, `SELECT name FROM organizations WHERE id = 'org-dump'`).Scan(&name); err != nil {
		t.Fatalf("read restored row: %v", err)
	}
	if name != "Dump" {
		t.Fatalf("unexpected restored row %q", name)
	}
}
