package projects

import (
	"testing"
	"time"
)

func TestCreatePolicyAndVerifyPassword(t *testing.T) {
	_, err := CreatePolicy("short", DefaultSessionSec)
	if err == nil {
		t.Fatal("expected error for short password")
	}

	validPass := "super-secure-production-password-2026"
	policy, err := CreatePolicy(validPass, DefaultSessionSec)
	if err != nil {
		t.Fatalf("unexpected error creating policy: %v", err)
	}

	if policy.Iterations != PBDKF2Iterations {
		t.Fatalf("expected %d iterations, got %d", PBDKF2Iterations, policy.Iterations)
	}
	if len(policy.Salt) != 64 {
		t.Fatalf("expected 64-char salt, got %d", len(policy.Salt))
	}
	if len(policy.Verifier) != 64 {
		t.Fatalf("expected 64-char verifier, got %d", len(policy.Verifier))
	}
	if len(policy.SigningKey) != 64 {
		t.Fatalf("expected 64-char signing key, got %d", len(policy.SigningKey))
	}

	if !VerifyPassword(validPass, policy) {
		t.Fatal("expected password verification to succeed")
	}
	if VerifyPassword("wrong-password-attempt", policy) {
		t.Fatal("expected wrong password to fail")
	}
}

func TestRevokeRouteSessions(t *testing.T) {
	validPass := "super-secure-production-password-2026"
	p1, _ := CreatePolicy(validPass, DefaultSessionSec)
	p2 := RevokeRouteSessions(p1)

	if p1.SigningKey == p2.SigningKey {
		t.Fatal("expected signing key to rotate")
	}
	if p1.Version == p2.Version {
		t.Fatal("expected version to rotate")
	}
	if p1.Verifier != p2.Verifier {
		t.Fatal("expected verifier to stay identical across session revokes")
	}

	// Verify old session token fails on revoked policy
	token := SignSessionToken(p1, "app.example.com", time.Now().Unix()+3600)
	if !VerifySessionToken(p1, "app.example.com", token) {
		t.Fatal("expected token to be valid for p1")
	}
	if VerifySessionToken(p2, "app.example.com", token) {
		t.Fatal("expected revoked policy p2 to reject p1 session token")
	}
}

func TestSessionTokens(t *testing.T) {
	policy, _ := CreatePolicy("super-secure-production-password-2026", DefaultSessionSec)
	host := "preview.example.com"

	validToken := SignSessionToken(policy, host, time.Now().Unix()+3600)
	if !VerifySessionToken(policy, host, validToken) {
		t.Fatal("expected valid session token to pass")
	}

	// Wrong host
	if VerifySessionToken(policy, "other.example.com", validToken) {
		t.Fatal("expected session token for different host to fail")
	}

	// Expired token
	expiredToken := SignSessionToken(policy, host, time.Now().Unix()-10)
	if VerifySessionToken(policy, host, expiredToken) {
		t.Fatal("expected expired session token to fail")
	}
}

func TestCSRFTokens(t *testing.T) {
	policy, _ := CreatePolicy("super-secure-production-password-2026", DefaultSessionSec)
	host := "preview.example.com"
	path := "/dashboard/settings"

	validCSRF := SignCSRFToken(policy, host, path, time.Now().Unix()+300)
	if !VerifyCSRFToken(policy, host, path, validCSRF) {
		t.Fatal("expected valid CSRF token to pass")
	}

	if VerifyCSRFToken(policy, host, "/different/path", validCSRF) {
		t.Fatal("expected CSRF token for different path to fail")
	}
}

func TestNormalizeHostname(t *testing.T) {
	cases := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"https://sub.domain.com/path", "sub.domain.com", false},
		{"http://SUB.DOMAIN.COM:8080/", "sub.domain.com", false},
		{"  my-app.codedock.run  ", "my-app.codedock.run", false},
		{"*.wildcard.com", "", true},
		{"", "", true},
	}

	for _, tc := range cases {
		got, err := NormalizeHostname(tc.input)
		if (err != nil) != tc.wantErr {
			t.Fatalf("NormalizeHostname(%q) err=%v, wantErr=%v", tc.input, err, tc.wantErr)
		}
		if got != tc.expected {
			t.Fatalf("NormalizeHostname(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
