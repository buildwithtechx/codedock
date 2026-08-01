package projects

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	deploymentservices "codedock.run/codedock/internal/services/deployments"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/telemetry"
)

type AppHandler struct {
	appService        *projectservices.AppService
	projectService    *projectservices.ProjectService
	deployer          *deploy.Deployer
	deploymentService *deploymentservices.DeploymentService
	envService        *projectservices.EnvironmentService
}

func NewAppHandler(s *projectservices.AppService, ps *projectservices.ProjectService, d *deploy.Deployer, ds *deploymentservices.DeploymentService, es *projectservices.EnvironmentService) *AppHandler {
	return &AppHandler{
		appService:        s,
		projectService:    ps,
		deployer:          d,
		deploymentService: ds,
		envService:        es,
	}
}

func (h *AppHandler) verifyProjectOwnership(c echo.Context, projectID string) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	if user.Role == "api" {
		tokenProjectID, ok := c.Get("project_id").(string)
		if ok && tokenProjectID != projectID {
			return utils.Error(c, http.StatusForbidden, "token does not have access to this project")
		}
	}

	project, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil || project == nil {
		return utils.Error(c, http.StatusNotFound, "project not found")
	}

	if !h.projectService.IsMemberOrOwner(c.Request().Context(), projectID, user.UserID, user.Role) {
		return utils.Error(c, http.StatusForbidden, "access denied")
	}
	return nil
}

func (h *AppHandler) Create(c echo.Context) error {
	envID := c.Param("id")
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "app service name is required")
	}
	if req.HealthCheckPath != "" && !isValidHealthCheckPath(req.HealthCheckPath) {
		return utils.Error(c, http.StatusBadRequest, "invalid health check path")
	}
	if err := h.verifyProjectOwnership(c, req.ProjectID); err != nil {
		return err
	}
	env, err := h.envService.GetEnvironment(c.Request().Context(), envID)
	if err != nil || env == nil {
		return utils.Error(c, http.StatusNotFound, "environment not found")
	}
	if env.ProjectID != req.ProjectID {
		return utils.Error(c, http.StatusBadRequest, "environment does not belong to specified project")
	}
	req.EnvironmentID = envID
	if req.InternalPort == 0 {
		req.InternalPort = 3000
	}
	if req.RuntimeMode == "" {
		req.RuntimeMode = models.RuntimeModeWeb
	}
	created, err := h.appService.CreateAppService(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	generatedDomain := utils.GenerateAppDomain(req.Name, "", "")
	parsedDomain, parseErr := url.Parse(generatedDomain)
	if parseErr != nil || parsedDomain.Hostname() == "" {
		slog.Warn("failed to parse generated app domain", "domain", generatedDomain)
		generatedDomain = strings.TrimPrefix(strings.TrimPrefix(generatedDomain, "https://"), "http://")
		generatedDomain = strings.Split(generatedDomain, "/")[0]
	} else {
		generatedDomain = parsedDomain.Hostname()
	}
	if _, err := h.envService.CreateDomain(c.Request().Context(), &models.DomainConfig{
		ServiceID:  created.ID,
		DomainName: generatedDomain,
	}); err != nil {
		slog.Warn("failed to create default domain", "error", err)
	}

	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	distinctID := "anonymous"
	if user != nil {
		distinctID = user.Email
	}
	sourceType := "github"
	if created.ImageRef != "" {
		sourceType = "docker_image"
	}
	telemetry.Track(distinctID, "app_created", map[string]any{
		"app_id": created.ID,
		"name":   created.Name,
		"type":   sourceType,
	})

	return utils.Created(c, "Created successfully", created)
}

