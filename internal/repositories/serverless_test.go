package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestServerlessRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, name) VALUES ('service-one', 'project-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed serverless dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewServerlessRepository(db)
	saved, err := repo.SaveCode(ctx, "service-one", "node22", "export default {}")
	if err != nil {
		t.Fatalf("save code: %v", err)
	}
	if saved.ID == "" || saved.ServiceID != "service-one" || saved.Runtime != "node22" || saved.CodeContent != "export default {}" {
		t.Fatalf("unexpected saved code: %+v", saved)
	}
	if saved.CreatedAt.IsZero() || saved.UpdatedAt.IsZero() {
		t.Fatal("expected code timestamps to roundtrip")
	}
	loaded, err := repo.GetCodeByServiceID(ctx, "service-one")
	if err != nil {
		t.Fatalf("get code: %v", err)
	}
	if loaded.ID != saved.ID || loaded.CodeContent != "export default {}" {
		t.Fatalf("code did not roundtrip: %+v", loaded)
	}
	resaved, err := repo.SaveCode(ctx, "service-one", "python313", "def handler(): pass")
	if err != nil {
		t.Fatalf("resave code: %v", err)
	}
	if resaved.ID != saved.ID || resaved.Runtime != "python313" || resaved.CodeContent != "def handler(): pass" {
		t.Fatalf("code update did not persist: %+v", resaved)
	}
	if _, err := repo.GetCodeByServiceID(ctx, "missing-service"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows for missing service, got %v", err)
	}
}
