package repositories

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func TestClusterUpgradeJournalRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed cluster upgrade dependencies: %v", err)
		}
	}
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewClusterRepo(db, vault)
	ctx := context.Background()
	cluster := &models.Cluster{ID: "cluster-one", ProjectID: "project-one", OrganizationID: "org-one", Name: "Cluster", Version: "v1.34.1+k3s1", Nodes: []models.ClusterNode{{ServerID: "node-one"}}, JoinToken: "fixture-cluster-secret"}
	if err := repo.Save(ctx, cluster, 0); err != nil {
		t.Fatalf("save cluster: %v", err)
	}
	plan := &models.ClusterPlan{OperationID: "op-one", PreviousVersion: "v1.34.1+k3s1", Cluster: *cluster, Action: "upgrade", Installer: "installer"}
	if err := repo.BeginUpgrade(ctx, plan); err != nil {
		t.Fatalf("begin upgrade: %v", err)
	}
	pending, err := repo.PendingUpgrades(ctx)
	if err != nil {
		t.Fatalf("list pending upgrades: %v", err)
	}
	if len(pending) != 1 || pending[0].Cluster.ID != "cluster-one" || pending[0].OperationID != "op-one" {
		t.Fatalf("unexpected pending upgrades: %+v", pending)
	}
	loaded, err := repo.Get(ctx, "cluster-one")
	if err != nil {
		t.Fatalf("get cluster: %v", err)
	}
	if err := repo.Save(ctx, loaded, loaded.Revision); err == nil {
		t.Fatal("expected save to fail while an upgrade journal exists")
	}
	if err := repo.FinishUpgrade(ctx, "cluster-one", "v1.35.0+k3s1", "READY", ""); err != nil {
		t.Fatalf("finish upgrade: %v", err)
	}
	upgraded, err := repo.Get(ctx, "cluster-one")
	if err != nil {
		t.Fatalf("get upgraded cluster: %v", err)
	}
	if upgraded.Version != "v1.35.0+k3s1" || upgraded.Status != "READY" || upgraded.Revision != loaded.Revision+1 {
		t.Fatalf("upgrade did not persist: %+v", upgraded)
	}
	pending, err = repo.PendingUpgrades(ctx)
	if err != nil {
		t.Fatalf("list pending upgrades after finish: %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("expected no pending upgrades, got %d", len(pending))
	}
	if err := repo.FinishUpgrade(ctx, "cluster-one", "v1.36.0+k3s1", "READY", ""); err == nil {
		t.Fatal("expected finish to fail without a journal")
	}
}
