package repositories

import (
	"context"
	"testing"

	"codedock.run/codedock/internal/models"
)

func TestSelfHostedRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewSelfHostedRepo(db)
	empty, err := repo.Load(ctx)
	if err != nil {
		t.Fatalf("load empty: %v", err)
	}
	if empty.JWTSecret != "" || empty.WildcardDomain != "" {
		t.Fatalf("expected empty configuration, got %+v", empty)
	}
	stored := &models.SelfHostedConfig{JWTSecret: "jwt", RefreshSecret: "refresh", TelemetrySalt: "salt", TLSEmail: "owner@example.com", WildcardDomain: "apps.example.com"}
	if err := repo.Save(ctx, stored); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := repo.Load(ctx)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if *loaded != *stored {
		t.Fatalf("mismatch: %+v", loaded)
	}
	stored.WildcardDomain = "other.example.com"
	if err := repo.Save(ctx, stored); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	updated, err := repo.Load(ctx)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if updated.WildcardDomain != "other.example.com" || updated.JWTSecret != "jwt" {
		t.Fatalf("overwrite did not persist: %+v", updated)
	}
}
