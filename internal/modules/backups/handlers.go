package backups

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ProjectResolver func(ctx context.Context, databaseID, serviceID string) string
type PermissionChecker func(ctx context.Context, projectID, userID, role, permission string) bool
type AuditLogger func(ctx context.Context, userID, action, resource, ip string)

type Handler struct {
	backupService     *Service
	projectResolver   ProjectResolver
	permissionChecker PermissionChecker
	auditLogger       AuditLogger
	operations        any
}

func NewHandler(backup *Service) *Handler {
	h := &Handler{
		backupService: backup,
	}
	if backup != nil && backup.backupRepo != nil {
		if pr, ok := backup.backupRepo.(interface {
			ResolveProjectID(context.Context, string, string) string
		}); ok {
			h.projectResolver = pr.ResolveProjectID
		}
		if pc, ok := backup.backupRepo.(interface {
			HasProjectPermission(context.Context, string, string, string, string) bool
		}); ok {
			h.permissionChecker = pc.HasProjectPermission
		}
		if al, ok := backup.backupRepo.(interface {
			LogAudit(context.Context, string, string, string, string)
		}); ok {
			h.auditLogger = al.LogAudit
		}
	}
	return h
}

func (h *Handler) SetSecurity(
	projRes ProjectResolver,
	permCheck PermissionChecker,
	audit AuditLogger,
) *Handler {
	h.projectResolver = projRes
	h.permissionChecker = permCheck
	h.auditLogger = audit
	return h
}

func (h *Handler) SetOperations(ops any) {
	h.operations = ops
}

func (h *Handler) List(c echo.Context) error {
	list, err := h.backupService.ListConfigs(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	var filtered []*BackupConfig
	for _, cfg := range list {
		if h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
			RedactBackupConfig(cfg)
			filtered = append(filtered, cfg)
		}
	}
	return utils.Success(c, "Operation successful", filtered)
}

func (h *Handler) Create(c echo.Context) error {
	var cfg BackupConfig
	if err := c.Bind(&cfg); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidPayload.Error())
	}

	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}

	if err := h.validatePolicy(c, &cfg); err != nil {
		return err
	}
	if err := h.backupService.CreateConfig(c.Request().Context(), &cfg); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), utils.AuditActorID(c), "backup.create", cfg.ID, c.RealIP())
	}
	RedactBackupConfig(&cfg)
	return utils.Created(c, "Created successfully", cfg)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")

	existing, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || existing == nil {
		if err != nil && !utils.IsNotFound(err) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup config")
		}
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}

	if !h.hasAdminAccess(c, existing.DatabaseID, existing.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}

	var req UpdateBackupRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidPayload.Error())
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
	if req.Timeout != nil {
		existing.Timeout = *req.Timeout
	}
	if req.RetentionDays != nil {
		existing.RetentionDays = *req.RetentionDays
	}
	if req.MaxBackups != nil {
		existing.MaxBackups = *req.MaxBackups
	}
	if req.MaxStorageGB != nil {
		existing.MaxStorageGB = *req.MaxStorageGB
	}
	if req.S3DestinationID != "" {
		existing.S3DestinationID = req.S3DestinationID
	}
	if req.DatabaseID != "" && req.DatabaseID != existing.DatabaseID {
		return utils.Error(c, http.StatusBadRequest, ErrProducerImmutable.Error())
	}

	if req.PreDeployment != nil {
		existing.PreDeployment = *req.PreDeployment
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

	if err := h.validatePolicy(c, existing); err != nil {
		return err
	}
	if err := h.backupService.UpdateConfig(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	RedactBackupConfig(existing)
	return utils.Success(c, "Updated successfully", existing)
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		if err != nil && !utils.IsNotFound(err) {
			return utils.Error(c, http.StatusInternalServerError, "failed to get backup config")
		}
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}

	if !h.hasAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientPermissions.Error())
	}
	RedactBackupConfig(cfg)
	return utils.Success(c, "Operation successful", cfg)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}

	if err := h.backupService.DeleteConfig(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), utils.AuditActorID(c), "backup.delete", id, c.RealIP())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Trigger(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}

	rec, err := h.backupService.TriggerBackup(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), utils.AuditActorID(c), "backup.trigger", id, c.RealIP())
	}
	return utils.Success(c, "Operation successful", rec)
}

func (h *Handler) Restore(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}

	cfg, err := h.backupService.GetConfig(c.Request().Context(), id)
	if err != nil || cfg == nil {
		return utils.Error(c, http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return utils.Error(c, http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}

	if err := h.backupService.RestoreBackup(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), utils.AuditActorID(c), "backup.restore", id, c.RealIP())
	}
	return utils.Success(c, "Restore initiated", nil)
}
