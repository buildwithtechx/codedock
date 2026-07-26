package deployments

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

type GitService struct {
	repo       repositories.GitRepository
	httpClient *http.Client
}

func NewGitService(r repositories.GitRepository) *GitService {
	return &GitService{
		repo:       r,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *GitService) SaveProvider(ctx context.Context, gp *models.GitProviderConfig) error {
	if gp == nil || gp.UserID == "" || gp.Provider == "" {
		return errors.New("valid git provider config with userId and provider required")
	}
	if gp.ID == "" {
		gp.ID = uuid.New().String()
	}
	gp.UpdatedAt = time.Now()
	if gp.CreatedAt.IsZero() {
		gp.CreatedAt = gp.UpdatedAt
	}
	return s.repo.SaveProvider(ctx, gp)
}

func (s *GitService) ConnectProvider(ctx context.Context, userID string, req *models.GitConnectRequest) (*models.GitProviderConfig, error) {
	switch req.Provider {
	case "github", "gitlab", "bitbucket", "gitea":
	default:
		return nil, errors.New("unsupported git provider")
	}
	if req.AccessToken == "" {
		return nil, errors.New("access token is required")
	}
	gp := &models.GitProviderConfig{
		UserID:      userID,
		Provider:    req.Provider,
		AccessToken: req.AccessToken,
		AccountName: req.AccountName,
	}
	if err := s.SaveProvider(ctx, gp); err != nil {
		return nil, fmt.Errorf("failed to save git provider: %w", err)
	}
	gp.AccessToken = ""
	return gp, nil
}

func (s *GitService) GetProvider(ctx context.Context, userID, provider string) (*models.GitProviderConfig, error) {
	if userID == "" || provider == "" {
		return nil, errors.New("userId and provider required")
	}
	return s.repo.GetProvider(ctx, userID, provider)
}

func (s *GitService) GetConnectedProviders(ctx context.Context, userID string) ([]map[string]any, error) {
	providers, err := s.repo.ListProvidersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	providerMap := make(map[string]*models.GitProviderConfig)
	for _, gp := range providers {
		providerMap[gp.Provider] = gp
	}
	var results []map[string]any
	for _, provider := range []string{"github", "gitlab", "bitbucket", "gitea"} {
		if gp, ok := providerMap[provider]; ok && gp != nil {
			results = append(results, map[string]any{
				"provider":    provider,
				"connected":   true,
				"accountName": gp.AccountName,
				"updatedAt":   gp.UpdatedAt,
			})
		} else {
			results = append(results, map[string]any{
				"provider":  provider,
				"connected": false,
			})
		}
	}
	return results, nil
}

func (s *GitService) GetAnyProviderByType(ctx context.Context, provider string) (*models.GitProviderConfig, error) {
	if provider == "" {
		return nil, errors.New("provider required")
	}
	return s.repo.GetAnyProviderByType(ctx, provider)
}

func (s *GitService) ListProvidersByUser(ctx context.Context, userID string) ([]*models.GitProviderConfig, error) {
	if userID == "" {
		return nil, errors.New("userId required")
	}
	return s.repo.ListProvidersByUser(ctx, userID)
}

func (s *GitService) DisconnectProvider(ctx context.Context, userID, provider string) error {
	if userID == "" || provider == "" {
		return errors.New("userId and provider required")
	}
	return s.repo.DeleteProvider(ctx, userID, provider)
}

func (s *GitService) DeleteProvider(ctx context.Context, userID, provider string) error {
	return s.DisconnectProvider(ctx, userID, provider)
}

func (s *GitService) CloneOrPullAppRepository(ctx context.Context, app *models.AppService, targetDir string, logWriter io.Writer) error {
	return s.SyncCodebase(ctx, app, targetDir, logWriter)
}

func (s *GitService) SyncCodebase(ctx context.Context, app *models.AppService, targetDir string, logWriter io.Writer) error {
	if app == nil {
		return errors.New("app service is nil")
	}
	repoURL := strings.TrimSpace(app.RepositoryURL)
	if repoURL == "" {
		return errors.New("repositoryUrl is not set for service")
	}
	branch := strings.TrimSpace(app.Branch)
	if branch == "" {
		branch = "main"
	}
	validBranch := regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
	if !validBranch.MatchString(branch) {
		return errors.New("invalid branch name")
	}

	cleanURL, token := s.getAuthTokenIfAvailable(ctx, repoURL)
	var gitConfigArgs []string
	if token != "" {
		authHeader := fmt.Sprintf("Authorization: Basic %s", base64.StdEncoding.EncodeToString([]byte("x-access-token:"+token)))
		gitConfigArgs = []string{"-c", fmt.Sprintf("http.extraheader=%s", authHeader)}
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📥 [GitService] Preparing to sync codebase from %s (branch: %s)...\n", cleanURL, branch)
	}
	gitDir := filepath.Join(targetDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🔄 [GitService] Existing local directory detected; pulling latest changes...\n")
		}
		fetchArgs := append(append([]string{}, gitConfigArgs...), "-C", targetDir, "fetch", "origin", branch)
		fetchCmd := exec.CommandContext(ctx, "git", fetchArgs...)
		if out, err := fetchCmd.CombinedOutput(); err != nil {
			return utils.NewDeploymentError(fmt.Sprintf("git fetch failed: %s", string(out)), err)
		}
		resetCmd := exec.CommandContext(ctx, "git", "-C", targetDir, "reset", "--hard", "origin/"+branch)
		if out, err := resetCmd.CombinedOutput(); err != nil {
			return utils.NewDeploymentError(fmt.Sprintf("git reset failed: %s", string(out)), err)
		}
		if logWriter != nil {
			fmt.Fprintf(logWriter, "✅ [GitService] Successfully updated local repository to latest commit on %s.\n", branch)
		}
		return nil
	}
	_ = os.RemoveAll(targetDir)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		return fmt.Errorf("failed to create build parent dir: %w", err)
	}

	cloneArgs := append(append([]string{}, gitConfigArgs...), "clone", "--depth", "1", "-b", branch, cleanURL, targetDir)
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [GitService] Running git clone --depth 1 -b %s...\n", branch)
	}
	cloneCmd := exec.CommandContext(ctx, "git", cloneArgs...)
	var stderr bytes.Buffer
	cloneCmd.Stderr = &stderr
	if err := cloneCmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "Remote branch") && branch == "main" {
			if logWriter != nil {
				fmt.Fprintf(logWriter, "⚠️ [GitService] Branch 'main' not found; retrying clone with repository default branch...\n")
			}
			_ = os.RemoveAll(targetDir)
			fallbackArgs := append(append([]string{}, gitConfigArgs...), "clone", "--depth", "1", cleanURL, targetDir)
			cloneCmd = exec.CommandContext(ctx, "git", fallbackArgs...)
			if errFallback := cloneCmd.Run(); errFallback != nil {
				return utils.NewDeploymentError(fmt.Sprintf("git clone failed: %s", stderr.String()), errFallback)
			}
		} else {
			return utils.NewDeploymentError(fmt.Sprintf("git clone failed: %s", stderr.String()), err)
		}
	}
	if logWriter != nil {
		fmt.Fprintf(logWriter, "✅ [GitService] Successfully cloned repository into %s.\n", targetDir)
	}
	return nil
}

func (s *GitService) getAuthTokenIfAvailable(ctx context.Context, repoURL string) (string, string) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return repoURL, ""
	}
	u.User = nil
	cleanURL := u.String()
	if u.Scheme != "https" {
		return cleanURL, ""
	}
	var provider string
	if strings.Contains(u.Host, "github.com") {
		provider = "github"
	} else {
		return cleanURL, ""
	}
	gp, err := s.repo.GetAnyProviderByType(ctx, provider)
	if err != nil || gp == nil || gp.AccessToken == "" {
		return cleanURL, ""
	}
	return cleanURL, gp.AccessToken
}
