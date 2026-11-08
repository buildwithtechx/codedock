package backups

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *Handler) getProjectID(ctx context.Context, databaseID, serviceID string) string {
	if h.projectResolver != nil {
		return h.projectResolver(ctx, databaseID, serviceID)
	}
	return ""
}

func (h *Handler) hasAccess(c echo.Context, databaseID, serviceID string) bool {
	userID, role, ok := utils.ExtractUser(c)
	if !ok {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	projectID := h.getProjectID(c.Request().Context(), databaseID, serviceID)
	if projectID == "" || h.permissionChecker == nil {
		return false
	}
	return h.permissionChecker(c.Request().Context(), projectID, userID, role, "")
}

func (h *Handler) hasAdminAccess(c echo.Context, databaseID, serviceID string) bool {
	userID, role, ok := utils.ExtractUser(c)
	if !ok {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	projectID := h.getProjectID(c.Request().Context(), databaseID, serviceID)
	if projectID == "" || h.permissionChecker == nil {
		return false
	}
	return h.permissionChecker(c.Request().Context(), projectID, userID, role, "admin")
}

func (h *Handler) backupProject(c echo.Context, cfg *BackupConfig) (string, error) {
	projectID := h.getProjectID(c.Request().Context(), cfg.DatabaseID, cfg.ServiceID)
	if projectID == "" && (cfg.ServiceID != "" || cfg.DatabaseID != "") {
		return "", echo.NewHTTPError(http.StatusNotFound, "resource not found")
	}
	return projectID, nil
}

func (h *Handler) validatePolicy(c echo.Context, cfg *BackupConfig) error {
	projectID, err := h.backupProject(c, cfg)
	if err != nil {
		return err
	}
	cfg.ProjectID = projectID

	userID, _, _ := utils.ExtractUser(c)
	if cfg.OwnerID == "" {
		cfg.OwnerID = userID
	}

	if cfg.PreDeployment && cfg.ServiceID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "pre-deployment backup policies are only supported on application services")
	}

	if cfg.S3Enabled && cfg.S3DestinationID != "" {
		dest, err := h.backupService.GetS3Destination(c.Request().Context(), cfg.S3DestinationID)
		if err != nil || dest == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid s3 destination")
		}
	}

	if cfg.SFTPEnabled && cfg.SFTPDestinationID != "" {
		dest, err := h.backupService.GetSFTPDestination(c.Request().Context(), cfg.SFTPDestinationID)
		if err != nil || dest == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid sftp destination")
		}
		if projectID != "" && dest.ProjectID != "" && dest.ProjectID != projectID {
			return echo.NewHTTPError(http.StatusForbidden, "sftp destination belongs to another project")
		}
	}

	if cfg.ParentBatchID != "" {
		batch, err := h.backupService.GetBatch(c.Request().Context(), cfg.ParentBatchID)
		if err != nil || batch == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid parent batch")
		}
		if projectID != "" && batch.ProjectID != "" && batch.ProjectID != projectID {
			return echo.NewHTTPError(http.StatusForbidden, "policy batch belongs to another project")
		}
		cfg.Schedule = batch.Schedule
		cfg.Timezone = batch.Timezone
		cfg.Timeout = batch.Timeout
	}

	return nil
}
