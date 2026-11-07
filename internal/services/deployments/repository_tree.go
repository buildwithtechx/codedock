package deployments

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"codedock/internal/models"
)

func (s *GitService) GetRepoTree(ctx context.Context, userID, repoURL, branch string) ([]models.RepoTreeEntry, error) {
	parsed, err := url.Parse(repoURL)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil {
		return nil, fmt.Errorf("use an HTTPS repository URL without credentials")
	}

	hostname := strings.ToLower(parsed.Hostname())
	if hostname == "github.com" && userID != "" {
		tree, err := s.fetchGitHubTree(ctx, userID, parsed.Path, branch)
		if err == nil && len(tree) > 0 {
			return tree, nil
		}
	}

	return s.fetchTreeViaGit(ctx, userID, parsed, branch)
}

func (s *GitService) fetchGitHubTree(ctx context.Context, userID, repoPath, branch string) ([]models.RepoTreeEntry, error) {
	credentials, err := s.repo.GetProvider(ctx, userID, "github")
	if err != nil || credentials == nil || credentials.AccessToken == "" {
		return nil, fmt.Errorf("no github token")
	}

	trimmed := strings.TrimPrefix(repoPath, "/")
	trimmed = strings.TrimSuffix(trimmed, ".git")
	targetBranch := branch
	if targetBranch == "" {
		targetBranch = "HEAD"
	}

	baseURL := os.Getenv("GITHUB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	apiURL := fmt.Sprintf("%s/repos/%s/git/trees/%s?recursive=1", strings.TrimRight(baseURL, "/"), trimmed, url.PathEscape(targetBranch))

	var response struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"` // "blob" or "tree"
		} `json:"tree"`
		Truncated bool `json:"truncated"`
	}

	err = s.fetchGitAPI(ctx, apiURL, credentials.AccessToken, map[string]string{
		"Accept": "application/vnd.github+json",
	}, &response)
	if err != nil {
		return nil, err
	}

	results := make([]models.RepoTreeEntry, 0, len(response.Tree))
	for _, item := range response.Tree {
		if strings.HasPrefix(item.Path, ".git/") || item.Path == ".git" {
			continue
		}
		results = append(results, models.RepoTreeEntry{
			Path: item.Path,
			Type: item.Type,
		})
	}
	return results, nil
}

func (s *GitService) fetchTreeViaGit(ctx context.Context, userID string, parsedURL *url.URL, branch string) ([]models.RepoTreeEntry, error) {
	directory, err := os.MkdirTemp("", "codedock-tree-")
	if err != nil {
		return nil, fmt.Errorf("create temp tree checkout: %w", err)
	}
	defer os.RemoveAll(directory)

	ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	args := []string{"-c", "protocol.file.allow=never", "-c", "http.followRedirects=false"}
	provider := map[string]string{"github.com": "github", "gitlab.com": "gitlab", "bitbucket.org": "bitbucket"}[strings.ToLower(parsedURL.Hostname())]
	if provider != "" && userID != "" {
		credentials, err := s.repo.GetProvider(ctx, userID, provider)
		if err == nil && credentials != nil && credentials.AccessToken != "" {
			username := map[string]string{"github": "x-access-token", "gitlab": "oauth2", "bitbucket": "x-token-auth"}[provider]
			header := "Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+credentials.AccessToken))
			args = append(args, "-c", "http.extraHeader="+header)
		}
	}

	args = append(args, "clone", "--depth", "1", "--single-branch")
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	checkout := filepath.Join(directory, "source")
	args = append(args, "--", parsedURL.String(), checkout)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_NOSYSTEM=1")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git checkout failed: %w", err)
	}

	var entries []models.RepoTreeEntry
	err = filepath.Walk(checkout, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		rel, err := filepath.Rel(checkout, path)
		if err != nil || rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		entryType := "blob"
		if info.IsDir() {
			entryType = "tree"
		}
		entries = append(entries, models.RepoTreeEntry{
			Path: rel,
			Type: entryType,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func DetectMonorepoCandidateApps(checkout string) []models.MonorepoCandidateApp {
	var candidates []models.MonorepoCandidateApp
	prefixes := []string{"apps", "packages", "services", "projects"}

	isAppDir := func(dir string) bool {
		markers := []string{"package.json", "Dockerfile", "go.mod", "requirements.txt", "Cargo.toml", "pyproject.toml"}
		for _, m := range markers {
			if info, err := os.Stat(filepath.Join(dir, m)); err == nil && !info.IsDir() {
				return true
			}
		}
		return false
	}

	inspectCandidate := func(relPath string, fullPath string) {
		insp, err := inspectRepositoryFiles(fullPath)
		if err != nil {
			return
		}
		name := filepath.Base(relPath)
		pkgData, err := os.ReadFile(filepath.Join(fullPath, "package.json"))
		if err == nil {
			var pkg struct {
				Name string `json:"name"`
			}
			if json.Unmarshal(pkgData, &pkg) == nil && pkg.Name != "" {
				name = strings.TrimPrefix(pkg.Name, "@")
				name = strings.ReplaceAll(name, "/", "-")
			}
		}

		dockerfile := ""
		if insp.DockerfilePath != "" {
			dockerfile = insp.DockerfilePath
		}

		candidates = append(candidates, models.MonorepoCandidateApp{
			Name:           name,
			RootDirectory:  filepath.ToSlash(relPath),
			Framework:      insp.Framework,
			InternalPort:   insp.InternalPort,
			BuildEngine:    insp.BuildEngine,
			PackageManager: insp.PackageManager,
			InstallCommand: insp.InstallCommand,
			BuildCommand:   insp.BuildCommand,
			StartCommand:   insp.StartCommand,
			DockerfilePath: dockerfile,
		})
	}

	for _, prefix := range prefixes {
		containerDir := filepath.Join(checkout, prefix)
		entries, err := os.ReadDir(containerDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			subPath := filepath.Join(prefix, entry.Name())
			fullSubPath := filepath.Join(checkout, subPath)
			if isAppDir(fullSubPath) {
				inspectCandidate(subPath, fullSubPath)
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].RootDirectory < candidates[j].RootDirectory
	})
	return candidates
}
