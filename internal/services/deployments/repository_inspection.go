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
	"strings"
	"time"

	"codedock.run/codedock/internal/models"
)

func (s *GitService) InspectRepository(ctx context.Context, userID string, request models.RepositoryInspectionRequest) (*models.RepositoryInspection, error) {
	parsed, err := url.Parse(request.RepositoryURL)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Port() != "" {
		return nil, fmt.Errorf("use an HTTPS repository URL without credentials")
	}
	providers := map[string]string{"github.com": "github", "gitlab.com": "gitlab", "bitbucket.org": "bitbucket"}
	provider, allowed := providers[strings.ToLower(parsed.Hostname())]
	if !allowed {
		return nil, fmt.Errorf("inspection supports GitHub, GitLab and Bitbucket HTTPS repositories; configure other repositories manually")
	}
	if userID == "" {
		return nil, fmt.Errorf("repository inspection requires a user")
	}
	branch := request.Branch
	if strings.HasPrefix(branch, "-") || strings.ContainsAny(branch, " \t\r\n") {
		return nil, fmt.Errorf("invalid branch")
	}
	directory, err := os.MkdirTemp("", "codedock-inspect-")
	if err != nil {
		return nil, fmt.Errorf("create inspection checkout: %w", err)
	}
	defer os.RemoveAll(directory)
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	args := []string{"-c", "protocol.file.allow=never", "-c", "http.followRedirects=false"}
	credentials, err := s.repo.GetProvider(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("load repository connection: %w", err)
	}
	if credentials != nil && credentials.AccessToken != "" {
		username := map[string]string{"github": "x-access-token", "gitlab": "oauth2", "bitbucket": "x-token-auth"}[provider]
		header := "Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+credentials.AccessToken))
		args = append(args, "-c", "http.extraHeader="+header)
	}
	args = append(args, "clone", "--depth", "1", "--single-branch")
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	checkout := filepath.Join(directory, "source")
	args = append(args, "--", parsed.String(), checkout)
	command := exec.CommandContext(ctx, "git", args...)
	command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_NOSYSTEM=1")
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("repository inspection failed; verify the URL, branch and connected account")
	}
	root, err := inspectionRoot(checkout, request.RootDirectory)
	if err != nil {
		return nil, err
	}
	return inspectRepositoryFiles(root)
}

func inspectionRoot(checkout, root string) (string, error) {
	root = strings.TrimPrefix(strings.ReplaceAll(root, "\\", "/"), "/")
	resolved, err := filepath.EvalSymlinks(filepath.Join(checkout, filepath.FromSlash(root)))
	if err != nil {
		return "", fmt.Errorf("root directory does not exist: %w", err)
	}
	relative, err := filepath.Rel(checkout, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("root directory must stay inside the repository")
	}
	return resolved, nil
}

func inspectRepositoryFiles(root string) (*models.RepositoryInspection, error) {
	result := &models.RepositoryInspection{Framework: "unknown", InternalPort: 3000, BuildEngine: "nixpacks", DockerfilePath: "Dockerfile", Warnings: []string{}}
	exists := func(name string) bool {
		info, err := os.Lstat(filepath.Join(root, name))
		return err == nil && info.Mode().IsRegular()
	}
	if exists("Dockerfile") {
		result.BuildEngine, result.Framework = "dockerfile", "Docker"
	}
	if exists("package.json") {
		info, err := os.Lstat(filepath.Join(root, "package.json"))
		if err != nil || info.Size() > 1024*1024 {
			return nil, fmt.Errorf("package.json exceeds inspection limits")
		}
		data, err := os.ReadFile(filepath.Join(root, "package.json"))
		if err != nil {
			return nil, fmt.Errorf("read package manifest: %w", err)
		}
		var manifest struct {
			Scripts         map[string]string `json:"scripts"`
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("invalid package.json: %w", err)
		}
		manager := "npm"
		if exists("pnpm-lock.yaml") {
			manager = "pnpm"
		} else if exists("yarn.lock") {
			manager = "yarn"
		} else if exists("bun.lock") || exists("bun.lockb") {
			manager = "bun"
		}
		result.PackageManager, result.InstallCommand = manager, manager+" install"
		if manager == "npm" && exists("package-lock.json") {
			result.InstallCommand = "npm ci"
		}
		if manifest.Scripts["build"] != "" {
			result.BuildCommand = manager + " run build"
		}
		if manifest.Scripts["start"] != "" {
			result.StartCommand = manager + " run start"
		}
		frameworks := []struct {
			dependency, name, output string
			port                     int
		}{{"next", "Next.js", "", 3000}, {"nuxt", "Nuxt", "", 3000}, {"@sveltejs/kit", "SvelteKit", "", 3000}, {"astro", "Astro", "dist", 4321}, {"vite", "Vite", "dist", 3000}, {"express", "Express", "", 3000}}
		if result.Framework != "Docker" {
			result.Framework = "Node.js"
		}
		for _, framework := range frameworks {
			if manifest.Dependencies[framework.dependency] != "" || manifest.DevDependencies[framework.dependency] != "" {
				result.Framework, result.StaticOutput, result.InternalPort = framework.name, framework.output, framework.port
				break
			}
		}
	} else if exists("go.mod") {
		result.Framework, result.BuildCommand, result.StartCommand, result.InternalPort = "Go", "go build -o app .", "./app", 8080
	} else if exists("requirements.txt") || exists("pyproject.toml") {
		result.Framework, result.InstallCommand, result.InternalPort = "Python", "pip install -r requirements.txt", 8000
		result.Warnings = append(result.Warnings, "Set the Python start command for your application.")
	}
	if result.Framework == "unknown" {
		result.Warnings = append(result.Warnings, "No supported framework detected. Enter commands and port manually.")
	}
	if result.StaticOutput != "" {
		result.InternalPort = 80
		result.StartCommand = ""
		result.Warnings = append(result.Warnings, "Confirm static output matches your framework configuration.")
	}
	return result, nil
}
