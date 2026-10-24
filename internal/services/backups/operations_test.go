package backups

import (
	"codedock/internal/models"
	"codedock/internal/repositories"
	"context"
	"strings"
	"testing"
)

type requiredPolicyStore struct {
	repositories.BackupRepository
	policies []*models.BackupConfig
}

func (s requiredPolicyStore) ListConfigs(context.Context) ([]*models.BackupConfig, error) {
	return s.policies, nil
}
func TestRequiredBackupBlocksDeploymentWhenDisabledOrRuntimeUnavailable(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		policy := &models.BackupConfig{ID: "policy", Name: "Required", ServiceID: "service", PreDeployment: true, BackupEnabled: enabled}
		service := NewBackupService(requiredPolicyStore{policies: []*models.BackupConfig{policy}}, nil, nil)
		err := service.BeforeDeployment(context.Background(), "service")
		if err == nil {
			t.Fatal("required backup allowed deployment without a verified archive")
		}
		if !enabled && !strings.Contains(err.Error(), "disabled") {
			t.Fatal("disabled policy cause lost", err)
		}
		if err := service.BeforeDeployment(context.Background(), "other-service"); err != nil {
			t.Fatal("unrelated policy blocked deployment", err)
		}
	}
}
