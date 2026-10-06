package backups

import (
	"codedock.run/codedock/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *BackupHandler) validatePolicy(c echo.Context, cfg *models.BackupConfig) error {
	if cfg.PreDeployment && cfg.ServiceID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "pre-deployment policy requires an application service")
	}
	if cfg.Timeout < 0 || cfg.Timeout > 86400 || cfg.RetentionDays < 0 || cfg.MaxBackups < 0 || cfg.MaxStorageGB < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid backup timeout or retention limits")
	}
	if cfg.DisableLocal && !cfg.S3Enabled {
		return echo.NewHTTPError(http.StatusBadRequest, "remote-only backups require an enabled storage destination")
	}
	if cfg.S3Enabled && cfg.S3DestinationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "select a storage destination")
	}
	if cfg.ServiceID == "" {
		return nil
	}
	if !h.hasAdminAccess(c, "", cfg.ServiceID) {
		return echo.NewHTTPError(http.StatusForbidden, "backup service permission required")
	}
	if cfg.DatabaseID == "" {
		if cfg.VolumeName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "service backup requires a database or named volume")
		}
		return nil
	}
	if cfg.VolumeName != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "select one backup producer")
	}
	database, err := h.dbService.GetDatabase(c.Request().Context(), cfg.DatabaseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "database not found")
	}
	service, err := h.appService.GetAppService(c.Request().Context(), cfg.ServiceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "application not found")
	}
	if database.ProjectID != service.ProjectID || !h.hasAdminAccess(c, database.ID, "") {
		return echo.NewHTTPError(http.StatusForbidden, "backup producer must be an authorized database in the service project")
	}
	return nil
}
