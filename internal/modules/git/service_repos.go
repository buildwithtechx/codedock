package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func (s *ProviderService) ListRepositories(ctx context.Context, userID, provider string) ([]GitRepository, error) {
	if provider == "" {
		providers, err := s.repo.ListProvidersByUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to list providers: %w", err)
		}
		var allRepos []GitRepository
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
	case ProviderGitHub:
		return s.listGitHubRepos(ctx, gp.AccessToken)
	case ProviderGitLab:
		return s.listGitLabRepos(ctx, gp.AccessToken)
	case ProviderBitbucket:
		return s.listBitbucketRepos(ctx, gp.AccessToken)
	case ProviderGitea:
		return s.listGiteaRepos(ctx, gp.AccessToken)
	default:
		return nil, errors.New("unsupported provider: " + provider)
	}
}

func (s *ProviderService) fetchGitAPI(ctx context.Context, reqURL, token string, headers map[string]string, target any) error {
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

func (s *ProviderService) listGitHubRepos(ctx context.Context, token string) ([]GitRepository, error) {
	baseURL := os.Getenv("GITHUB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	reqURL := strings.TrimRight(baseURL, "/") + "/user/repos?per_page=100&sort=updated"
	var ghRepos []GitHubRepoDTO

	err := s.fetchGitAPI(ctx, reqURL, token, map[string]string{
		"Accept": "application/vnd.github+json",
	}, &ghRepos)
	if err != nil {
		return nil, fmt.Errorf("github: %w", err)
	}

	results := make([]GitRepository, len(ghRepos))
	for i, r := range ghRepos {
		results[i] = GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Private,
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.DefaultBranch,
		}
	}
	return results, nil
}

func (s *ProviderService) listGitLabRepos(ctx context.Context, token string) ([]GitRepository, error) {
	baseURL := os.Getenv("GITLAB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	reqURL := strings.TrimRight(baseURL, "/") + "/api/v4/projects?membership=true&per_page=100&order_by=updated_at"
	var glRepos []GitLabRepoDTO

	err := s.fetchGitAPI(ctx, reqURL, token, nil, &glRepos)
	if err != nil {
		return nil, fmt.Errorf("gitlab: %w", err)
	}

	results := make([]GitRepository, len(glRepos))
	for i, r := range glRepos {
		branch := r.DefaultBranch
		if branch == "" {
			branch = "main"
		}
		results[i] = GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.PathWithNamespace,
			Private:       r.Visibility == "private",
			CloneURL:      r.HTTPURL,
			HTMLURL:       r.WebURL,
			DefaultBranch: branch,
		}
	}
	return results, nil
}

func (s *ProviderService) listBitbucketRepos(ctx context.Context, token string) ([]GitRepository, error) {
	baseURL := os.Getenv("BITBUCKET_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.bitbucket.org"
	}
	reqURL := strings.TrimRight(baseURL, "/") + "/2.0/repositories?role=member&pagelen=100"
	var listResp BitbucketRepoListDTO

	err := s.fetchGitAPI(ctx, reqURL, token, nil, &listResp)
	if err != nil {
		return nil, fmt.Errorf("bitbucket: %w", err)
	}

	results := make([]GitRepository, len(listResp.Values))
	for i, r := range listResp.Values {
		cloneURL := ""
		for _, c := range r.Links.Clone {
			if c.Name == "https" {
				cloneURL = c.Href
				break
			}
		}
		branch := r.Mainbranch.Name
		if branch == "" {
			branch = "main"
		}
		results[i] = GitRepository{
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.IsPriv,
			CloneURL:      cloneURL,
			HTMLURL:       r.Links.HTML.Href,
			DefaultBranch: branch,
		}
	}
	return results, nil
}

func (s *ProviderService) listGiteaRepos(ctx context.Context, token string) ([]GitRepository, error) {
	baseURL := os.Getenv("GITEA_BASE_URL")
	if baseURL == "" {
		return nil, errors.New("gitea base url not configured")
	}
	reqURL := strings.TrimRight(baseURL, "/") + "/api/v1/user/repos?limit=100"
	var repos []GiteaRepoDTO

	err := s.fetchGitAPI(ctx, reqURL, token, nil, &repos)
	if err != nil {
		return nil, fmt.Errorf("gitea: %w", err)
	}

	results := make([]GitRepository, len(repos))
	for i, r := range repos {
		results[i] = GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Private,
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.DefaultBranch,
		}
	}
	return results, nil
}
