package repositories

import (
	"context"
	"testing"
)

func TestDeviceTokenRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewDeviceTokenRepo(db)

	token, err := repo.RegisterToken(ctx, "usr-1", "ExponentPushToken[abc-123]", "expo")
	if err != nil {
		t.Fatalf("register token: %v", err)
	}
	if token.ID == "" || token.Token != "ExponentPushToken[abc-123]" || token.Platform != "expo" {
		t.Fatalf("unexpected token returned: %+v", token)
	}

	// Idempotent upsert with different user or platform
	updated, err := repo.RegisterToken(ctx, "usr-2", "ExponentPushToken[abc-123]", "ios")
	if err != nil {
		t.Fatalf("upsert token: %v", err)
	}
	if updated.UserID != "usr-2" || updated.Platform != "ios" {
		t.Fatalf("expected updated attributes, got %+v", updated)
	}

	userTokens, err := repo.ListTokensByUser(ctx, "usr-2")
	if err != nil {
		t.Fatalf("list user tokens: %v", err)
	}
	if len(userTokens) != 1 || userTokens[0].Token != "ExponentPushToken[abc-123]" {
		t.Fatalf("expected 1 user token, got %+v", userTokens)
	}

	allTokens, err := repo.ListAllTokens(ctx)
	if err != nil {
		t.Fatalf("list all tokens: %v", err)
	}
	if len(allTokens) != 1 {
		t.Fatalf("expected 1 total token, got %d", len(allTokens))
	}

	if err := repo.DeleteToken(ctx, "usr-2", "ExponentPushToken[abc-123]"); err != nil {
		t.Fatalf("delete token: %v", err)
	}

	remaining, err := repo.ListAllTokens(ctx)
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("expected 0 tokens after delete, got %d", len(remaining))
	}
}
