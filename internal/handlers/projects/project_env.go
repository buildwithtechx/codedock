package projects

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	projectservices "codedock/internal/services/projects"
	"codedock/internal/utils"
)

type ProjectEnvHandler struct {
	envService *projectservices.EnvironmentService
	appService *projectservices.AppService
}

func NewProjectEnvHandler(s *projectservices.EnvironmentService, as *projectservices.AppService) *ProjectEnvHandler {
	return &ProjectEnvHandler{envService: s, appService: as}
}

func (h *ProjectEnvHandler) GetVars(c echo.Context) error {
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

func (h *ProjectEnvHandler) SetVars(c echo.Context) error {
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

type MigrateCommandEnvironmentRequest struct {
	Environment string            `json:"environment"`
	ServiceID   string            `json:"serviceId,omitempty"`
	Commands    map[string]string `json:"commands"`
}

type MigrateCommandEnvironmentResponse struct {
	Keys     []string          `json:"keys"`
	Commands map[string]string `json:"commands"`
}

func (h *ProjectEnvHandler) MigrateCommands(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}

	var req MigrateCommandEnvironmentRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if len(req.Commands) == 0 {
		return utils.Error(c, http.StatusBadRequest, "Select a command to migrate")
	}

	cleanedCommands := make(map[string]string)
	variables := make(map[string]string)
	var keys []string

	for field, cmd := range req.Commands {
		parsed := utils.InspectCommandEnvironment(cmd)
		if !parsed.Migratable {
			return utils.Error(c, http.StatusBadRequest, fmt.Sprintf("%s contains shell expressions or has no literal variables to migrate", field))
		}
		for k, v := range parsed.Variables {
			if existingVal, exists := variables[k]; exists && existingVal != v {
				return utils.Error(c, http.StatusConflict, fmt.Sprintf("Commands use different values for %s. Resolve the difference before migrating.", k))
			}
			if _, exists := variables[k]; !exists {
				keys = append(keys, k)
			}
			variables[k] = v
		}
		cleanedCommands[field] = parsed.Command
	}

	ctx := c.Request().Context()
	for k, v := range variables {
		if err := h.envService.SetVar(ctx, projectID, k, v); err != nil {
			return utils.Error(c, http.StatusInternalServerError, fmt.Sprintf("save variable %s: %v", k, err))
		}
	}

	if req.ServiceID != "" && h.appService != nil {
		if svc, err := h.appService.GetAppService(ctx, req.ServiceID); err == nil && svc != nil && svc.ProjectID == projectID {
			if val, ok := cleanedCommands["installCommand"]; ok {
				svc.InstallCommand = val
			}
			if val, ok := cleanedCommands["buildCommand"]; ok {
				svc.BuildCommand = val
			}
			if val, ok := cleanedCommands["startCommand"]; ok {
				svc.StartCommand = val
			}
			_ = h.appService.UpdateAppService(ctx, svc)
		}
	}

	return utils.Success(c, "Variables migrated successfully", MigrateCommandEnvironmentResponse{
		Keys:     keys,
		Commands: cleanedCommands,
	})
}
