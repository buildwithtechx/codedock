package integration_test

import (
	"context"
	"testing"

	"github.com/docker/docker/client"

	"codedock/internal/app"
	cdkdb "codedock/internal/db"
	"codedock/internal/testdb"
	"codedock/internal/utils"
)

func TestCodedockBackendInitialization(t *testing.T) {
	t.Setenv("CODEDOCK_JWT_SECRET", "testsecret")

	dataDir := t.TempDir()
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	db := testdb.Open(t)

	if err := cdkdb.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	server, err := app.New(app.Config{
		DB:           db,
		Vault:        vlt,
		DockerClient: dockerClient,
		DataDir:      dataDir,
	})
	if err != nil {
		t.Fatalf("failed to initialize server: %v", err)
	}

	if server == nil {
		t.Fatalf("expected server to be non-nil")
	}

	if err := db.PingContext(context.Background()); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}
}
