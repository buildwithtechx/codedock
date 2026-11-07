package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/engine/deploy"
	"codedock/internal/http/middleware"
	authservices "codedock/internal/services/auth"
	deploymentservices "codedock/internal/services/deployments"
	projectservices "codedock/internal/services/projects"
	"codedock/internal/services/runtimes"
)

type AppHandler struct {
	Runtime           *runtimes.Service
	appService        *projectservices.AppService
	projectService    *projectservices.ProjectService
	deployer          *deploy.Deployer
	deploymentService *deploymentservices.DeploymentService
	envService        *projectservices.EnvironmentService
	auditService      *authservices.AuditService
}

func NewAppHandler(s *projectservices.AppService, ps *projectservices.ProjectService, d *deploy.Deployer, ds *deploymentservices.DeploymentService, es *projectservices.EnvironmentService, auditService *authservices.AuditService) *AppHandler {
	return &AppHandler{
		appService:        s,
		projectService:    ps,
		deployer:          d,
		deploymentService: ds,
		envService:        es,
		auditService:      auditService,
	}
}

func (h *AppHandler) verifyProjectOwnership(c echo.Context, projectID string) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	if user.Role == "api" {
		tokenProjectID, ok := c.Get("project_id").(string)
		if ok && tokenProjectID != projectID {
			return echo.NewHTTPError(http.StatusForbidden, "token does not have access to this project")
		}
	}

	project, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil || project == nil {
		return echo.NewHTTPError(http.StatusNotFound, "project not found")
	}

	if !h.projectService.IsMemberOrOwner(c.Request().Context(), projectID, user.UserID, user.Role) {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}
	return nil
}
