package git

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type DeploymentTrigger interface {
	TriggerDeployment(ctx context.Context, serviceID, environmentID, projectID, branch, trigger, buildLogs string) (string, error)
}

type ServiceLookup interface {
	GetService(ctx context.Context, serviceID string) (*ServiceInfo, error)
	ListServices(ctx context.Context) ([]*ServiceInfo, error)
}

type ServiceInfo struct {
	ID               string
	ProjectID        string
	EnvironmentID    string
	Name             string
	RepositoryURL    string
	Branch           string
	RootDirectory    string
	BuildCommand     string
	RuntimeMode      string
	EnablePRPreviews bool
	DeployToken      string
}

type WebhookDeliveryService struct {
	repo WebhookDeliveryRepository
}

func NewWebhookDeliveryService(repo WebhookDeliveryRepository) *WebhookDeliveryService {
	return &WebhookDeliveryService{repo: repo}
}

func (s *WebhookDeliveryService) ClaimGithub(
	ctx context.Context,
	deliveryID string,
	event string,
	clientIP *string,
	userAgent *string,
) (ClaimResult, error) {
	input := WebhookDeliveryInput{
		DeliveryID: &deliveryID,
		Event:      event,
		Source:     WebhookSourceGitHub,
		ClientIP:   clientIP,
		UserAgent:  userAgent,
		Outcome:    WebhookOutcomeReceived,
	}
	return s.repo.ClaimGithub(ctx, input)
}

func (s *WebhookDeliveryService) Record(ctx context.Context, input WebhookDeliveryInput) (string, error) {
	return s.repo.Record(ctx, input)
}

func (s *WebhookDeliveryService) MarkProcessed(ctx context.Context, id string, res WebhookDeliveryResult) error {
	return s.repo.MarkProcessed(ctx, id, res)
}

func (s *WebhookDeliveryService) ListByProject(
	ctx context.Context,
	projectID string,
	limit int,
	cursor string,
) ([]WebhookDelivery, string, error) {
	return s.repo.ListByProject(ctx, projectID, limit, cursor)
}

func (s *WebhookDeliveryService) Prune(ctx context.Context, olderThan time.Duration) (int64, error) {
	return s.repo.Prune(ctx, olderThan)
}

type WebhookDispatcher struct {
	deliveryService   *WebhookDeliveryService
	deploymentTrigger DeploymentTrigger
	serviceLookup     ServiceLookup
	gitAppsService    *GitAppsService
}

func NewWebhookDispatcher(
	deliveryService *WebhookDeliveryService,
	deploymentTrigger DeploymentTrigger,
	serviceLookup ServiceLookup,
	gitAppsService *GitAppsService,
) *WebhookDispatcher {
	return &WebhookDispatcher{
		deliveryService:   deliveryService,
		deploymentTrigger: deploymentTrigger,
		serviceLookup:     serviceLookup,
		gitAppsService:    gitAppsService,
	}
}

func (d *WebhookDispatcher) VerifySignature(ctx context.Context, bodyBytes []byte, signature string) (bool, error) {
	if signature == "" {
		return false, nil
	}
	if d.gitAppsService == nil {
		return true, nil
	}
	apps, err := d.gitAppsService.ListGithubApps(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to list github apps: %w", err)
	}
	for _, app := range apps {
		if app.WebhookSecret != "" && VerifyHMAC(bodyBytes, app.WebhookSecret, signature) {
			return true, nil
		}
	}
	return false, nil
}

func VerifyHMAC(payload []byte, secret, signature string) bool {
	if secret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}

func (d *WebhookDispatcher) DispatchPush(
	ctx context.Context,
	deliveryID string,
	targetServiceID string,
	payload GithubWebhookPayload,
	clientIP *string,
	userAgent *string,
) (*PushDispatchResult, error) {
	var claim ClaimResult
	if deliveryID != "" && d.deliveryService != nil {
		var err error
		claim, err = d.deliveryService.ClaimGithub(ctx, deliveryID, WebhookEventPush, clientIP, userAgent)
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

	var targetServices []*ServiceInfo
	if targetServiceID != "" && d.serviceLookup != nil {
		svc, err := d.serviceLookup.GetService(ctx, targetServiceID)
		if err == nil && svc != nil {
			targetServices = append(targetServices, svc)
		}
	} else if d.serviceLookup != nil {
		allServices, err := d.serviceLookup.ListServices(ctx)
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
		impact := ResolveProjectChangeImpact(svc.RootDirectory, svc.BuildCommand, svc.RuntimeMode, changedFiles, commitMsg)
		if !impact.Affected {
			slog.Info("skipping auto-deploy: changes do not affect service", "serviceId", svc.ID, "name", svc.Name, "root", svc.RootDirectory, "reason", impact.Reason)
			if d.deliveryService != nil {
				_, _ = d.deliveryService.Record(ctx, WebhookDeliveryInput{
					ProjectID:  &svc.ProjectID,
					Source:     WebhookSourceGitHub,
					DeliveryID: &deliveryID,
					Event:      WebhookEventPush,
					Outcome:    WebhookOutcomeSkipped,
					Summary: map[string]any{
						"reason": impact.Reason,
						"root":   svc.RootDirectory,
					},
					ClientIP:  clientIP,
					UserAgent: userAgent,
				})
			}
			continue
		}

		if d.deploymentTrigger != nil {
			logs := fmt.Sprintf("Initiating automated build from %s branch %s (delivery %s)...\n", svc.RepositoryURL, branch, deliveryID)
			depID, err := d.deploymentTrigger.TriggerDeployment(ctx, svc.ID, svc.EnvironmentID, svc.ProjectID, branch, "Git Webhook Push", logs)
			if err != nil {
				result.FailedCount++
				errMsg := fmt.Sprintf("service %s: %v", svc.Name, err)
				result.Errors = append(result.Errors, errMsg)
				if d.deliveryService != nil {
					_, _ = d.deliveryService.Record(ctx, WebhookDeliveryInput{
						ProjectID:  &svc.ProjectID,
						Source:     WebhookSourceGitHub,
						DeliveryID: &deliveryID,
						Event:      WebhookEventPush,
						Outcome:    WebhookOutcomeFailed,
						Error:      &errMsg,
						ClientIP:   clientIP,
						UserAgent:  userAgent,
					})
				}
				continue
			}

			result.HandledCount++
			result.HandledProjIDs = append(result.HandledProjIDs, svc.ProjectID)
			if d.deliveryService != nil {
				_, _ = d.deliveryService.Record(ctx, WebhookDeliveryInput{
					ProjectID:  &svc.ProjectID,
					Source:     WebhookSourceGitHub,
					DeliveryID: &deliveryID,
					Event:      WebhookEventPush,
					Outcome:    WebhookOutcomeDispatched,
					ActionRef:  &depID,
					Summary: map[string]any{
						"serviceId": svc.ID,
						"branch":    branch,
					},
					ClientIP:  clientIP,
					UserAgent: userAgent,
				})
			}
		}
	}

	return result, nil
}
