package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
)

type EnvironmentHandler struct {
	envService     *services.EnvironmentService
	projectService *services.ProjectService
}

func NewEnvironmentHandler(s *services.EnvironmentService, proj *services.ProjectService) *EnvironmentHandler {
	return &EnvironmentHandler{
		envService:     s,
		projectService: proj,
	}
}

func (h *EnvironmentHandler) hasAccess(ctx echo.Context, projectID string) bool {
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

	hasPerm := h.projectService.HasPermission(ctx.Request().Context(), projectID, userClaims.UserID, userClaims.Role, models.MemberPermissionAdmin)
	return hasPerm
}

func (h *EnvironmentHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	envs, err := h.envService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", envs)
}

func (h *EnvironmentHandler) Create(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	var env models.EnvironmentConfig
	if err := c.Bind(&env); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	env.ProjectID = projectID
	if env.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "environment name is required")
	}
	created, err := h.envService.CreateEnvironment(c.Request().Context(), &env)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *EnvironmentHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	env, err := h.envService.GetEnvironment(c.Request().Context(), id)
	if err == nil && env != nil {
		if !h.hasAccess(c, env.ProjectID) {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions")
		}
	}

	if err := h.envService.DeleteEnvironment(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
