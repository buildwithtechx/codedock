package repositories

import (
	"codedock/internal/models"
	"context"
	"testing"
)

func TestPersonalTokenLookupAndRevocation(t *testing.T) {
	db := openPGTestDB(t)
	if _, err := db.Exec(`INSERT INTO users(id,email,name,password_hash) VALUES('user','user@example.com','User','hash')`); err != nil {
		t.Fatal(err)
	}
	repo := NewUserRepo(db)
	ctx := context.Background()
	pat := &models.PersonalAccessToken{ID: "pat", UserID: "user", Name: "automation", TokenHash: "digest", Prefix: "vpt_1234", AccessLevel: "read", ProjectScope: "all"}
	if err := repo.CreatePAT(ctx, pat); err != nil {
		t.Fatal(err)
	}
	found, err := repo.GetPATByHash(ctx, "digest")
	if err != nil || found == nil || found.UserID != "user" || found.ExpiresAt != nil || found.AccessLevel != "read" || found.ProjectScope != "all" {
		t.Fatalf("lookup failed: %v", err)
	}
	if err := repo.DeletePAT(ctx, "pat", "user"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetPATByHash(ctx, "digest"); err == nil {
		t.Fatal("revoked token remained usable")
	}
}

func TestPersonalTokenResolvesNestedProjectResources(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations(id,name) VALUES('org','Org')`,
		`INSERT INTO project_apps(id,organization_id,name,slug) VALUES('app','org','App','app')`,
		`INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','app','org','Project','project')`,
		`INSERT INTO app_services(id,project_id,name) VALUES('service','project','Service')`,
		`INSERT INTO scheduled_tasks(id,service_id,name,schedule,command) VALUES('task','service','Task','manual','true')`,
		`INSERT INTO service_webhooks(id,service_id,url) VALUES('hook','service','https://example.com')`,
		`INSERT INTO registries(id,project_id,name,registry_url) VALUES('registry','project','Registry','https://example.com')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewPersonalTokenResources(db)
	for kind, id := range map[string]string{"scheduled-tasks": "task", "webhooks": "hook", "registries": "registry"} {
		project, err := repo.ProjectForResource(context.Background(), kind, id)
		if err != nil || project != "project" {
			t.Fatalf("%s: %s %v", kind, project, err)
		}
	}
}
