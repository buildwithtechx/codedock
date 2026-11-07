package repositories

import (
	"codedock/internal/models"
	"codedock/internal/utils"
	"context"
	"strings"
	"testing"
)

func TestRuntimeRecoveryJournalEncryptionAndRevisionConflicts(t *testing.T) {
	db := openPGTestDB(t)
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

func TestDesiredWorkloadCommitIsAtomicAndEncrypted(t *testing.T) {
	db := openPGTestDB(t)
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
	if err := repo.Save(ctx, &models.ServiceRuntime{ServiceID: "service", ProjectID: "project", Target: models.RuntimeTarget{Kind: "kubernetes"}}, 0); err != nil {
		t.Fatal(err)
	}
	if err := repo.Begin(ctx, "service", 1, "journal"); err != nil {
		t.Fatal(err)
	}
	desired := &models.DesiredRuntime{Revision: 1, Manifest: "private-registry-credentials"}
	if err := repo.CommitDesired(ctx, "service", 2, desired); err == nil {
		t.Fatal("stale desired state committed")
	}
	pending, err := repo.Pending(ctx)
	if err != nil || len(pending) != 1 {
		t.Fatal("stale commit lost recovery journal", err)
	}
	if err := repo.CommitDesired(ctx, "service", 1, desired); err != nil {
		t.Fatal(err)
	}
	var encrypted string
	if err := db.QueryRow(`SELECT encrypted_workload FROM runtime_desired WHERE service_id='service'`).Scan(&encrypted); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(encrypted, desired.Manifest) {
		t.Fatal("desired secrets stored in plaintext")
	}
	loaded, err := repo.Desired(ctx, "service")
	if err != nil || loaded.Manifest != desired.Manifest {
		t.Fatal("desired state cannot reload", err)
	}
	pending, err = repo.Pending(ctx)
	if err != nil || len(pending) != 0 {
		t.Fatal("committed journal retained", err)
	}
}

func TestRuntimeDesiredIDsAndSyncKinds(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{`INSERT INTO organizations(id,name) VALUES('org','Org')`, `INSERT INTO project_apps(id,organization_id,name,slug) VALUES('app','org','App','app')`, `INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','app','org','Project','project')`, `INSERT INTO app_services(id,project_id,name) VALUES('service-one','project','One')`, `INSERT INTO app_services(id,project_id,name) VALUES('service-two','project','Two')`} {
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
	if err := repo.Save(ctx, &models.ServiceRuntime{ServiceID: "service-one", ProjectID: "project", Target: models.RuntimeTarget{Kind: "docker"}}, 0); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, &models.ServiceRuntime{ServiceID: "service-two", ProjectID: "project", Target: models.RuntimeTarget{Kind: "kubernetes"}}, 0); err != nil {
		t.Fatal(err)
	}
	ids, err := repo.DesiredIDs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Fatal("desired ids listed before any commit", ids)
	}
	if err := repo.Begin(ctx, "service-one", 1, "journal"); err != nil {
		t.Fatal(err)
	}
	if err := repo.CommitDesired(ctx, "service-one", 1, &models.DesiredRuntime{Revision: 1, Manifest: "manifest"}); err != nil {
		t.Fatal(err)
	}
	ids, err = repo.DesiredIDs(ctx)
	if err != nil || len(ids) != 1 || ids[0] != "service-one" {
		t.Fatal("committed desired workload missing", ids, err)
	}
	if _, err := db.Exec(`UPDATE service_runtimes SET runtime_kind='bare' WHERE service_id='service-two'`); err != nil {
		t.Fatal(err)
	}
	if err := repo.SyncKinds(ctx); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.Get(ctx, "service-two")
	if err != nil || loaded.RuntimeKind != "kubernetes" {
		t.Fatal("stale runtime kind was not repaired", err)
	}
	untouched, err := repo.Get(ctx, "service-one")
	if err != nil || untouched.RuntimeKind != "docker" {
		t.Fatal("matching runtime kind was rewritten", err)
	}
}
