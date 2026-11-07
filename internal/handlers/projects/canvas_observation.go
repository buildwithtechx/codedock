package projects

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *CanvasHandler) ObserveResource(c echo.Context) error {
	canvas, err := h.canvasService.GetEnvironmentCanvas(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "environment unavailable")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil || !h.projectService.HasPermission(c.Request().Context(), canvas.Environment.ProjectID, user.UserID, models.UserRole(user.Role), "") {
		return utils.Error(c, http.StatusForbidden, "project access denied")
	}
	observation, err := h.canvasService.ReadObservation(c.Request().Context(), canvas, c.QueryParam("node"))
	if err != nil {
		return utils.Error(c, http.StatusServiceUnavailable, err.Error())
	}
	return utils.Success(c, "Runtime observation retrieved", observation)
}
