package monitoring

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type MetricsHandler struct {
	service MetricsOperations
}

func NewMetricsHandler(s MetricsOperations) *MetricsHandler {
	return &MetricsHandler{service: s}
}

func (h *MetricsHandler) GetHistoricalMetrics(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	rangeParam := c.QueryParam("range")
	if rangeParam == "" {
		rangeParam = "24h"
	}

	data, err := h.service.GetServiceMetrics(c.Request().Context(), serviceID, rangeParam)
	if err != nil {
		if errors.Is(err, ErrObservabilityDisabled) {
			return utils.Error(c, http.StatusServiceUnavailable, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Historical metrics fetched", data)
}

type LogHandler struct {
	service LogOperations
}

func NewLogHandler(s LogOperations) *LogHandler {
	return &LogHandler{service: s}
}

func (h *LogHandler) GetHistoricalLogs(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	rangeParam := c.QueryParam("range")
	if rangeParam == "" {
		rangeParam = "24h"
	}

	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 1000
	}

	logs, err := h.service.GetHistoricalLogs(c.Request().Context(), serviceID, rangeParam, limit)
	if err != nil {
		if errors.Is(err, ErrObservabilityDisabled) {
			return utils.Error(c, http.StatusServiceUnavailable, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Historical logs fetched", logs)
}
