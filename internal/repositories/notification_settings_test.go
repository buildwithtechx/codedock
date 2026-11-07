package repositories

import (
	"context"
	"testing"
)

func TestNotificationSettingsRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewNotificationSettingsRepo(db)
	cfg, err := repo.GetNotificationSettings(ctx)
	if err != nil {
		t.Fatalf("get defaults: %v", err)
	}
	if cfg.ID != "global" || !cfg.NotificationAlerts || cfg.SMTPPort != 587 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	cfg.SlackEnabled = true
	cfg.SlackWebhookURL = "https://hooks.slack.test/x"
	cfg.SMTPPort = 2525
	if err := repo.UpdateNotificationSettings(ctx, cfg); err != nil {
		t.Fatalf("update: %v", err)
	}
	updated, err := repo.GetNotificationSettings(ctx)
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if !updated.SlackEnabled || updated.SlackWebhookURL != "https://hooks.slack.test/x" || updated.SMTPPort != 2525 {
		t.Fatalf("update did not persist: %+v", updated)
	}
}
