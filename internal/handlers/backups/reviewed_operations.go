package backups

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/services/operations"
	"codedock/internal/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *BackupHandler) SetOperations(service *operations.Service) { h.operations = service }
func (h *BackupHandler) authorizedRecord(c echo.Context, id string) (*models.BackupConfig, string, error) {
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
func (h *BackupHandler) backupProject(c echo.Context, cfg *models.BackupConfig) (string, error) {
	if cfg.ServiceID != "" {
		app, err := h.appService.GetAppService(c.Request().Context(), cfg.ServiceID)
		if err != nil {
			return "", echo.NewHTTPError(http.StatusNotFound, "application not found")
		}
		return app.ProjectID, nil
	}
	if cfg.DatabaseID == "global" || cfg.DatabaseID == "" {
		return "", nil
	}
	database, err := h.dbService.GetDatabase(c.Request().Context(), cfg.DatabaseID)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusNotFound, "database not found")
	}
	return database.ProjectID, nil
}
func (h *BackupHandler) PrepareRestore(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	var request models.RestoreReviewRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restore target")
	}
	_, project, err := h.authorizedRecord(c, c.Param("recordId"))
	if err != nil {
		return err
	}
	if request.TargetDatabaseID != "" {
		target, err := h.dbService.GetDatabase(c.Request().Context(), request.TargetDatabaseID)
		if err != nil || target.ProjectID != project || !h.hasAdminAccess(c, target.ID, "") {
			return echo.NewHTTPError(http.StatusForbidden, "restore target must be an authorized database in the same project")
		}
	}
	review, err := h.backupService.PrepareRestore(c.Request().Context(), user.UserID, project, c.Param("recordId"), request.TargetDatabaseID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review restore effects before applying", review)
}
func (h *BackupHandler) ApplyRestore(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	op, err := h.operations.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op.UserID != user.UserID || op.Kind != "restore" {
		return echo.NewHTTPError(http.StatusNotFound, "restore operation not found")
	}
	var target models.RestoreTarget
	if err := json.Unmarshal([]byte(op.Payload), &target); err != nil {
		return err
	}
	if _, _, err := h.authorizedRecord(c, target.RecordID); err != nil {
		return err
	}
	if target.DatabaseID != "" && !h.hasAdminAccess(c, target.DatabaseID, "") {
		return echo.NewHTTPError(http.StatusForbidden, "restore target permission required")
	}
	var request models.OperationApply
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	if err := h.backupService.ApplyRestore(c.Request().Context(), user.UserID, op.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]any{"status": "success", "data": op})
}
func (h *BackupHandler) ProtectRecord(c echo.Context) error {
	if _, _, err := h.authorizedRecord(c, c.Param("recordId")); err != nil {
		return err
	}
	var request models.BackupProtectionRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid protection expiry")
	}
	if err := h.backupService.ProtectRecord(c.Request().Context(), c.Param("recordId"), request.ProtectedUntil); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Backup protection saved", nil)
}
func (h *BackupHandler) StartRun(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	cfg, err := h.backupService.GetConfig(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "backup policy not found")
	}
	if !h.hasAdminAccess(c, cfg.DatabaseID, cfg.ServiceID) {
		return echo.NewHTTPError(http.StatusForbidden, "backup admin permission required")
	}
	project, err := h.backupProject(c, cfg)
	if err != nil {
		return err
	}
	op, err := h.backupService.StartRun(c.Request().Context(), user.UserID, project, cfg.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]any{"status": "success", "data": op})
}
