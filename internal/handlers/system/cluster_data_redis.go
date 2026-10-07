package system

import (
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *ClusterDataHandler) ReviewRedisSnapshot(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	result, err := h.service.ReviewRedisSnapshot(c.Request().Context(), user.UserID, c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review Redis snapshot", result)
}

func (h *ClusterDataHandler) ApplyRedisSnapshot(c echo.Context) error {
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
	if err != nil || op.UserID != user.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "snapshot operation not found")
	}
	if err := h.service.ApplyRedisSnapshot(c.Request().Context(), user.UserID, op.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Snapshot applied", op)
}

func (h *ClusterDataHandler) ListRedisSnapshots(c echo.Context) error {
	result, err := h.service.ListRedisSnapshots(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Redis snapshots", result)
}

func (h *ClusterDataHandler) ReviewRedisRestore(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request struct {
		SnapshotID string `json:"snapshotId"`
		TargetID   string `json:"targetId"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restore request")
	}
	result, err := h.service.ReviewRedisRestore(c.Request().Context(), user.UserID, c.Param("id"), c.Param("clusterId"), c.Param("databaseId"), request.SnapshotID, request.TargetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review Redis restore", result)
}

func (h *ClusterDataHandler) ApplyRedisRestore(c echo.Context) error {
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
	if err != nil || op.UserID != user.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "restore operation not found")
	}
	if err := h.service.ApplyRedisRestore(c.Request().Context(), user.UserID, op.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Restore applied", op)
}

func (h *ClusterDataHandler) ReviewRedisFailover(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	result, err := h.service.ReviewRedisFailover(c.Request().Context(), user.UserID, c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review Redis failover", result)
}

func (h *ClusterDataHandler) ApplyRedisFailover(c echo.Context) error {
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
	if err != nil || op.UserID != user.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "failover operation not found")
	}
	if err := h.service.ApplyRedisFailover(c.Request().Context(), user.UserID, op.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Failover applied", op)
}
