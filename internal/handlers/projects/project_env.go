package projects

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"

	projectservices "codedock/internal/services/projects"
)

type ProjectEnvHandler struct {
	envService *projectservices.EnvironmentService
}

func NewProjectEnvHandler(s *projectservices.EnvironmentService) *ProjectEnvHandler {
	return &ProjectEnvHandler{envService: s}
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
