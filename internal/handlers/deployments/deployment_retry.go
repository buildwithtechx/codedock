package deployments

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	"codedock/internal/utils"
)

const (
	ErrDeploymentBuildUnavailable = "DEPLOYMENT_BUILD_UNAVAILABLE"
)

func (h *DeploymentHandler) Retry(c echo.Context) error {
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

	imageRef := strings.TrimSpace(targetDep.ImageRef)
	if imageRef == "" {
		return c.JSON(http.StatusConflict, map[string]any{
			"error":   "This completed build is no longer available for retry. Use Redeploy to rebuild it.",
			"code":    ErrDeploymentBuildUnavailable,
			"success": false,
		})
	}

	var app *models.AppService
	if targetDep.ServiceID != "" {
		app, _ = h.deploymentService.GetService(c.Request().Context(), targetDep.ServiceID)
	}

	// Verify image exists in Docker daemon before scheduling retry
	if app != nil {
		cli, release, dockerErr := h.deploymentService.DockerClientForApp(c.Request().Context(), app)
		if dockerErr != nil || cli == nil {
			return c.JSON(http.StatusConflict, map[string]any{
				"error":   "This completed build is no longer available for retry. Use Redeploy to rebuild it.",
				"code":    ErrDeploymentBuildUnavailable,
				"success": false,
			})
		}
		defer release()

		checkCtx, checkCancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer checkCancel()

		if _, _, inspectErr := cli.ImageInspectWithRaw(checkCtx, imageRef); inspectErr != nil {
			return c.JSON(http.StatusConflict, map[string]any{
				"error":   "This completed build is no longer available for retry. Use Redeploy to rebuild it.",
				"code":    ErrDeploymentBuildUnavailable,
				"success": false,
			})
		}
	}

	newDep := &models.Deployment{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        models.DeploymentStatusPending,
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Retry deployment " + targetDep.ID,
		Branch:        targetDep.Branch,
		ImageRef:      imageRef,
		Trigger:       "retry",
		BuildLogs:     "⚡ Retrying deployment using verified build artifact: " + imageRef + "\n",
	}

	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), newDep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "deployment.retry",
		Resource:  newDep.ServiceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"deploymentId":  created.ID,
			"retrySourceId": targetDep.ID,
			"imageRef":      imageRef,
		},
	})

	return utils.Accepted(c, "Deployment retry created", created)
}
