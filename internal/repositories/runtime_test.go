package repositories

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"strings"
	"testing"
)

func TestRuntimeRecoveryJournalEncryptionAndRevisionConflicts(t *testing.T) {
	db := openTestDB(t)
	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{`INSERT INTO organizations(id,name) VALUES('org','Org')`, `INSERT INTO project_apps(id,organization_id,name,slug) VALUES('app','org','App','app')`, `INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','app','org','Project','project')`, `INSERT INTO app_services(id,project_id,name) VALUES('service','project','Service')`} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRuntimeRepo(db, vault)
	ctx := context.Background()
	runtime := &models.ServiceRuntime{ServiceID: "service", ProjectID: "project", Target: models.RuntimeTarget{Kind: "kubernetes", ClusterID: "cluster"}}
	if err := repo.Save(ctx, runtime, 0); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.Get(ctx, "service")
	if err != nil || loaded.Revision != 1 || loaded.Target.ClusterID != "cluster" {
		t.Fatal("runtime cannot reload", err)
	}
	if err := repo.Begin(ctx, "service", 1, "private-recovery-secret"); err != nil {
		t.Fatal(err)
	}
	var raw string
	if err := db.QueryRow(`SELECT encrypted_journal FROM service_runtimes WHERE service_id='service'`).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(raw, "private-recovery-secret") {
		t.Fatal("journal stored in plaintext")
	}
	if err := repo.Save(ctx, runtime, 1); err == nil {
		t.Fatal("pending recovery target replaced")
	}
	if err := repo.Begin(ctx, "service", 1, "replacement"); err == nil {
		t.Fatal("active operation replaced")
	}
	loaded, err = repo.Get(ctx, "service")
	if err != nil || loaded.Journal != "private-recovery-secret" {
		t.Fatal("journal cannot reload", err)
	}
	if err := repo.Observe(ctx, "service", "FAILED", "", true); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, runtime, 1); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, runtime, 1); err == nil {
		t.Fatal("stale revision accepted")
	}
}
