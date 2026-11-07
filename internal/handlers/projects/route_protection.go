package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	projectservices "codedock/internal/services/projects"
	"codedock/internal/utils"
)

type RouteProtectionHandler struct {
	service        *projectservices.RouteProtectionService
	projectService *projectservices.ProjectService
}

func NewRouteProtectionHandler(service *projectservices.RouteProtectionService, projectService *projectservices.ProjectService) *RouteProtectionHandler {
	return &RouteProtectionHandler{
		service:        service,
		projectService: projectService,
	}
}

func (h *RouteProtectionHandler) List(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}

	result, err := h.service.List(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Route protection status observed", result)
}

func (h *RouteProtectionHandler) Change(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}

	var input models.RouteProtectionInput
	if err := c.Bind(&input); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid route protection input")
	}
	if input.Hostname == "" || input.Action == "" {
		return utils.Error(c, http.StatusBadRequest, "hostname and action are required")
	}

	status, err := h.service.Change(c.Request().Context(), projectID, input)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Route protection updated", status)
}
