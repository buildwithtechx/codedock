package networking

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParentWildcardDomain(t *testing.T) {
	tests := []struct {
		domain   string
		expected string
	}{
		{"api.example.com", "*.example.com"},
		{"deep.sub.example.com", "*.sub.example.com"},
		{"*.example.com", ""},
		{"example.com", ""},
		{"localhost", ""},
		{"", ""},
		{"APP.TEST.IO", "*.test.io"},
	}

	for _, tc := range tests {
		got := ParentWildcardDomain(tc.domain)
		if got != tc.expected {
			t.Errorf("ParentWildcardDomain(%q) = %q; want %q", tc.domain, got, tc.expected)
		}
	}
}

func TestResolveCertDomain(t *testing.T) {
	mockCerts := map[string]bool{
		"/certs/*.example.com/fullchain.pem":     true,
		"/certs/direct.domain.com/fullchain.pem": true,
	}

	exists := func(p string) bool {
		return mockCerts[p]
	}

	// Direct match
	if got := ResolveCertDomain("/certs", "direct.domain.com", exists); got != "direct.domain.com" {
		t.Errorf("expected direct.domain.com, got %q", got)
	}

	// Subdomain inheriting wildcard
	if got := ResolveCertDomain("/certs", "api.example.com", exists); got != "*.example.com" {
		t.Errorf("expected *.example.com, got %q", got)
	}

	// No cert found fallback
	if got := ResolveCertDomain("/certs", "other.xyz", exists); got != "other.xyz" {
		t.Errorf("expected other.xyz, got %q", got)
	}
}

func TestCleanupObsoleteCertbotRenewalConf(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "renewal_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	confFile := filepath.Join(tmpDir, "example.com.conf")
	if err := os.WriteFile(confFile, []byte("renewalparams"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := CleanupObsoleteCertbotRenewalConf(tmpDir, "example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(confFile); !os.IsNotExist(err) {
		t.Errorf("expected renewal conf to be deleted")
	}

	// Non-existent file should not error
	if err := CleanupObsoleteCertbotRenewalConf(tmpDir, "nonexistent.com"); err != nil {
		t.Errorf("unexpected error on nonexistent conf: %v", err)
	}
}
