package networking

import (
	"strings"
	"testing"
)

const opener = "Certbot failed to authenticate some domains (authenticator: webroot)."

func TestSummarizeACMEFailure(t *testing.T) {
	t.Run("HTTP 502/504 gateway failure", func(t *testing.T) {
		out := opener + "\n  Detail: Invalid response from http://app.example.test/.well-known/acme-challenge/abc: 502 Bad Gateway"
		summary := SummarizeACMEFailure(out, "app.example.test")
		if !strings.Contains(summary, "HTTP 502") {
			t.Errorf("expected summary to contain HTTP 502, got: %s", summary)
		}
		if !strings.Contains(summary, "49180") {
			t.Errorf("expected summary to mention port 49180, got: %s", summary)
		}
		if strings.Contains(summary, "Port 80 is occupied") {
			t.Errorf("summary should not infer port 80 takeover from 502, got: %s", summary)
		}
	})

	t.Run("Does not diagnose gateway from challenge token digits", func(t *testing.T) {
		out := opener + "\n  Detail: Invalid response from http://502.example.test/.well-known/acme-challenge/token-504: 404"
		summary := SummarizeACMEFailure(out, "502.example.test")
		if strings.Contains(summary, "HTTP 502") || strings.Contains(summary, "HTTP 504") {
			t.Errorf("summary falsely identified 502/504: %s", summary)
		}
		if !strings.Contains(summary, "HTTP 404") {
			t.Errorf("expected 404 diagnosis, got: %s", summary)
		}
	})

	t.Run("Port 80 bind failure", func(t *testing.T) {
		out := "Problem binding to port 80: Could not bind to IPv4 or IPv6."
		summary := SummarizeACMEFailure(out, "app.example.test")
		if !strings.Contains(summary, "Port 80 is occupied") {
			t.Errorf("expected port 80 collision diagnosis, got: %s", summary)
		}
	})

	t.Run("Redirected to HTTPS gateway failure", func(t *testing.T) {
		out := opener + "\n  Detail: Invalid response from https://app.example.test/.well-known/acme-challenge/abc: 502"
		summary := SummarizeACMEFailure(out, "app.example.test")
		if !strings.Contains(summary, "redirected to HTTPS") {
			t.Errorf("expected note about HTTPS redirection, got: %s", summary)
		}
	})
}
