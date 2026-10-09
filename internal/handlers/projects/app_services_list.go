package projects

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/utils"
)

func (h *AppHandler) ListByOrganization(c echo.Context) error {
	organizationID := c.QueryParam("organizationId")
	if organizationID == "" {
		return utils.Error(c, http.StatusBadRequest, "organizationId is required")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if !h.projectService.HasOrgPermission(c.Request().Context(), organizationID, user.UserID, user.Role, "") {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions for this organization")
	}
	apps, err := h.appService.ListByOrganization(c.Request().Context(), organizationID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for _, app := range apps {
		if app != nil {
			app.DeployToken = ""
		}
	}
	return utils.Success(c, "Operation successful", apps)
}

func (h *AppHandler) ListByEnvironment(c echo.Context) error {
	envID := c.Param("id")
	apps, err := h.appService.ListByEnvironment(c.Request().Context(), envID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != models.UserRoleAdmin && user.Role != models.UserRoleOwner {
		var filtered []*models.AppService
		for _, app := range apps {
			if h.projectService.IsMemberOrOwner(c.Request().Context(), app.ProjectID, user.UserID, user.Role) {
				filtered = append(filtered, app)
			}
		}
		return utils.Success(c, "Operation successful", filtered)
	}
	return utils.Success(c, "Operation successful", apps)
}

func (h *AppHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		projectID = c.Param("projectId")
	}
	if err := h.verifyProjectOwnership(c, projectID); err != nil {
		return err
	}
	apps, err := h.appService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	targetEnv := c.QueryParam("environmentId")
	if targetEnv == "" {
		targetEnv = c.QueryParam("environment_id")
	}
	envName := c.QueryParam("env")
	if envName == "" {
		envName = c.QueryParam("environment")
	}

	if targetEnv == "" && envName != "" && h.envService != nil {
		envs, _ := h.envService.ListByProject(c.Request().Context(), projectID)
		for _, e := range envs {
			if strings.EqualFold(e.Name, envName) {
				targetEnv = e.ID
				break
			}
		}
	}

	if targetEnv != "" {
		var filtered []*models.AppService
		for _, app := range apps {
			if app.EnvironmentID == targetEnv {
				filtered = append(filtered, app)
			}
		}
		apps = filtered
	}

	appID := c.QueryParam("appId")
	if appID == "" {
		appID = c.QueryParam("app_id")
	}
	if appID != "" {
		var filtered []*models.AppService
		for _, app := range apps {
			if app.AppID == appID {
				filtered = append(filtered, app)
			}
		}
		apps = filtered
	}

	return utils.Success(c, "Operation successful", apps)
}
