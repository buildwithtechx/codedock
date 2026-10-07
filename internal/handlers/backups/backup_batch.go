package backups

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type PolicyBatchPayload struct {
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	Timezone    string `json:"timezone"`
	Timeout     int    `json:"timeout"`
}

func (h *BackupHandler) projectAdmin(c echo.Context, projectID string) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	if userClaims.Role == models.UserRoleAdmin || userClaims.Role == models.UserRoleOwner {
		return nil
	}
	if !h.projectService.HasPermission(c.Request().Context(), projectID, userClaims.UserID, userClaims.Role, models.MemberPermissionAdmin) {
		return echo.NewHTTPError(http.StatusForbidden, "project admin access required")
	}
	return nil
}

func (h *BackupHandler) ListBatches(c echo.Context) error {
	projectID := c.Param("id")
	if err := h.projectAdmin(c, projectID); err != nil {
		return err
	}
	list, err := h.backupService.ListBatches(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *BackupHandler) CreateBatch(c echo.Context) error {
	projectID := c.Param("id")
	if err := h.projectAdmin(c, projectID); err != nil {
		return err
	}
	var req PolicyBatchPayload
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	batch := models.BackupPolicyBatch{
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
	return utils.Created(c, "Created successfully", batch)
}

func (h *BackupHandler) TriggerBatch(c echo.Context) error {
	batchID := c.Param("batchId")
	if batchID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing batch id")
	}
	batch, err := h.backupService.GetBatch(c.Request().Context(), batchID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "policy batch not found")
	}
	if err := h.projectAdmin(c, batch.ProjectID); err != nil {
		return err
	}
	records, err := h.backupService.TriggerBatch(c.Request().Context(), batchID)
	if err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Accepted(c, "Batch triggered", records)
}
