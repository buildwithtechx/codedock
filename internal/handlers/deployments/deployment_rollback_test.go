package deployments

import (
	"testing"

	"codedock/internal/models"
)

func TestHasRedactedConfiguration(t *testing.T) {
	cases := []struct {
		name     string
		values   []string
		expected bool
	}{
		{
			name:     "valid clean configuration",
			values:   []string{"postgresql://app:release-secret@db:5432/app", "PORT=3000"},
			expected: false,
		},
		{
			name:     "scrubbed dsn with asterisks",
			values:   []string{"Connecting to postgresql://***@db:5432/app"},
			expected: true,
		},
		{
			name:     "bare asterisk redaction",
			values:   []string{"***"},
			expected: true,
		},
		{
			name:     "quoted secret mask",
			values:   []string{`DATABASE_URL="***"`},
			expected: true,
		},
		{
			name:     "assignment mask",
			values:   []string{"SECRET_TOKEN=***"},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := HasRedactedConfiguration(tc.values...)
			if got != tc.expected {
				t.Errorf("HasRedactedConfiguration(%v) = %v; want %v", tc.values, got, tc.expected)
			}
		})
	}
}

func TestIsRestorableStatus(t *testing.T) {
	valid := []models.DeploymentStatus{
		models.DeploymentStatusReady,
		models.DeploymentStatusActive,
		models.DeploymentStatus("SUCCESS"),
		models.DeploymentStatus("PARTIAL_FAILURE"),
	}
	for _, s := range valid {
		if !isRestorableStatus(s) {
			t.Errorf("expected status %s to be restorable", s)
		}
	}

	invalid := []models.DeploymentStatus{
		models.DeploymentStatusPending,
		models.DeploymentStatusFailed,
		models.DeploymentStatusBuilding,
		models.DeploymentStatus("CANCELLED"),
	}
	for _, s := range invalid {
		if isRestorableStatus(s) {
			t.Errorf("expected status %s NOT to be restorable", s)
		}
	}
}

func TestRollbackErrorCodes(t *testing.T) {
	if RollbackErrorArtifactGone != "ROLLBACK_ARTIFACT_GONE" {
		t.Errorf("expected RollbackErrorArtifactGone to be ROLLBACK_ARTIFACT_GONE, got %s", RollbackErrorArtifactGone)
	}
	if RollbackErrorNotReady != "ROLLBACK_NOT_READY" {
		t.Errorf("expected RollbackErrorNotReady to be ROLLBACK_NOT_READY, got %s", RollbackErrorNotReady)
	}
	if RollbackErrorAlreadyActive != "ROLLBACK_ALREADY_ACTIVE" {
		t.Errorf("expected RollbackErrorAlreadyActive to be ROLLBACK_ALREADY_ACTIVE, got %s", RollbackErrorAlreadyActive)
	}
	if RollbackErrorConfigurationRedacted != "ROLLBACK_CONFIGURATION_REDACTED" {
		t.Errorf("expected RollbackErrorConfigurationRedacted to be ROLLBACK_CONFIGURATION_REDACTED, got %s", RollbackErrorConfigurationRedacted)
	}
}
