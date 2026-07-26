package projects

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"codedock.run/codedock/internal/models"
)

type gitProviderFetcher interface {
	GetAnyProviderByType(ctx context.Context, providerType string) (*models.GitProviderConfig, error)
}

type EnvSuggestionService struct {
	gitService gitProviderFetcher
}

func NewEnvSuggestionService(gitService gitProviderFetcher) *EnvSuggestionService {
	return &EnvSuggestionService{
		gitService: gitService,
	}
}

type EnvExampleVariableSuggestion struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Label      string `json:"label"`
	SourcePath string `json:"sourcePath"`
}

func (s *EnvSuggestionService) SuggestEnvVars(ctx context.Context, repoURL, branch, rootDir string) ([]EnvExampleVariableSuggestion, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, fmt.Errorf("invalid repository url: %w", err)
	}

	if !strings.Contains(u.Host, "github.com") {
		return nil, errors.New("env suggestions only supported for github repositories")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, errors.New("invalid github repository format")
	}
	owner := parts[0]
	repo := strings.TrimSuffix(parts[1], ".git")

	if branch == "" {
		branch = "main"
	}

	pathsToTry := []string{".env.example"}
	rootDir = strings.Trim(rootDir, "/")
	if rootDir != "" {
		pathsToTry = []string{fmt.Sprintf("%s/.env.example", rootDir), ".env.example"}
	}

	var content string
	var foundPath string

	for _, path := range pathsToTry {
		c, err := s.fetchGitHubFile(ctx, owner, repo, branch, path)
		if err == nil && c != "" {
			content = c
			foundPath = path
			break
		}
	}

	if content == "" {
		return []EnvExampleVariableSuggestion{}, nil
	}

	return s.parseEnvExampleText(content, foundPath), nil
}

func (s *EnvSuggestionService) fetchGitHubFile(ctx context.Context, owner, repo, branch, path string) (string, error) {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, path, branch)

	gp, _ := s.gitService.GetAnyProviderByType(ctx, "github")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if gp != nil && gp.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+gp.AccessToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var result struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(result.Content)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}

	return result.Content, nil
}

func (s *EnvSuggestionService) parseEnvExampleText(input, sourcePath string) []EnvExampleVariableSuggestion {
	suggestions := []EnvExampleVariableSuggestion{}
	seen := map[string]bool{}

	lines := strings.Split(input, "\n")
	keyRegex := regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		if !keyRegex.MatchString(strings.ToUpper(key)) {
			continue
		}

		if seen[strings.ToUpper(key)] {
			continue
		}

		value := strings.TrimSpace(line[idx+1:])
		if len(value) >= 2 && ((strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) || (strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`))) {
			value = value[1 : len(value)-1]
		}

		suggestions = append(suggestions, EnvExampleVariableSuggestion{
			Key:        key,
			Value:      value,
			Label:      fmt.Sprintf("%s variable", sourcePath),
			SourcePath: sourcePath,
		})
		seen[strings.ToUpper(key)] = true
	}

	return suggestions
}
