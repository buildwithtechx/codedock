package monitoring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AnalyticsHandler struct {
	service AnalyticsService
}

func NewAnalyticsHandler(service AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func extractUser(c echo.Context) (string, error) {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return "", utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	return userID, nil
}

func (h *AnalyticsHandler) Periods(c echo.Context) error {
	periods, err := h.service.Periods(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Aggregation periods retrieved", periods)
}

func (h *AnalyticsHandler) Dashboard(c echo.Context) error {
	dashboard, err := h.service.Dashboard(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Analytics dashboard retrieved", dashboard)
}

func (h *AnalyticsHandler) Summary(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	summary, err := h.service.Summary(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("from"), c.QueryParam("to"), c.QueryParam("period"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Analytics summary retrieved", summary)
}

func (h *AnalyticsHandler) Overview(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	overview, err := h.service.Overview(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("domain"), c.QueryParam("from"), c.QueryParam("to"), c.QueryParam("period"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Analytics overview retrieved", overview)
}

func (h *AnalyticsHandler) Geo(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	geo, err := h.service.Geo(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("from"), c.QueryParam("to"), c.QueryParam("period"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Geographic analytics retrieved", geo)
}

func (h *AnalyticsHandler) Live(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	to := time.Now().UTC().Format("2006-01-02T15:04")
	from := time.Now().UTC().Add(-15 * time.Minute).Format("2006-01-02T15:04")
	overview, err := h.service.Overview(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("domain"), from, to, "1")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Live traffic retrieved", overview)
}

func (h *AnalyticsHandler) GetPaths(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	enabled, err := h.service.GetPaths(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Path analytics setting retrieved", map[string]bool{"enabled": enabled})
}

func (h *AnalyticsHandler) SetPaths(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	var req PathsToggleRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	if err := h.service.SetPaths(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), req.Enabled); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Path analytics setting updated", map[string]bool{"enabled": req.Enabled})
}

func (h *AnalyticsHandler) DeploymentStats(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	stats, err := h.service.DeploymentStats(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Deployment statistics retrieved", stats)
}

func (h *AnalyticsHandler) Usage(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	usage, err := h.service.Usage(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("period"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Project resource usage retrieved", usage)
}

func (h *AnalyticsHandler) Resources(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	resources, err := h.service.Resources(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("period"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Detailed resource consumption retrieved", resources)
}

func (h *AnalyticsHandler) UsageHistory(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	history, err := h.service.UsageHistory(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("period"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Historical resource usage retrieved", history)
}

func (h *AnalyticsHandler) Container(c echo.Context) error {
	userID, err := extractUser(c)
	if err != nil {
		return err
	}
	stats, err := h.service.Container(c.Request().Context(), userID, c.Param("id"), c.Param("projectId"), c.QueryParam("serviceId"), c.QueryParam("period"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Container telemetry retrieved", stats)
}

func (h *AnalyticsHandler) UsageStream(c echo.Context) error {
	_, err := extractUser(c)
	if err != nil {
		return err
	}
	writer, flush, ok := startEventStream(c)
	if !ok {
		return utils.Error(c, http.StatusInternalServerError, "streaming unsupported")
	}
	projectID := c.Param("projectId")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-ticker.C:
			usage, err := h.service.CurrentUsage(c.Request().Context(), projectID)
			if err != nil {
				writeStreamEvent(writer, flush, "ERROR", err.Error())
				return nil
			}
			encoded, _ := json.Marshal(usage)
			writeStreamEvent(writer, flush, "USAGE", string(encoded))
		}
	}
}

func startEventStream(c echo.Context) (http.ResponseWriter, http.Flusher, bool) {
	writer := c.Response().Writer
	flush, ok := writer.(http.Flusher)
	if !ok {
		return nil, nil, false
	}
	header := writer.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	return writer, flush, true
}

func writeStreamEvent(writer http.ResponseWriter, flush http.Flusher, event, data string) {
	fmt.Fprintf(writer, "event: %s\ndata: %s\n\n", event, data)
	flush.Flush()
}
