package deployments

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	deploymentservices "codedock/internal/services/deployments"
	projectservices "codedock/internal/services/projects"
	"codedock/internal/utils"
)

type WebhookHandler struct {
	gitService        *deploymentservices.GitService
	projectService    *projectservices.ProjectService
	appService        *projectservices.AppService
	deploymentService *deploymentservices.DeploymentService
	prPreviewService  *deploymentservices.PRPreviewService
	gitAppsService    *deploymentservices.GitAppsService
	dispatcher        *deploymentservices.WebhookDispatcher
	deliveryService   *deploymentservices.WebhookDeliveryService
}

func NewWebhookHandler(
	gitService *deploymentservices.GitService,
	projectService *projectservices.ProjectService,
	appService *projectservices.AppService,
	deploymentService *deploymentservices.DeploymentService,
	prPreviewService *deploymentservices.PRPreviewService,
	gitAppsService *deploymentservices.GitAppsService,
	dispatcher *deploymentservices.WebhookDispatcher,
	deliveryService *deploymentservices.WebhookDeliveryService,
) *WebhookHandler {
	return &WebhookHandler{
		gitService:        gitService,
		projectService:    projectService,
		appService:        appService,
		deploymentService: deploymentService,
		prPreviewService:  prPreviewService,
		gitAppsService:    gitAppsService,
		dispatcher:        dispatcher,
		deliveryService:   deliveryService,
	}
}

func (h *WebhookHandler) HandleServiceWebhook(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	appSvc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || appSvc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	token := c.QueryParam("token")
	if appSvc.DeployToken == "" || token != appSvc.DeployToken {
		return utils.Error(c, http.StatusUnauthorized, "invalid or missing deploy token")
	}

	clientIP := c.RealIP()
	userAgent := c.Request().UserAgent()

	depID := h.deployServiceAsync(appSvc)
	if h.deliveryService != nil {
		_, _ = h.deliveryService.Record(c.Request().Context(), models.WebhookDeliveryInput{
			ProjectID:  &appSvc.ProjectID,
			Source:     "deploy_token",
			Event:      "deploy",
			Outcome:    "dispatched",
			ActionRef:  &depID,
			ClientIP:   &clientIP,
			UserAgent:  &userAgent,
			StatusCode: intPtr(202),
		})
	}

	return utils.Accepted(c, fmt.Sprintf("triggering background build & rollout for service %s", appSvc.Name), nil)
}

func intPtr(v int) *int {
	return &v
}

func (h *WebhookHandler) deployServiceAsync(appSvc *models.AppService) string {
	depID := uuid.NewString()
	go func() {
		ctx := context.Background()
		dep := &models.Deployment{
			ID:            depID,
			ServiceID:     appSvc.ID,
			EnvironmentID: appSvc.EnvironmentID,
			ProjectID:     appSvc.ProjectID,
			Status:        "BUILDING",
			Branch:        appSvc.Branch,
			Trigger:       "Git Webhook Push",
			BuildLogs:     fmt.Sprintf("Initiating automated build from %s branch %s...\n", appSvc.RepositoryURL, appSvc.Branch),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		var err error
		dep, err = h.deploymentService.CreateDeployment(ctx, dep)
		if err != nil {
			slog.Error("failed to create deployment from webhook", "error", err)
			return
		}
		h.deploymentService.ExecuteDeploymentAsync(dep)
	}()
	return depID
}

func (h *WebhookHandler) HandleGitHubWebhook(c echo.Context) error {
	event := c.Request().Header.Get("X-GitHub-Event")
	if event == "" {
		return utils.Error(c, http.StatusBadRequest, "missing X-GitHub-Event header")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "failed to read body")
	}

	signature := c.Request().Header.Get("X-Hub-Signature-256")
	valid, err := h.dispatcher.VerifySignature(c.Request().Context(), bodyBytes, signature)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to verify signature")
	}
	if !valid {
		return utils.Error(c, http.StatusUnauthorized, "invalid webhook signature")
	}

	deliveryID := c.Request().Header.Get("X-GitHub-Delivery")
	clientIP := c.RealIP()
	userAgent := c.Request().UserAgent()
	serviceID := c.Param("serviceId")

	var payload models.GithubWebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if event == "ping" {
		return utils.Success(c, "Pong", map[string]string{"message": "Pong"})
	}

	if event == "push" {
		result, err := h.dispatcher.DispatchPush(c.Request().Context(), deliveryID, serviceID, payload, &clientIP, &userAgent)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		if result.Duplicate {
			return utils.Success(c, "Duplicate delivery ignored", map[string]string{"message": "Duplicate delivery ignored"})
		}
		if result.FailedCount > 0 {
			msg := fmt.Sprintf("Triggered %d, %d failed: %s", result.HandledCount, result.FailedCount, strings.Join(result.Errors, "; "))
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": msg,
			})
		}
		return utils.Success(c, fmt.Sprintf("Triggered %d deployments", result.HandledCount), map[string]any{
			"handledCount": result.HandledCount,
		})
	}

	if event == "pull_request" && serviceID != "" {
		return h.handlePullRequest(c, serviceID, payload)
	}

	return utils.Success(c, "Operation successful", map[string]string{"message": "Event ignored"})
}

func (h *WebhookHandler) handlePullRequest(c echo.Context, serviceID string, payload models.GithubWebhookPayload) error {
	appSvc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
	if err != nil || appSvc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	switch payload.Action {
	case "opened", "synchronize", "reopened":
		if !appSvc.EnablePRPreviews {
			return utils.Success(c, "PR previews are disabled for this service", nil)
		}
		prInfo := deploymentservices.PullRequestInfo{
			Number:   payload.Number,
			State:    payload.PullRequest.State,
			Draft:    payload.PullRequest.Draft,
			HeadSHA:  payload.PullRequest.Head.Sha,
			HeadRef:  payload.PullRequest.Head.Ref,
			HeadRepo: payload.PullRequest.Head.Repo.FullName,
			BaseRef:  payload.PullRequest.Base.Ref,
			BaseRepo: payload.PullRequest.Base.Repo.FullName,
		}
		gitOwner, gitRepo := deploymentservices.ParseRepoOwnerAndName(appSvc.RepositoryURL)
		if eligible, reason := deploymentservices.IsPreviewEligible(prInfo, gitOwner, gitRepo, appSvc.Branch); !eligible {
			return utils.Success(c, "PR is not eligible for preview: "+reason, nil)
		}
		go func() {
			ctx := context.Background()
			opts := deploymentservices.DeployPRPreviewOpts{
				AppID:      serviceID,
				PRNumber:   payload.Number,
				CommitHash: payload.PullRequest.Head.Sha,
				Branch:     payload.PullRequest.Head.Ref,
			}
			_, _ = h.prPreviewService.DeployPRPreview(ctx, opts)
		}()
		return utils.Accepted(c, "Deploying PR preview", nil)
	case "closed":
		go func() {
			ctx := context.Background()
			_ = h.prPreviewService.DestroyPRPreview(ctx, serviceID, payload.Number)
		}()
		return utils.Accepted(c, "Destroying PR preview", nil)
	}
	return utils.Success(c, "PR action ignored", nil)
}

func (h *WebhookHandler) ListProjectDeliveries(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id required")
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	cursor := c.QueryParam("cursor")

	deliveries, nextCursor, err := h.deliveryService.ListByProject(c.Request().Context(), projectID, limit, cursor)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Deliveries retrieved", map[string]any{
		"deliveries": deliveries,
		"nextCursor": nextCursor,
	})
}
