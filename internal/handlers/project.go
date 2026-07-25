package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/telemetry"
)

type ProjectHandler struct {
	projectService         *services.ProjectService
	projectSettingsService *services.ProjectSettingsService
}

func NewProjectHandler(s *services.ProjectService, pss *services.ProjectSettingsService) *ProjectHandler {
	return &ProjectHandler{projectService: s, projectSettingsService: pss}
}

func (h *ProjectHandler) ListProjects(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		return utils.Error(c, http.StatusBadRequest, "organizationId is required")
	}

	projects, total, err := h.projectService.ListProjectsByOrganization(c.Request().Context(), orgID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Paginated(c, "Operation successful", projects, total, page, limit)
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req models.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.OrganizationID == "none" {
		req.OrganizationID = ""
	}
	if req.ServerID == "local" {
		req.ServerID = ""
	}

	userClaims, ok := c.Get("user").(*models.UserClaims)

	if ok && userClaims != nil {
		if userClaims.PlanType != "pro" {
			count, err := h.projectService.CountProjectsByUser(c.Request().Context(), userClaims.UserID)
			if err != nil {
				return utils.Error(c, http.StatusInternalServerError, "failed to check project limits")
			}
			if count >= 2 {
				return utils.Error(c, http.StatusPaymentRequired, "hobby plan is limited to 2 projects. please upgrade to pro")
			}
		}
	}

	p, err := h.projectService.CreateProjectFromRequest(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if ok && userClaims != nil {
		telemetry.Track(userClaims.Email, "project_created", map[string]any{
			"project_id": p.ID,
			"name":       p.Name,
		})
	} else {
		telemetry.Track("anonymous", "project_created", map[string]any{
			"project_id": p.ID,
			"name":       p.Name,
		})
	}

	return utils.Created(c, "Created successfully", p)
}

func (h *ProjectHandler) GetProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	p, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
	}
	return utils.Success(c, "Operation successful", p)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	_, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
	}
	if err := h.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}
