package environments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ProjectPermissionChecker interface {
	HasPermission(ctx context.Context, projectID, userID, role, permission string) bool
}

type EnvironmentHandler struct {
	envService    *EnvironmentService
	permChecker   ProjectPermissionChecker
	suggestionSvc *EnvSuggestionService
}

func NewEnvironmentHandler(s *EnvironmentService, pc ProjectPermissionChecker, ss *EnvSuggestionService) *EnvironmentHandler {
	return &EnvironmentHandler{
		envService:    s,
		permChecker:   pc,
		suggestionSvc: ss,
	}
}

func (h *EnvironmentHandler) hasAccess(ctx echo.Context, projectID string) bool {
	userID, role, ok := utils.ExtractUser(ctx)
	if !ok {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	if projectID == "" || h.permChecker == nil {
		return false
	}
	return h.permChecker.HasPermission(ctx.Request().Context(), projectID, userID, role, "admin")
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
	var env EnvironmentConfig
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
	if err != nil || env == nil {
		return utils.Error(c, http.StatusNotFound, "environment not found")
	}
	if !h.hasAccess(c, env.ProjectID) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	if err := h.envService.DeleteEnvironment(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *EnvironmentHandler) GetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	vars, err := h.envService.GetVars(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if vars == nil {
		vars = map[string]string{}
	}
	return utils.Success(c, "Operation successful", vars)
}

func (h *EnvironmentHandler) SetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	var raw map[string]any
	if err := c.Bind(&raw); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	vars := make(map[string]string)
	if nested, ok := raw["variables"].(map[string]any); ok {
		for k, v := range nested {
			vars[k] = fmt.Sprint(v)
		}
	} else {
		for k, v := range raw {
			vars[k] = fmt.Sprint(v)
		}
	}
	for k, v := range vars {
		if err := h.envService.SetVar(c.Request().Context(), projectID, k, v); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
	}
	return utils.Success(c, "Operation successful", vars)
}

func (h *EnvironmentHandler) SuggestVars(c echo.Context) error {
	if h.suggestionSvc == nil {
		return utils.Success(c, "Operation successful", []EnvExampleVariableSuggestion{})
	}
	repoURL := c.QueryParam("repoUrl")
	branch := c.QueryParam("branch")
	rootDir := c.QueryParam("rootDir")
	if repoURL == "" {
		return utils.Error(c, http.StatusBadRequest, "repoUrl parameter required")
	}
	suggestions, err := h.suggestionSvc.SuggestEnvVars(c.Request().Context(), repoURL, branch, rootDir)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Operation successful", suggestions)
}
