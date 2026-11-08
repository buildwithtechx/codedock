package git

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/utils"
)

var hex40Regex = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)

func ParseRepoOwnerAndName(rawURL string) (string, string) {
	trimmed := strings.TrimSuffix(strings.TrimSpace(rawURL), ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2], parts[len(parts)-1]
	}
	return "", ""
}

func IsPreviewEligible(pr PullRequestInfo, targetOwner, targetRepo, targetBranch string) (bool, string) {
	if strings.ToLower(pr.State) != "open" {
		return false, "PR is not open"
	}
	if pr.Draft {
		return false, "PR is a draft"
	}
	if pr.BaseRef != targetBranch {
		return false, fmt.Sprintf("PR base branch %q does not match project branch %q", pr.BaseRef, targetBranch)
	}

	expectedRepo := strings.ToLower(strings.TrimSpace(targetOwner) + "/" + strings.TrimSuffix(strings.TrimSpace(targetRepo), ".git"))
	if strings.ToLower(pr.BaseRepo) != expectedRepo {
		return false, fmt.Sprintf("PR base repo %q does not match project repo %q", pr.BaseRepo, expectedRepo)
	}
	if pr.HeadRepo != "" && strings.ToLower(pr.HeadRepo) != expectedRepo {
		return false, fmt.Sprintf("PR head repo %q does not match project repo %q", pr.HeadRepo, expectedRepo)
	}
	if !hex40Regex.MatchString(pr.HeadSHA) {
		return false, "PR head commit SHA must be a valid 40-character hex string"
	}

	return true, ""
}

func ValidatePreviewCommands(commands ...string) error {
	for _, cmd := range commands {
		if cmd == "" {
			continue
		}
		res := utils.InspectCommandEnvironment(cmd)
		if len(res.Keys) > 0 {
			return fmt.Errorf("command %q embeds environment variable assignments (%v); assign them in the Environment tab instead", cmd, res.Keys)
		}
	}
	return nil
}

func GeneratePreviewDomain(baseDomain, magicDomain, appName, projectID string, prNumber int) string {
	hash := sha256.Sum256(fmt.Appendf(nil, "%s:%d", projectID, prNumber))
	suffix := hex.EncodeToString(hash[:])[:6]
	cleanApp := strings.ToLower(strings.ReplaceAll(appName, "_", "-"))

	if baseDomain != "" {
		return fmt.Sprintf("%s-pr-%d-%s.%s", cleanApp, prNumber, suffix, strings.TrimPrefix(baseDomain, "."))
	}
	if magicDomain == "" {
		magicDomain = "sslip.io"
	}
	return fmt.Sprintf("%s-pr-%d-%s.%s", cleanApp, prNumber, suffix, strings.TrimPrefix(magicDomain, "."))
}

type PRPreviewService struct {
	repo          PRPreviewRepository
	serviceLookup ServiceLookup
}

func NewPRPreviewService(repo PRPreviewRepository, serviceLookup ServiceLookup) *PRPreviewService {
	return &PRPreviewService{
		repo:          repo,
		serviceLookup: serviceLookup,
	}
}

func (s *PRPreviewService) ListByApp(ctx context.Context, appID string) ([]*PRPreview, error) {
	return s.repo.GetByApp(ctx, appID)
}

func (s *PRPreviewService) CreatePreviewRecord(ctx context.Context, appID, projectID, branch, commitHash, previewDomain string, prNumber int) (*PRPreview, error) {
	now := time.Now().UTC()
	preview := &PRPreview{
		ID:            uuid.NewString(),
		ServiceID:     appID,
		ProjectID:     projectID,
		PRNumber:      prNumber,
		Branch:        branch,
		CommitHash:    commitHash,
		Status:        PRPreviewStatusBuilding,
		PreviewDomain: previewDomain,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.repo.Create(ctx, preview); err != nil {
		return nil, err
	}
	return preview, nil
}

func (s *PRPreviewService) UpdateStatus(ctx context.Context, previewID, status, containerID string) error {
	previews, err := s.repo.GetByApp(ctx, previewID)
	if err != nil || len(previews) == 0 {
		return nil
	}
	p := previews[0]
	p.Status = status
	p.ContainerID = containerID
	p.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, p)
}

func (s *PRPreviewService) DestroyPRPreview(ctx context.Context, appID string, prNumber int) error {
	previews, err := s.repo.GetByAppAndPR(ctx, appID, prNumber)
	if err != nil {
		return err
	}
	for _, p := range previews {
		_ = s.repo.Delete(ctx, p.ID)
	}
	return nil
}

func (s *PRPreviewService) UpdateCommitStatus(ctx context.Context, repoURL, commitHash, previewDomain string) {
	if repoURL == "" || commitHash == "" {
		return
	}
	repoParts := strings.Split(strings.TrimSuffix(repoURL, ".git"), "/")
	if len(repoParts) < 2 {
		return
	}
	owner := repoParts[len(repoParts)-2]
	repo := repoParts[len(repoParts)-1]
	token := ""
	if strings.Contains(repoURL, "@") {
		parts := strings.Split(repoURL, "@")
		authParts := strings.Split(parts[0], "://")
		if len(authParts) == 2 {
			creds := strings.Split(authParts[1], ":")
			if len(creds) == 2 {
				token = creds[1]
			} else {
				token = creds[0]
			}
		}
	}
	if token == "" || !strings.Contains(repoURL, "github.com") {
		return
	}
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/statuses/%s", owner, repo, commitHash)
	payload := map[string]string{
		"state":       "success",
		"target_url":  "https://" + previewDomain,
		"description": "PR Preview is ready",
		"context":     "codedock/pr-preview",
	}
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}
