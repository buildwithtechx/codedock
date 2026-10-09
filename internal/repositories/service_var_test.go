package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestServiceVarRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO environments (id, project_id, name, created_at, updated_at) VALUES ('env-one', 'project-one', 'production', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		`INSERT INTO app_services (id, project_id, environment_id, name) VALUES ('service-one', 'project-one', 'env-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed service var dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewServiceVarRepo(db)
	variable := &models.Variable{ServiceID: "service-one", EnvironmentID: "env-one", Key: "API_KEY", Value: "secret", IsSecret: true}
	if err := repo.Create(ctx, variable); err != nil {
		t.Fatalf("create var: %v", err)
	}
	if variable.ID == "" || variable.CreatedAt.IsZero() || variable.UpdatedAt.IsZero() {
		t.Fatal("expected generated id and timestamps")
	}
	got, err := repo.GetByID(ctx, variable.ID)
	if err != nil {
		t.Fatalf("get var: %v", err)
	}
	if got.Key != "API_KEY" || got.Value != "secret" || !got.IsSecret || got.EnvironmentID != "env-one" {
		t.Fatalf("unexpected var row: %+v", got)
	}
	upsert := &models.Variable{ServiceID: "service-one", EnvironmentID: "env-one", Key: "API_KEY", Value: "rotated", IsSecret: false}
	if err := repo.Create(ctx, upsert); err != nil {
		t.Fatalf("upsert var: %v", err)
	}
	list, err := repo.ListByService(ctx, "service-one")
	if err != nil {
		t.Fatalf("list vars: %v", err)
	}
	if len(list) != 1 || list[0].Value != "rotated" || list[0].IsSecret {
		t.Fatalf("upsert did not persist: %+v", list)
	}
	stored := list[0]
	stored.Key = "RENAMED_KEY"
	stored.Value = "final"
	stored.IsSecret = true
	if err := repo.Update(ctx, stored); err != nil {
		t.Fatalf("update var: %v", err)
	}
	updated, err := repo.GetByID(ctx, stored.ID)
	if err != nil {
		t.Fatalf("get updated var: %v", err)
	}
	if updated.Key != "RENAMED_KEY" || updated.Value != "final" || !updated.IsSecret {
		t.Fatalf("update did not persist: %+v", updated)
	}
	if err := repo.Delete(ctx, stored.ID); err != nil {
		t.Fatalf("delete var: %v", err)
	}
	if _, err := repo.GetByID(ctx, stored.ID); err == nil {
		t.Fatal("expected deleted var lookup to fail")
	}
}
