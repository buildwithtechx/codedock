package e2e_test

import (
	"net/http"
	"testing"
)

func TestE2ESystemAndServiceSettings(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	res, body, err := h.get("/api/system/public", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/system/public to succeed, got status %d", res.StatusCode)
	}

	h.post("/api/auth/signup", map[string]string{
		"email":    "sys_settings_admin@codedock.local",
		"password": "Password123!",
		"name":     "Settings Admin",
	}, nil)

	res, body, err = h.get("/api/system/stats", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/system/stats to succeed, got status %d", res.StatusCode)
	}

	res, body, err = h.get("/api/notifications", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/notifications to succeed, got status %d", res.StatusCode)
	}

	res, body, err = h.put("/api/notifications", map[string]any{
		"slackWebhookUrl": "https://hooks.slack.com/services/TEST/B000/XXXX",
		"notifyOnDeploy":  true,
		"notifyOnBackup":  true,
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected PUT /api/notifications to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, body, err = h.get("/api/ai", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/ai to succeed, got status %d", res.StatusCode)
	}

	res, body, err = h.put("/api/ai", map[string]any{
		"provider": "openai",
		"apiKey":   "sk-proj-testkey1234567890",
		"model":    "gpt-4o",
		"enabled":  true,
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected PUT /api/ai to succeed, got status %d, body %v", res.StatusCode, body)
	}
}
