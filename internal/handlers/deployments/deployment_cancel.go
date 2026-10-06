package deployments

import (
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *DeploymentHandler) Cancel(c echo.Context) error {
	deployment, err := h.deploymentService.GetDeployment(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}
	if err := h.verifyProjectAdmin(c, deployment.ProjectID); err != nil {
		return err
	}
	if err := h.deploymentService.CancelDeployment(deployment.ID); err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Accepted(c, "Cancellation requested; restoring previous containers if activation has started", nil)
}
