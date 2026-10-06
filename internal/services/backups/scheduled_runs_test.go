package backups

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"context"
	"errors"
	"strings"
	"testing"
)

type scheduledPolicyStore struct {
	repositories.BackupRepository
	config *models.BackupConfig
}

func (s scheduledPolicyStore) GetConfigByID(context.Context, string) (*models.BackupConfig, error) {
	return s.config, nil
}
func TestScheduledRunRequiresCurrentOwnerAuthorization(t *testing.T) {
	cfg := &models.BackupConfig{ID: "policy", OwnerID: "owner", ProjectID: "project", BackupEnabled: true, Status: models.BackupConfigStatusActive, Schedule: "0 2 * * *"}
	service := NewBackupService(scheduledPolicyStore{config: cfg}, nil, nil)
	called := false
	service.RunAuthorization = func(_ context.Context, user, project string) error {
		called = true
		if user != "owner" || project != "project" {
			t.Fatal("scheduled run used wrong identity")
		}
		return errors.New("permission revoked")
	}
	if err := service.StartScheduledRun(context.Background(), cfg.ID); err == nil || !strings.Contains(err.Error(), "revoked") || !called {
		t.Fatal("revoked owner was allowed to start backup", err)
	}
	called = false
	cfg.OwnerID = ""
	if err := service.StartScheduledRun(context.Background(), cfg.ID); err == nil || called {
		t.Fatal("ownerless scheduled policy executed")
	}
	cfg.OwnerID = "owner"
	cfg.BackupEnabled = false
	if err := service.StartScheduledRun(context.Background(), cfg.ID); err == nil || called {
		t.Fatal("disabled policy executed")
	}
	cfg.BackupEnabled = true
	cfg.Schedule = "manual"
	if err := service.StartScheduledRun(context.Background(), cfg.ID); err == nil || called {
		t.Fatal("manual policy scheduled")
	}
}
