package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	systemservices "codedock.run/codedock/internal/services/system"
)

type UpdaterHandler struct {
	updaterService *systemservices.UpdaterService
}

func NewUpdaterHandler(s *systemservices.UpdaterService) *UpdaterHandler {
	return &UpdaterHandler{updaterService: s}
}

func (h *UpdaterHandler) GetUpdateStatus(c echo.Context) error {
	if h.updaterService == nil {
		return utils.Error(c, http.StatusInternalServerError, "updater service not initialized")
	}
	status := h.updaterService.GetStatus()
	return utils.Success(c, "Operation successful", status)
}

func (h *UpdaterHandler) CheckUpdate(c echo.Context) error {
	if h.updaterService == nil {
		return utils.Error(c, http.StatusInternalServerError, "updater service not initialized")
	}
	if _, err := h.updaterService.CheckForUpdates(c.Request().Context()); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	status := h.updaterService.GetStatus()
	return utils.Success(c, "Operation successful", status)
}

func (h *UpdaterHandler) DeployUpdate(c echo.Context) error {
	if h.updaterService == nil {
		return utils.Error(c, http.StatusInternalServerError, "updater service not initialized")
	}
	if err := h.updaterService.DeployUpdate(c.Request().Context()); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Accepted(c, "update deployment triggered", nil)
}
