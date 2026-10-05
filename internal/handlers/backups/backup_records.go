package backups

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func (h *BackupHandler) ListAllRecords(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	recs, err := h.backupService.ListAllRecords(c.Request().Context(), limit)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	filtered := make([]*models.BackupRecord, 0, len(recs))
	for _, record := range recs {
		cfg, err := h.backupService.GetConfig(c.Request().Context(), record.BackupConfigID)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, "failed to check backup access")
		}
		if cfg != nil && h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			filtered = append(filtered, record)
		}
	}
	return utils.Success(c, "Operation successful", filtered)
}

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

	var rec *models.BackupRecord
	var cfg *models.BackupConfig

	directRec, err := h.backupService.GetRecord(c.Request().Context(), id)
	if err == nil && directRec != nil {
		rec = directRec
		cfg, _ = h.backupService.GetConfig(c.Request().Context(), rec.BackupConfigID)
	} else {
		directCfg, cfgErr := h.backupService.GetConfig(c.Request().Context(), id)
		if cfgErr == nil && directCfg != nil {
			cfg = directCfg
			records, _ := h.backupService.ListRecordsByConfig(c.Request().Context(), id)
			for _, r := range records {
				if r.Status == models.BackupRecordStatusCompleted {
					rec = r
					break
				}
			}
		}
	}

	if rec == nil {
		return utils.Error(c, http.StatusNotFound, "backup record not found")
	}
	if cfg == nil {
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}

	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient admin permissions to restore this backup")
	}

	err = h.backupService.RestoreBackup(c.Request().Context(), rec.ID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Backup successfully restored", nil)
}

func (h *BackupHandler) ListRecordsByDatabase(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id")
	}
	if !h.hasAccess(c, id, "") {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}
	records, err := h.backupService.ListRecordsByDatabase(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", records)
}

func (h *BackupHandler) TriggerDatabaseBackup(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing database id")
	}
	if !h.hasAdminAccess(c, id, "") {
		return utils.Error(c, http.StatusForbidden, "insufficient admin permissions")
	}
	cfg, err := h.backupService.GetConfigByDatabaseID(c.Request().Context(), id)
	if err != nil || cfg == nil {
		db, dbErr := h.dbService.GetDatabase(c.Request().Context(), id)
		if dbErr != nil || db == nil {
			return utils.Error(c, http.StatusNotFound, "database not found")
		}
		newCfg := models.BackupConfig{
			DatabaseID:    id,
			Name:          db.Name + "-backup",
			Description:   "Manual snapshots for " + db.Name,
			BackupEnabled: true,
			Schedule:      "manual",
			RetentionDays: 7,
			Status:        models.BackupConfigStatusInactive,
		}
		destinations, _ := h.backupService.ListS3Destinations(c.Request().Context())
		for _, d := range destinations {
			if d.IsDefault {
				newCfg.S3DestinationID = d.ID
				newCfg.S3Enabled = true
				break
			}
		}
		if err := h.backupService.CreateConfig(c.Request().Context(), &newCfg); err != nil {
			return utils.Error(c, http.StatusInternalServerError, "failed to create backup configuration: "+err.Error())
		}
		cfg = &newCfg
	}
	rec, err := h.backupService.TriggerBackup(c.Request().Context(), cfg.ID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Backup triggered successfully", rec)
}
