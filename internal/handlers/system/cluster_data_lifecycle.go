package system

import (
	"codedock/internal/http/middleware"
	"codedock/internal/services/clusterdata"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *ClusterDataHandler) ReviewLifecycle(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request clusterdata.LifecycleRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid lifecycle request")
	}
	result, err := h.service.ReviewLifecycle(c.Request().Context(), user.UserID, c.Param("id"), c.Param("clusterId"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review database lifecycle", result)
}

func (h *ClusterDataHandler) ApplyLifecycle(c echo.Context) error {
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
	if err != nil || op.Kind != "cluster-data-lifecycle" || op.UserID != user.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "lifecycle operation not found")
	}
	if err := h.service.ApplyLifecycle(c.Request().Context(), user.UserID, op.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Lifecycle applied", op)
}
