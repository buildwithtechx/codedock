package projects

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	services "codedock/internal/services/projects"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ComposeStackHandler struct {
	stacks       *services.ComposeStackService
	projects     *services.ProjectService
	environments *services.EnvironmentService
}

func NewComposeStackHandler(stacks *services.ComposeStackService, projects *services.ProjectService, environments *services.EnvironmentService) *ComposeStackHandler {
	return &ComposeStackHandler{stacks: stacks, projects: projects, environments: environments}
}

func (h *ComposeStackHandler) authorize(c echo.Context, write bool) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	permission := models.MemberPermission("")
	if write {
		permission = models.MemberPermissionAdmin
	}
	if !h.projects.HasPermission(c.Request().Context(), c.Param("id"), user.UserID, models.UserRole(user.Role), permission) {
		return echo.NewHTTPError(http.StatusForbidden, "project access denied")
	}
	project, err := h.projects.GetProject(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "project not found")
	}
	if project.ServerID != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Compose stacks currently require local Docker; this SSH target is unsupported")
	}
	return nil
}

func (h *ComposeStackHandler) request(c echo.Context) (*models.ComposeStackRequest, error) {
	if err := h.authorize(c, true); err != nil {
		return nil, err
	}
	var request models.ComposeStackRequest
	if err := c.Bind(&request); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid stack request")
	}
	environment, err := h.environments.GetEnvironment(c.Request().Context(), request.EnvironmentID)
	if err != nil || environment == nil || environment.ProjectID != c.Param("id") {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "environment must belong to this project")
	}
	return &request, nil
}

func (h *ComposeStackHandler) Review(c echo.Context) error {
	request, err := h.request(c)
	if err != nil {
		return err
	}
	review, err := h.stacks.Review(c.Request().Context(), *request)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Stack reviewed", review)
}

func (h *ComposeStackHandler) Save(c echo.Context) error {
	request, err := h.request(c)
	if err != nil {
		return err
	}
	stack, err := h.stacks.Save(c.Request().Context(), c.Param("id"), *request)
	if err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Stack configuration saved", stack)
}

func (h *ComposeStackHandler) List(c echo.Context) error {
	if err := h.authorize(c, false); err != nil {
		return err
	}
	stacks, err := h.stacks.ListObserved(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Stack configuration unavailable")
	}
	return utils.Success(c, "Stacks retrieved", stacks)
}

func (h *ComposeStackHandler) Deploy(c echo.Context) error {
	if err := h.authorize(c, true); err != nil {
		return err
	}
	if err := h.stacks.Deploy(c.Request().Context(), c.Param("id"), c.Param("stackId")); err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Accepted(c, "Stack deployment started", nil)
}

func (h *ComposeStackHandler) Cancel(c echo.Context) error {
	if err := h.authorize(c, true); err != nil {
		return err
	}
	stacks, err := h.stacks.List(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "Stack unavailable")
	}
	for _, stack := range stacks {
		if stack.ID == c.Param("stackId") {
			if err := h.stacks.Cancel(stack.ProjectID, stack.ID); err != nil {
				return utils.Error(c, http.StatusConflict, err.Error())
			}
			return utils.Accepted(c, "Cancellation requested; inspect per-service states before retrying", nil)
		}
	}
	return utils.Error(c, http.StatusNotFound, "stack not found")
}
