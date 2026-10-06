package system

import (
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/services/operations"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type OperationHandler struct{ service *operations.Service }

func NewOperationHandler(service *operations.Service) *OperationHandler {
	return &OperationHandler{service: service}
}
func (h *OperationHandler) Get(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	op, err := h.service.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op.UserID != user.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "operation not found")
	}
	return utils.Success(c, "Operation observed", op)
}
func (h *OperationHandler) List(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if _, token := c.Get("api_scopes").([]string); token && c.QueryParam("projectId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "token requests must select a project")
	}
	result, err := h.service.List(c.Request().Context(), c.QueryParam("projectId"), user.UserID)
	if err != nil {
		return err
	}
	return utils.Success(c, "Operations observed", result)
}
func (h *OperationHandler) Cancel(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	op, err := h.service.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op.UserID != user.UserID || (op.Kind != "backup" && op.Kind != "restore") {
		return echo.NewHTTPError(http.StatusNotFound, "backup operation not found")
	}
	if err := h.service.Cancel(c.Request().Context(), c.Param("operationId"), user.UserID); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cancellation requested; partial changes may remain", nil)
}
