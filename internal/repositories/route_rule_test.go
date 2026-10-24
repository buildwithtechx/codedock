package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestRouteRuleRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, name) VALUES ('service-one', 'project-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed route rule dependencies: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewRouteRuleRepository(db)
	rule := &models.RouteRule{ServiceID: "service-one", Name: "limit", Enabled: true, RuleType: models.RouteRuleTypeRateLimit, SpecJSON: `{"average":10}`}
	if err := repo.Create(ctx, rule); err != nil {
		t.Fatalf("create rule: %v", err)
	}
	if rule.ID == "" || rule.CreatedAt.IsZero() {
		t.Fatal("expected generated id and timestamps")
	}
	loaded, err := repo.GetByID(ctx, rule.ID)
	if err != nil {
		t.Fatalf("get rule: %v", err)
	}
	if !loaded.Enabled || loaded.RuleType != models.RouteRuleTypeRateLimit || loaded.SpecJSON != `{"average":10}` {
		t.Fatalf("unexpected rule row: %+v", loaded)
	}
	listed, err := repo.ListByService(ctx, "service-one")
	if err != nil {
		t.Fatalf("list rules: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != rule.ID {
		t.Fatalf("unexpected rule list: %+v", listed)
	}
	name := "limit-renamed"
	enabled := false
	spec := `{"average":5}`
	if err := repo.Update(ctx, rule.ID, &name, &enabled, &spec); err != nil {
		t.Fatalf("update rule: %v", err)
	}
	updated, err := repo.GetByID(ctx, rule.ID)
	if err != nil {
		t.Fatalf("get updated rule: %v", err)
	}
	if updated.Name != "limit-renamed" || updated.Enabled || updated.SpecJSON != spec {
		t.Fatalf("update did not persist: %+v", updated)
	}
	if err := repo.Delete(ctx, rule.ID); err != nil {
		t.Fatalf("delete rule: %v", err)
	}
	remaining, err := repo.ListByService(ctx, "service-one")
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("expected deleted rule to be gone, got %+v", remaining)
	}
}
