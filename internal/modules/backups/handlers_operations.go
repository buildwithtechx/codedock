package backups

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *Handler) ListBatches(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	batches, err := h.backupService.ListBatches(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", batches)
}

func (h *Handler) CreateBatch(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	var req CreatePolicyBatchRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidPayload.Error())
	}
	batch := BackupPolicyBatch{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Schedule:    req.Schedule,
		Timezone:    req.Timezone,
		Timeout:     req.Timeout,
	}
	if err := h.backupService.CreateBatch(c.Request().Context(), &batch); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.batch_create", batch.ID, c.RealIP())
	}
	return utils.Created(c, "Backup batch created successfully", batch)
}

func (h *Handler) TriggerBatch(c echo.Context) error {
	batchID := c.Param("batchId")
	if batchID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing batch id parameter")
	}
	summary, err := h.backupService.TriggerBatch(c.Request().Context(), batchID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.batch_trigger", batchID, c.RealIP())
	}
	return utils.Success(c, "Backup batch triggered", summary)
}

func (h *Handler) authorizedRecord(c echo.Context, id string) (*BackupConfig, string, error) {
	record, err := h.backupService.GetRecord(c.Request().Context(), id)
	if err != nil {
		return nil, "", echo.NewHTTPError(http.StatusNotFound, "backup record not found")
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), record.BackupConfigID)
	if err != nil {
		return nil, "", echo.NewHTTPError(http.StatusNotFound, "backup policy not found")
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return nil, "", echo.NewHTTPError(http.StatusForbidden, "backup admin permission required")
	}
	project, err := h.backupProject(c, cfg)
	return cfg, project, err
}

func (h *Handler) StartRun(c echo.Context) error {
	userID, _, ok := extractUser(c)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrUnauthorized.Error())
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return echo.NewHTTPError(http.StatusForbidden, ErrInsufficientAdminPermissions.Error())
	}
	project, err := h.backupProject(c, cfg)
	if err != nil {
		return err
	}
	op, err := h.backupService.StartRun(c.Request().Context(), userID, project, cfg.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]any{"status": "success", "data": op})
}

func (h *Handler) PrepareRestore(c echo.Context) error {
	userID, _, ok := extractUser(c)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrUnauthorized.Error())
	}
	var req RestoreReviewRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restore target")
	}
	_, project, err := h.authorizedRecord(c, c.Param("recordId"))
	if err != nil {
		return err
	}
	if req.Database != "" {
		targetProj := h.getProjectID(c.Request().Context(), req.Database, "")
		if targetProj != project || !h.hasAdminAccess(c, req.Database, "") {
			return echo.NewHTTPError(http.StatusForbidden, "restore target must be an authorized database in the same project")
		}
	}
	review, err := h.backupService.PrepareRestore(c.Request().Context(), userID, project, c.Param("recordId"), req.Database)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review restore effects before applying", review)
}

func (h *Handler) ApplyRestore(c echo.Context) error {
	userID, _, ok := extractUser(c)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrUnauthorized.Error())
	}
	ops, ok := h.operations.(interface {
		Get(ctx context.Context, id string) (any, error)
	})
	if !ok || ops == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "operations service unavailable")
	}
	op, err := ops.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op == nil {
		return echo.NewHTTPError(http.StatusNotFound, "restore operation not found")
	}
	opUser := extractStructField(op, "UserID")
	opKind := extractStructField(op, "Kind")
	opPayload := extractStructField(op, "Payload")
	opID := extractStructField(op, "ID")
	if opUser != userID || opKind != "restore" {
		return echo.NewHTTPError(http.StatusNotFound, "restore operation not found")
	}
	var target RestoreTarget
	if err := json.Unmarshal([]byte(opPayload), &target); err != nil {
		return err
	}
	if _, _, err := h.authorizedRecord(c, target.RecordID); err != nil {
		return err
	}
	if target.DatabaseID != "" && !h.hasAdminAccess(c, target.DatabaseID, "") {
		return echo.NewHTTPError(http.StatusForbidden, "restore target permission required")
	}
	var req ApplyRestoreRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	if err := h.backupService.ApplyRestore(c.Request().Context(), userID, opID, req.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]any{"status": "success", "data": op})
}

func (h *Handler) ProtectRecord(c echo.Context) error {
	if _, _, err := h.authorizedRecord(c, c.Param("recordId")); err != nil {
		return err
	}
	var req ProtectRecordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid protection expiry")
	}
	if err := h.backupService.ProtectRecord(c.Request().Context(), c.Param("recordId"), req.Until); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Backup protection saved", nil)
}

func (h *Handler) authorizeVolumeRestore(c echo.Context) (string, error) {
	id := c.Param("recordId")
	record, err := h.backupService.GetRecord(c.Request().Context(), id)
	if err != nil && !utils.IsNotFound(err) {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "failed to load backup record")
	}
	if err != nil || record == nil {
		return "", echo.NewHTTPError(http.StatusNotFound, ErrBackupRecordNotFound.Error())
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), record.BackupConfigID)
	if err != nil && !utils.IsNotFound(err) {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "failed to load backup configuration")
	}
	if err != nil || cfg == nil {
		return "", echo.NewHTTPError(http.StatusNotFound, ErrBackupConfigNotFound.Error())
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return "", echo.NewHTTPError(http.StatusForbidden, "admin access required for volume restoration")
	}
	return id, nil
}

func (h *Handler) VolumeRestoreTarget(c echo.Context) error {
	id, err := h.authorizeVolumeRestore(c)
	if err != nil {
		return err
	}
	serverID, volumeName, err := h.backupService.VolumeRestoreTarget(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Volume restore target validated", VolumeRestoreTargetResponse{
		TargetServerID: serverID,
		TargetVolume:   volumeName,
	})
}

func (h *Handler) RestoreVolume(c echo.Context) error {
	return echo.NewHTTPError(http.StatusGone, "Use the reviewed restore prepare/apply workflow")
}

func (h *Handler) CancelVolumeRestore(c echo.Context) error {
	id, err := h.authorizeVolumeRestore(c)
	if err != nil {
		return err
	}
	if !h.backupService.CancelVolumeRestore(id) {
		return utils.Error(c, http.StatusConflict, "no active volume restore for this record")
	}
	return utils.Success(c, "Volume restore interruption requested; existing data may be partially overwritten", nil)
}
