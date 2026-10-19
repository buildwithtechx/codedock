package deployments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/models"
	projectservices "codedock/internal/services/projects"
)

type WebhookDispatcher struct {
	deliveryService   *WebhookDeliveryService
	deploymentService *DeploymentService
	appService        *projectservices.AppService
	gitAppsService    *GitAppsService
	prPreviewService  *PRPreviewService
}

func NewWebhookDispatcher(
	deliveryService *WebhookDeliveryService,
	deploymentService *DeploymentService,
	appService *projectservices.AppService,
	gitAppsService *GitAppsService,
	prPreviewService *PRPreviewService,
) *WebhookDispatcher {
	return &WebhookDispatcher{
		deliveryService:   deliveryService,
		deploymentService: deploymentService,
		appService:        appService,
		gitAppsService:    gitAppsService,
		prPreviewService:  prPreviewService,
	}
}

func (d *WebhookDispatcher) VerifySignature(ctx context.Context, bodyBytes []byte, signature string) (bool, error) {
	if signature == "" {
		return false, nil
	}
	apps, err := d.gitAppsService.ListGithubApps(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to list github apps: %w", err)
	}
	for _, app := range apps {
		if app.WebhookSecret != "" && verifyHMAC(bodyBytes, app.WebhookSecret, signature) {
			return true, nil
		}
	}
	return false, nil
}

func verifyHMAC(payload []byte, secret, signature string) bool {
	if secret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}

type PushDispatchResult struct {
	Duplicate      bool
	HandledCount   int
	FailedCount    int
	Errors         []string
	HandledProjIDs []string
}

func (d *WebhookDispatcher) DispatchPush(
	ctx context.Context,
	deliveryID string,
	targetServiceID string,
	payload models.GithubWebhookPayload,
	clientIP *string,
	userAgent *string,
) (*PushDispatchResult, error) {
	var claim models.ClaimResult
	if deliveryID != "" {
		var err error
		claim, err = d.deliveryService.ClaimGithub(ctx, deliveryID, "push", clientIP, userAgent)
		if err != nil {
			slog.Warn("failed to claim github delivery, proceeding fail-open", "error", err)
		} else if !claim.Claimed {
			return &PushDispatchResult{Duplicate: true}, nil
		}
	}

	handledSet := make(map[string]bool)
	for _, id := range claim.HandledProjectIDs {
		handledSet[id] = true
	}

	repoFullName := strings.ToLower(payload.Repository.FullName)
	repoName := strings.ToLower(payload.Repository.Name)

	var targetServices []*models.AppService
	if targetServiceID != "" {
		svc, err := d.appService.GetAppService(ctx, targetServiceID)
		if err == nil && svc != nil {
			targetServices = append(targetServices, svc)
		}
	} else {
		allServices, err := d.appService.ListAllServices(ctx, "")
		if err == nil {
			for _, svc := range allServices {
				svcURL := strings.ToLower(svc.RepositoryURL)
				if (repoFullName != "" && strings.Contains(svcURL, repoFullName)) ||
					(repoName != "" && strings.Contains(svcURL, repoName)) {
					targetServices = append(targetServices, svc)
				}
			}
		}
	}

	result := &PushDispatchResult{
		HandledProjIDs: claim.HandledProjectIDs,
	}

	for _, svc := range targetServices {
		branch := svc.Branch
		if branch == "" {
			branch = payload.Repository.DefaultBranch
		}
		if branch == "" {
			branch = "main"
		}
		expectedRef := "refs/heads/" + branch
		if payload.Ref != expectedRef {
			continue
		}

		if handledSet[svc.ProjectID] {
			slog.Info("skipping already handled project from previous delivery attempt", "projectId", svc.ProjectID, "serviceId", svc.ID)
			continue
		}

		changedFiles := ExtractChangedFilesFromPush(payload)
		commitMsg := ""
		if payload.HeadCommit != nil {
			commitMsg = payload.HeadCommit.Message
		}
		impact := ResolveProjectChangeImpact(svc, changedFiles, commitMsg)
		if !impact.Affected {
			slog.Info("skipping auto-deploy: changes do not affect service", "serviceId", svc.ID, "name", svc.Name, "root", svc.RootDirectory, "reason", impact.Reason)
			_, _ = d.deliveryService.Record(ctx, models.WebhookDeliveryInput{
				ProjectID:  &svc.ProjectID,
				Source:     "github",
				DeliveryID: &deliveryID,
				Event:      "push",
				Outcome:    "skipped",
				Summary: map[string]any{
					"reason": impact.Reason,
					"root":   svc.RootDirectory,
				},
				ClientIP:  clientIP,
				UserAgent: userAgent,
			})
			continue
		}

		dep := &models.Deployment{
			ID:            uuid.NewString(),
			ServiceID:     svc.ID,
			EnvironmentID: svc.EnvironmentID,
			ProjectID:     svc.ProjectID,
			Status:        "BUILDING",
			Branch:        branch,
			Trigger:       "Git Webhook Push",
			BuildLogs:     fmt.Sprintf("Initiating automated build from %s branch %s (delivery %s)...\n", svc.RepositoryURL, branch, deliveryID),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}

		createdDep, err := d.deploymentService.CreateDeployment(ctx, dep)
		if err != nil {
			result.FailedCount++
			errMsg := fmt.Sprintf("service %s: %v", svc.Name, err)
			result.Errors = append(result.Errors, errMsg)
			_, _ = d.deliveryService.Record(ctx, models.WebhookDeliveryInput{
				ProjectID:  &svc.ProjectID,
				Source:     "github",
				DeliveryID: &deliveryID,
				Event:      "push",
				Outcome:    "failed",
				Error:      &errMsg,
				ClientIP:   clientIP,
				UserAgent:  userAgent,
			})
			continue
		}

		d.deploymentService.ExecuteDeploymentAsync(createdDep)
		result.HandledCount++
		handledSet[svc.ProjectID] = true
		result.HandledProjIDs = append(result.HandledProjIDs, svc.ProjectID)

		_, _ = d.deliveryService.Record(ctx, models.WebhookDeliveryInput{
			ProjectID:  &svc.ProjectID,
			Source:     "github",
			DeliveryID: &deliveryID,
			Event:      "push",
			Outcome:    "dispatched",
			ActionRef:  &createdDep.ID,
			ClientIP:   clientIP,
			UserAgent:  userAgent,
		})
	}

	if claim.ID != "" {
		if result.FailedCount > 0 {
			errSummary := strings.Join(result.Errors, "; ")
			statusCode := 500
			_ = d.deliveryService.MarkProcessed(ctx, claim.ID, models.WebhookDeliveryResult{
				Outcome:    "failed",
				StatusCode: &statusCode,
				Error:      &errSummary,
				Summary: map[string]any{
					"handledProjectIds": result.HandledProjIDs,
					"repo":              payload.Repository.FullName,
					"ref":               payload.Ref,
				},
			})
		} else {
			statusCode := 200
			_ = d.deliveryService.MarkProcessed(ctx, claim.ID, models.WebhookDeliveryResult{
				Outcome:    "received",
				StatusCode: &statusCode,
				Summary: map[string]any{
					"handledProjectIds": result.HandledProjIDs,
					"repo":              payload.Repository.FullName,
					"ref":               payload.Ref,
				},
			})
		}
	}

	return result, nil
}
