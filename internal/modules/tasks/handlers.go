package tasks

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ScheduledTaskAccessChecker interface {
	HasServiceAccess(ctx context.Context, userID, role, serviceID string) bool
	HasProjectAccess(ctx context.Context, userID, role, projectID string) bool
}

type ScheduledTaskHandler struct {
	service       ScheduledTaskServiceFull
	accessChecker ScheduledTaskAccessChecker
}

func NewScheduledTaskHandler(
	s ScheduledTaskServiceFull,
	ac ScheduledTaskAccessChecker,
) *ScheduledTaskHandler {
	return &ScheduledTaskHandler{
		service:       s,
		accessChecker: ac,
	}
}

func (h *ScheduledTaskHandler) hasAccess(c echo.Context, serviceID string) bool {
	if h.accessChecker == nil {
		return true
	}
	userID, role, ok := utils.ExtractUser(c)
	if !ok {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	return h.accessChecker.HasServiceAccess(c.Request().Context(), userID, role, serviceID)
}

func (h *ScheduledTaskHandler) hasProjectAccess(c echo.Context, projectID string) bool {
	if h.accessChecker == nil {
		return true
	}
	userID, role, ok := utils.ExtractUser(c)
	if !ok {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	return h.accessChecker.HasProjectAccess(c.Request().Context(), userID, role, projectID)
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

	if h.service == nil {
		return utils.Success(c, "Operation successful", []ScheduledTask{})
	}

	var tasks []ScheduledTask
	var err error

	if serviceID != "" {
		tasks, err = h.service.ListByService(c.Request().Context(), serviceID)
	} else {
		tasks, err = h.service.ListByProject(c.Request().Context(), projectID)
	}

	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", tasks)
}

func (h *ScheduledTaskHandler) Create(c echo.Context) error {
	var j ScheduledTask
	if err := c.Bind(&j); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if !h.hasAccess(c, j.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	if h.service == nil {
		return utils.Created(c, "Created successfully", j)
	}

	created, err := h.service.CreateScheduledTask(c.Request().Context(), &j)
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
	if h.service == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	j, err := h.service.GetScheduledTask(c.Request().Context(), id)
	if err != nil || j == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}

	if !h.hasAccess(c, j.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	return utils.Success(c, "Operation successful", j)
}

func (h *ScheduledTaskHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	existing, err := h.service.GetScheduledTask(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}

	if !h.hasAccess(c, existing.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	var req UpdateScheduledTaskRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Schedule != nil {
		existing.Schedule = *req.Schedule
	}
	if req.Command != nil {
		existing.Command = *req.Command
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	if err := h.service.UpdateScheduledTask(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Updated successfully", existing)
}

func (h *ScheduledTaskHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	j, err := h.service.GetScheduledTask(c.Request().Context(), id)
	if err != nil || j == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	if !h.hasAccess(c, j.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.DeleteScheduledTask(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ScheduledTaskHandler) Run(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if h.service == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	j, err := h.service.GetScheduledTask(c.Request().Context(), id)
	if err != nil || j == nil {
		return utils.Error(c, http.StatusNotFound, "scheduled task not found")
	}
	if !h.hasAccess(c, j.ServiceID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	out, err := h.service.ExecuteScheduledTask(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "executed", "output": out})
}
