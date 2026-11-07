package repositories

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
)

func TestEnvironmentRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
		`INSERT INTO app_services (id, project_id, environment_id, name) VALUES ('service-one', 'project-one', 'env-one', 'Service')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed project: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewEnvironmentRepo(db)
	env := &models.EnvironmentConfig{ID: "env-one", ProjectID: "project-one", Name: "production", IsDefault: true}
	if err := repo.Create(ctx, env); err != nil {
		t.Fatalf("create environment: %v", err)
	}
	if env.CreatedAt.IsZero() || env.UpdatedAt.IsZero() {
		t.Fatal("expected timestamps to roundtrip")
	}
	got, err := repo.Get(ctx, "env-one")
	if err != nil {
		t.Fatalf("get environment: %v", err)
	}
	if got.Name != "production" || !got.IsDefault || got.ProjectID != "project-one" {
		t.Fatalf("unexpected environment row: %+v", got)
	}
	if !got.CreatedAt.Truncate(time.Second).Equal(env.CreatedAt.Truncate(time.Second)) {
		t.Fatal("created_at did not roundtrip")
	}
	envs, err := repo.ListByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list environments: %v", err)
	}
	if len(envs) != 1 || envs[0].ID != "env-one" {
		t.Fatalf("expected 1 environment, got %d", len(envs))
	}
	domains := NewDomainRepo(db)
	domain := &models.DomainConfig{ServiceID: "service-one", DomainName: "example.com", SSLCertStatus: "pending", PathPrefix: "/"}
	if err := domains.Create(ctx, domain); err != nil {
		t.Fatalf("create domain: %v", err)
	}
	if domain.ID == "" {
		t.Fatal("expected generated domain id")
	}
	stored, err := domains.GetByID(ctx, domain.ID)
	if err != nil {
		t.Fatalf("get domain: %v", err)
	}
	if stored == nil || stored.DomainName != "example.com" || stored.DNSProvisionStatus != "pending" {
		t.Fatalf("unexpected domain row: %+v", stored)
	}
	byService, err := domains.ListByService(ctx, "service-one")
	if err != nil {
		t.Fatalf("list domains by service: %v", err)
	}
	if len(byService) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(byService))
	}
	every, err := domains.ListAll(ctx)
	if err != nil {
		t.Fatalf("list all domains: %v", err)
	}
	if len(every) != 1 {
		t.Fatalf("expected 1 domain overall, got %d", len(every))
	}
	if err := domains.UpdateDNSProvisionStatus(ctx, domain.ID, "provisioned", "cloudflare", "203.0.113.7"); err != nil {
		t.Fatalf("update dns status: %v", err)
	}
	provisioned, err := domains.GetByID(ctx, domain.ID)
	if err != nil || provisioned == nil || provisioned.DNSProvisionStatus != "provisioned" || provisioned.DNSProvisionedIP != "203.0.113.7" {
		t.Fatalf("dns status did not persist: %+v %v", provisioned, err)
	}
	if err := domains.Delete(ctx, domain.ID); err != nil {
		t.Fatalf("delete domain: %v", err)
	}
	gone, err := domains.GetByID(ctx, domain.ID)
	if err != nil || gone != nil {
		t.Fatalf("expected deleted domain to be gone, got %+v %v", gone, err)
	}
	if err := repo.Delete(ctx, "env-one"); err != nil {
		t.Fatalf("delete environment: %v", err)
	}
	if _, err := repo.Get(ctx, "env-one"); err == nil {
		t.Fatal("expected deleted environment lookup to fail")
	}
}
