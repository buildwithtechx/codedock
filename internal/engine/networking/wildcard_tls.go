package networking

import (
	"os"
	"path/filepath"
	"strings"
)

// ParentWildcardDomain returns the wildcard representation for a subdomain.
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

// ResolveCertDomain resolves direct certificate or parent wildcard certificate.
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

// CleanupObsoleteCertbotRenewalConf removes stale Certbot renewal config.
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
