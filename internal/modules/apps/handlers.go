package apps

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AccessChecker interface {
	HasPermission(ctx context.Context, projectID, userID, role, permission string) bool
	AccessibleServiceIDs(ctx context.Context, userID, role string) (map[string]bool, error)
}

type EnvironmentLookup interface {
	GetProjectID(ctx context.Context, environmentID string) (string, error)
}

type Handler struct {
	service   *Service
	access    AccessChecker
	envLookup EnvironmentLookup
	deployer  LifecycleDeployer
	runtime   RuntimeEngine
}

func NewHandler(s *Service, ac AccessChecker, el EnvironmentLookup) *Handler {
	return &Handler{
		service:   s,
		access:    ac,
		envLookup: el,
	}
}

func (h *Handler) hasAccess(ctx echo.Context, projectID string, adminOnly bool) bool {
	userID, role, ok := utils.ExtractUser(ctx)
	if !ok {
		return false
	}
	if role == "admin" || role == "owner" {
		return true
	}
	if projectID == "" || h.access == nil {
		return false
	}
	perm := ""
	if adminOnly {
		perm = "admin"
	}
	return h.access.HasPermission(ctx.Request().Context(), projectID, userID, role, perm)
}

func (h *Handler) ListAll(c echo.Context) error {
	userID, role, ok := utils.ExtractUser(c)
	if !ok {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	search := c.QueryParam("search")
	apps, err := h.service.ListAllServices(c.Request().Context(), search)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if role != "admin" && role != "owner" && h.access != nil {
		allowed, err := h.access.AccessibleServiceIDs(c.Request().Context(), userID, role)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		filtered := make([]*AppService, 0, len(apps))
		for _, a := range apps {
			if _, ok := allowed[a.ID]; ok {
				filtered = append(filtered, a)
			}
		}
		apps = filtered
	}

	return utils.Success(c, "Operation successful", apps)
}

func (h *Handler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}
	if !h.hasAccess(c, projectID, false) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	apps, err := h.service.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", apps)
}

func (h *Handler) ListByEnvironment(c echo.Context) error {
	envID := c.Param("id")
	if envID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing environment id")
	}

	apps, err := h.service.ListByEnvironment(c.Request().Context(), envID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", apps)
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing app service id")
	}

	app, err := h.service.GetAppService(c.Request().Context(), id)
	if err != nil || app == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if !h.hasAccess(c, app.ProjectID, false) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	return utils.Success(c, "Operation successful", app)
}

func (h *Handler) Create(c echo.Context) error {
	envID := c.Param("id")
	var req AppService
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "app service name is required")
	}
	if req.ProjectID == "" && h.envLookup != nil {
		pID, _ := h.envLookup.GetProjectID(c.Request().Context(), envID)
		req.ProjectID = pID
	}
	if req.ProjectID == "" {
		return utils.Error(c, http.StatusBadRequest, "projectId is required")
	}
	if !h.hasAccess(c, req.ProjectID, true) {
		return utils.Error(c, http.StatusForbidden, "project admin access required")
	}

	req.EnvironmentID = envID
	if req.ID != "" {
		if _, err := uuid.Parse(req.ID); err != nil {
			return utils.Error(c, http.StatusBadRequest, "application ID must be a UUID")
		}
	}
	userID := utils.ExtractUserID(c)
	req.GitUserID = userID

	created, err := h.service.CreateAppService(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing app service id")
	}

	existing, err := h.service.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if !h.hasAccess(c, existing.ProjectID, true) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	var req UpdateAppServiceRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Branch != nil {
		existing.Branch = *req.Branch
	}
	if req.RootDirectory != nil {
		existing.RootDirectory = *req.RootDirectory
	}
	if req.InstallCommand != nil {
		existing.InstallCommand = *req.InstallCommand
	}
	if req.BuildCommand != nil {
		existing.BuildCommand = *req.BuildCommand
	}
	if req.StartCommand != nil {
		existing.StartCommand = *req.StartCommand
	}
	if req.DockerfilePath != nil {
		existing.DockerfilePath = *req.DockerfilePath
	}
	if req.BuildEngine != nil {
		existing.BuildEngine = *req.BuildEngine
	}
	if req.InternalPort != nil {
		existing.InternalPort = *req.InternalPort
	}
	if req.Domain != nil {
		existing.Domain = *req.Domain
	}
	if req.StaticOutput != nil {
		existing.StaticOutput = *req.StaticOutput
	}
	if req.HealthCheckPath != nil {
		existing.HealthCheckPath = *req.HealthCheckPath
	}
	if req.CPULimit != nil {
		existing.CPULimit = *req.CPULimit
	}
	if req.MemoryLimit != nil {
		existing.MemoryLimit = *req.MemoryLimit
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	if err := h.service.UpdateAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Updated successfully", existing)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing app service id")
	}

	existing, err := h.service.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if !h.hasAccess(c, existing.ProjectID, true) {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.DeleteAppService(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
