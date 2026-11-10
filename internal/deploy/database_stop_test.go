package deploy

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"codedock/internal/modules/databases"
)

type stoppedDatabaseStore struct {
	DatabaseDeployerStore
	status    databases.DatabaseStatus
	container string
}

func (s *stoppedDatabaseStore) GetDatabase(string) (*databases.Database, error) {
	return &databases.Database{ID: "db", Name: "database", Status: databases.DatabaseStatusRunning, ContainerID: "dry-run-container"}, nil
}
func (s *stoppedDatabaseStore) UpdateDatabaseStatus(_ string, status databases.DatabaseStatus, container string) error {
	s.status = status
	s.container = container
	return nil
}

func TestDatabaseStopDryRunPersistsStatus(t *testing.T) {
	if os.Getenv("CODEDOCK_TEST_STOP_DRY_RUN") != "1" {
		command := exec.Command(os.Args[0], "-test.run=^TestDatabaseStopDryRunPersistsStatus$")
		command.Env = append(os.Environ(), "DEPLOY_DRY_RUN=true", "CODEDOCK_TEST_STOP_DRY_RUN=1")
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("dry-run stop failed: %v: %s", err, output)
		}
		return
	}
	store := &stoppedDatabaseStore{}
	deployer := NewDatabaseDeployer(nil, store)
	deployer.SetVolumeOperations(NewVolumeGate())
	if err := deployer.Stop(context.Background(), "db"); err != nil {
		t.Fatal(err)
	}
	if store.status != databases.DatabaseStatusStopped || store.container != "dry-run-container" {
		t.Fatal("dry-run stop did not persist stopped status and preserve its container identity")
	}
}
