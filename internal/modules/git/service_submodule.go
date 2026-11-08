package git

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	scpURLPattern       = regexp.MustCompile(`^[^@\s/]+@([^:/\s]+):\/?([^/\s]+)\/([^/\s]+?)(?:\.git)?\/?$`)
	bareRepoNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,100}$`)
)

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

	if matches := scpURLPattern.FindStringSubmatch(href); len(matches) == 4 {
		if !sameHost(matches[1], hostLc) {
			return ""
		}
		if strings.ToLower(matches[2]) != ownerLc {
			return ""
		}
		return bareName(matches[3])
	}

	parsed, err := url.Parse(href)
	if err != nil {
		return ""
	}
	if !sameHost(parsed.Hostname(), hostLc) {
		return ""
	}
	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(segments) != 2 {
		return ""
	}
	if strings.ToLower(segments[0]) != ownerLc {
		return ""
	}
	return bareName(segments[1])
}

func BuildGitConfigArgs(cleanURL, token string) []string {
	gitConfigArgs := []string{"-c", "http.followRedirects=false", "-c", "protocol.file.allow=never"}
	if token != "" {
		username := "x-access-token"
		if strings.Contains(cleanURL, "https://gitlab.com/") {
			username = "oauth2"
		}
		if strings.Contains(cleanURL, "https://bitbucket.org/") {
			username = "x-token-auth"
		}
		authHeader := fmt.Sprintf("Authorization: Basic %s", base64.StdEncoding.EncodeToString([]byte(username+":"+token)))
		gitConfigArgs = append(gitConfigArgs, "-c", fmt.Sprintf("http.extraheader=%s", authHeader))
	}
	return gitConfigArgs
}

func SyncAndInitSubmodules(ctx context.Context, targetDir string, gitConfigArgs []string, logWriter io.Writer) error {
	gitmodulesPath := filepath.Join(targetDir, ".gitmodules")
	if _, err := os.Stat(gitmodulesPath); os.IsNotExist(err) {
		return nil
	}
	if logWriter != nil {
		fmt.Fprintf(logWriter, "📦 [Git] Updating submodules...\n")
	}
	subArgs := append(append([]string{}, gitConfigArgs...), "-C", targetDir, "submodule", "update", "--init", "--recursive", "--depth", "1")
	cmd := exec.CommandContext(ctx, "git", subArgs...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "⚠️ [Git] Submodule update warning: %s\n", string(out))
		}
		return nil
	}
	if logWriter != nil {
		fmt.Fprintf(logWriter, "✅ [Git] Submodules synchronized.\n")
	}
	return nil
}
