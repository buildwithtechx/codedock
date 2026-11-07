package deployments

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"codedock/internal/utils"
)

const MaxSubmoduleRepos = 25

var scpURLPattern = regexp.MustCompile(`^[^@\s/]+@([^:/\s]+):\/?([^/\s]+)\/([^/\s]+?)(?:\.git)?\/?$`)
var bareRepoNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,100}$`)

// ParseSameOwnerSubmoduleRepos parses .gitmodules content and returns bare repository
// names belonging to the same organization/owner on the same host.
func ParseSameOwnerSubmoduleRepos(gitmodules, owner, parentRepo, host string) []string {
	if host == "" {
		host = "github.com"
	}
	ownerLc := strings.ToLower(owner)
	parentLc := strings.ToLower(strings.TrimSuffix(parentRepo, ".git"))
	hostLc := strings.ToLower(host)

	found := make(map[string]bool)
	var repos []string

	scanner := bufio.NewScanner(strings.NewReader(gitmodules))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if !strings.HasPrefix(strings.ToLower(line), "url") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		rawVal := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		repoName := resolveSameOwnerRepo(rawVal, ownerLc, parentLc, hostLc)
		if repoName != "" && repoName != parentLc && !found[repoName] {
			found[repoName] = true
			repos = append(repos, repoName)
			if len(repos) >= MaxSubmoduleRepos {
				break
			}
		}
	}

	return repos
}

func sameHost(h1, h2 string) bool {
	h1 = strings.ToLower(h1)
	h2 = strings.ToLower(h2)
	return h1 == h2 || (h2 == "github.com" && h1 == "www.github.com")
}

func bareName(repo string) string {
	name := strings.ToLower(strings.TrimSuffix(repo, ".git"))
	if bareRepoNamePattern.MatchString(name) && name != "." && name != ".." {
		return name
	}
	return ""
}

func resolveSameOwnerRepo(rawURL, ownerLc, parentLc, hostLc string) string {
	href := strings.TrimSpace(rawURL)
	if href == "" {
		return ""
	}

	// Handle relative submodule URLs (e.g. ../sub.git or ./sub.git)
	if strings.HasPrefix(href, "./") || strings.HasPrefix(href, "../") {
		base, err := url.Parse(fmt.Sprintf("https://%s/%s/%s/", hostLc, ownerLc, parentLc))
		if err != nil {
			return ""
		}
		rel, err := url.Parse(href)
		if err != nil {
			return ""
		}
		href = base.ResolveReference(rel).String()
	}

	// Handle SCP-style URLs: git@github.com:owner/repo.git
	if matches := scpURLPattern.FindStringSubmatch(href); len(matches) == 4 {
		if !sameHost(matches[1], hostLc) {
			return ""
		}
		if strings.ToLower(matches[2]) == ownerLc {
			return bareName(matches[3])
		}
		return ""
	}

	// Handle HTTPS / SSH URLs
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	if !sameHost(u.Hostname(), hostLc) {
		return ""
	}
	cleanPath := strings.Trim(path.Clean(u.Path), "/")
	parts := strings.Split(cleanPath, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return ""
	}
	if strings.ToLower(parts[0]) != ownerLc {
		return ""
	}
	return bareName(parts[1])
}

// SyncAndInitSubmodules configures safe auth and syncs/initializes submodules if present.
func SyncAndInitSubmodules(ctx context.Context, targetDir string, gitConfigArgs []string, logWriter io.Writer) error {
	gitmodulesPath := filepath.Join(targetDir, ".gitmodules")
	if _, err := os.Stat(gitmodulesPath); os.IsNotExist(err) {
		return nil
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📦 [GitService] .gitmodules detected; synchronizing submodules...\n")
	}

	// Safe protocol restrictions: prevent command execution via malicious git protocol schemes
	baseArgs := append([]string{}, gitConfigArgs...)
	baseArgs = append(baseArgs, "-C", targetDir)

	syncCmd := exec.CommandContext(ctx, "git", append(baseArgs, "submodule", "sync", "--recursive")...)
	if out, err := syncCmd.CombinedOutput(); err != nil {
		return utils.NewDeploymentError(fmt.Sprintf("git submodule sync failed: %s", string(out)), err)
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "🔄 [GitService] Updating and cloning recursive submodules...\n")
	}

	updateCmd := exec.CommandContext(ctx, "git", append(baseArgs, "submodule", "update", "--init", "--recursive")...)
	if out, err := updateCmd.CombinedOutput(); err != nil {
		return utils.NewDeploymentError(fmt.Sprintf("git submodule update failed: %s", string(out)), err)
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "✅ [GitService] Submodules initialized and checked out successfully.\n")
	}

	return nil
}
