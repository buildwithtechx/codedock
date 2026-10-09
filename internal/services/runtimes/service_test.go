package runtimes

import (
	"codedock/internal/engine/deploy"
	"codedock/internal/models"
	"testing"
)

func TestRuntimeLocksConflictWithClusterMaintenance(t *testing.T) {
	gate := deploy.NewVolumeGate()
	service := &Service{gate: gate}
	cluster := &models.Cluster{ID: "cluster", Nodes: []models.ClusterNode{{ServerID: "one"}, {ServerID: "two"}}}
	release, err := gate.AcquireVolume("server:two")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if _, err := service.lockCluster(cluster); err == nil {
		t.Fatal("deployment raced cluster maintenance")
	}
	for _, name := range []string{"cluster:cluster", "server:one"} {
		cleanup, err := gate.AcquireVolume(name)
		if err != nil {
			t.Fatal("partial lock leaked", err)
		}
		cleanup()
	}
}