func (h *AppHandler) ListByEnvironment(c echo.Context) error {
	envID := c.Param("id")
	apps, err := h.appService.ListByEnvironment(c.Request().Context(), envID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user != nil && user.Role != "admin" {
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
	if err := h.verifyProjectOwnership(c, projectID); err != nil {
		return err
	}
	apps, err := h.appService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", apps)
}

func (h *AppHandler) Get(c echo.Context) error {
	id := c.Param("id")
	svc, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, svc.ProjectID); err != nil {
		return err
	}
	return utils.Success(c, "Operation successful", svc)
}

func (h *AppHandler) Update(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	var req models.AppService
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.HealthCheckPath != "" && !isValidHealthCheckPath(req.HealthCheckPath) {
		return utils.Error(c, http.StatusBadRequest, "invalid health check path")
	}
	targetProjectID := existing.ProjectID
	if req.ProjectID != "" && req.ProjectID != existing.ProjectID {
		if err := h.verifyProjectOwnership(c, req.ProjectID); err != nil {
			return err
		}
		targetProjectID = req.ProjectID
	}
	targetEnvID := existing.EnvironmentID
	if req.EnvironmentID != "" {
		targetEnvID = req.EnvironmentID
	}
	if targetEnvID != "" {
		env, err := h.envService.GetEnvironment(c.Request().Context(), targetEnvID)
		if err != nil || env == nil {
			return utils.Error(c, http.StatusNotFound, "environment not found")
		}
		if env.ProjectID != targetProjectID {
			return utils.Error(c, http.StatusBadRequest, "environment does not belong to target project")
		}
	}
	existing.Name = req.Name
	existing.ProjectID = targetProjectID
	existing.EnvironmentID = targetEnvID
	existing.RepositoryURL = req.RepositoryURL
	existing.Branch = req.Branch
	existing.RootDirectory = req.RootDirectory
	existing.BuildCommand = req.BuildCommand
	existing.StartCommand = req.StartCommand
	existing.InstallCommand = req.InstallCommand
	existing.DockerfilePath = req.DockerfilePath
	existing.BuildEngine = req.BuildEngine
	existing.InternalPort = req.InternalPort
	existing.RuntimeMode = req.RuntimeMode
	existing.Domain = req.Domain
	existing.StaticOutput = req.StaticOutput
	existing.HealthCheckPath = req.HealthCheckPath
	existing.ContainerID = req.ContainerID
	existing.Status = req.Status
	existing.CPULimit = req.CPULimit
	existing.MemoryLimit = req.MemoryLimit
	existing.DeployToken = req.DeployToken
	existing.MaintenanceMode = req.MaintenanceMode
	existing.EnablePRPreviews = req.EnablePRPreviews
	existing.ImageRef = req.ImageRef
	existing.RegistryID = req.RegistryID
	if err := h.appService.UpdateAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", existing)
}

func (h *AppHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}

	if err := h.deployer.StopAppService(c.Request().Context(), existing); err != nil {
		slog.Warn("failed to stop app service", "error", err)
	}

	if err := h.appService.DeleteAppService(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AppHandler) StopService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.deployer.StopAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	existing.Status = models.AppServiceStatusStopped
	_ = h.appService.UpdateAppService(c.Request().Context(), existing)
	return utils.Success(c, "Service stopped successfully", existing)
}

func (h *AppHandler) RedeployService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}

	newDep := &models.Deployment{
		ServiceID:     existing.ID,
		EnvironmentID: existing.EnvironmentID,
		ProjectID:     existing.ProjectID,
		Status:        "BUILDING",
		CommitMessage: "Manual Redeploy",
		Branch:        existing.Branch,
		Trigger:       "Manual Redeploy",
	}
	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), newDep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)
	return utils.Accepted(c, "Redeployment triggered", created)
}

func (h *AppHandler) RestartService(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.deployer.RestartAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	existing.Status = models.AppServiceStatusRunning
	_ = h.appService.UpdateAppService(c.Request().Context(), existing)
	return utils.Success(c, "Service restarted successfully", existing)
}

func isValidHealthCheckPath(path string) bool {
	if path == "" {
		return true
	}
	if path[0] != '/' {
		return false
	}
	for _, ch := range path {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '/' || ch == '-' || ch == '_' || ch == '.' || ch == '?' || ch == '=' || ch == '&') {
			return false
		}
	}
	return true
}
