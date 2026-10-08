package integration_test

import (
	"context"
	"testing"

	"github.com/docker/docker/client"

	codedockhttp "codedock.run/codedock/internal/http"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/testdb"
	"codedock.run/codedock/internal/utils"
)

func TestCodedockBackendInitialization(t *testing.T) {
	t.Setenv("CODEDOCK_JWT_SECRET", "testsecret")

	dataDir := t.TempDir()
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	db := testdb.Open(t)

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	server, err := codedockhttp.NewServer(db, vlt, nil, nil, dockerClient, dataDir)
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
