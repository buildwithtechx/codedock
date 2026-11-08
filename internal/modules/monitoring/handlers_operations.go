package monitoring

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type OperationHandler struct {
	service OperationsService
}

func NewOperationHandler(service OperationsService) *OperationHandler {
	return &OperationHandler{service: service}
}

func (h *OperationHandler) Get(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return echo.NewHTTPError(http.StatusNotFound, "operation not found")
	}
	op, err := h.service.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op == nil || op.UserID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "operation not found")
	}
	return utils.Success(c, "Operation observed", op)
}

func (h *OperationHandler) List(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if _, token := c.Get("api_scopes").([]string); token && c.QueryParam("projectId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "token requests must select a project")
	}
	if h.service == nil {
		return utils.Success(c, "Operations observed", []Operation{})
	}
	result, err := h.service.List(c.Request().Context(), c.QueryParam("projectId"), userID)
	if err != nil {
		return err
	}
	return utils.Success(c, "Operations observed", result)
}

func (h *OperationHandler) Cancel(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if h.service == nil {
		return echo.NewHTTPError(http.StatusNotFound, "backup operation not found")
	}
	op, err := h.service.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op == nil || op.UserID != userID || (op.Kind != "backup" && op.Kind != "restore") {
		return echo.NewHTTPError(http.StatusNotFound, "backup operation not found")
	}
	if err := h.service.Cancel(c.Request().Context(), c.Param("operationId"), userID); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cancellation requested; partial changes may remain", nil)
}
