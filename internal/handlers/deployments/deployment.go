package deployments

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
	deploymentservices "codedock.run/codedock/internal/services/deployments"
	projectservices "codedock.run/codedock/internal/services/projects"
)

type DeploymentHandler struct {
	deploymentService *deploymentservices.DeploymentService
	appService        *projectservices.AppService
	auditService      *authservices.AuditService
	aiAnalysis        *projectservices.AIAnalysisService
	prPreviewService  *deploymentservices.PRPreviewService
	projectService    *projectservices.ProjectService
}

func NewDeploymentHandler(ds *deploymentservices.DeploymentService, as *projectservices.AppService, audit *authservices.AuditService, aiAnalysis *projectservices.AIAnalysisService, prp *deploymentservices.PRPreviewService, ps *projectservices.ProjectService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: ds,
		appService:        as,
		auditService:      audit,
		aiAnalysis:        aiAnalysis,
		prPreviewService:  prp,
		projectService:    ps,
	}
}

func (h *DeploymentHandler) verifyProjectOwnership(c echo.Context, projectID string) error {
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

func (h *DeploymentHandler) verifyProjectAdmin(c echo.Context, projectID string) error {
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

	if !h.projectService.HasPermission(c.Request().Context(), projectID, user.UserID, models.UserRole(user.Role), models.MemberPermissionAdmin) {
		return utils.Error(c, http.StatusForbidden, "admin access required")
	}
	return nil
}

func (h *DeploymentHandler) ListServiceDeployments(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	if err := h.verifyProjectOwnership(c, svc.ProjectID); err != nil {
		return err
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	deps, total, err := h.deploymentService.ListByService(c.Request().Context(), serviceID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, "Deployments retrieved", deps, total, page, limit)
}

func (h *DeploymentHandler) ListOrganizationDeployments(c echo.Context) error {
	organizationID := c.QueryParam("organizationId")
	if organizationID == "" {
		return utils.Error(c, http.StatusBadRequest, "organizationId is required")
	}
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if !h.projectService.HasOrgPermission(c.Request().Context(), organizationID, user.UserID, models.UserRole(user.Role), "") {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions for this organization")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	filter := models.DeploymentListFilter{
		OrganizationID: organizationID,
		ProjectID:      c.QueryParam("projectId"),
		ServiceID:      c.QueryParam("serviceId"),
		Status:         c.QueryParam("status"),
		Search:         c.QueryParam("search"),
		Limit:          limit,
		Offset:         (page - 1) * limit,
	}
	deployments, total, err := h.deploymentService.ListByOrganization(c.Request().Context(), filter)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, "Deployments retrieved", deployments, total, page, limit)
}

func (h *DeploymentHandler) Trigger(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	dep := &models.Deployment{
		ServiceID:     serviceID,
		EnvironmentID: svc.EnvironmentID,
		ProjectID:     svc.ProjectID,
		Status:        "BUILDING",
		Branch:        svc.Branch,
		Trigger:       "Manual Deploy",
		BuildLogs:     "Initiating build...\n",
	}
	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), dep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "deployment.trigger",
		Resource:  serviceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"deploymentId": created.ID,
		},
	})

	return utils.Accepted(c, "Deployment created", created)
}

func (h *DeploymentHandler) Rollback(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	targetDep, err := h.deploymentService.GetDeployment(c.Request().Context(), id)
	if err != nil || targetDep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}

	if err := h.verifyProjectAdmin(c, targetDep.ProjectID); err != nil {
		return err
	}

	newDep := &models.Deployment{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        "BUILDING",
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Rollback to " + targetDep.ID,
		Branch:        targetDep.Branch,
		Trigger:       "Rollback",
		BuildLogs:     "Rolling back to deployment " + targetDep.ID + "...\n",
	}
	created, err := h.deploymentService.CreateDeployment(c.Request().Context(), newDep)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.deploymentService.ExecuteDeploymentAsync(created)

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "deployment.rollback",
		Resource:  newDep.ServiceID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"deploymentId": created.ID,
		},
	})

	return utils.Accepted(c, "Rollback created", created)
}

func (h *DeploymentHandler) GetLogs(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	dep, err := h.deploymentService.GetDeployment(c.Request().Context(), id)
	if err != nil || dep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}
	if err := h.verifyProjectOwnership(c, dep.ProjectID); err != nil {
		return err
	}

	return utils.Success(c, "Logs fetched successfully", map[string]string{
		"id":        dep.ID,
		"buildLogs": dep.BuildLogs,
		"status":    string(dep.Status),
	})
}

func (h *DeploymentHandler) GetMetrics(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "serviceId is required")
	}

	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	if err := h.verifyProjectOwnership(c, svc.ProjectID); err != nil {
		return err
	}

	health, err := h.deploymentService.GetMetrics(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	now := time.Now().UTC()
	metrics := []map[string]any{
		{
			"timestamp":  now.Format(time.RFC3339),
			"cpuPercent": health.CPUUsagePercentage,
			"memoryMB":   float64(health.MemoryUsageBytes) / 1024 / 1024,
			"status":     health.Status,
			"uptime":     health.UptimeSeconds,
		},
	}
	return utils.Success(c, "Operation successful", metrics)
}

func (h *DeploymentHandler) ExplainFailure(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	dep, err := h.deploymentService.GetDeployment(c.Request().Context(), id)
	if err != nil || dep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}

	if err := h.verifyProjectOwnership(c, dep.ProjectID); err != nil {
		return err
	}

	explanation, err := h.aiAnalysis.ExplainDeploymentFailure(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "AI Analysis completed", explanation)
}

func (h *DeploymentHandler) ListPRPreviews(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}
	if err := h.verifyProjectOwnership(c, svc.ProjectID); err != nil {
		return err
	}

	previews, err := h.prPreviewService.ListByApp(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", previews)
}
