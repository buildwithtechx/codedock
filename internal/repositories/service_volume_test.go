package repositories

import (
	"context"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
)

func TestServiceVolumeRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, environment_id, name) VALUES ('service-one', 'project-one', 'env-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed volume dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewServiceVolumeRepo(db)
	now := time.Now().UTC().Truncate(time.Second)
	volume := &models.ServiceVolume{ID: "vol-one", ServiceID: "service-one", HostPath: "/data", ContainerPath: "/app/data", CreatedAt: now}
	if err := repo.Create(ctx, volume); err != nil {
		t.Fatalf("create volume: %v", err)
	}
	got, err := repo.GetByID(ctx, "vol-one")
	if err != nil {
		t.Fatalf("get volume: %v", err)
	}
	if got.ServiceID != "service-one" || got.HostPath != "/data" || got.ContainerPath != "/app/data" {
		t.Fatalf("unexpected volume row: %+v", got)
	}
	if !got.CreatedAt.Truncate(time.Second).Equal(now) {
		t.Fatalf("created_at did not roundtrip: %v", got.CreatedAt)
	}
	list, err := repo.ListByService(ctx, "service-one")
	if err != nil {
		t.Fatalf("list volumes: %v", err)
	}
	if len(list) != 1 || list[0].ID != "vol-one" {
		t.Fatalf("expected 1 volume, got %d", len(list))
	}
	if err := repo.Delete(ctx, "vol-one"); err != nil {
		t.Fatalf("delete volume: %v", err)
	}
	if _, err := repo.GetByID(ctx, "vol-one"); err == nil {
		t.Fatal("expected deleted volume lookup to fail")
	}
}
