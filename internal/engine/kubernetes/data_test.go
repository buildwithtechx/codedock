package kubernetes

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func dataFixture() *models.ClusterDataPlan {
	return &models.ClusterDataPlan{Record: models.ClusterData{ID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", ProjectID: "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", Spec: models.ClusterDataSpec{ID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", EnvironmentID: "cccccccc-cccc-4ccc-8ccc-cccccccccccc", Name: "database", Engine: "postgres", Image: "ghcr.io/cloudnative-pg/postgresql:17.6", Instances: 3, StorageGiB: 10, Synchronous: true}}, Password: "generated"}
}
func TestDataModesRequireDistinctNodesAndExactImages(t *testing.T) {
	plan := dataFixture()
	if err := ValidateDataSpec(plan.Record.Spec, 3); err != nil {
		t.Fatal(err)
	}
	plan.Record.Spec.Synchronous = false
	if err := ValidateDataSpec(plan.Record.Spec, 3); err == nil {
		t.Fatal("replicated PostgreSQL accepted without synchronous durability")
	}
	plan.Record.Spec.Synchronous = true
	if err := ValidateDataSpec(plan.Record.Spec, 2); err == nil {
		t.Fatal("replication accepted without enough nodes")
	}
	plan.Record.Spec.Image = "postgres:latest"
	if err := ValidateDataSpec(plan.Record.Spec, 3); err == nil {
		t.Fatal("unpinned PostgreSQL image accepted")
	}
	plan.Record.Spec.Engine = "redis"
	plan.Record.Spec.Image = "redis:7.4.2"
	plan.Record.Spec.S3DestinationID = "destination"
	if err := ValidateDataSpec(plan.Record.Spec, 3); err == nil {
		t.Fatal("Redis advertised PostgreSQL backup compatibility")
	}
}
func TestPostgresRecoveryUsesSeparatePrefixAndNewCredentials(t *testing.T) {
	plan := dataFixture()
	plan.RestoreTime = "2026-10-06T18:00:00Z"
	destination := &models.S3Destination{Bucket: "backups", PathPrefix: "prefix", AccessKeyID: "key", SecretAccessKey: "secret"}
	manifest, err := DataManifest(plan, destination, destination, "db-source")
	if err != nil {
		t.Fatal(err)
	}
	var list struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(manifest), &list); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range list.Items {
		if item["kind"] != "Cluster" {
			continue
		}
		spec := item["spec"].(map[string]any)
		recovery := spec["bootstrap"].(map[string]any)["recovery"].(map[string]any)
		if recovery["owner"] != "app" || recovery["secret"] == nil {
			t.Fatal("recovery did not reset generated credentials")
		}
		found = true
	}
	if !found || !strings.Contains(manifest, "db-source") || !strings.Contains(manifest, "db-"+plan.Record.ID) {
		t.Fatal("source and destination backup prefixes not separated")
	}
}
func TestDatabaseApplyRefusesForeignResources(t *testing.T) {
	desired := map[string]any{"metadata": map[string]any{"name": "database", "labels": map[string]any{"codedock.run/project": "owner", "codedock.run/database": "database"}}}
	actual := map[string]any{"metadata": map[string]any{"labels": map[string]any{"codedock.run/project": "foreign", "codedock.run/database": "database"}}}
	if err := matchingDataOwner(actual, desired); err == nil {
		t.Fatal("foreign database resource adopted")
	}
}

type foreignDataCommands struct{ mutated bool }

func (r *foreignDataCommands) Kubectl(ctx context.Context, node models.ClusterNode, args []string, input string) (string, error) {
	if input != "" {
		r.mutated = true
	}
	return `{"metadata":{"uid":"foreign-uid","labels":{"codedock.run/project":"foreign","codedock.run/database":"database"}}}`, nil
}
func TestDatabaseOwnershipFailurePreventsSubmission(t *testing.T) {
	commands := &foreignDataCommands{}
	runtime := NewWorkloadRuntime(commands)
	manifest := `{"items":[{"apiVersion":"v1","kind":"Secret","metadata":{"name":"database","labels":{"codedock.run/project":"owner","codedock.run/database":"database"}}}]}`
	if err := runtime.ApplyOwnedData(context.Background(), models.ClusterNode{}, manifest, func(string, string) error { return nil }); err == nil || commands.mutated {
		t.Fatal("foreign resource was submitted for mutation")
	}
}

func TestDatabaseObservationRejectsMissingReadyInstances(t *testing.T) {
	commands := &foreignDataCommands{}
	runtime := NewWorkloadRuntime(commands)
	record := dataFixture().Record
	if err := runtime.DataObserved(context.Background(), models.ClusterNode{}, &record); err == nil {
		t.Fatal("foreign or incomplete runtime appeared healthy")
	}
}
func TestRedisShardingBuildsPerShardStatefulSets(t *testing.T) {
	plan := dataFixture()
	plan.Record.Spec.Engine = "redis"
	plan.Record.Spec.Image = "redis:7.4.2"
	plan.Record.Spec.Instances = 1
	plan.Record.Spec.Shards = 3
	plan.Record.Spec.Synchronous = false
	if err := ValidateDataSpec(plan.Record.Spec, 3); err != nil {
		t.Fatal(err)
	}
	manifest, err := DataManifest(plan, nil, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	var list struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(manifest), &list); err != nil {
		t.Fatal(err)
	}
	statefuls, jobs := 0, 0
	for _, item := range list.Items {
		if item["kind"] == "StatefulSet" {
			statefuls++
		}
		if item["kind"] == "Job" {
			jobs++
		}
	}
	if statefuls != 3 || jobs != 1 {
		t.Fatal("sharded Redis did not declare per-shard state and bootstrap")
	}
	plan.Record.Spec.Shards = 9
	if err := ValidateDataSpec(plan.Record.Spec, 9); err == nil {
		t.Fatal("excessive Redis shards accepted")
	}
	plan.Record.Spec.Shards = 0
	plan.Record.Spec.Engine = "postgres"
	if err := ValidateDataSpec(plan.Record.Spec, 3); err == nil {
		t.Fatal("PostgreSQL sharding accepted")
	}
}
