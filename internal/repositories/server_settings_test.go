package repositories

import (
	"context"
	"testing"
)

func TestServerSettingsCreatesAndUpdatesDefaults(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewSettingsRepo(db, nil)
	cfg, err := repo.GetServerSettings(ctx)
	if err != nil {
		t.Fatalf("get defaults: %v", err)
	}
	if cfg.ID != "global" || cfg.SiteName != "" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	cfg.SiteName = "Acme Dock"
	cfg.RegistrationEnabled = true
	cfg.TelemetryEnabled = false
	if err := repo.UpdateServerSettings(ctx, cfg); err != nil {
		t.Fatalf("update: %v", err)
	}
	updated, err := repo.GetServerSettings(ctx)
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if updated.SiteName != "Acme Dock" || !updated.RegistrationEnabled || updated.TelemetryEnabled {
		t.Fatalf("update did not persist: %+v", updated)
	}
}
