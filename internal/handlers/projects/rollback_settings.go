package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	deploymentservices "codedock/internal/services/deployments"
	"codedock/internal/utils"
)

type RollbackSettingsHandler struct {
	retentionService *deploymentservices.RetentionService
}

func NewRollbackSettingsHandler(rs *deploymentservices.RetentionService) *RollbackSettingsHandler {
	return &RollbackSettingsHandler{retentionService: rs}
}

func (h *RollbackSettingsHandler) GetRollbackCapacity(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}
	capacity, err := h.retentionService.GetRollbackCapacity(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Rollback capacity retrieved", capacity)
}

func (h *RollbackSettingsHandler) UpdateRollbackSettings(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}
	var req models.UpdateRollbackSettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	if err := h.retentionService.UpdateRollbackSettings(c.Request().Context(), projectID, req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	capacity, err := h.retentionService.GetRollbackCapacity(c.Request().Context(), projectID)
	if err != nil {
		return utils.Success(c, "Rollback settings updated", nil)
	}
	return utils.Success(c, "Rollback settings updated", capacity)
}
