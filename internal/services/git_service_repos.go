package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"codedock.run/codedock/internal/models"
)

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
	baseURL := os.Getenv("GITHUB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	reqURL := strings.TrimRight(baseURL, "/") + "/user/repos?per_page=100&sort=updated"
	var ghRepos []models.GitHubRepoDTO

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
			DefaultBranch: r.DefaultBranch,
		})
	}
	return results, nil
}

func (s *GitService) listGitLabRepos(ctx context.Context, token string) ([]models.GitRepository, error) {
	baseURL := os.Getenv("GITLAB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	reqURL := strings.TrimRight(baseURL, "/") + "/api/v4/projects?membership=true&per_page=100&order_by=updated_at"
	var glRepos []models.GitLabRepoDTO

	err := s.fetchGitAPI(ctx, reqURL, token, nil, &glRepos)
	if err != nil {
		return nil, fmt.Errorf("gitlab: %w", err)
	}

	var results []models.GitRepository
	for _, r := range glRepos {
		results = append(results, models.GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.PathWithNamespace,
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
	var bbResponse models.BitbucketRepoListDTO

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
			ID:            int64(i + 1),
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
	baseURL := os.Getenv("GITEA_BASE_URL")
	if baseURL == "" {
		baseURL = "https://gitea.com"
	}
	reqURL := strings.TrimRight(baseURL, "/") + "/api/v1/user/repos"
	var gtRepos []models.GiteaRepoDTO

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
			DefaultBranch: r.DefaultBranch,
		})
	}
	return results, nil
}
