package clusters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *DataHandler) ReviewBinding(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request BindingReviewRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid binding request")
	}
	result, err := h.service.ReviewBinding(c.Request().Context(), userID, c.Param("id"), c.Param("clusterId"), request.AppID, request.DatabaseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review database binding", result)
}

func (h *DataHandler) ApplyBinding(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request BindingApplyRequest
	if err := c.Bind(&request); err != nil || request.OperationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.service.ApplyBinding(c.Request().Context(), userID, request.OperationID, request.Confirmation)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Binding applied", op)
}

func (h *DataHandler) ReviewLifecycle(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request LifecycleReviewRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid lifecycle request")
	}
	result, err := h.service.ReviewLifecycle(c.Request().Context(), userID, c.Param("id"), c.Param("clusterId"), request.DatabaseID, request.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review database lifecycle", result)
}

func (h *DataHandler) ApplyLifecycle(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request LifecycleApplyRequest
	if err := c.Bind(&request); err != nil || request.OperationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.service.ApplyLifecycle(c.Request().Context(), userID, request.OperationID, request.Confirmation)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Lifecycle applied", op)
}

func (h *DataHandler) ReviewRedisSnapshot(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	result, err := h.service.ReviewRedisSnapshot(c.Request().Context(), userID, c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review Redis snapshot", result)
}

func (h *DataHandler) ApplyRedisSnapshot(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request RedisSnapshotApplyRequest
	if err := c.Bind(&request); err != nil || request.OperationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.service.ApplyRedisSnapshot(c.Request().Context(), userID, request.OperationID, request.Confirmation)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Snapshot applied", op)
}

func (h *DataHandler) ListRedisSnapshots(c echo.Context) error {
	result, err := h.service.ListRedisSnapshots(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Redis snapshots", result)
}

func (h *DataHandler) ReviewRedisRestore(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request RedisRestoreReviewRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restore request")
	}
	result, err := h.service.ReviewRedisRestore(c.Request().Context(), userID, c.Param("id"), c.Param("clusterId"), c.Param("databaseId"), request.SnapshotID, request.TargetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review Redis restore", result)
}

func (h *DataHandler) ApplyRedisRestore(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request RedisRestoreApplyRequest
	if err := c.Bind(&request); err != nil || request.OperationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.service.ApplyRedisRestore(c.Request().Context(), userID, request.OperationID, request.Confirmation)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Restore applied", op)
}

func (h *DataHandler) ReviewRedisFailover(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	result, err := h.service.ReviewRedisFailover(c.Request().Context(), userID, c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review Redis failover", result)
}

func (h *DataHandler) ApplyRedisFailover(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request RedisFailoverApplyRequest
	if err := c.Bind(&request); err != nil || request.OperationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.service.ApplyRedisFailover(c.Request().Context(), userID, request.OperationID, request.Confirmation)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Failover applied", op)
}
