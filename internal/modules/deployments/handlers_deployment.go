package deployments

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type DeploymentHandler struct {
	service       *DeploymentService
	appLookup     AppServiceLookup
	projectLookup ProjectLookup
	aiExplainer   AIExplainer
	prLookup      PRPreviewLookup
}

func NewDeploymentHandler(
	service *DeploymentService,
	appLookup AppServiceLookup,
	projectLookup ProjectLookup,
	aiExplainer AIExplainer,
	prLookup PRPreviewLookup,
) *DeploymentHandler {
	return &DeploymentHandler{
		service:       service,
		appLookup:     appLookup,
		projectLookup: projectLookup,
		aiExplainer:   aiExplainer,
		prLookup:      prLookup,
	}
}

func (h *DeploymentHandler) ListOrganizationDeployments(c echo.Context) error {
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.Param("id")
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	days, _ := strconv.Atoi(c.QueryParam("days"))

	filter := DeploymentListFilter{
		OrganizationID: orgID,
		ProjectID:      c.QueryParam("projectId"),
		ServiceID:      c.QueryParam("serviceId"),
		Status:         c.QueryParam("status"),
		Search:         c.QueryParam("search"),
		Days:           days,
		Limit:          limit,
		Offset:         offset,
	}

	items, total, err := h.service.ListByOrganization(c.Request().Context(), filter)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Deployments retrieved", map[string]any{
		"items": items,
		"total": total,
	})
}

func (h *DeploymentHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	dep, err := h.service.GetDeployment(c.Request().Context(), id)
	if err != nil || dep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}

	return utils.Success(c, "Deployment retrieved", dep)
}

func (h *DeploymentHandler) ListProjectDeployments(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing project id")
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	filter := DeploymentListFilter{
		ProjectID: projectID,
		Limit:     limit,
		Offset:    offset,
	}
	items, total, err := h.service.ListByOrganization(c.Request().Context(), filter)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Project deployments retrieved", map[string]any{
		"items": items,
		"total": total,
	})
}

func (h *DeploymentHandler) ListServiceDeployments(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	deps, total, err := h.service.ListByService(c.Request().Context(), serviceID, limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Service deployments retrieved", map[string]any{
		"items": deps,
		"total": total,
	})
}

func (h *DeploymentHandler) GetLogs(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	dep, err := h.service.GetDeployment(c.Request().Context(), id)
	if err != nil || dep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}

	return utils.Success(c, "Logs retrieved", LogsResponse{
		DeploymentID: dep.ID,
		BuildLogs:    dep.BuildLogs,
	})
}

func (h *DeploymentHandler) ExplainFailure(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	dep, err := h.service.GetDeployment(c.Request().Context(), id)
	if err != nil || dep == nil {
		return utils.Error(c, http.StatusNotFound, "deployment not found")
	}

	if h.aiExplainer != nil {
		explanation, err := h.aiExplainer.ExplainDeploymentFailure(c.Request().Context(), id)
		if err == nil && explanation != nil {
			return utils.Success(c, "AI Analysis completed", explanation)
		}
	}

	return utils.Success(c, "Analysis completed", FailureExplanation{
		DeploymentID: dep.ID,
		RootCause:    dep.ResourceIssue,
		Explanation:  "Deployment failed during execution. Inspect build logs for detailed error stack.",
	})
}

func (h *DeploymentHandler) GetMetrics(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	health, err := h.service.GetMetrics(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Metrics retrieved", health)
}

func (h *DeploymentHandler) ListPRPreviews(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	if h.prLookup != nil {
		previews, err := h.prLookup.ListByApp(c.Request().Context(), serviceID)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return utils.Success(c, "PR previews retrieved", previews)
	}

	return utils.Success(c, "PR previews retrieved", []any{})
}
