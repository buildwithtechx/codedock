package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/utils"
)

type CanvasHandler struct {
	canvasService  *projectservices.CanvasService
	projectService *projectservices.ProjectService
}

func NewCanvasHandler(s *projectservices.CanvasService, ps *projectservices.ProjectService) *CanvasHandler {
	return &CanvasHandler{canvasService: s, projectService: ps}
}

func (h *CanvasHandler) ListCanvasSummaries(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	organizationID := c.QueryParam("organizationId")
	if organizationID != "" && user != nil && h.projectService != nil {
		if !h.projectService.HasOrgPermission(c.Request().Context(), organizationID, user.UserID, models.UserRole(user.Role), "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this organization")
		}
	}
	summaries, err := h.canvasService.ListSummaries(c.Request().Context(), organizationID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if summaries == nil {
		summaries = make([]models.CanvasSummary, 0)
	}
	if user != nil && user.Role != "admin" && h.projectService != nil {
		filtered := make([]models.CanvasSummary, 0, len(summaries))
		for _, s := range summaries {
			if h.projectService.HasPermission(c.Request().Context(), s.ID, user.UserID, models.UserRole(user.Role), "") {
				filtered = append(filtered, s)
			}
		}
		summaries = filtered
	}
	return utils.Success(c, "Operation successful", summaries)
}

func (h *CanvasHandler) GetCanvasSummary(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" && h.projectService != nil {
		if !h.projectService.HasPermission(c.Request().Context(), id, user.UserID, models.UserRole(user.Role), "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this canvas summary")
		}
	}
	summary, err := h.canvasService.GetSummary(c.Request().Context(), id)
	if err != nil || summary == nil {
		return utils.Error(c, http.StatusNotFound, "canvas summary not found")
	}
	return utils.Success(c, "Operation successful", summary)
}

func (h *CanvasHandler) GetEnvironmentCanvas(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	canvas, err := h.canvasService.GetEnvironmentCanvas(c.Request().Context(), id)
	if err != nil || canvas == nil {
		return utils.Error(c, http.StatusNotFound, "environment canvas not found")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" && h.projectService != nil && canvas.Environment != nil {
		if !h.projectService.HasPermission(c.Request().Context(), canvas.Environment.ProjectID, user.UserID, models.UserRole(user.Role), "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this environment canvas")
		}
	}
	return utils.Success(c, "Operation successful", canvas)
}
