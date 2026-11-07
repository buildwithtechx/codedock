package repositories

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
)

func TestPRPreviewRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, environment_id, name) VALUES ('service-one', 'project-one', 'env-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed preview dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewPRPreviewRepository(db)
	now := time.Now().UTC().Truncate(time.Second)
	preview := &models.PRPreview{
		ID:            "pr-one",
		ServiceID:     "service-one",
		ProjectID:     "project-one",
		PRNumber:      42,
		Branch:        "feature",
		CommitHash:    "abc1234",
		Status:        models.PRPreviewStatusPending,
		PreviewDomain: "pr42.example.com",
		ContainerID:   "container-one",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := repo.Create(ctx, preview); err != nil {
		t.Fatalf("create preview: %v", err)
	}
	byApp, err := repo.GetByApp(ctx, "service-one")
	if err != nil {
		t.Fatalf("get by app: %v", err)
	}
	if len(byApp) != 1 || byApp[0].ID != "pr-one" {
		t.Fatalf("expected 1 preview, got %d", len(byApp))
	}
	if byApp[0].Branch != "feature" || byApp[0].Status != models.PRPreviewStatusPending {
		t.Fatalf("unexpected preview row: %+v", byApp[0])
	}
	if !byApp[0].CreatedAt.Truncate(time.Second).Equal(now) {
		t.Fatalf("created_at did not roundtrip: %v", byApp[0].CreatedAt)
	}
	byPR, err := repo.GetByAppAndPR(ctx, "service-one", 42)
	if err != nil {
		t.Fatalf("get by app and pr: %v", err)
	}
	if len(byPR) != 1 {
		t.Fatalf("expected 1 preview for PR 42, got %d", len(byPR))
	}
	other, err := repo.GetByAppAndPR(ctx, "service-one", 43)
	if err != nil {
		t.Fatalf("get by app and other pr: %v", err)
	}
	if len(other) != 0 {
		t.Fatalf("expected 0 previews for PR 43, got %d", len(other))
	}
	preview.Status = models.PRPreviewStatusReady
	preview.PreviewDomain = "ready.example.com"
	preview.UpdatedAt = now
	if err := repo.Update(ctx, preview); err != nil {
		t.Fatalf("update preview: %v", err)
	}
	updated, err := repo.GetByApp(ctx, "service-one")
	if err != nil || len(updated) != 1 || updated[0].Status != models.PRPreviewStatusReady || updated[0].PreviewDomain != "ready.example.com" {
		t.Fatalf("update did not persist: %+v %v", updated, err)
	}
	if err := repo.Delete(ctx, "pr-one"); err != nil {
		t.Fatalf("delete preview: %v", err)
	}
	empty, err := repo.GetByApp(ctx, "service-one")
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected 0 previews after delete, got %d", len(empty))
	}
	if err := repo.Delete(ctx, ""); err == nil {
		t.Fatal("expected empty id delete to fail")
	}
}
