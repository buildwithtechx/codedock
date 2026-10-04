package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/utils"
)

type ProjectAppHandler struct {
	appService     *projectservices.ProjectAppService
	projectService *projectservices.ProjectService
}

func NewProjectAppHandler(appService *projectservices.ProjectAppService, projectService *projectservices.ProjectService) *ProjectAppHandler {
	return &ProjectAppHandler{
		appService:     appService,
		projectService: projectService,
	}
}

func (h *ProjectAppHandler) List(c echo.Context) error {
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.QueryParam("organization_id")
	}
	if orgID == "" {
		return utils.Error(c, http.StatusBadRequest, "organizationId is required")
	}

	apps, err := h.appService.ListByOrganization(c.Request().Context(), orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", apps)
}

func (h *ProjectAppHandler) Create(c echo.Context) error {
	var req models.CreateProjectAppRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "name is required")
	}
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.QueryParam("organization_id")
	}
	if orgID == "" {
		return utils.Error(c, http.StatusBadRequest, "organizationId is required")
	}

	app := &models.ProjectApp{
		OrganizationID: orgID,
		Name:           req.Name,
		Slug:           req.Slug,
		GitProvider:    req.GitProvider,
		GitOwner:       req.GitOwner,
		GitRepo:        req.GitRepo,
		GitURL:         req.GitURL,
		InstallationID: req.InstallationID,
	}

	created, err := h.appService.Create(c.Request().Context(), app)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *ProjectAppHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		id = c.Param("appId")
	}
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing app id")
	}

	app, err := h.appService.GetByID(c.Request().Context(), id)
	if err != nil || app == nil {
		return utils.Error(c, http.StatusNotFound, "project app not found")
	}
	return utils.Success(c, "Operation successful", app)
}

func (h *ProjectAppHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		id = c.Param("appId")
	}
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing app id")
	}

	var req models.UpdateProjectAppRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	app, err := h.appService.GetByID(c.Request().Context(), id)
	if err != nil || app == nil {
		return utils.Error(c, http.StatusNotFound, "project app not found")
	}

	if req.Name != "" {
		app.Name = req.Name
	}
	if req.Slug != "" {
		app.Slug = req.Slug
	}
	if req.Favicon != "" {
		app.Favicon = req.Favicon
	}
	if req.GitProvider != "" {
		app.GitProvider = req.GitProvider
	}
	if req.GitOwner != "" {
		app.GitOwner = req.GitOwner
	}
	if req.GitRepo != "" {
		app.GitRepo = req.GitRepo
	}
	if req.GitURL != "" {
		app.GitURL = req.GitURL
	}

	if err := h.appService.Update(c.Request().Context(), app); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Updated successfully", app)
}

func (h *ProjectAppHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		id = c.Param("appId")
	}
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing app id")
	}

	if err := h.appService.Delete(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Deleted successfully", map[string]string{"status": "deleted"})
}
