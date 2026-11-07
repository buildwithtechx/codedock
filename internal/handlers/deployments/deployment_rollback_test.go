package deployments

import (
	"context"
	"testing"

	"codedock/internal/models"
	deploymentservices "codedock/internal/services/deployments"
	"codedock/internal/utils"
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
			got := utils.HasRedactedConfiguration(tc.values...)
			if got != tc.expected {
				t.Errorf("HasRedactedConfiguration(%v) = %v; want %v", tc.values, got, tc.expected)
			}
		})
	}
}

func TestValidateRollback(t *testing.T) {
	svc := &deploymentservices.RetentionService{}

	// Not ready
	res := svc.ValidateRollback(context.Background(), &models.Deployment{Status: "PENDING"}, nil)
	if res.Eligible || res.ErrorCode != RollbackErrorNotReady {
		t.Errorf("expected RollbackErrorNotReady, got %+v", res)
	}

	// Already active
	res = svc.ValidateRollback(context.Background(), &models.Deployment{
		Status:      "READY",
		ContainerID: "cont-123",
	}, &models.AppService{ContainerID: "cont-123"})
	if res.Eligible || res.ErrorCode != RollbackErrorAlreadyActive {
		t.Errorf("expected RollbackErrorAlreadyActive, got %+v", res)
	}

	// Redacted
	res = svc.ValidateRollback(context.Background(), &models.Deployment{
		Status:    "READY",
		BuildLogs: "secret: ***",
	}, nil)
	if res.Eligible || res.ErrorCode != RollbackErrorConfigurationRedacted {
		t.Errorf("expected RollbackErrorConfigurationRedacted, got %+v", res)
	}

	// Artifact gone and no commit hash
	res = svc.ValidateRollback(context.Background(), &models.Deployment{
		Status: "READY",
	}, nil)
	if res.Eligible || res.ErrorCode != RollbackErrorArtifactGone {
		t.Errorf("expected RollbackErrorArtifactGone, got %+v", res)
	}

	// Rebuild fallback when artifact gone but commit hash exists
	res = svc.ValidateRollback(context.Background(), &models.Deployment{
		Status:     "READY",
		CommitHash: "abc1234",
	}, nil)
	if !res.Eligible || res.UseArtifact || res.RollbackImage != "" {
		t.Errorf("expected eligible rebuild fallback, got %+v", res)
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
