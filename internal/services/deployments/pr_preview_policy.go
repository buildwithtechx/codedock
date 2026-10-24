package deployments

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"codedock/internal/models"
	"codedock/internal/utils"
)

var hex40Regex = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)

type PullRequestInfo struct {
	Number   int    `json:"number"`
	State    string `json:"state"`
	Draft    bool   `json:"draft"`
	HeadSHA  string `json:"headSha"`
	HeadRef  string `json:"headRef"`
	HeadRepo string `json:"headRepo"`
	BaseRef  string `json:"baseRef"`
	BaseRepo string `json:"baseRepo"`
}

// ParseRepoOwnerAndName extracts the owner and repo from a git repository URL.
func ParseRepoOwnerAndName(rawURL string) (string, string) {
	trimmed := strings.TrimSuffix(strings.TrimSpace(rawURL), ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2], parts[len(parts)-1]
	}
	return "", ""
}

// normalizeRepoIdentity returns a canonical lowercase repository identity.
func normalizeRepoIdentity(source, owner, repo string) string {
	src := strings.ToLower(strings.TrimSpace(source))
	o := strings.ToLower(strings.TrimSpace(owner))
	r := strings.ToLower(strings.TrimSpace(repo))
	r = strings.TrimSuffix(r, ".git")
	if src == "" {
		src = "github"
	}
	return fmt.Sprintf("%s:%s/%s", src, o, r)
}

// IsPreviewEligible checks the PR eligibility gate:
// - must be open
// - must not be draft
// - base ref must match project branch
// - head and base repos must match project repo (case-insensitive)
// - head SHA must be a valid 40-character hex string
func IsPreviewEligible(pr PullRequestInfo, targetOwner, targetRepo, targetBranch string) (bool, string) {
	if strings.ToLower(pr.State) != "open" {
		return false, "PR is not open"
	}
	if pr.Draft {
		return false, "PR is a draft"
	}
	if pr.BaseRef != targetBranch {
		return false, fmt.Sprintf("PR base branch %q does not match project branch %q", pr.BaseRef, targetBranch)
	}

	expectedRepo := strings.ToLower(strings.TrimSpace(targetOwner) + "/" + strings.TrimSuffix(strings.TrimSpace(targetRepo), ".git"))
	if strings.ToLower(pr.BaseRepo) != expectedRepo {
		return false, fmt.Sprintf("PR base repo %q does not match project repo %q", pr.BaseRepo, expectedRepo)
	}
	if pr.HeadRepo != "" && strings.ToLower(pr.HeadRepo) != expectedRepo {
		return false, fmt.Sprintf("PR head repo %q does not match project repo %q", pr.HeadRepo, expectedRepo)
	}
	if !hex40Regex.MatchString(pr.HeadSHA) {
		return false, "PR head commit SHA must be a valid 40-character hex string"
	}

	return true, ""
}

// validatePreviewVolumes ensures only isolated named volumes are mounted,
// rejecting host bind mounts (starting with / or containing . or ..).
func validatePreviewVolumes(volumes []models.ServiceVolume) error {
	for _, v := range volumes {
		source := strings.TrimSpace(v.HostPath)
		if strings.HasPrefix(source, "/") || strings.HasPrefix(source, "./") || source == "." || source == ".." || strings.Contains(source, `\`) {
			return fmt.Errorf("PR previews require isolated named volumes; host bind mounts (%s) must be removed from the preview template", source)
		}
	}
	return nil
}

// validatePreviewCommands checks that commands do not embed environment variables,
// enforcing that embedded variables are moved to Environment first.
func validatePreviewCommands(commands ...string) error {
	for _, cmd := range commands {
		if cmd == "" {
			continue
		}
		res := utils.InspectCommandEnvironment(cmd)
		if len(res.Keys) > 0 {
			return fmt.Errorf("move embedded command variables (%s) to Environment before enabling previews; configure preview values separately", strings.Join(res.Keys, ", "))
		}
	}
	return nil
}

// generatePreviewDomain generates an isolated domain for the preview environment.
func generatePreviewDomain(baseDomain, magicDomain, appName, projectID string, prNumber int) string {
	hash := sha256.Sum256(fmt.Appendf(nil, "%s:%d", projectID, prNumber))
	suffix := hex.EncodeToString(hash[:])[:6]
	cleanApp := strings.ToLower(strings.ReplaceAll(appName, "_", "-"))

	if baseDomain != "" {
		return fmt.Sprintf("%s-pr-%d-%s.%s", cleanApp, prNumber, suffix, strings.TrimPrefix(baseDomain, "."))
	}
	if magicDomain == "" {
		magicDomain = "sslip.io"
	}
	return fmt.Sprintf("%s-pr-%d-%s.%s", cleanApp, prNumber, suffix, strings.TrimPrefix(magicDomain, "."))
}
