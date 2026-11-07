package system

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	systemservices "codedock/internal/services/system"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
)

type LogHandler struct {
	logService *systemservices.LogService
}

func NewLogHandler(ls *systemservices.LogService) *LogHandler {
	return &LogHandler{logService: ls}
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

	end := time.Now()
	var start time.Time

	switch rangeParam {
	case "7d":
		start = end.Add(-7 * 24 * time.Hour)
	case "24h":
		start = end.Add(-24 * time.Hour)
	case "1h":
		start = end.Add(-1 * time.Hour)
	default:
		start = end.Add(-24 * time.Hour)
	}

	opts := systemservices.HistoricalLogsOpts{
		ServiceID: serviceID,
		Start:     start,
		End:       end,
		Limit:     limit,
	}
	logs, err := h.logService.GetHistoricalLogs(c.Request().Context(), opts)
	if err != nil {
		if errors.Is(err, systemservices.ErrObservabilityDisabled) {
			return utils.Error(c, http.StatusServiceUnavailable, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Historical logs fetched", logs)
}
