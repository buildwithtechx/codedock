package repositories

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
)

func TestProjectAppRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	if _, err := db.Exec(`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`); err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	ctx := context.Background()
	repo := NewProjectAppRepo(db)
	app := &models.ProjectApp{OrganizationID: "org-one", Name: "App", Slug: "app"}
	if err := repo.Create(ctx, app); err != nil {
		t.Fatalf("create app: %v", err)
	}
	if app.ID == "" || app.GitProvider != "github" {
		t.Fatalf("expected generated id and github default, got %+v", app)
	}
	if app.CreatedAt.IsZero() || app.UpdatedAt.IsZero() {
		t.Fatal("expected timestamps to roundtrip")
	}
	got, err := repo.GetByID(ctx, app.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Name != "App" || got.Slug != "app" {
		t.Fatalf("unexpected app row: %+v", got)
	}
	bySlug, err := repo.GetBySlug(ctx, "org-one", "app")
	if err != nil {
		t.Fatalf("get by slug: %v", err)
	}
	if bySlug.ID != app.ID {
		t.Fatal("slug lookup returned wrong app")
	}
	list, err := repo.ListByOrganization(ctx, "org-one")
	if err != nil {
		t.Fatalf("list apps: %v", err)
	}
	if len(list) != 1 || list[0].ID != app.ID {
		t.Fatalf("expected 1 app, got %d", len(list))
	}
	app.Name = "Renamed"
	if err := repo.Update(ctx, app); err != nil {
		t.Fatalf("update app: %v", err)
	}
	updated, err := repo.GetByID(ctx, app.ID)
	if err != nil {
		t.Fatalf("get updated app: %v", err)
	}
	if updated.Name != "Renamed" {
		t.Fatalf("update did not persist: %+v", updated)
	}
	if err := repo.Delete(ctx, app.ID); err != nil {
		t.Fatalf("delete app: %v", err)
	}
	if _, err := repo.GetByID(ctx, app.ID); err == nil {
		t.Fatal("expected deleted app lookup to fail")
	}
	list, err = repo.ListByOrganization(ctx, "org-one")
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 apps after delete, got %d", len(list))
	}
}
