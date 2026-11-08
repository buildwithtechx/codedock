package git

import (
	"path/filepath"
	"strings"
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

func ExtractChangedFilesFromPush(payload GithubWebhookPayload) []string {
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
	rootDirectory string,
	buildCommand string,
	runtimeMode string,
	files []string,
	commitMessage string,
) ProjectChangeImpact {
	if strings.Contains(strings.ToLower(commitMessage), "[force-deploy]") {
		return ProjectChangeImpact{Affected: true, Reason: "force-deploy-token"}
	}

	if runtimeMode == "docker" || strings.ToLower(buildCommand) == "docker" {
		return ProjectChangeImpact{Affected: true, Reason: "custom-build-context"}
	}

	root := normalizePath(rootDirectory)
	if root == "" || root == "." {
		return ProjectChangeImpact{Affected: true, Reason: "repository-root"}
	}

	if len(files) == 0 {
		return ProjectChangeImpact{Affected: true, Reason: "unknown-changes"}
	}

	for _, file := range files {
		clean := normalizePath(file)
		if globalWorkspaceFiles[clean] {
			return ProjectChangeImpact{Affected: true, Reason: "global-manifest"}
		}
		if clean == "codedock.json" || clean == "openship.json" {
			return ProjectChangeImpact{Affected: true, Reason: "repository-config"}
		}
	}

	rootPrefix := root + "/"
	for _, file := range files {
		clean := normalizePath(file)
		if clean == root || strings.HasPrefix(clean, rootPrefix) {
			return ProjectChangeImpact{Affected: true, Reason: "project-files"}
		}
	}

	return ProjectChangeImpact{
		Affected: false,
		Reason:   "unrelated-workspaces",
	}
}
