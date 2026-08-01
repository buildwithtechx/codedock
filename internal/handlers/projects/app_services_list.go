package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/utils"
)

func (h *AppHandler) ListByOrganization(c echo.Context) error {
	organizationID := c.QueryParam("organizationId")
	if organizationID == "" {
		return utils.Error(c, http.StatusBadRequest, "organizationId is required")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if !h.projectService.HasOrgPermission(c.Request().Context(), organizationID, user.UserID, user.Role, "") {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions for this organization")
	}
	apps, err := h.appService.ListByOrganization(c.Request().Context(), organizationID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", apps)
}
