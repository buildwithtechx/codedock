package systemdb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestSupervisorLifecycle(t *testing.T) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("docker unavailable: %v", err)
	}
	ctx := context.Background()
	if _, err := docker.Ping(ctx); err != nil {
		t.Skipf("docker daemon unreachable: %v", err)
	}
	supervisor := NewSupervisor(docker, "codedock-postgres-test", t.TempDir(), "postgres:16", 5544)
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_ = docker.ContainerStop(cleanupCtx, "codedock-postgres-test", container.StopOptions{})
		_ = docker.ContainerRemove(cleanupCtx, "codedock-postgres-test", container.RemoveOptions{Force: true})
	})
	url, err := supervisor.EnsureRunning(ctx)
	if err != nil {
		t.Fatalf("ensure running: %v", err)
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatalf("open provisioned postgres: %v", err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping provisioned postgres: %v", err)
	}
	again, err := supervisor.EnsureRunning(ctx)
	if err != nil {
		t.Fatalf("second ensure running: %v", err)
	}
	if again != url {
		t.Fatal("expected stable database URL across runs")
	}
	var user, dbname string
	if err := db.QueryRowContext(ctx, `SELECT current_user, current_database()`).Scan(&user, &dbname); err != nil {
		t.Fatalf("identify database: %v", err)
	}
	if user != DBUser || dbname != DBName {
		t.Fatalf("expected %s/%s, got %s/%s", DBUser, DBName, user, dbname)
	}
}
