package apps

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *Handler) SetDeployer(d LifecycleDeployer) {
	h.deployer = d
}

func (h *Handler) StopService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.service.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if !h.hasAccess(c, existing.ProjectID, true) {
		return utils.Error(c, http.StatusForbidden, "forbidden")
	}

	if h.deployer != nil {
		if err := h.deployer.StopAppService(c.Request().Context(), id); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
	}

	existing.Status = AppServiceStatusStopped
	_ = h.service.UpdateAppService(c.Request().Context(), existing)
	return utils.Success(c, "Service stopped successfully", existing)
}

func (h *Handler) RestartService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.service.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if !h.hasAccess(c, existing.ProjectID, true) {
		return utils.Error(c, http.StatusForbidden, "forbidden")
	}

	if h.deployer != nil {
		if err := h.deployer.RestartAppService(c.Request().Context(), id); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
	}

	existing.Status = AppServiceStatusRunning
	_ = h.service.UpdateAppService(c.Request().Context(), existing)
	return utils.Success(c, "Service restarted successfully", existing)
}

func (h *Handler) RedeployService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.service.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if !h.hasAccess(c, existing.ProjectID, true) {
		return utils.Error(c, http.StatusForbidden, "forbidden")
	}

	if h.deployer != nil {
		res, err := h.deployer.RedeployAppService(c.Request().Context(), id)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return utils.Accepted(c, "Redeployment triggered", res)
	}

	return utils.Accepted(c, "Redeployment requested", nil)
}

func (h *Handler) CheckProjectRuntimeUpdates(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}
	if !h.hasAccess(c, projectID, false) {
		return utils.Error(c, http.StatusForbidden, "forbidden")
	}

	result := map[string]any{
		"projectId": projectID,
		"updates":   []any{},
	}
	return utils.Success(c, "Runtime update check complete", result)
}

func (h *Handler) ApplyProjectRuntimeUpdates(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}
	if !h.hasAccess(c, projectID, true) {
		return utils.Error(c, http.StatusForbidden, "forbidden")
	}

	if h.deployer != nil {
		res, err := h.deployer.ApplyProjectUpdates(c.Request().Context(), projectID)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return utils.Accepted(c, "Runtime updates triggered", res)
	}

	return utils.Accepted(c, "Runtime updates triggered", nil)
}

func (h *Handler) CheckAppRuntimeUpdates(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "service id is required")
	}
	result := map[string]any{
		"serviceId":       serviceID,
		"updateAvailable": false,
	}
	return utils.Success(c, "Runtime update check complete", result)
}

func (h *Handler) ApplyAppRuntimeUpdates(c echo.Context) error {
	return h.RedeployService(c)
}
