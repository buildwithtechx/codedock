package networking

import (
	"os"
	"path/filepath"
	"strings"
)

// ParentWildcardDomain returns the wildcard representation for a subdomain (e.g. "*.example.com" for "api.example.com").
// Returns empty string if the domain is already a wildcard, has no dots, or is an apex domain.
func ParentWildcardDomain(domain string) string {
	clean := strings.TrimSpace(strings.ToLower(domain))
	if clean == "" || strings.HasPrefix(clean, "*.") {
		return ""
	}
	firstDot := strings.Index(clean, ".")
	if firstDot <= 0 {
		return ""
	}
	parent := clean[firstDot+1:]
	if !strings.Contains(parent, ".") {
		return ""
	}
	return "*." + parent
}

// ResolveCertDomain returns the directory name containing a valid certificate.
// If direct cert doesn't exist, checks if a parent wildcard cert exists.
func ResolveCertDomain(certDir, domain string, exists func(string) bool) string {
	clean := strings.TrimSpace(strings.ToLower(domain))
	if clean == "" {
		return ""
	}
	if exists == nil {
		exists = func(p string) bool {
			_, err := os.Stat(p)
			return err == nil
		}
	}

	directCert := filepath.Join(certDir, clean, "fullchain.pem")
	if exists(directCert) {
		return clean
	}

	wildcard := ParentWildcardDomain(clean)
	if wildcard != "" {
		wildcardCert := filepath.Join(certDir, wildcard, "fullchain.pem")
		if exists(wildcardCert) {
			return wildcard
		}
	}

	return clean
}

// CleanupObsoleteCertbotRenewalConf removes the Certbot renewal configuration file
// when a custom or manual certificate replaces an ACME lineage, preventing renewal crashes.
func CleanupObsoleteCertbotRenewalConf(renewalDir, domain string) error {
	clean := strings.TrimSpace(strings.ToLower(domain))
	if clean == "" || renewalDir == "" {
		return nil
	}
	confPath := filepath.Join(renewalDir, clean+".conf")
	if _, err := os.Stat(confPath); err == nil {
		return os.Remove(confPath)
	}
	return nil
}
