package repositories

import (
	"context"
	"testing"
)

func TestAISettingsRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewAISettingsRepo(db)
	cfg, err := repo.GetAISettings(ctx)
	if err != nil {
		t.Fatalf("get defaults: %v", err)
	}
	if cfg.ID != "global" || cfg.DefaultProvider != "none" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	cfg.DefaultProvider = "openai"
	cfg.OpenAIModel = "gpt-4o"
	if err := repo.UpdateAISettings(ctx, cfg); err != nil {
		t.Fatalf("update: %v", err)
	}
	updated, err := repo.GetAISettings(ctx)
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if updated.DefaultProvider != "openai" || updated.OpenAIModel != "gpt-4o" {
		t.Fatalf("update did not persist: %+v", updated)
	}
}
