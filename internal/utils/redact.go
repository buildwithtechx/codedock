package utils

import (
	"regexp"
	"strings"
)

var redactionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|secret|key|token|auth|dsn)\s*[:=]\s*["']?\*{3,}`),
	regexp.MustCompile(`:\/\/[^:@\s]+:\*{3,}@`),
	regexp.MustCompile(`["'][^"']*\*{4,}[^"']*["']`),
	regexp.MustCompile(`\*{6,}`),
}

// HasRedactedConfiguration detects unrecoverable redacted secrets (e.g. "://***@" or "***")
// in configuration dumps, connection strings, or build logs.
func HasRedactedConfiguration(values ...string) bool {
	if len(values) == 0 {
		return false
	}
	for _, val := range values {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			continue
		}
		if strings.Contains(val, "://***@") || trimmed == "***" || strings.Contains(val, "=\"***\"") || strings.Contains(val, "=***") {
			return true
		}
		for _, re := range redactionPatterns {
			if re.MatchString(val) {
				return true
			}
		}
	}
	return false
}
