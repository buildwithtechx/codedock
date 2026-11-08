package backups

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *Handler) ListAllRecords(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := DefaultLimit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	_, role, ok := extractUser(c)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, ErrUnauthorized.Error())
	}
	var records []*BackupRecord
	var err error
	if role == "admin" || role == "owner" {
		records, err = h.backupService.ListAllRecords(c.Request().Context(), limit)
	} else {
		configs, configErr := h.backupService.ListConfigs(c.Request().Context())
		if configErr != nil {
			return utils.Error(c, http.StatusInternalServerError, "failed to check backup access")
		}
		access := make(map[string]bool)
		ids := make([]string, 0, len(configs))
		for _, cfg := range configs {
			key := cfg.DatabaseID + ":" + cfg.ServiceID
			allowed, found := access[key]
			if !found {
				allowed = h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID)
				access[key] = allowed
			}
			if allowed {
				ids = append(ids, cfg.ID)
			}
		}
		records, err = h.backupService.ListRecordsByConfigs(c.Request().Context(), ids, limit)
	}
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to list backup records")
	}
	return utils.Success(c, "Operation successful", records)
}

func (h *Handler) ListRecords(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}
	if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientPermissions.Error())
	}
	records, err := h.backupService.ListRecords(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", records)
}

func (h *Handler) DownloadRecord(c echo.Context) error {
	recordID := c.Param("recordId")
	if recordID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id parameter")
	}

	record, err := h.backupService.GetRecord(c.Request().Context(), recordID)
	if err != nil || record == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupRecordNotFound.Error())
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), record.BackupConfigID)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}

	if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientPermissions.Error())
	}

	if record.FilePath != "" {
		if fi, err := os.Stat(record.FilePath); err == nil && !fi.IsDir() {
			filename := filepath.Base(record.FilePath)
			contentType := mime.TypeByExtension(filepath.Ext(filename))
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=\""+filename+"\"")
			c.Response().Header().Set(echo.HeaderContentType, contentType)
			return c.File(record.FilePath)
		}
	}

	url, err := h.backupService.GetRecordDownloadURL(c.Request().Context(), recordID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if url != "" {
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}

	return utils.Error(c, http.StatusNotFound, "backup file not found")
}

func (h *Handler) DeleteRecord(c echo.Context) error {
	recordID := c.Param("recordId")
	if recordID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id parameter")
	}

	record, err := h.backupService.GetRecord(c.Request().Context(), recordID)
	if err != nil || record == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupRecordNotFound.Error())
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), record.BackupConfigID)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}

	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}

	if err := h.backupService.DeleteRecord(c.Request().Context(), recordID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.record_delete", recordID, c.RealIP())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListRecordsByDatabase(c echo.Context) error {
	dbID := c.Param("id")
	if dbID == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}

	if !h.hasAccess(c, dbID, "") {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientPermissions.Error())
	}

	records, err := h.backupService.ListRecordsByDatabase(c.Request().Context(), dbID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", records)
}

func (h *Handler) TriggerDatabaseBackup(c echo.Context) error {
	dbID := c.Param("id")
	if dbID == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}

	if !h.hasAdminAccess(c, dbID, "") {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}

	cfg, err := h.backupService.backupRepo.GetConfigByDatabaseID(c.Request().Context(), dbID)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, "no backup configuration found for database")
	}

	rec, err := h.backupService.TriggerBackup(c.Request().Context(), cfg.ID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.trigger", cfg.ID, c.RealIP())
	}
	return utils.Success(c, "Backup triggered successfully", rec)
}
