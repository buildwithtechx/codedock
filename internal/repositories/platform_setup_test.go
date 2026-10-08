package repositories

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"strings"
	"testing"
)

func TestStackEncryptionRevisionAndRestartRecovery(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations(id,name) VALUES('org','Org')`,
		`INSERT INTO project_apps(id,organization_id,name,slug) VALUES('application','org','Application','application')`,
		`INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','application','org','Project','project')`,
		`INSERT INTO environments(id,project_id,name,created_at,updated_at) VALUES('environment','project','Production',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	vault, err := utils.NewVault(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := NewComposeStackRepo(db, vault)
	ctx := context.Background()
	stack := &models.ComposeStack{ID: "stack", ProjectID: "project", EnvironmentID: "environment", Name: "Stack", Config: `{"services":{"web":{"image":"nginx","environment":{"PASSWORD":"fixture-secret"}}}}`}
	if err := repo.Save(ctx, stack, 0); err != nil {
		t.Fatal(err)
	}
	var stored string
	if err := db.QueryRow(`SELECT encrypted_config FROM compose_stacks WHERE id='stack'`).Scan(&stored); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stored, "fixture-secret") {
		t.Fatal("stack secrets were persisted as plaintext")
	}
	loaded, err := repo.Get(ctx, "project", "stack")
	if err != nil || loaded.Config != stack.Config || loaded.Revision != 1 {
		t.Fatalf("stack could not be restored: %v", err)
	}
	if err := repo.Save(ctx, stack, 0); err == nil {
		t.Fatal("duplicate creation overwrote the saved stack")
	}
	if err := repo.Save(ctx, stack, 1); err != nil {
		t.Fatal(err)
	}
	if err := repo.Claim(ctx, "project", "stack", 1); err == nil {
		t.Fatal("deployment claimed a stale revision")
	}
	if err := repo.Claim(ctx, "project", "stack", 2); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, stack, 2); err == nil {
		t.Fatal("configuration changed during activation")
	}
	if err := repo.Recover(ctx); err != nil {
		t.Fatal(err)
	}
	loaded, err = repo.Get(ctx, "project", "stack")
	if err != nil || loaded.Status != "INTERRUPTED" || loaded.Config != stack.Config {
		t.Fatal("interrupted stack did not preserve its reviewed configuration")
	}
}

func TestTopologyRejectsStaleCrossEnvironmentAndCyclicEdits(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations(id,name) VALUES('org','Org')`,
		`INSERT INTO project_apps(id,organization_id,name,slug) VALUES('application','org','Application','application')`,
		`INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','application','org','Project','project')`,
		`INSERT INTO environments(id,project_id,name,created_at,updated_at) VALUES('environment','project','Production',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP),('other','project','Preview',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
		`INSERT INTO app_services(id,project_id,environment_id,name) VALUES('a','project','environment','A'),('b','project','environment','B'),('outside','project','other','Outside')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	ctx := context.Background()
	repo := NewCanvasRepo(db, NewEnvironmentRepo(db))
	canvas, err := repo.GetEnvironmentCanvas(ctx, "environment")
	if err != nil {
		t.Fatal(err)
	}
	request := models.TopologyApplyRequest{Revision: canvas.Revision, Dependencies: []models.CanvasEdge{{Source: "app-a", Target: "app-b", Kind: "dependency"}}}
	if err := repo.ApplyTopology(ctx, "environment", request); err != nil {
		t.Fatal(err)
	}
	if err := repo.ApplyTopology(ctx, "environment", request); err == nil {
		t.Fatal("stale edit overwrote changed topology")
	}
	canvas, err = repo.GetEnvironmentCanvas(ctx, "environment")
	if err != nil {
		t.Fatal(err)
	}
	request.Revision = canvas.Revision
	request.Dependencies = append(request.Dependencies, models.CanvasEdge{Source: "app-b", Target: "app-a", Kind: "dependency"})
	if err := repo.ApplyTopology(ctx, "environment", request); err == nil {
		t.Fatal("cycle accepted")
	}
	request.Dependencies = []models.CanvasEdge{{Source: "app-outside", Target: "app-a", Kind: "dependency"}}
	if err := repo.ApplyTopology(ctx, "environment", request); err == nil {
		t.Fatal("cross-environment dependency accepted")
	}
	variable := &models.Variable{ServiceID: "b", EnvironmentID: "environment", Key: "URL", Value: "initial", ExpectedTopologyRevision: canvas.Revision}
	variables := NewServiceVarRepo(db)
	if err := variables.Create(ctx, variable); err != nil {
		t.Fatal(err)
	}
	variable.Value = "overwrite"
	if err := variables.Create(ctx, variable); err == nil {
		t.Fatal("stale binding review overwrote a newer value")
	}
	saved, err := variables.ListByService(ctx, "b")
	if err != nil || len(saved) != 1 || saved[0].Value != "initial" {
		t.Fatal("failed binding edit changed persisted data")
	}
}
