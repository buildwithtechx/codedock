package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func (h *AppHandler) CheckProjectRuntimeUpdates(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	result, err := h.appService.CheckRuntimeUpdates(c.Request().Context(), projectID, "")
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to check runtime updates: "+err.Error())
	}
	return utils.Success(c, "Runtime update check complete", result)
}

func (h *AppHandler) ApplyProjectRuntimeUpdates(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	services, err := h.appService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to list services: "+err.Error())
	}

	var triggered []*models.Deployment
	for _, svc := range services {
		dep := &models.Deployment{
			ServiceID:     svc.ID,
			EnvironmentID: svc.EnvironmentID,
			ProjectID:     svc.ProjectID,
			Status:        "BUILDING",
			CommitMessage: "Redeploy with latest base image",
			Branch:        svc.Branch,
			Trigger:       "Runtime Image Update",
		}
		created, err := h.deploymentService.CreateDeployment(c.Request().Context(), dep)
		if err == nil && created != nil {
			h.deploymentService.ExecuteDeploymentAsync(created)
			triggered = append(triggered, created)
		}
	}

	return utils.Accepted(c, "Runtime updates triggered", triggered)
}

func (h *AppHandler) CheckAppRuntimeUpdates(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "service id is required")
	}

	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	result, err := h.appService.CheckRuntimeUpdates(c.Request().Context(), svc.ProjectID, serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to check runtime updates: "+err.Error())
	}
	return utils.Success(c, "Runtime update check complete", result)
}

func (h *AppHandler) ApplyAppRuntimeUpdates(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "service id is required")
	}

	existing, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	dep := &models.Deployment{
		ServiceID:     existing.ID,
		EnvironmentID: existing.EnvironmentID,
		ProjectID:     existing.ProjectID,
		Status:        "BUILDING",
		CommitMessage: "Redeploy with latest base image",
		Branch:        existing.Branch,
		Trigger:       "Runtime Image Update",
	}
	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), dep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)
	return utils.Accepted(c, "Redeployment triggered with latest base image", created)
}
