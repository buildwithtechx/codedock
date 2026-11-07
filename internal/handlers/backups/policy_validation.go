package backups

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func (h *BackupHandler) validatePolicy(c echo.Context, cfg *models.BackupConfig) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	project, err := h.backupProject(c, cfg)
	if err != nil {
		return err
	}
	cfg.OwnerID = user.UserID
	cfg.ProjectID = project
	if strings.TrimSpace(cfg.Name) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "backup policy name is required")
	}

	if cfg.PreDeployment && cfg.ProjectID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "pre-deployment policy requires a project resource")
	}
	if cfg.Timeout < 0 || cfg.Timeout > 86400 || cfg.RetentionDays < 0 || cfg.MaxBackups < 0 || cfg.MaxStorageGB < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid backup timeout or retention limits")
	}
	if cfg.DisableLocal && !cfg.S3Enabled && !cfg.SFTPEnabled {
		return echo.NewHTTPError(http.StatusBadRequest, "remote-only backups require an enabled storage destination")
	}
	if cfg.S3Enabled && cfg.S3DestinationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "select a storage destination")
	}
	if cfg.SFTPEnabled && cfg.SFTPDestinationID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "select an SFTP destination")
	}
	if cfg.Incremental && cfg.VolumeName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "incremental backups require a named volume producer")
	}
	if (cfg.QuiesceCommand != "" || cfg.UnquiesceCommand != "") && cfg.ServiceID == "" && cfg.DatabaseID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "quiesce commands require a service or database producer")
	}
	if cfg.FileSourcePath != "" {
		if cfg.DatabaseID != "" || cfg.VolumeName != "" || cfg.CustomBackupCommand != "" {
			return echo.NewHTTPError(http.StatusBadRequest, "select one backup producer")
		}
		if len(cfg.FileSourcePath) > 512 || strings.Contains(cfg.FileSourcePath, "..") {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid file source path")
		}
		return nil
	}
	if cfg.CustomBackupCommand != "" {
		if cfg.ServiceID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "custom producers require a service container")
		}
		if !h.hasAdminAccess(c, "", cfg.ServiceID) {
			return echo.NewHTTPError(http.StatusForbidden, "backup service permission required")
		}
		return nil
	}
	if cfg.ParentBatchID != "" && cfg.ServiceID != "" && cfg.DatabaseID == "" && cfg.VolumeName == "" {
		if !h.hasAdminAccess(c, "", cfg.ServiceID) {
			return echo.NewHTTPError(http.StatusForbidden, "backup service permission required")
		}
		return nil
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
