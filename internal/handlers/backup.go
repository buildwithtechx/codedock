package handlers

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
)

type BackupHandler struct {
	backupService  *services.BackupService
	appService     *services.AppService
	dbService      *services.DatabaseService
	projectService *services.ProjectService
}

func NewBackupHandler(backup *services.BackupService, app *services.AppService, db *services.DatabaseService, proj *services.ProjectService) *BackupHandler {
	return &BackupHandler{
		backupService:  backup,
		appService:     app,
		dbService:      db,
		projectService: proj,
	}
}

func (h *BackupHandler) hasAccess(c echo.Context, databaseID, serviceID string) bool {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return false
	}
	if userClaims.Role == models.UserRoleAdmin || userClaims.Role == models.UserRoleOwner {
		return true
	}

	var projectID string
	if databaseID != "" {
		db, err := h.dbService.GetDatabase(c.Request().Context(), databaseID)
		if err == nil && db != nil {
			projectID = db.ProjectID
		}
	} else if serviceID != "" {
		app, err := h.appService.GetAppService(c.Request().Context(), serviceID)
		if err == nil && app != nil {
			projectID = app.ProjectID
		}
	}
	if projectID == "" {
		return false
	}

	return h.projectService.HasPermission(c.Request().Context(), projectID, userClaims.UserID, userClaims.Role, "")
}

func (h *BackupHandler) List(c echo.Context) error {
	list, err := h.backupService.ListConfigs(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	var filtered []*models.BackupConfig
	for _, cfg := range list {
		if h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			cfg.DbPassword = "********"
			filtered = append(filtered, cfg)
		}
	}
	return utils.Success(c, "Operation successful", filtered)
}

func (h *BackupHandler) Create(c echo.Context) error {
	var cfg models.BackupConfig
	if err := c.Bind(&cfg); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions to create backup for this resource")
	}

	if err := h.backupService.CreateConfig(c.Request().Context(), &cfg); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", cfg)
}

func (h *BackupHandler) Update(c echo.Context) error {
	id := c.Param("id")

	existing, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || existing == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup config")
		}
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}

	if !h.hasAccess(c, existing.DatabaseID, existing.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	var req struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		DbUser          string `json:"dbUser"`
		DbPassword      string `json:"dbPassword"`
		Schedule        string `json:"schedule"`
		Timezone        string `json:"timezone"`
		Timeout         int    `json:"timeout"`
		RetentionDays   int    `json:"retentionDays"`
		MaxBackups      int    `json:"maxBackups"`
		MaxStorageGB    int    `json:"maxStorageGB"`
		S3DestinationID string `json:"s3DestinationId"`
		DatabaseID      string `json:"databaseId"`
		BackupEnabled   *bool  `json:"backupEnabled"`
		S3Enabled       *bool  `json:"s3Enabled"`
		DisableLocal    *bool  `json:"disableLocal"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.DbUser != "" {
		existing.DbUser = req.DbUser
	}
	if req.DbPassword != "" {
		existing.DbPassword = req.DbPassword
	}
	if req.Schedule != "" {
		existing.Schedule = req.Schedule
	}
	if req.Timezone != "" {
		existing.Timezone = req.Timezone
	}
	if req.Timeout != 0 {
		existing.Timeout = req.Timeout
	}
	if req.RetentionDays != 0 {
		existing.RetentionDays = req.RetentionDays
	}
	if req.MaxBackups != 0 {
		existing.MaxBackups = req.MaxBackups
	}
	if req.MaxStorageGB != 0 {
		existing.MaxStorageGB = req.MaxStorageGB
	}
	if req.S3DestinationID != "" {
		existing.S3DestinationID = req.S3DestinationID
	}
	if req.DatabaseID != "" {
		existing.DatabaseID = req.DatabaseID
	}

	if req.BackupEnabled != nil {
		existing.BackupEnabled = *req.BackupEnabled
	}
	if req.S3Enabled != nil {
		existing.S3Enabled = *req.S3Enabled
	}
	if req.DisableLocal != nil {
		existing.DisableLocal = *req.DisableLocal
	}

	if err := h.backupService.UpdateConfig(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	existing.DbPassword = "********"
	return utils.Success(c, "Updated successfully", existing)
}

func (h *BackupHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		var notFoundErr *utils.NotFoundError
		if err != nil && !errors.As(err, &notFoundErr) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup config")
		}
		return utils.Error(c, http.StatusNotFound, "backup config not found")
	}

	if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}
	cfg.DbPassword = "********"
	return utils.Success(c, "Operation successful", cfg)
}

func (h *BackupHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err == nil && cfg != nil {
		if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	}

	if err := h.backupService.DeleteConfig(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *BackupHandler) Trigger(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err == nil && cfg != nil {
		if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	}

	rec, err := h.backupService.TriggerBackup(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", rec)
}

func (h *BackupHandler) ListRecords(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err == nil && cfg != nil {
		if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
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
	if err == nil && cfg != nil {
		if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
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
	if err == nil && cfg != nil {
		if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
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

func (h *BackupHandler) ListS3Destinations(c echo.Context) error {
	list, err := h.backupService.ListS3Destinations(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *BackupHandler) CreateS3Destination(c echo.Context) error {
	var dest models.S3Destination
	if err := c.Bind(&dest); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.backupService.CreateS3Destination(c.Request().Context(), &dest); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", dest)
}

func (h *BackupHandler) DeleteS3Destination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}
	if err := h.backupService.DeleteS3Destination(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
