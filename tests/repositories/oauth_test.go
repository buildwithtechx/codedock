package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"

	_ "modernc.org/sqlite"
)

func TestOAuthRepositoryOnFreshDatabase(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := repositories.NewOAuthRepo(db)
	providers, err := repo.ListProviders(context.Background())
	if err != nil {
		t.Fatalf("list oauth providers: %v", err)
	}
	if len(providers) != 0 {
		t.Fatalf("expected no oauth providers, got %d", len(providers))
	}

	provider := &models.OAuthProviderConfig{
		ID:           "github",
		ProviderName: "github",
		Enabled:      true,
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "http://localhost:8080/api/auth/oauth/github/callback",
		BaseURL:      "https://github.com",
	}
	if err := repo.SaveProvider(context.Background(), provider); err != nil {
		t.Fatalf("save oauth provider: %v", err)
	}

	providers, err = repo.ListProviders(context.Background())
	if err != nil {
		t.Fatalf("list saved oauth providers: %v", err)
	}
	if len(providers) != 1 || providers[0].ProviderName != "github" {
		t.Fatalf("unexpected oauth providers: %+v", providers)
	}
}
