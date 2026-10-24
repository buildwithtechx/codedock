package deployments

import (
	"path/filepath"
	"strings"

	"codedock/internal/models"
)

type ProjectChangeImpact struct {
	Affected bool   `json:"affected"`
	Reason   string `json:"reason"`
}

var globalWorkspaceFiles = map[string]bool{
	"package.json":        true,
	"pnpm-lock.yaml":      true,
	"yarn.lock":           true,
	"package-lock.json":   true,
	"pnpm-workspace.yaml": true,
	"turbo.json":          true,
	"turbo.jsonc":         true,
	"lerna.json":          true,
	"nx.json":             true,
	"bun.lockb":           true,
	"bun.lock":            true,
	"deno.json":           true,
	"deno.lock":           true,
	"go.work":             true,
	"go.work.sum":         true,
	"Cargo.toml":          true,
	"Cargo.lock":          true,
}

func ExtractChangedFilesFromPush(payload models.GithubWebhookPayload) []string {
	fileMap := make(map[string]bool)
	if payload.HeadCommit != nil {
		for _, f := range payload.HeadCommit.Added {
			fileMap[normalizePath(f)] = true
		}
		for _, f := range payload.HeadCommit.Modified {
			fileMap[normalizePath(f)] = true
		}
		for _, f := range payload.HeadCommit.Removed {
			fileMap[normalizePath(f)] = true
		}
	}
	for _, c := range payload.Commits {
		for _, f := range c.Added {
			fileMap[normalizePath(f)] = true
		}
		for _, f := range c.Modified {
			fileMap[normalizePath(f)] = true
		}
		for _, f := range c.Removed {
			fileMap[normalizePath(f)] = true
		}
	}
	var files []string
	for f := range fileMap {
		if f != "" {
			files = append(files, f)
		}
	}
	return files
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	path = filepath.ToSlash(path)
	return strings.TrimPrefix(path, "/")
}

func ResolveProjectChangeImpact(
	svc *models.AppService,
	files []string,
	commitMessage string,
) ProjectChangeImpact {
	// 1. Force escape token in commit message
	if strings.Contains(strings.ToLower(commitMessage), "[force-deploy]") {
		return ProjectChangeImpact{Affected: true, Reason: "force-deploy-token"}
	}

	// 2. Custom build / Docker context always deploys
	if svc.RuntimeMode == "docker" || strings.ToLower(svc.BuildCommand) == "docker" {
		return ProjectChangeImpact{Affected: true, Reason: "custom-build-context"}
	}

	root := normalizePath(svc.RootDirectory)
	// 3. If project is at root or no rootDirectory configured, fail-open (always deploy)
	if root == "" || root == "." {
		return ProjectChangeImpact{Affected: true, Reason: "repository-root"}
	}

	// 4. If no files changed list is known (e.g. truncated / empty payload), fail-open
	if len(files) == 0 {
		return ProjectChangeImpact{Affected: true, Reason: "unknown-changes"}
	}

	// 5. Check if any global workspace manifest file changed at repo root
	for _, file := range files {
		clean := normalizePath(file)
		if globalWorkspaceFiles[clean] {
			return ProjectChangeImpact{Affected: true, Reason: "global-manifest"}
		}
		if clean == "codedock.json" || clean == "openship.json" {
			return ProjectChangeImpact{Affected: true, Reason: "repository-config"}
		}
	}

	// 6. Check if any file under the service root directory changed
	rootPrefix := root + "/"
	for _, file := range files {
		clean := normalizePath(file)
		if clean == root || strings.HasPrefix(clean, rootPrefix) {
			return ProjectChangeImpact{Affected: true, Reason: "project-files"}
		}
	}

	// 7. If none of the files touch this service root or any global manifest,
	// and changes are strictly in other paths, skip deployment!
	return ProjectChangeImpact{
		Affected: false,
		Reason:   "unrelated-workspaces",
	}
}
