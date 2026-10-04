package deployments

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
	"codedock.run/codedock/internal/utils"
)

func (h *DeploymentHandler) TriggerProject(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId parameter")
	}

	if err := h.verifyProjectAdmin(c, projectID); err != nil {
		return err
	}

	services, err := h.appService.ListByProject(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	envID := c.QueryParam("environmentId")
	if envID == "" {
		envID = c.QueryParam("environment_id")
	}

	var triggered []*models.Deployment
	for _, svc := range services {
		if envID != "" && svc.EnvironmentID != envID {
			continue
		}
		dep := &models.Deployment{
			ServiceID:     svc.ID,
			EnvironmentID: svc.EnvironmentID,
			ProjectID:     svc.ProjectID,
			Status:        "BUILDING",
			Branch:        svc.Branch,
			Trigger:       "Manual Project Deploy",
			BuildLogs:     "Initiating build...\n",
		}
		created, err := h.deploymentService.CreateDeployment(c.Request().Context(), dep)
		if err == nil {
			h.deploymentService.ExecuteDeploymentAsync(created)
			triggered = append(triggered, created)
		}
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    "system",
		Action:    "project.deployment.trigger",
		Resource:  projectID,
		IPAddress: c.RealIP(),
		Details: map[string]string{
			"count": fmt.Sprintf("%d", len(triggered)),
		},
	})

	return utils.Accepted(c, "Deployments triggered", triggered)
}

func (h *DeploymentHandler) ListProjectDeployments(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId parameter")
	}

	if err := h.verifyProjectOwnership(c, projectID); err != nil {
		return err
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
		ProjectID: projectID,
		ServiceID: c.QueryParam("serviceId"),
		Status:    c.QueryParam("status"),
		Search:    c.QueryParam("search"),
		Limit:     limit,
		Offset:    (page - 1) * limit,
	}

	deps, total, err := h.deploymentService.ListByOrganization(c.Request().Context(), filter)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Paginated(c, "Deployments retrieved", deps, total, page, limit)
}
