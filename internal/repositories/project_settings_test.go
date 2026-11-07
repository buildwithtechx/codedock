package repositories

import (
	"context"
	"testing"
	"time"

	"codedock/internal/models"
)

func TestProjectSettingsTokenRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	for _, query := range []string{
		`INSERT INTO organizations (id, name) VALUES ('org-one', 'Organization')`,
		`INSERT INTO project_apps (id, organization_id, name, slug) VALUES ('app-one', 'org-one', 'App', 'app')`,
		`INSERT INTO projects (id, app_id, organization_id, name, slug) VALUES ('project-one', 'app-one', 'org-one', 'Project', 'project')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("seed project: %v", err)
		}
	}
	ctx := context.Background()
	repo := NewProjectSettingsRepo(db)
	now := time.Now().UTC().Truncate(time.Second)
	expires := now.Add(24 * time.Hour)
	token := &models.ProjectToken{
		ID:          "token-one",
		ProjectID:   "project-one",
		Name:        "ci",
		TokenPrefix: "cdk_",
		Scopes:      []string{"read", "write"},
		IPAllowlist: []string{"10.0.0.1"},
		ExpiresAt:   &expires,
		CreatedAt:   now,
	}
	if err := repo.CreateToken(ctx, token, "hash-secret"); err != nil {
		t.Fatalf("create token: %v", err)
	}
	list, err := repo.ListTokensByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list tokens: %v", err)
	}
	if len(list) != 1 || list[0].ID != "token-one" {
		t.Fatalf("expected 1 token, got %d", len(list))
	}
	stored := list[0]
	if stored.Name != "ci" || stored.TokenPrefix != "cdk_" {
		t.Fatalf("unexpected token row: %+v", stored)
	}
	if len(stored.Scopes) != 2 || stored.Scopes[0] != "read" || stored.Scopes[1] != "write" {
		t.Fatalf("scopes did not roundtrip: %v", stored.Scopes)
	}
	if len(stored.IPAllowlist) != 1 || stored.IPAllowlist[0] != "10.0.0.1" {
		t.Fatalf("ip allowlist did not roundtrip: %v", stored.IPAllowlist)
	}
	if stored.ExpiresAt == nil || !stored.ExpiresAt.Truncate(time.Second).Equal(expires) {
		t.Fatalf("expires_at did not roundtrip: %v", stored.ExpiresAt)
	}
	if !stored.CreatedAt.Truncate(time.Second).Equal(now) {
		t.Fatalf("created_at did not roundtrip: %v", stored.CreatedAt)
	}
	byHash, err := repo.GetTokenByHash(ctx, "hash-secret")
	if err != nil {
		t.Fatalf("get by hash: %v", err)
	}
	if byHash.ID != "token-one" {
		t.Fatal("hash lookup returned wrong token")
	}
	if err := repo.UpdateTokenLastUsed(ctx, "token-one"); err != nil {
		t.Fatalf("update last used: %v", err)
	}
	if err := repo.DeleteToken(ctx, "token-one", "project-one"); err != nil {
		t.Fatalf("delete token: %v", err)
	}
	list, err = repo.ListTokensByProject(ctx, "project-one")
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 tokens after delete, got %d", len(list))
	}
	if err := repo.DeleteToken(ctx, "token-one", "project-one"); err == nil {
		t.Fatal("expected second delete to fail")
	}
}
