package deploy

import (
	"codedock.run/codedock/internal/models"
	"context"
	"github.com/docker/docker/client"
	"strings"
	"testing"
)

func TestVolumeGateFencesDatabaseStartup(t *testing.T) {
	t.Setenv("CODEDOCK_DRY_RUN", "false")
	gate := NewVolumeGate()
	release, err := gate.AcquireVolume("codedock-db-data-db")
	if err != nil {
		t.Fatal(err)
	}
	cli, err := client.NewClientWithOpts(client.WithHost("http://127.0.0.1:1"), client.WithVersion("1.47"))
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	deployer := NewDatabaseDeployer(cli, nil)
	deployer.SetVolumeOperations(gate)
	_, err = deployer.SpinUp(context.Background(), &models.Database{ID: "db", Name: "database"})
	if err == nil || !strings.Contains(err.Error(), "busy") {
		t.Fatalf("database startup escaped restore fence: %v", err)
	}
	if err := deployer.ImportData(context.Background(), &models.Database{ID: "db", Name: "database"}, "source"); err == nil || !strings.Contains(err.Error(), "busy") {
		t.Fatalf("database import escaped restore fence: %v", err)
	}
	other, err := gate.AcquireVolume("other")
	if err != nil {
		t.Fatal(err)
	}
	other()
	release()
	release()
	next, err := gate.AcquireVolume("codedock-db-data-db")
	if err != nil {
		t.Fatal(err)
	}
	next()
}
