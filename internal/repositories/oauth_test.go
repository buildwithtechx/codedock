package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
)

func TestOAuthRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewOAuthRepo(db)
	p := &models.OAuthProviderConfig{
		ID:           uuid.NewString(),
		ProviderName: "github",
		Enabled:      true,
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "https://app.example.com/oauth/callback",
		BaseURL:      "https://github.com",
		Tenant:       "acme",
	}
	if err := repo.SaveProvider(ctx, p); err != nil {
		t.Fatalf("save provider: %v", err)
	}
	if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
		t.Fatal("expected save to stamp created and updated times")
	}
	byID, err := repo.GetProvider(ctx, p.ID)
	if err != nil {
		t.Fatalf("get provider by id: %v", err)
	}
	if byID.ProviderName != "github" || !byID.Enabled || byID.ClientID != "client-id" || byID.Tenant != "acme" {
		t.Fatalf("unexpected provider row: %+v", byID)
	}
	if !byID.CreatedAt.Truncate(time.Second).Equal(p.CreatedAt.Truncate(time.Second)) {
		t.Fatalf("created_at mismatch: got %v want %v", byID.CreatedAt, p.CreatedAt)
	}
	byName, err := repo.GetProvider(ctx, "github")
	if err != nil {
		t.Fatalf("get provider by name: %v", err)
	}
	if byName.ID != p.ID {
		t.Fatal("name lookup returned wrong provider")
	}
	providers, err := repo.ListProviders(ctx)
	if err != nil {
		t.Fatalf("list providers: %v", err)
	}
	if len(providers) != 1 || providers[0].ID != p.ID {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}
	p.Enabled = false
	p.ClientID = "rotated-id"
	if err := repo.SaveProvider(ctx, p); err != nil {
		t.Fatalf("update provider: %v", err)
	}
	updated, err := repo.GetProvider(ctx, p.ID)
	if err != nil {
		t.Fatalf("get updated provider: %v", err)
	}
	if updated.Enabled || updated.ClientID != "rotated-id" {
		t.Fatalf("update did not persist: %+v", updated)
	}
	if _, err := repo.GetProvider(ctx, "missing"); err == nil {
		t.Fatal("expected missing provider to error")
	}
	users := NewUserRepo(db)
	user := &models.User{Email: "totp@example.com", Name: "TOTP", PasswordHash: "hash", Role: models.UserRoleMember, IsActive: true}
	if err := users.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	secret, codes, err := repo.GetUserTOTPSecret(ctx, user.ID)
	if err != nil {
		t.Fatalf("get empty totp secret: %v", err)
	}
	if secret != "" || len(codes) != 0 {
		t.Fatalf("expected empty totp material, got %q %v", secret, codes)
	}
	if err := repo.UpdateUserTOTP(ctx, user.ID, true, "SECRET", []string{"a", "b", "c"}); err != nil {
		t.Fatalf("update totp: %v", err)
	}
	secret, codes, err = repo.GetUserTOTPSecret(ctx, user.ID)
	if err != nil {
		t.Fatalf("get totp secret: %v", err)
	}
	if secret != "SECRET" || len(codes) != 3 || codes[0] != "a" || codes[2] != "c" {
		t.Fatalf("totp did not roundtrip: %q %v", secret, codes)
	}
	if err := repo.UpdateUserTOTP(ctx, user.ID, false, "", nil); err != nil {
		t.Fatalf("clear totp: %v", err)
	}
	secret, codes, err = repo.GetUserTOTPSecret(ctx, user.ID)
	if err != nil {
		t.Fatalf("get cleared totp secret: %v", err)
	}
	if secret != "" || len(codes) != 0 {
		t.Fatalf("expected cleared totp material, got %q %v", secret, codes)
	}
}
