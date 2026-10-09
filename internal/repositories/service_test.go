package repositories

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
)

func TestAppServiceRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO environments (id, project_id, name, created_at, updated_at) VALUES ('env-one', 'project-one', 'production', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed service dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewAppServiceRepo(db)
	svc := &models.AppService{
		ProjectID:        "project-one",
		EnvironmentID:    "env-one",
		Name:             "web",
		RepositoryURL:    "https://example.com/app.git",
		EnablePRPreviews: true,
		MaintenanceMode:  true,
	}
	if err := repo.Create(ctx, svc); err != nil {
		t.Fatalf("create service: %v", err)
	}
	if svc.ID == "" || svc.Status != "building" || svc.InternalPort != 3000 {
		t.Fatalf("expected generated id with defaults, got %+v", svc)
	}
	got, err := repo.GetByID(ctx, svc.ID)
	if err != nil {
		t.Fatalf("get service: %v", err)
	}
	if got.Name != "web" || got.AppID != "app-one" || !got.EnablePRPreviews || !got.MaintenanceMode {
		t.Fatalf("unexpected service row: %+v", got)
	}
	if !got.CreatedAt.Truncate(time.Second).Equal(svc.CreatedAt.Truncate(time.Second)) {
		t.Fatal("created_at did not roundtrip")
	}
	byEnv, err := repo.ListByEnvironment(ctx, "env-one")
	if err != nil {
		t.Fatalf("list by environment: %v", err)
	}
	if len(byEnv) != 1 || byEnv[0].ID != svc.ID {
		t.Fatalf("expected 1 service by environment, got %d", len(byEnv))
	}
	byProject, err := repo.ListByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list by project: %v", err)
	}
	if len(byProject) != 1 {
		t.Fatalf("expected 1 service by project, got %d", len(byProject))
	}
	byOrg, err := repo.ListByOrganization(ctx, "org-one")
	if err != nil {
		t.Fatalf("list by organization: %v", err)
	}
	if len(byOrg) != 1 {
		t.Fatalf("expected 1 service by organization, got %d", len(byOrg))
	}
	every, err := repo.ListAll(ctx)
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(every) != 1 {
		t.Fatalf("expected 1 service overall, got %d", len(every))
	}
	got.Name = "web-renamed"
	got.MaintenanceMode = false
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update service: %v", err)
	}
	updated, err := repo.GetByID(ctx, svc.ID)
	if err != nil {
		t.Fatalf("get updated service: %v", err)
	}
	if updated.Name != "web-renamed" || updated.MaintenanceMode || !updated.EnablePRPreviews {
		t.Fatalf("update did not persist: %+v", updated)
	}
	hook := &models.Webhook{ServiceID: svc.ID, URL: "https://hooks.example.com/x", EventTypes: []string{"deploy", "rollback"}, IncludePREnvironments: true}
	if err := repo.CreateWebhook(ctx, hook); err != nil {
		t.Fatalf("create webhook: %v", err)
	}
	hooks, err := repo.ListWebhooksByService(ctx, svc.ID)
	if err != nil {
		t.Fatalf("list webhooks: %v", err)
	}
	if len(hooks) != 1 || hooks[0].ID != hook.ID {
		t.Fatalf("expected 1 webhook, got %d", len(hooks))
	}
	if len(hooks[0].EventTypes) != 2 || !hooks[0].IncludePREnvironments || hooks[0].CreatedAt.IsZero() {
		t.Fatalf("webhook did not roundtrip: %+v", hooks[0])
	}
	if err := repo.DeleteWebhook(ctx, hook.ID, svc.ID); err != nil {
		t.Fatalf("delete webhook: %v", err)
	}
	if err := repo.DeleteWebhook(ctx, hook.ID, svc.ID); err == nil {
		t.Fatal("expected second webhook delete to fail")
	}
	drain := &models.LogDrain{ServiceID: svc.ID, ProjectID: "project-one", DrainType: models.LogDrainTypeWebhook, EndpointURL: "https://logs.example.com/x", AuthToken: "secret"}
	if err := repo.CreateLogDrain(ctx, drain); err != nil {
		t.Fatalf("create log drain: %v", err)
	}
	drains, err := repo.ListLogDrainsByService(ctx, svc.ID)
	if err != nil {
		t.Fatalf("list log drains: %v", err)
	}
	if len(drains) != 1 || drains[0].ID != drain.ID {
		t.Fatalf("expected 1 log drain, got %d", len(drains))
	}
	if drains[0].EndpointURL != "https://logs.example.com/x" || drains[0].AuthToken != "secret" || drains[0].CreatedAt.IsZero() || drains[0].UpdatedAt.IsZero() {
		t.Fatalf("log drain did not roundtrip: %+v", drains[0])
	}
	if err := repo.DeleteLogDrain(ctx, drain.ID, svc.ID); err != nil {
		t.Fatalf("delete log drain: %v", err)
	}
	if err := repo.DeleteLogDrain(ctx, drain.ID, svc.ID); err == nil {
		t.Fatal("expected second log drain delete to fail")
	}
	if _, err := db.Exec(`INSERT INTO topology_dependencies (environment_id, source, target) VALUES ('env-one', 'db-x', $1)`, "app-"+svc.ID); err != nil {
		t.Fatalf("seed dependency: %v", err)
	}
	deps, err := repo.DeploymentDependencies(ctx, updated)
	if err != nil {
		t.Fatalf("deployment dependencies: %v", err)
	}
	if len(deps) != 1 || deps[0] != "db-x" {
		t.Fatalf("unexpected dependencies: %v", deps)
	}
	if err := repo.Delete(ctx, svc.ID); err != nil {
		t.Fatalf("delete service: %v", err)
	}
	if _, err := repo.GetByID(ctx, svc.ID); err == nil {
		t.Fatal("expected deleted service lookup to fail")
	}
}
