package repositories

import (
	"context"
	"errors"
	"strings"
	"testing"

	"codedock/internal/models"
)

type gitTestVault struct{}

func (gitTestVault) Encrypt(plaintext string) (string, error) {
	return "enc:" + plaintext, nil
}

func (gitTestVault) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, "enc:") {
		return "", errors.New("not encrypted")
	}
	return strings.TrimPrefix(ciphertext, "enc:"), nil
}

func TestGitProviderRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	if _, err := db.Exec(`INSERT INTO users (id, email, name, password_hash) VALUES ('user-one', 'git@example.com', 'Git', 'hash')`); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	ctx := context.Background()
	repo := NewGitRepo(db, gitTestVault{})
	provider := &models.GitProviderConfig{UserID: "user-one", Provider: "github", AccessToken: "token-one", AccountName: "octo"}
	if err := repo.SaveProvider(ctx, provider); err != nil {
		t.Fatalf("save provider: %v", err)
	}
	if provider.ID == "" || provider.CreatedAt.IsZero() {
		t.Fatal("expected generated id and timestamps")
	}
	var stored string
	if err := db.QueryRowContext(ctx, `SELECT encrypted_access_token FROM user_git_providers WHERE id = $1`, provider.ID).Scan(&stored); err != nil {
		t.Fatalf("read stored token: %v", err)
	}
	if stored == "" || stored == "token-one" {
		t.Fatalf("expected encrypted token at rest, got %q", stored)
	}
	loaded, err := repo.GetProvider(ctx, "user-one", "github")
	if err != nil {
		t.Fatalf("get provider: %v", err)
	}
	if loaded == nil || loaded.AccessToken != "token-one" || loaded.AccountName != "octo" {
		t.Fatalf("unexpected provider row: %+v", loaded)
	}
	any, err := repo.GetAnyProviderByType(ctx, "github")
	if err != nil {
		t.Fatalf("get any provider: %v", err)
	}
	if any == nil || any.AccessToken != "token-one" {
		t.Fatalf("unexpected any-provider row: %+v", any)
	}
	fallback, err := repo.GetProvider(ctx, "", "github")
	if err != nil {
		t.Fatalf("get provider with empty user: %v", err)
	}
	if fallback == nil || fallback.ID != provider.ID {
		t.Fatalf("expected fallback to any provider, got %+v", fallback)
	}
	listed, err := repo.ListProvidersByUser(ctx, "user-one")
	if err != nil {
		t.Fatalf("list providers: %v", err)
	}
	if len(listed) != 1 || listed[0].AccessToken != "token-one" {
		t.Fatalf("unexpected provider list: %+v", listed)
	}
	provider.AccessToken = "token-two"
	provider.AccountName = "octo-two"
	if err := repo.SaveProvider(ctx, provider); err != nil {
		t.Fatalf("resave provider: %v", err)
	}
	updated, err := repo.GetProvider(ctx, "user-one", "github")
	if err != nil {
		t.Fatalf("get updated provider: %v", err)
	}
	if updated == nil || updated.AccessToken != "token-two" || updated.AccountName != "octo-two" {
		t.Fatalf("upsert did not persist: %+v", updated)
	}
	if err := repo.DeleteProvider(ctx, "user-one", "github"); err != nil {
		t.Fatalf("delete provider: %v", err)
	}
	deleted, err := repo.GetProvider(ctx, "user-one", "github")
	if err != nil {
		t.Fatalf("get deleted provider: %v", err)
	}
	if deleted != nil {
		t.Fatalf("expected deleted provider to be gone, got %+v", deleted)
	}
}
