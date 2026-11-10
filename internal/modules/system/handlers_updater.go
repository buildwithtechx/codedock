package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type UpdaterHandler struct {
	service UpdaterOperations
}

func NewUpdaterHandler(s UpdaterOperations) *UpdaterHandler {
	return &UpdaterHandler{service: s}
}

func (h *UpdaterHandler) GetUpdateStatus(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "updater service not initialized")
	}
	status := h.service.GetStatus()
	return utils.Success(c, "Operation successful", status)
}

func (h *UpdaterHandler) CheckUpdate(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "updater service not initialized")
	}
	status, err := h.service.CheckUpdate(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", status)
}

func (h *UpdaterHandler) ApplyUpdate(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "updater service not initialized")
	}
	status, err := h.service.ApplyUpdate(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", status)
}

func (h *UpdaterHandler) DeployUpdate(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "updater service not initialized")
	}
	if _, err := h.service.ApplyUpdate(c.Request().Context()); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Accepted(c, "update deployment triggered", nil)
}
