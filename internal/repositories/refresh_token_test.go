package repositories

import (
	"context"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
)

func TestRefreshTokenLifecycle(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	users := NewUserRepo(db)
	user := &models.User{Email: "refresh@example.com", Name: "Refresh", PasswordHash: "hash", Role: models.UserRoleMember, IsActive: true}
	if err := users.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	repo := NewRefreshTokenRepo(db)
	hash := HashToken("token-one")
	if err := repo.StoreToken(ctx, user.ID, hash, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("store token: %v", err)
	}
	revoked, err := repo.IsRevoked(ctx, hash)
	if err != nil {
		t.Fatalf("check revocation: %v", err)
	}
	if revoked {
		t.Fatal("fresh token reported revoked")
	}
	if err := repo.RevokeToken(ctx, hash); err != nil {
		t.Fatalf("revoke token: %v", err)
	}
	revoked, err = repo.IsRevoked(ctx, hash)
	if err != nil {
		t.Fatalf("recheck revocation: %v", err)
	}
	if !revoked {
		t.Fatal("revoked token reported active")
	}
	if err := repo.StoreToken(ctx, user.ID, hash, time.Now().Add(2*time.Hour)); err != nil {
		t.Fatalf("re-store token: %v", err)
	}
	revoked, err = repo.IsRevoked(ctx, hash)
	if err != nil {
		t.Fatalf("check after re-store: %v", err)
	}
	if revoked {
		t.Fatal("re-stored token reported revoked")
	}
	second := HashToken("token-two")
	if err := repo.StoreToken(ctx, user.ID, second, time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("store second token: %v", err)
	}
	if err := repo.RevokeAllForUser(ctx, user.ID); err != nil {
		t.Fatalf("revoke all: %v", err)
	}
	revoked, err = repo.IsRevoked(ctx, second)
	if err != nil {
		t.Fatalf("check second revocation: %v", err)
	}
	if !revoked {
		t.Fatal("second token reported active after revoke-all")
	}
	expired := HashToken("token-expired")
	if err := repo.StoreToken(ctx, user.ID, expired, time.Now().Add(-time.Hour)); err != nil {
		t.Fatalf("store expired token: %v", err)
	}
	if err := repo.PruneExpired(ctx); err != nil {
		t.Fatalf("prune expired: %v", err)
	}
	revoked, err = repo.IsRevoked(ctx, expired)
	if err != nil {
		t.Fatalf("check pruned token: %v", err)
	}
	if revoked {
		t.Fatal("pruned token still present")
	}
}
