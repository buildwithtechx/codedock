package deployments

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	deploymentservices "codedock/internal/services/deployments"
	"codedock/internal/utils"
)

const (
	RollbackErrorNotReady              = deploymentservices.RollbackErrorNotReady
	RollbackErrorAlreadyActive         = deploymentservices.RollbackErrorAlreadyActive
	RollbackErrorArtifactGone          = deploymentservices.RollbackErrorArtifactGone
	RollbackErrorConfigurationRedacted = deploymentservices.RollbackErrorConfigurationRedacted
)

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

	var currentApp *models.AppService
	if targetDep.ServiceID != "" {
		if app, err := h.deploymentService.GetService(c.Request().Context(), targetDep.ServiceID); err == nil && app != nil {
			currentApp = app
		}
	}

	var validation *deploymentservices.RollbackValidationResult
	if h.retentionService != nil {
		validation = h.retentionService.ValidateRollback(c.Request().Context(), targetDep, currentApp)
	} else {
		validation = (&deploymentservices.RetentionService{}).ValidateRollback(c.Request().Context(), targetDep, currentApp)
	}

	if !validation.Eligible {
		return c.JSON(validation.StatusCode, map[string]any{
			"error":   validation.ErrorMessage,
			"code":    validation.ErrorCode,
			"success": false,
		})
	}

	rollbackImageRef := validation.RollbackImage

	newDep := &models.Deployment{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        "BUILDING",
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Rollback to " + targetDep.ID,
		Branch:        targetDep.Branch,
		ImageRef:      rollbackImageRef,
		Trigger:       "Rollback",
		BuildLogs:     "Rolling back to deployment " + targetDep.ID + "...\n",
	}

	if rollbackImageRef != "" && currentApp != nil {
		currentApp.ImageRef = rollbackImageRef
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

func (h *DeploymentHandler) TogglePin(c echo.Context) error {
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

	pinned, err := h.deploymentService.TogglePin(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Deployment pin toggled", models.PinDeploymentResponse{
		ID:     id,
		Pinned: pinned,
	})
}
