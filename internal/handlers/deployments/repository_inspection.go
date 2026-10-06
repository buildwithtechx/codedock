package deployments

import (
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *GitHandler) InspectRepository(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var request models.RepositoryInspectionRequest
	if err := c.Bind(&request); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid inspection request")
	}
	result, err := h.gitService.InspectRepository(c.Request().Context(), user.UserID, request)
	if err != nil {
		return utils.Error(c, http.StatusUnprocessableEntity, err.Error())
	}
	return utils.Success(c, "Repository inspected", result)
}
