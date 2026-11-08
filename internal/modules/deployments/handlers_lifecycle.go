package deployments

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AuditLogger interface {
	LogAction(ctx any, action, resource, ip string, details map[string]string)
}

type LifecycleHandler struct {
	service   *DeploymentService
	appLookup AppServiceLookup
	audit     AuditLogger
}

func NewLifecycleHandler(service *DeploymentService, appLookup AppServiceLookup, audit AuditLogger) *LifecycleHandler {
	return &LifecycleHandler{
		service:   service,
		appLookup: appLookup,
		audit:     audit,
	}
}

func (h *LifecycleHandler) Trigger(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	var req TriggerDeploymentRequest
	_ = c.Bind(&req)

	if h.appLookup == nil {
		return utils.Error(c, http.StatusInternalServerError, "app service lookup unavailable")
	}

	app, err := h.appLookup.GetAppService(c.Request().Context(), serviceID)
	if err != nil || app == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	branch := app.Branch
	if req.Branch != nil && *req.Branch != "" {
		branch = *req.Branch
	}

	logs := fmt.Sprintf("Initiating manual deployment for %s branch %s...\n", app.Name, branch)
	depID, err := h.service.TriggerDeployment(c.Request().Context(), app.ID, app.EnvironmentID, app.ProjectID, branch, "Manual", logs)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.audit != nil {
		h.audit.LogAction(c.Request().Context(), "deployment.create", depID, c.RealIP(), map[string]string{
			"serviceId": app.ID,
			"branch":    branch,
		})
	}

	return utils.Accepted(c, "Deployment triggered", DeployResponse{
		DeploymentID: depID,
		Status:       string(DeploymentStatusBuilding),
		Message:      "Deployment initiated",
	})
}

func (h *LifecycleHandler) TriggerProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}

	// Trigger deployment on project scope
	depID, err := h.service.TriggerDeployment(c.Request().Context(), "", "", projectID, "main", "Project Manual", "Initiating project rollout...\n")
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Accepted(c, "Project deployment triggered", DeployResponse{
		DeploymentID: depID,
		Status:       string(DeploymentStatusBuilding),
		Message:      "Project deployment initiated",
	})
}

func (h *LifecycleHandler) Rollback(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	var req RollbackRequest
	_ = c.Bind(&req)

	dep, err := h.service.Rollback(c.Request().Context(), id, req.Reason)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.audit != nil {
		h.audit.LogAction(c.Request().Context(), "deployment.rollback", dep.ID, c.RealIP(), map[string]string{
			"targetId": id,
		})
	}

	return utils.Accepted(c, "Rollback initiated", RollbackResponse{
		DeploymentID: dep.ID,
		PreviousID:   id,
		Status:       string(dep.Status),
	})
}

func (h *LifecycleHandler) Retry(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	dep, err := h.service.Retry(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.audit != nil {
		h.audit.LogAction(c.Request().Context(), "deployment.retry", dep.ID, c.RealIP(), map[string]string{
			"sourceId": id,
		})
	}

	return utils.Accepted(c, "Deployment retry initiated", DeployResponse{
		DeploymentID: dep.ID,
		Status:       string(dep.Status),
		Message:      "Retry queued",
	})
}

func (h *LifecycleHandler) TogglePin(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	pinned, err := h.service.TogglePin(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Pin toggled", PinToggleResponse{
		DeploymentID: id,
		Pinned:       pinned,
	})
}

func (h *LifecycleHandler) Cancel(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	err := h.service.Cancel(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Deployment cancelled", CancelResponse{
		DeploymentID: id,
		Cancelled:    true,
	})
}

func (h *LifecycleHandler) GetAdmission(c echo.Context) error {
	status := h.service.GetAdmissionStatus()
	return utils.Success(c, "Admission status retrieved", status)
}
