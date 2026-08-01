package e2e_test

import (
	"net/http"
	"testing"
)

func TestEDNSProviderCredentialsAreConfiguredAndMasked(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	res, body, err := h.post("/api/auth/signup", map[string]string{
		"email":    "dns-provider-owner@codedock.local",
		"password": "Password123!",
		"name":     "DNS Provider Owner",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("signup failed: status=%d body=%v err=%v", res.StatusCode, body, err)
	}

	res, body, err = h.put("/api/settings", map[string]string{
		"spaceshipApiKey":    "spaceship-test-key",
		"spaceshipApiSecret": "spaceship-test-secret",
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("provider settings update failed: status=%d body=%v err=%v", res.StatusCode, body, err)
	}
	settings, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected settings response data, got %T", body["data"])
	}
	if settings["spaceshipApiSecret"] != "********" {
		t.Fatalf("expected Spaceship secret to be masked in update response, got %v", settings["spaceshipApiSecret"])
	}

	res, body, err = h.get("/api/settings", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("provider settings read failed: status=%d body=%v err=%v", res.StatusCode, body, err)
	}
	settings, ok = body["data"].(map[string]any)
	if !ok || settings["spaceshipApiSecret"] != "********" {
		t.Fatalf("expected Spaceship secret to remain masked, got %v", settings["spaceshipApiSecret"])
	}
}
