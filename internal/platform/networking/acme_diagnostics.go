package networking

import (
	"fmt"
	"regexp"
	"strings"
)

const AcmeHTTP01Port = 49180

var (
	acmeChallengeURLRegex = regexp.MustCompile(`(?i)https?://\S*\.well-known/acme-challenge/\S*`)
	acmeStatusRegex       = regexp.MustCompile(`(?i)invalid response from https?://\S+:\s*(\d{3})\b`)
	acmeBindFailureRegex  = regexp.MustCompile(`(?i)bind: address already in use|could not bind to IPv4 or IPv6|failed to bind to :80\b`)
	acmeTimeoutRegex      = regexp.MustCompile(`(?i)timeout during connect|connection refused|failed to connect|timed? ?out`)
	acmeNotFoundRegex     = regexp.MustCompile(`(?i)404|unauthorized|not found`)
)

// SummarizeACMEFailure classifies Certbot and Let's Encrypt challenge errors.
// It avoids diagnosing a port 80 takeover collision when a gateway returns HTTP 502/504.
func SummarizeACMEFailure(output, domain string) string {
	challengeURL := acmeChallengeURLRegex.FindString(output)
	redirectedToHTTPS := strings.HasPrefix(strings.ToLower(challengeURL), "https:")

	statusMatches := acmeStatusRegex.FindStringSubmatch(output)
	var challengeStatus string
	if len(statusMatches) > 1 {
		challengeStatus = statusMatches[1]
	}

	if acmeBindFailureRegex.MatchString(output) {
		return fmt.Sprintf("Port 80 is occupied by another process on this server. Another web server or container is actively bound to :80. Stop the process holding port 80 before requesting a certificate for %s.", domain)
	}

	if challengeStatus == "502" || challengeStatus == "504" {
		redirectNote := ". "
		if redirectedToHTTPS {
			redirectNote = " after being redirected to HTTPS. "
		}
		return fmt.Sprintf(
			"The certificate challenge for %s returned HTTP %s%sA gateway could not reach its challenge responder. Check forwarding for /.well-known/acme-challenge/ and that Codedock's temporary listener on port %d is reachable from the edge, then retry. This response alone does not establish a port-80 conflict.",
			domain, challengeStatus, redirectNote, AcmeHTTP01Port,
		)
	}

	if challengeStatus == "525" || strings.Contains(strings.ToLower(output), "525") {
		return fmt.Sprintf("Cloudflare answered the challenge for %s with HTTP 525. Set SSL/TLS to Full (not Full Strict) or pause proxy mode until issuance finishes.", domain)
	}

	if acmeTimeoutRegex.MatchString(output) {
		return fmt.Sprintf("Port 80 for %s isn't reachable from the internet — a cloud security group or firewall is likely blocking inbound traffic on port 80.", domain)
	}

	if acmeNotFoundRegex.MatchString(output) {
		return fmt.Sprintf("The ACME challenge token was not found (HTTP 404). Verify that %s resolves to this server's IP address.", domain)
	}

	return fmt.Sprintf("Let's Encrypt could not verify the challenge for %s. Check DNS and server connectivity.", domain)
}
