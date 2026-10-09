package projects

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AppHandler) GetRuntime(c echo.Context) error {
	result, err := h.Runtime.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Runtime configuration", result)
}
func (h *AppHandler) RuntimeOperation(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	result, err := h.Runtime.Operation(c.Request().Context(), user.UserID, c.Param("id"), c.QueryParam("operationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "runtime operation not found")
	}
	return utils.Success(c, "Runtime operation", result)
}
func (h *AppHandler) ReviewRuntime(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request models.RuntimeReviewRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid runtime configuration")
	}
	result, err := h.Runtime.Review(c.Request().Context(), user.UserID, c.Param("id"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review destination before applying", result)
}
func (h *AppHandler) ApplyRuntime(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request models.RuntimeOperationApply
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid runtime confirmation")
	}
	if err := h.Runtime.ApplyForService(c.Request().Context(), user.UserID, c.Param("id"), request.OperationID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Accepted(c, "Runtime configuration applying", nil)
}
func (h *AppHandler) ObserveRuntime(c echo.Context) error {
	result, err := h.Runtime.Observe(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Runtime observed", result)
}
func (h *AppHandler) RuntimeLogs(c echo.Context) error {
	result, err := h.Runtime.Logs(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Runtime logs", result)
}
func (h *AppHandler) RuntimeExec(c echo.Context) error {
	var request models.RuntimeExecRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid command")
	}
	result, err := h.Runtime.Exec(c.Request().Context(), c.Param("id"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Command completed", result)
}
