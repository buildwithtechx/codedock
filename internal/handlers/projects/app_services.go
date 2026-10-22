package projects

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"codedock/internal/http/middleware"
	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	"codedock/internal/telemetry"
	"codedock/internal/utils"
)

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
	userClaims := middleware.GetUserClaimsFromContext(c.Request().Context())
	if userClaims == nil || !h.projectService.HasPermission(c.Request().Context(), req.ProjectID, userClaims.UserID, models.UserRole(userClaims.Role), models.MemberPermissionAdmin) {
		return utils.Error(c, http.StatusForbidden, "project admin access required")
	}
	if req.InternalPort < 0 || req.InternalPort > 65535 {
		return utils.Error(c, http.StatusBadRequest, "invalid internal port")
	}
	if req.ImageRef != "" && req.RepositoryURL != "" {
		return utils.Error(c, http.StatusBadRequest, "choose exactly one repository or image")
	}
	if req.CPULimit < 0 || req.MemoryLimit < 0 {
		return utils.Error(c, http.StatusBadRequest, "resource limits cannot be negative")
	}
	if req.SwapDisabled && req.MemoryLimit <= 0 {
		return utils.Error(c, http.StatusBadRequest, "disabling swap requires a positive memory limit")
	}
	if req.RuntimeMode == "docker" {
		req.RuntimeMode = models.RuntimeModeWeb
	}
	if req.RuntimeMode != "" && req.RuntimeMode != models.RuntimeModeWeb && req.RuntimeMode != models.RuntimeModeWorker {
		return utils.Error(c, http.StatusBadRequest, "unsupported runtime")
	}
	env, err := h.envService.GetEnvironment(c.Request().Context(), envID)
	if err != nil || env == nil {
		return utils.Error(c, http.StatusNotFound, "environment not found")
	}
	if env.ProjectID != req.ProjectID {
		return utils.Error(c, http.StatusBadRequest, "environment does not belong to specified project")
	}
	req.EnvironmentID = envID
	if req.ID != "" {
		if _, err := uuid.Parse(req.ID); err != nil {
			return utils.Error(c, http.StatusBadRequest, "application ID must be a UUID")
		}
	}
	req.GitUserID = userClaims.UserID
	if req.InternalPort == 0 {
		req.InternalPort = 3000
	}
	version, err := normalizeRuntimeVersion(req.RuntimeVersion)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	req.RuntimeVersion = version
	if req.RuntimeMode == "" {
		req.RuntimeMode = models.RuntimeModeWeb
	}
	if err := h.prepareApplicationDomain(c.Request().Context(), &req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	created, err := h.appService.CreateApplicationSetup(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if created.RuntimeMode == models.RuntimeModeWeb {
		if err := h.ensureApplicationDomain(c.Request().Context(), created); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
	}

	for _, role := range req.Roles {
		if role.Command == "" {
			continue
		}
		sibling := models.AppService{
			ProjectID:              created.ProjectID,
			EnvironmentID:          created.EnvironmentID,
			Name:                   created.Name + "-" + role.Name,
			RepositoryURL:          created.RepositoryURL,
			ImageRef:               created.ImageRef,
			Branch:                 created.Branch,
			RootDirectory:          created.RootDirectory,
			Icon:                   created.Icon,
			GitUserID:              created.GitUserID,
			RuntimeMode:            models.RuntimeModeWorker,
			InstallCommand:         created.InstallCommand,
			BuildCommand:           created.BuildCommand,
			StartCommand:           role.Command,
			DockerfilePath:         created.DockerfilePath,
			BuildEngine:            created.BuildEngine,
			InternalPort:           0,
			Optimize:               created.Optimize,
			Stack:                  created.Stack,
			StopGracePeriodSeconds: roleGraceSeconds(role.StopGracePeriod),
			RuntimeVersion:         created.RuntimeVersion,
		}
		if _, err := h.appService.CreateAppService(c.Request().Context(), &sibling); err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
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
	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "app.create",
		Resource:  created.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"name": created.Name, "project": req.ProjectID},
	})

	return utils.Created(c, "Created successfully", created)
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
	release, err := h.guardAppOperation(c, id)
	if err != nil {
		return err
	}
	defer release()
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
	userClaims := middleware.GetUserClaimsFromContext(c.Request().Context())
	if userClaims != nil {
		existing.GitUserID = userClaims.UserID
	}
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
	if req.SwapDisabled && req.MemoryLimit <= 0 {
		return utils.Error(c, http.StatusBadRequest, "disabling swap requires a positive memory limit")
	}
	existing.CPULimit = req.CPULimit
	existing.MemoryLimit = req.MemoryLimit
	existing.SwapDisabled = req.SwapDisabled
	existing.DeployToken = req.DeployToken
	existing.MaintenanceMode = req.MaintenanceMode
	existing.Optimize = req.Optimize
	existing.Stack = req.Stack
	existing.StopGracePeriodSeconds = req.StopGracePeriodSeconds
	version, err := normalizeRuntimeVersion(req.RuntimeVersion)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	existing.RuntimeVersion = version
	existing.EnablePRPreviews = req.EnablePRPreviews
	existing.ImageRef = req.ImageRef
	existing.RegistryID = req.RegistryID
	if err := h.appService.UpdateAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    utils.AuditActorID(c),
		Action:    "app.update",
		Resource:  id,
		IPAddress: c.RealIP(),
	})
	return utils.Success(c, "Operation successful", existing)
}
