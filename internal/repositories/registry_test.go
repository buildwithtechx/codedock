package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestRegistryRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed registry dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewRegistryRepository(db)
	registry := &models.Registry{ProjectID: "project-one", Name: "Docker Hub", RegistryURL: "https://index.docker.io", Username: "octo", PasswordToken: "token"}
	if err := repo.Create(ctx, registry); err != nil {
		t.Fatalf("create registry: %v", err)
	}
	if registry.ID == "" || registry.CreatedAt.IsZero() {
		t.Fatal("expected generated id and timestamps")
	}
	loaded, err := repo.Get(ctx, registry.ID)
	if err != nil {
		t.Fatalf("get registry: %v", err)
	}
	if loaded == nil || loaded.Name != "Docker Hub" || loaded.ProjectID != "project-one" {
		t.Fatalf("unexpected registry row: %+v", loaded)
	}
	listed, err := repo.ListByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list registries: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != registry.ID {
		t.Fatalf("unexpected registry list: %+v", listed)
	}
	if err := repo.Delete(ctx, registry.ID); err != nil {
		t.Fatalf("delete registry: %v", err)
	}
	deleted, err := repo.Get(ctx, registry.ID)
	if err != nil {
		t.Fatalf("get deleted registry: %v", err)
	}
	if deleted != nil {
		t.Fatalf("expected deleted registry to be gone, got %+v", deleted)
	}
}
