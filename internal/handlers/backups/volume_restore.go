package backups

import (
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *BackupHandler) authorizeVolumeRestore(c echo.Context) (string, error) {
	id := c.Param("recordId")
	record, err := h.backupService.GetRecord(c.Request().Context(), id)
	if err != nil && !utils.IsNotFound(err) {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "failed to load backup record")
	}
	if err != nil || record == nil {
		return "", echo.NewHTTPError(http.StatusNotFound, "backup record not found")
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), record.BackupConfigID)
	if err != nil && !utils.IsNotFound(err) {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "failed to load backup configuration")
	}
	if err != nil || cfg == nil {
		return "", echo.NewHTTPError(http.StatusNotFound, "backup configuration not found")
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return "", echo.NewHTTPError(http.StatusForbidden, "admin access required for volume restoration")
	}
	return id, nil
}

func (h *BackupHandler) VolumeRestoreTarget(c echo.Context) error {
	id, err := h.authorizeVolumeRestore(c)
	if err != nil {
		return err
	}
	target, err := h.backupService.ValidateVolumeRestore(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Volume restore target validated", target)
}

func (h *BackupHandler) RestoreVolume(c echo.Context) error {
	return echo.NewHTTPError(http.StatusGone, "Use the reviewed restore prepare/apply workflow")
}

func (h *BackupHandler) CancelVolumeRestore(c echo.Context) error {
	id, err := h.authorizeVolumeRestore(c)
	if err != nil {
		return err
	}
	if !h.backupService.CancelVolumeRestore(id) {
		return utils.Error(c, http.StatusConflict, "no active volume restore for this record")
	}
	return utils.Success(c, "Volume restore interruption requested; existing data may be partially overwritten", nil)
}
