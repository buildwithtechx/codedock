package deployments

import (
	"context"
	"net/http"
	"strings"

	"codedock/internal/utils"
)

const (
	RollbackErrorNotReady              = "ROLLBACK_NOT_READY"
	RollbackErrorAlreadyActive         = "ROLLBACK_ALREADY_ACTIVE"
	RollbackErrorArtifactGone          = "ROLLBACK_ARTIFACT_GONE"
	RollbackErrorConfigurationRedacted = "ROLLBACK_CONFIGURATION_REDACTED"
)

type RollbackValidationResult struct {
	Eligible      bool   `json:"eligible"`
	StatusCode    int    `json:"statusCode"`
	ErrorCode     string `json:"errorCode,omitempty"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
	UseArtifact   bool   `json:"useArtifact"`
	RollbackImage string `json:"rollbackImage,omitempty"`
}

type RetentionService struct {
	deploymentRepo DeploymentRepository
}

func NewRetentionService(deploymentRepo DeploymentRepository) *RetentionService {
	return &RetentionService{
		deploymentRepo: deploymentRepo,
	}
}

func isRestorableStatus(status DeploymentStatus) bool {
	upper := strings.ToUpper(string(status))
	return upper == "READY" || upper == "ACTIVE" || upper == "SUCCESS"
}

func (s *RetentionService) ValidateRollback(ctx context.Context, targetDep *Deployment, currentApp *AppServiceSummary) *RollbackValidationResult {
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

	isArtifactAvailable := targetDep.ArtifactRetainedAt != nil && targetDep.ImageRef != ""

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

func (s *RetentionService) PruneProject(ctx context.Context, projectID string, maxRetained int) (int, error) {
	if maxRetained <= 0 {
		maxRetained = DefaultRetentionCount
	}
	deployments, err := s.deploymentRepo.ListForRetention(ctx, projectID)
	if err != nil {
		return 0, err
	}

	pruned := 0
	eligibleCount := 0
	for _, dep := range deployments {
		if dep.Pinned {
			continue
		}
		eligibleCount++
		if eligibleCount > maxRetained {
			if dep.ArtifactRetainedAt != nil {
				_ = s.deploymentRepo.SetArtifactRetainedAt(ctx, dep.ID, nil)
				pruned++
			}
		}
	}
	return pruned, nil
}
