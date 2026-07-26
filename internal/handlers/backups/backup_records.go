package backups

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"
)

func (h *BackupHandler) ListRecords(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}
	if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	recs, err := h.backupService.ListRecordsByConfig(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", recs)
}

func (h *BackupHandler) DownloadRecord(c echo.Context) error {
	id := c.Param("id")
	recordID := c.Param("recordId")
	if id == "" || recordID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id or recordId parameter")
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient admin permissions to download backup record")
	}

	rec, err := h.backupService.GetRecord(c.Request().Context(), recordID)
	if err != nil {
		var notFound *utils.NotFoundError
		if !errors.As(err, &notFound) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup record")
		}
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec == nil {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec.BackupConfigID != id {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}

	if rec.FilePath == "" {
		return utils.Error(c, http.StatusNotFound, "local backup file not available")
	}
	return c.Attachment(rec.FilePath, filepath.Base(rec.FilePath))
}

func (h *BackupHandler) DeleteRecord(c echo.Context) error {
	id := c.Param("id")
	recordID := c.Param("recordId")
	if id == "" || recordID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id or recordId parameter")
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient admin permissions to delete backup record")
	}

	rec, err := h.backupService.GetRecord(c.Request().Context(), recordID)
	if err != nil {
		var notFound *utils.NotFoundError
		if !errors.As(err, &notFound) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup record")
		}
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec == nil {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}
	if rec.BackupConfigID != id {
		return utils.Error(c, http.StatusNotFound, "record not found")
	}

	if err := h.backupService.DeleteRecord(c.Request().Context(), recordID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *BackupHandler) Restore(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id parameter")
	}

	rec, err := h.backupService.GetRecord(c.Request().Context(), id)
	if err != nil {
		var notFound *utils.NotFoundError
		if !errors.As(err, &notFound) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup record: "+err.Error())
		}
		return utils.Error(c, http.StatusNotFound, "backup record not found")
	}
	if rec == nil {
		return utils.Error(c, http.StatusNotFound, "backup record not found")
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), rec.BackupConfigID)
	if err != nil {
		var notFound *utils.NotFoundError
		if !errors.As(err, &notFound) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup config: "+err.Error())
		}
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}
	if cfg == nil {
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}

	if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions to restore this backup")
	}

	err = h.backupService.RestoreBackup(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Backup successfully restored", nil)
}
