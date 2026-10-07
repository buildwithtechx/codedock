package system

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services/analytics"
	"codedock.run/codedock/internal/utils"
)

type AnalyticsHandler struct {
	analytics *analytics.Service
}

func NewAnalyticsHandler(service *analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{analytics: service}
}

func analyticsClaims(c echo.Context) (*models.UserClaims, error) {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return nil, utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	return userClaims, nil
}

func (h *AnalyticsHandler) Summary(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	summary, err := h.analytics.Summary(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"), c.QueryParam("domain"), c.QueryParam("from"), c.QueryParam("to"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Traffic summary retrieved", summary)
}

func (h *AnalyticsHandler) Overview(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	overview, err := h.analytics.Overview(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"), c.QueryParam("domain"), c.QueryParam("from"), c.QueryParam("to"), c.QueryParam("step"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Traffic overview retrieved", overview)
}

func (h *AnalyticsHandler) Geo(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	geo, err := h.analytics.Geo(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"), c.QueryParam("from"), c.QueryParam("to"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Traffic geography retrieved", geo)
}

func (h *AnalyticsHandler) GetPaths(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	enabled, err := h.analytics.PathsEnabled(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Paths collection retrieved", map[string]bool{"enabled": enabled})
}

func (h *AnalyticsHandler) SetPaths(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if err := h.analytics.SetPathsEnabled(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"), req.Enabled); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Paths collection updated", map[string]bool{"enabled": req.Enabled})
}

func (h *AnalyticsHandler) Periods(c echo.Context) error {
	if _, err := analyticsClaims(c); err != nil {
		return err
	}
	return utils.Success(c, "Periods retrieved", h.analytics.Periods())
}

func (h *AnalyticsHandler) DeploymentStats(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	days, _ := strconv.Atoi(c.QueryParam("days"))
	stats, err := h.analytics.DeploymentStats(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"), days)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Deployment stats retrieved", stats)
}

func (h *AnalyticsHandler) Dashboard(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	dashboard, err := h.analytics.Dashboard(c.Request().Context(), claims.UserID, c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Dashboard retrieved", dashboard)
}

func (h *AnalyticsHandler) Usage(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.requireProjectRead(c, claims); err != nil {
		return err
	}
	usage, err := h.analytics.Usage().Current(c.Request().Context(), c.Param("projectId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Usage retrieved", usage)
}

func (h *AnalyticsHandler) Resources(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.requireProjectRead(c, claims); err != nil {
		return err
	}
	resources, err := h.analytics.Usage().Resources(c.Request().Context(), c.Param("projectId"))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Resources retrieved", resources)
}

func (h *AnalyticsHandler) UsageHistory(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.requireProjectRead(c, claims); err != nil {
		return err
	}
	hours, _ := strconv.Atoi(c.QueryParam("hours"))
	history, err := h.analytics.Usage().History(c.Request().Context(), c.QueryParam("serviceId"), hours)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Usage history retrieved", history)
}

func (h *AnalyticsHandler) Container(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.requireProjectRead(c, claims); err != nil {
		return err
	}
	entry, err := h.analytics.Usage().Container(c.Request().Context(), c.Param("projectId"), c.QueryParam("serviceId"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Container usage retrieved", entry)
}

func (h *AnalyticsHandler) UsageStream(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	if err := h.requireProjectRead(c, claims); err != nil {
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
			usage, err := h.analytics.Usage().Current(c.Request().Context(), projectID)
			if err != nil {
				writeStreamEvent(writer, flush, "ERROR", err.Error())
				return nil
			}
			encoded, _ := json.Marshal(usage)
			writeStreamEvent(writer, flush, "USAGE", string(encoded))
		}
	}
}

func (h *AnalyticsHandler) requireProjectRead(c echo.Context, claims *models.UserClaims) error {
	_, err := h.analytics.Summary(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"), "", "", "")
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return nil
}

func (h *AnalyticsHandler) Live(c echo.Context) error {
	claims, err := analyticsClaims(c)
	if err != nil {
		return err
	}
	to := time.Now().UTC().Format("2006-01-02T15:04")
	from := time.Now().UTC().Add(-15 * time.Minute).Format("2006-01-02T15:04")
	overview, err := h.analytics.Overview(c.Request().Context(), claims.UserID, c.Param("id"), c.Param("projectId"), c.QueryParam("domain"), from, to, "1")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Live traffic retrieved", overview)
}
