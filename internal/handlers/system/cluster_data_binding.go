package system

import (
	"codedock/internal/http/middleware"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *ClusterDataHandler) ReviewBinding(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request struct {
		AppID      string `json:"appId"`
		DatabaseID string `json:"databaseId"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid binding request")
	}
	if h.apps == nil || h.vars == nil {
		return echo.NewHTTPError(http.StatusConflict, "binding workflow unavailable")
	}
	result, err := h.service.ReviewBinding(c.Request().Context(), user.UserID, c.Param("id"), c.Param("clusterId"), request.AppID, request.DatabaseID, h.apps, h.vars)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review database binding", result)
}

func (h *ClusterDataHandler) ApplyBinding(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request struct {
		OperationID  string `json:"operationId"`
		Confirmation string `json:"confirmation"`
	}
	if err := c.Bind(&request); err != nil || request.OperationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.operations.Get(c.Request().Context(), request.OperationID)
	if err != nil || op.Kind != "cluster-binding" || op.UserID != user.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "binding operation not found")
	}
	if h.apps == nil || h.vars == nil {
		return echo.NewHTTPError(http.StatusConflict, "binding workflow unavailable")
	}
	if err := h.service.ApplyBinding(c.Request().Context(), user.UserID, op.ID, request.Confirmation, h.apps, h.vars); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Binding applied", op)
}
