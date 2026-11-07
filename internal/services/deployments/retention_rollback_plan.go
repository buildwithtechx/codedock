package deployments

import (
	"context"
	"net/http"

	"codedock/internal/models"
	"codedock/internal/utils"
)

const (
	RollbackErrorNotReady              = "ROLLBACK_NOT_READY"
	RollbackErrorAlreadyActive         = "ROLLBACK_ALREADY_ACTIVE"
	RollbackErrorArtifactGone          = "ROLLBACK_ARTIFACT_GONE"
	RollbackErrorConfigurationRedacted = "ROLLBACK_CONFIGURATION_REDACTED"
)

// RollbackValidationResult contains the domain validation decision for a rollback attempt.
type RollbackValidationResult struct {
	Eligible      bool
	StatusCode    int
	ErrorCode     string
	ErrorMessage  string
	UseArtifact   bool
	RollbackImage string
}

// ValidateRollback evaluates whether a past target deployment can be restored,
// verifying status readiness, non-active status, non-redacted configuration, and artifact availability.
func (s *RetentionService) ValidateRollback(ctx context.Context, targetDep *models.Deployment, currentApp *models.AppService) *RollbackValidationResult {
	if targetDep == nil || !isRestorableStatus(targetDep.Status) {
		return &RollbackValidationResult{
			Eligible:     false,
			StatusCode:   http.StatusBadRequest,
			ErrorCode:    RollbackErrorNotReady,
			ErrorMessage: "Only successful deployments can be rolled back to.",
		}
	}

	if currentApp != nil && targetDep.ContainerID != "" && currentApp.ContainerID == targetDep.ContainerID {
		return &RollbackValidationResult{
			Eligible:     false,
			StatusCode:   http.StatusConflict,
			ErrorCode:    RollbackErrorAlreadyActive,
			ErrorMessage: "This deployment is already active.",
		}
	}

	if utils.HasRedactedConfiguration(targetDep.BuildLogs, targetDep.CommitMessage) {
		return &RollbackValidationResult{
			Eligible:   false,
			StatusCode: http.StatusConflict,
			ErrorCode:  RollbackErrorConfigurationRedacted,
			ErrorMessage: "This release's saved configuration may contain credentials removed by an older version. " +
				"Rollback was not started and the active release was left unchanged. " +
				"Choose an intact release, or redeploy the original source with the required credentials to create a new restore point.",
		}
	}

	var isArtifactAvailable bool
	if s != nil {
		available, _ := s.IsArtifactAvailable(ctx, targetDep.ProjectID, targetDep)
		isArtifactAvailable = available
	} else {
		isArtifactAvailable = targetDep.ArtifactRetainedAt != nil && targetDep.ImageRef != ""
	}

	if !isArtifactAvailable && targetDep.CommitHash == "" {
		return &RollbackValidationResult{
			Eligible:     false,
			StatusCode:   http.StatusConflict,
			ErrorCode:    RollbackErrorArtifactGone,
			ErrorMessage: "This deployment can't be restored: its artifact is no longer retained and it has no commit to rebuild from.",
		}
	}

	rollbackImage := ""
	if isArtifactAvailable {
		rollbackImage = targetDep.ImageRef
	}

	return &RollbackValidationResult{
		Eligible:      true,
		UseArtifact:   isArtifactAvailable,
		RollbackImage: rollbackImage,
	}
}
