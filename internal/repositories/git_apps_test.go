package repositories

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
)

func TestGithubAppRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewGitAppRepo(db, gitTestVault{})
	app := &models.GithubApp{ID: "app-one", Name: "CI", AppID: "123", InstallationID: "456", ClientID: "cid", ClientSecret: "secret", WebhookSecret: "hook", PrivateKey: "key", IsPublic: true}
	if err := repo.SaveGithubApp(ctx, app); err != nil {
		t.Fatalf("save app: %v", err)
	}
	loaded, err := repo.GetGithubApp(ctx, "app-one")
	if err != nil {
		t.Fatalf("get app: %v", err)
	}
	if loaded.ClientSecret != "secret" || loaded.WebhookSecret != "hook" || loaded.PrivateKey != "key" || !loaded.IsPublic {
		t.Fatalf("unexpected app row: %+v", loaded)
	}
	if loaded.CreatedAt.IsZero() || loaded.UpdatedAt.IsZero() {
		t.Fatal("expected timestamps to roundtrip")
	}
	listed, err := repo.ListGithubApps(ctx)
	if err != nil {
		t.Fatalf("list apps: %v", err)
	}
	if len(listed) != 1 || listed[0].ClientSecret != "********" || listed[0].WebhookSecret != "********" || listed[0].PrivateKey != "********" {
		t.Fatalf("expected masked list, got %+v", listed)
	}
	app.Name = "CI Renamed"
	app.ClientSecret = "********"
	app.WebhookSecret = ""
	app.PrivateKey = "********"
	if err := repo.SaveGithubApp(ctx, app); err != nil {
		t.Fatalf("resave app: %v", err)
	}
	updated, err := repo.GetGithubApp(ctx, "app-one")
	if err != nil {
		t.Fatalf("get updated app: %v", err)
	}
	if updated.Name != "CI Renamed" || updated.ClientSecret != "secret" || updated.WebhookSecret != "hook" || updated.PrivateKey != "key" {
		t.Fatalf("masked secrets were not preserved: %+v", updated)
	}
	if err := repo.DeleteGithubApp(ctx, "app-one"); err != nil {
		t.Fatalf("delete app: %v", err)
	}
	if _, err := repo.GetGithubApp(ctx, "app-one"); err == nil {
		t.Fatal("expected deleted app to be gone")
	}
}
