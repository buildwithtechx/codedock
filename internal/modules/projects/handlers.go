package projects

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/users"
	"codedock/internal/telemetry"
	"codedock/internal/utils"
)

type Handler struct {
	projectService         *Service
	projectSettingsService *ProjectSettingsService
	orgProvider            OrganizationProvider
}

func NewHandler(s *Service, pss *ProjectSettingsService, op OrganizationProvider) *Handler {
	return &Handler{
		projectService:         s,
		projectSettingsService: pss,
		orgProvider:            op,
	}
}

func (h *Handler) ListProjects(c echo.Context) error {
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
	userID, roleStr, ok := utils.ExtractUser(c)
	if !ok || userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	role := users.UserRole(roleStr)

	if orgID == "" && h.orgProvider != nil {
		org, err := h.orgProvider.GetOrCreateDefaultOrganization(c.Request().Context(), userID)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		if org != nil {
			orgID = org.ID
		}
	}

	if orgID != "" && !h.projectService.HasOrgPermission(c.Request().Context(), orgID, userID, role, "") {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions for this organization")
	}

	projects, total, err := h.projectService.ListProjectsByOrganization(c.Request().Context(), orgID, c.QueryParam("search"), limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Paginated(c, "Operation successful", projects, total, page, limit)
}

func (h *Handler) CreateProject(c echo.Context) error {
	var req CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.OrganizationID == "none" {
		req.OrganizationID = ""
	}
	if req.ServerID == "local" {
		req.ServerID = ""
	}

	userID, roleStr, ok := utils.ExtractUser(c)
	_, userEmail := utils.ExtractClaims(c)
	role := users.UserRole(roleStr)
	planType := utils.ExtractStructField(c.Get("user"), "PlanType")

	if ok && userID != "" {
		if req.OrganizationID == "" && h.orgProvider != nil {
			organization, err := h.orgProvider.GetOrCreateDefaultOrganization(c.Request().Context(), userID)
			if err != nil {
				return utils.Error(c, http.StatusInternalServerError, err.Error())
			}
			req.OrganizationID = organization.ID
		}

		if req.OrganizationID != "" {
			if !h.projectService.HasOrgPermission(c.Request().Context(), req.OrganizationID, userID, role, "") {
				return utils.Error(c, http.StatusForbidden, "you do not have permission to create a project in this organization")
			}
		}

		if planType != "pro" {
			count, err := h.projectService.CountProjectsByUser(c.Request().Context(), userID)
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

	trackEmail := "anonymous"
	if ok && userEmail != "" {
		trackEmail = userEmail
	}
	telemetry.Track(trackEmail, "project_created", map[string]any{
		"project_id": p.ID,
		"name":       p.Name,
	})

	return utils.Created(c, "Created successfully", p)
}

func (h *Handler) GetProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	p, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil || p == nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	return utils.Success(c, "Operation successful", p)
}

func (h *Handler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id parameter")
	}
	p, err := h.projectService.GetProject(c.Request().Context(), id)
	if err != nil || p == nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}
	if err := h.projectService.DeleteProject(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", map[string]string{"status": "deleted"})
}
