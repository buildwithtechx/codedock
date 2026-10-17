package deployments

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *DeploymentHandler) ExplainFailure(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	dep, err := h.deploymentService.GetDeployment(c.Request().Context(), id)
	if err != nil || dep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}

	if err := h.verifyProjectOwnership(c, dep.ProjectID); err != nil {
		return err
	}

	explanation, err := h.aiAnalysis.ExplainDeploymentFailure(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "AI Analysis completed", explanation)
}

func (h *DeploymentHandler) ListPRPreviews(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	if err := h.verifyProjectOwnership(c, svc.ProjectID); err != nil {
		return err
	}

	previews, err := h.prPreviewService.ListByApp(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", previews)
}
