package git

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type WebhookHandler struct {
	dispatcher      *WebhookDispatcher
	deliveryService *WebhookDeliveryService
	prService       *PRPreviewService
	serviceLookup   ServiceLookup
}

func NewWebhookHandler(
	dispatcher *WebhookDispatcher,
	deliveryService *WebhookDeliveryService,
	prService *PRPreviewService,
	serviceLookup ServiceLookup,
) *WebhookHandler {
	return &WebhookHandler{
		dispatcher:      dispatcher,
		deliveryService: deliveryService,
		prService:       prService,
		serviceLookup:   serviceLookup,
	}
}

func (h *WebhookHandler) HandleServiceWebhook(c echo.Context) error {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}
	if h.serviceLookup == nil {
		return utils.Error(c, http.StatusInternalServerError, "service lookup not available")
	}
	svc, err := h.serviceLookup.GetService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	token := c.QueryParam("token")
	if svc.DeployToken == "" || token != svc.DeployToken {
		return utils.Error(c, http.StatusUnauthorized, "invalid or missing deploy token")
	}

	clientIP := c.RealIP()
	userAgent := c.Request().UserAgent()

	if h.dispatcher != nil && h.dispatcher.deploymentTrigger != nil {
		go func() {
			ctx := context.Background()
			logs := fmt.Sprintf("Initiating automated build from %s branch %s (deploy token)...\n", svc.RepositoryURL, svc.Branch)
			depID, _ := h.dispatcher.deploymentTrigger.TriggerDeployment(ctx, svc.ID, svc.EnvironmentID, svc.ProjectID, svc.Branch, "Deploy Token", logs)
			if h.deliveryService != nil {
				stCode := 202
				_, _ = h.deliveryService.Record(ctx, WebhookDeliveryInput{
					ProjectID:  &svc.ProjectID,
					Source:     WebhookSourceDeployToken,
					Event:      WebhookEventDeploy,
					Outcome:    WebhookOutcomeDispatched,
					ActionRef:  &depID,
					ClientIP:   &clientIP,
					UserAgent:  &userAgent,
					StatusCode: &stCode,
				})
			}
		}()
	}

	return utils.Accepted(c, fmt.Sprintf("triggering background build & rollout for service %s", svc.Name), nil)
}

func (h *WebhookHandler) HandleGitHubWebhook(c echo.Context) error {
	event := c.Request().Header.Get("X-GitHub-Event")
	if event == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingWebhookEvent.Error())
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
		return utils.Error(c, http.StatusUnauthorized, ErrInvalidSignature.Error())
	}

	deliveryID := c.Request().Header.Get("X-GitHub-Delivery")
	clientIP := c.RealIP()
	userAgent := c.Request().UserAgent()
	serviceID := c.Param("serviceId")

	var payload GithubWebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if event == WebhookEventPing {
		return utils.Success(c, "Pong", MessageResponse{Message: "Pong"})
	}

	if event == WebhookEventPush {
		result, err := h.dispatcher.DispatchPush(c.Request().Context(), deliveryID, serviceID, payload, &clientIP, &userAgent)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		if result.Duplicate {
			return utils.Success(c, "Duplicate delivery ignored", MessageResponse{Message: "Duplicate delivery ignored"})
		}
		if result.FailedCount > 0 {
			msg := fmt.Sprintf("Triggered %d, %d failed: %s", result.HandledCount, result.FailedCount, strings.Join(result.Errors, "; "))
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": msg,
			})
		}
		return utils.Success(c, fmt.Sprintf("Triggered %d deployments", result.HandledCount), result)
	}

	if event == WebhookEventPullRequest && serviceID != "" {
		return h.handlePullRequest(c, serviceID, payload)
	}

	return utils.Success(c, "Operation successful", MessageResponse{Message: "Event ignored"})
}

func (h *WebhookHandler) handlePullRequest(c echo.Context, serviceID string, payload GithubWebhookPayload) error {
	if h.serviceLookup == nil {
		return utils.Error(c, http.StatusInternalServerError, "service lookup unavailable")
	}
	svc, err := h.serviceLookup.GetService(c.Request().Context(), serviceID)
	if err != nil || svc == nil {
		return utils.Error(c, http.StatusNotFound, "service not found")
	}

	switch payload.Action {
	case "opened", "synchronize", "reopened":
		if !svc.EnablePRPreviews {
			return utils.Success(c, "PR previews are disabled for this service", nil)
		}
		prInfo := PullRequestInfo{
			Number:   payload.Number,
			State:    payload.PullRequest.State,
			Draft:    payload.PullRequest.Draft,
			HeadSHA:  payload.PullRequest.Head.Sha,
			HeadRef:  payload.PullRequest.Head.Ref,
			HeadRepo: payload.PullRequest.Head.Repo.FullName,
			BaseRef:  payload.PullRequest.Base.Ref,
			BaseRepo: payload.PullRequest.Base.Repo.FullName,
		}
		gitOwner, gitRepo := ParseRepoOwnerAndName(svc.RepositoryURL)
		if eligible, reason := IsPreviewEligible(prInfo, gitOwner, gitRepo, svc.Branch); !eligible {
			return utils.Success(c, "PR is not eligible for preview: "+reason, nil)
		}
		if h.prService != nil {
			go func() {
				ctx := context.Background()
				previewDomain := GeneratePreviewDomain("", "", svc.Name, svc.ProjectID, payload.Number)
				preview, _ := h.prService.CreatePreviewRecord(ctx, serviceID, svc.ProjectID, payload.PullRequest.Head.Ref, payload.PullRequest.Head.Sha, previewDomain, payload.Number)
				if preview != nil && h.dispatcher != nil && h.dispatcher.deploymentTrigger != nil {
					logs := fmt.Sprintf("Deploying PR preview for #%d (commit %s)...\n", payload.Number, payload.PullRequest.Head.Sha)
					_, _ = h.dispatcher.deploymentTrigger.TriggerDeployment(ctx, serviceID, svc.EnvironmentID, svc.ProjectID, payload.PullRequest.Head.Ref, "PR Preview", logs)
				}
			}()
		}
		return utils.Accepted(c, "Deploying PR preview", nil)
	case "closed":
		if h.prService != nil {
			go func() {
				ctx := context.Background()
				_ = h.prService.DestroyPRPreview(ctx, serviceID, payload.Number)
			}()
		}
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

	return utils.Success(c, "Deliveries retrieved", DeliveriesResponse{
		Deliveries: deliveries,
		NextCursor: nextCursor,
	})
}
