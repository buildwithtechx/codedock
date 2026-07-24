package services

import (
	"bytes"
	"context"
	"encoding/json"
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

func (s *GitService) ListRepositories(ctx context.Context, userID, provider string) ([]models.GitRepository, error) {
	if provider == "" {
		providers, err := s.repo.ListProvidersByUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to list providers: %w", err)
		}
		var allRepos []models.GitRepository
		for _, p := range providers {
			repos, err := s.ListRepositories(ctx, userID, p.Provider)
			if err == nil {
				allRepos = append(allRepos, repos...)
			}
		}
		return allRepos, nil
	}

	gp, err := s.repo.GetProvider(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to load git credentials: %w", err)
	}
	if gp == nil || gp.AccessToken == "" {
		return nil, fmt.Errorf("user is not authenticated with %s", provider)
	}
	switch provider {
	case "github":
		return s.listGitHubRepos(ctx, gp.AccessToken)
	case "gitlab":
		return s.listGitLabRepos(ctx, gp.AccessToken)
	case "bitbucket":
		return s.listBitbucketRepos(ctx, gp.AccessToken)
	case "gitea":
		return s.listGiteaRepos(ctx, gp.AccessToken)
	default:
		return nil, errors.New("unsupported provider: " + provider)
	}
}

func (s *GitService) fetchGitAPI(ctx context.Context, reqURL, token string, headers map[string]string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("api request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api returned status %d: %s", resp.StatusCode, string(body))
	}
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func (s *GitService) listGitHubRepos(ctx context.Context, token string) ([]models.GitRepository, error) {
	reqURL := "https://api.github.com/user/repos?per_page=100&sort=updated"
	var ghRepos []struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		CloneURL string `json:"clone_url"`
		HTMLURL  string `json:"html_url"`
		Default  string `json:"default_branch"`
	}

	err := s.fetchGitAPI(ctx, reqURL, token, map[string]string{
		"Accept": "application/vnd.github+json",
	}, &ghRepos)

	if err != nil {
		return nil, fmt.Errorf("github: %w", err)
	}

	var results []models.GitRepository
	for _, r := range ghRepos {
		results = append(results, models.GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Private,
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.Default,
		})
	}
	return results, nil
}

func (s *GitService) listGitLabRepos(ctx context.Context, token string) ([]models.GitRepository, error) {
	reqURL := "https://gitlab.com/api/v4/projects?membership=true&per_page=100&order_by=updated_at"
	var glRepos []struct {
		ID            int64  `json:"id"`
		Name          string `json:"name"`
		PathWithNames string `json:"path_with_namespace"`
		Visibility    string `json:"visibility"`
		HTTPURL       string `json:"http_url_to_repo"`
		WebURL        string `json:"web_url"`
		DefaultBranch string `json:"default_branch"`
	}

	err := s.fetchGitAPI(ctx, reqURL, token, nil, &glRepos)
	if err != nil {
		return nil, fmt.Errorf("gitlab: %w", err)
	}

	var results []models.GitRepository
	for _, r := range glRepos {
		results = append(results, models.GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.PathWithNames,
			Private:       r.Visibility != "public",
			CloneURL:      r.HTTPURL,
			HTMLURL:       r.WebURL,
			DefaultBranch: r.DefaultBranch,
		})
	}
	return results, nil
}

func (s *GitService) listBitbucketRepos(ctx context.Context, token string) ([]models.GitRepository, error) {
	reqURL := "https://api.bitbucket.org/2.0/repositories?role=member"
	var bbResponse struct {
		Values []struct {
			UUID     string `json:"uuid"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			IsPriv   bool   `json:"is_private"`
			Links    struct {
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Clone []struct {
					Name string `json:"name"`
					Href string `json:"href"`
				} `json:"clone"`
			} `json:"links"`
			Mainbranch struct {
				Name string `json:"name"`
			} `json:"mainbranch"`
		} `json:"values"`
	}

	err := s.fetchGitAPI(ctx, reqURL, token, nil, &bbResponse)
	if err != nil {
		return nil, fmt.Errorf("bitbucket: %w", err)
	}

	var results []models.GitRepository
	for i, r := range bbResponse.Values {
		cloneURL := ""
		for _, link := range r.Links.Clone {
			if link.Name == "https" {
				cloneURL = link.Href
			}
		}
		results = append(results, models.GitRepository{
			ID:            int64(i + 1), // Bitbucket uses UUIDs, fake it for ID
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.IsPriv,
			CloneURL:      cloneURL,
			HTMLURL:       r.Links.HTML.Href,
			DefaultBranch: r.Mainbranch.Name,
		})
	}
	return results, nil
}

func (s *GitService) listGiteaRepos(ctx context.Context, token string) ([]models.GitRepository, error) {
	reqURL := "https://gitea.com/api/v1/user/repos" // assuming public gitea.com for simplicity
	var gtRepos []struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		CloneURL string `json:"clone_url"`
		HTMLURL  string `json:"html_url"`
		Default  string `json:"default_branch"`
	}

	err := s.fetchGitAPI(ctx, reqURL, token, nil, &gtRepos)
	if err != nil {
		return nil, fmt.Errorf("gitea: %w", err)
	}

	var results []models.GitRepository
	for _, r := range gtRepos {
		results = append(results, models.GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Private,
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.Default,
		})
	}
	return results, nil
}

func (s *GitService) CloneOrPullAppRepository(ctx context.Context, app *models.AppService, targetDir string, logWriter io.Writer) error {
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
	authURL := s.injectAuthTokenIfAvailable(ctx, repoURL)
	if logWriter != nil {
		fmt.Fprintf(logWriter, "📥 [GitService] Preparing to sync codebase from %s (branch: %s)...\n", repoURL, branch)
	}
	gitDir := filepath.Join(targetDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🔄 [GitService] Existing local directory detected; pulling latest changes...\n")
		}
		fetchCmd := exec.CommandContext(ctx, "git", "-C", targetDir, "fetch", "origin", branch)
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
	cloneArgs := []string{"clone", "--depth", "1", "-b", branch, authURL, targetDir}
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
			cloneCmd = exec.CommandContext(ctx, "git", "clone", "--depth", "1", authURL, targetDir)
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

func (s *GitService) injectAuthTokenIfAvailable(ctx context.Context, repoURL string) string {
	u, err := url.Parse(repoURL)
	if err != nil || u.Scheme != "https" {
		return repoURL
	}
	var provider string
	if strings.Contains(u.Host, "github.com") {
		provider = "github"
	} else {
		return repoURL
	}
	gp, err := s.repo.GetAnyProviderByType(ctx, provider)
	if err != nil || gp == nil || gp.AccessToken == "" {
		return repoURL
	}
	switch provider {
	case "github":
		u.User = url.UserPassword("x-access-token", gp.AccessToken)
	}
	return u.String()
}
