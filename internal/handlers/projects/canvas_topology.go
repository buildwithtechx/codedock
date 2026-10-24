package projects

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *CanvasHandler) ApplyTopology(c echo.Context) error {
	canvas, err := h.canvasService.GetEnvironmentCanvas(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "environment not found")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil || !h.projectService.HasPermission(c.Request().Context(), canvas.Environment.ProjectID, user.UserID, models.UserRole(user.Role), models.MemberPermissionAdmin) {
		return utils.Error(c, http.StatusForbidden, "project admin access required")
	}
	var request models.TopologyApplyRequest
	if err := c.Bind(&request); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid topology request")
	}
	if err := h.canvasService.ApplyTopology(c.Request().Context(), canvas.Environment.ID, request); err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Dependencies saved; deployment checks these prerequisites before starting the application", nil)
}
