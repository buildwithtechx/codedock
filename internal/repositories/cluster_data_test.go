package repositories

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestClusterDatabaseSecretsAreEncryptedAndHiddenFromListing(t *testing.T) {
	db := openTestDB(t)
	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{`INSERT INTO organizations(id,name) VALUES('org','Org')`, `INSERT INTO project_apps(id,organization_id,name,slug) VALUES('app','org','App','app')`, `INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','app','org','Project','project')`, `INSERT INTO clusters(id,project_id,organization_id,name,version,nodes_json,encrypted_token,updated_at) VALUES('cluster','project','org','Cluster','v1.34.1+k3s1','[]','','now')`} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewClusterDataRepo(db, vault)
	ctx := context.Background()
	plan := &models.ClusterDataPlan{Record: models.ClusterData{ID: "database", ClusterID: "cluster", ProjectID: "project", Spec: models.ClusterDataSpec{ID: "database", Name: "Database", Engine: "redis", Image: "redis:7.4.2"}}, Password: "private-password", Manifest: "private-manifest"}
	if err := repo.Create(ctx, plan); err != nil {
		t.Fatal(err)
	}
	saved, err := repo.Get(ctx, plan.Record.ID)
	if err != nil || saved.Password != plan.Password || saved.Manifest != plan.Manifest {
		t.Fatal("encrypted config cannot reload", err)
	}
	if strings.Contains(saved.Record.Config, plan.Password) || strings.Contains(saved.Record.Config, plan.Manifest) {
		t.Fatal("plaintext secrets persisted")
	}
	records, err := repo.List(ctx, "cluster")
	if err != nil {
		t.Fatal(err)
	}
	public, err := json.Marshal(records)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(public), plan.Password) || strings.Contains(string(public), plan.Manifest) {
		t.Fatal("listing disclosed private database configuration")
	}
	if err := repo.Observe(ctx, plan.Record.ID, "APPLYING", ""); err != nil {
		t.Fatal(err)
	}
	if err := repo.Recover(ctx); err != nil {
		t.Fatal(err)
	}
	saved, err = repo.Get(ctx, plan.Record.ID)
	if err != nil || saved.Record.Status != "INTERRUPTED" {
		t.Fatal("interrupted creation appeared ready", err)
	}
}
