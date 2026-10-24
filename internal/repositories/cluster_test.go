package repositories

import (
	"codedock/internal/models"
	"codedock/internal/utils"
	"context"
	"strings"
	"testing"
)

func TestClusterSecretsRevisionsAndInterruptedRecovery(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{`INSERT INTO organizations(id,name) VALUES('org','Org')`, `INSERT INTO project_apps(id,organization_id,name,slug) VALUES('app','org','App','app')`, `INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','app','org','Project','project')`} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewClusterRepo(db, vault)
	ctx := context.Background()
	cluster := &models.Cluster{ID: "cluster", ProjectID: "project", OrganizationID: "org", Name: "Cluster", Version: "v1.34.1+k3s1", Nodes: []models.ClusterNode{{ServerID: "node"}}, JoinToken: "fixture-cluster-secret"}
	if err := repo.Save(ctx, cluster, 0); err != nil {
		t.Fatal(err)
	}
	var encrypted string
	if err := db.QueryRow(`SELECT encrypted_token FROM clusters WHERE id='cluster'`).Scan(&encrypted); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(encrypted, cluster.JoinToken) {
		t.Fatal("cluster credentials stored in plaintext")
	}
	loaded, err := repo.Get(ctx, cluster.ID)
	if err != nil || loaded.JoinToken != cluster.JoinToken || loaded.Revision != 1 {
		t.Fatal("cluster cannot be reloaded", err)
	}
	if err := repo.Save(ctx, loaded, 1); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, loaded, 1); err == nil {
		t.Fatal("stale desired revision accepted")
	}
	if err := repo.Observe(ctx, cluster.ID, "APPLYING", ""); err != nil {
		t.Fatal(err)
	}
	if err := repo.Recover(ctx); err != nil {
		t.Fatal(err)
	}
	loaded, err = repo.Get(ctx, cluster.ID)
	if err != nil || loaded.Status != "INTERRUPTED" || loaded.Error == "" {
		t.Fatal("interrupted cluster reported ready", err)
	}
	listed, err := repo.List(ctx, "project")
	if err != nil || len(listed) != 1 || listed[0].JoinToken != "" {
		t.Fatal("cluster list leaked credentials", err)
	}
}
