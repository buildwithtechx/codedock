package deployments

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	"codedock/internal/utils"
)

const (
	RollbackErrorNotReady              = "ROLLBACK_NOT_READY"
	RollbackErrorAlreadyActive         = "ROLLBACK_ALREADY_ACTIVE"
	RollbackErrorArtifactGone          = "ROLLBACK_ARTIFACT_GONE"
	RollbackErrorConfigurationRedacted = "ROLLBACK_CONFIGURATION_REDACTED"
)

// HasRedactedConfiguration detects unrecoverable redacted secrets (e.g. "://***@" or "***")
// in saved frozen release configurations or logs.
func HasRedactedConfiguration(values ...string) bool {
	for _, val := range values {
		trimmed := strings.TrimSpace(val)
		if strings.Contains(val, "://***@") || trimmed == "***" || strings.Contains(val, "=\"***\"") || strings.Contains(val, "=***") {
			return true
		}
	}
	return false
}

func isRestorableStatus(status models.DeploymentStatus) bool {
	s := strings.ToUpper(string(status))
	return s == "READY" || s == "ACTIVE" || s == "SUCCESS" || s == "PARTIAL_FAILURE"
}

func (h *DeploymentHandler) Rollback(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	targetDep, err := h.deploymentService.GetDeployment(c.Request().Context(), id)
	if err != nil || targetDep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}

	if err := h.verifyProjectAdmin(c, targetDep.ProjectID); err != nil {
		return err
	}

	if !isRestorableStatus(targetDep.Status) {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error":   "Only successful deployments can be rolled back to.",
			"code":    RollbackErrorNotReady,
			"success": false,
		})
	}

	var currentApp *models.AppService
	if targetDep.ServiceID != "" {
		if app, err := h.deploymentService.GetService(c.Request().Context(), targetDep.ServiceID); err == nil && app != nil {
			currentApp = app
		}
	}

	if currentApp != nil && targetDep.ContainerID != "" && currentApp.ContainerID == targetDep.ContainerID {
		return c.JSON(http.StatusConflict, map[string]any{
			"error":   "This deployment is already active.",
			"code":    RollbackErrorAlreadyActive,
			"success": false,
		})
	}

	if HasRedactedConfiguration(targetDep.BuildLogs, targetDep.CommitMessage) {
		return c.JSON(http.StatusConflict, map[string]any{
			"error": "This release's saved configuration may contain credentials removed by an older version. " +
				"Rollback was not started and the active release was left unchanged. " +
				"Choose an intact release, or redeploy the original source with the required credentials to create a new restore point.",
			"code":    RollbackErrorConfigurationRedacted,
			"success": false,
		})
	}

	newDep := &models.Deployment{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        "BUILDING",
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Rollback to " + targetDep.ID,
		Branch:        targetDep.Branch,
		ImageRef:      targetDep.ImageRef,
		Trigger:       "Rollback",
		BuildLogs:     "Rolling back to deployment " + targetDep.ID + "...\n",
	}

	if targetDep.ImageRef != "" && currentApp != nil {
		currentApp.ImageRef = targetDep.ImageRef
		if updateErr := h.deploymentService.UpdateService(c.Request().Context(), currentApp); updateErr != nil {
			return utils.Error(c, http.StatusInternalServerError, updateErr.Error())
		}
	}

	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), newDep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "deployment.rollback",
		Resource:  newDep.ServiceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"deploymentId": created.ID,
		},
	})

	return utils.Accepted(c, "Rollback created", created)
}
