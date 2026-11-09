package clusters

import (
	"codedock/internal/engine/deploy"
	"testing"
)

func TestClusterLocksReleasePartialAcquisitions(t *testing.T) {
	gate := deploy.NewVolumeGate()
	service := &ClusterService{gate: gate}
	release, err := gate.AcquireVolume("server:second")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	cluster := &Cluster{ID: "cluster", Nodes: []ClusterNode{{ServerID: "first"}, {ServerID: "second"}}}
	if _, err := service.acquireTargets(cluster); err == nil {
		t.Fatal("conflicting server lock accepted")
	}
	cleanup, err := gate.AcquireVolume("cluster:cluster")
	if err != nil {
		t.Fatal("partial cluster lock leaked", err)
	}
	cleanup()
	cleanup, err = gate.AcquireVolume("server:first")
	if err != nil {
		t.Fatal("partial node lock leaked", err)
	}
	cleanup()
}
func TestClusterNodeChangesCannotDropOrReplaceExistingNodes(t *testing.T) {
	previous := &Cluster{Version: "v1.34.1+k3s1", Nodes: []ClusterNode{{ServerID: "first", PrivateIP: "192.168.3.4"}}}
	next := *previous
	next.Nodes = []ClusterNode{{ServerID: "replacement"}}
	if err := preserveExistingNodes(previous, &next, "join"); err == nil {
		t.Fatal("existing node replaced")
	}
	next = *previous
	next.Version = "v1.35.1+k3s1"
	if err := preserveExistingNodes(previous, &next, "join"); err == nil {
		t.Fatal("join also upgraded the cluster")
	}
	next = *previous
	next.Nodes = append(append([]ClusterNode{}, previous.Nodes...), ClusterNode{ServerID: "second"})
	if err := preserveExistingNodes(previous, &next, "join"); err != nil {
		t.Fatal(err)
	}
	next.Revision++
	if clusterSnapshot(previous) == clusterSnapshot(&next) {
		t.Fatal("changed cluster retained stale confirmation snapshot")
	}
}
