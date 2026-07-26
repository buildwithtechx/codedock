package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
)

type ScheduledTaskHandler struct {
	scheduledTaskService *services.ScheduledTaskService
	appService           *services.AppService
	projectService       *services.ProjectService
}

func NewScheduledTaskHandler(s *services.ScheduledTaskService, app *services.AppService, proj *services.ProjectService) *ScheduledTaskHandler {
	return &ScheduledTaskHandler{
		scheduledTaskService: s,
		appService:           app,
		projectService:       proj,
	}
}

func (h *ScheduledTaskHandler) hasAccess(ctx echo.Context, serviceID string) bool {
	userClaims, ok := ctx.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return false
	}
	if userClaims.Role == models.UserRoleAdmin || userClaims.Role == models.UserRoleOwner {
		return true
	}

	if serviceID == "" {
		return false
	}

	svc, err := h.appService.GetAppService(ctx.Request().Context(), serviceID)
	if err != nil || svc == nil || svc.ProjectID == "" {
		return false
	}

	hasPerm := h.projectService.HasPermission(ctx.Request().Context(), svc.ProjectID, userClaims.UserID, userClaims.Role, "")
	return hasPerm
}

func (h *ScheduledTaskHandler) hasProjectAccess(ctx echo.Context, projectID string) bool {
	userClaims, ok := ctx.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return false
	}
	if userClaims.Role == models.UserRoleAdmin || userClaims.Role == models.UserRoleOwner {
		return true
	}

	if projectID == "" {
		return false
	}

	hasPerm := h.projectService.HasPermission(ctx.Request().Context(), projectID, userClaims.UserID, userClaims.Role, "")
	return hasPerm
}

func (h *ScheduledTaskHandler) ListProjectScheduledTasks(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	serviceID := c.QueryParam("serviceId")

	if projectID == "" && serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "must provide projectId or serviceId")
	}

	if serviceID != "" {
		if !h.hasAccess(c, serviceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	} else if projectID != "" {
		if !h.hasProjectAccess(c, projectID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	}

	var tasks []models.ScheduledTask
	var err error

	if serviceID != "" {
		tasks, err = h.scheduledTaskService.ListScheduledTasksByService(c.Request().Context(), serviceID)
	} else {
		tasks, err = h.scheduledTaskService.ListScheduledTasksByProject(c.Request().Context(), projectID)
	}

	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", tasks)
}

func (h *ScheduledTaskHandler) Create(c echo.Context) error {
	var j models.ScheduledTask
	if err := c.Bind(&j); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if !h.hasAccess(c, j.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	created, err := h.scheduledTaskService.CreateScheduledTask(c.Request().Context(), &j)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *ScheduledTaskHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	j, err := h.scheduledTaskService.GetScheduledTask(c.Request().Context(), id)
	if err != nil || j == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}

	if !h.hasAccess(c, j.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	return utils.Success(c, "Operation successful", j)
}

func (h *ScheduledTaskHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	j, err := h.scheduledTaskService.GetScheduledTask(c.Request().Context(), id)
	if err == nil && j != nil {
		if !h.hasAccess(c, j.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	}

	if err := h.scheduledTaskService.DeleteScheduledTask(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ScheduledTaskHandler) Run(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	j, err := h.scheduledTaskService.GetScheduledTask(c.Request().Context(), id)
	if err == nil && j != nil {
		if !h.hasAccess(c, j.ServiceID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	}

	out, err := h.scheduledTaskService.ExecuteScheduledTask(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "executed", "output": out})
}
