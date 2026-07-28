package e2e_test

import (
	"net/http"
	"testing"
)

func TestE2EProjectTokensAndSharedEnvs(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	h.post("/api/auth/signup", map[string]string{
		"email":    "proj_admin@codedock.local",
		"password": "Password123!",
		"name":     "Project Admin",
	}, nil)

	res, body, err := h.post("/api/projects", map[string]string{
		"name":        "E2E Shared Project",
		"description": "Project for testing shared envs and API tokens",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create project to succeed, got status %d, body %v", res.StatusCode, body)
	}

	projData, _ := body["data"].(map[string]any)
	projID, _ := projData["id"].(string)
	if projID == "" {
		t.Fatalf("expected valid project id, got empty")
	}

	res, body, err = h.put("/api/projects/"+projID+"/env", map[string]string{
		"SHARED_API_KEY": "secret_key_123",
		"PUBLIC_APP_URL": "https://myapp.example.com",
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected PUT /api/projects/:id/env to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, body, err = h.get("/api/projects/"+projID+"/env", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/projects/:id/env to succeed, got status %d", res.StatusCode)
	}

	res, body, err = h.post("/api/projects/"+projID+"/tokens", map[string]string{
		"name":  "CI/CD Deploy Token",
		"scope": "env:read",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create project token to succeed, got status %d, body %v", res.StatusCode, body)
	}

	tokenData, _ := body["data"].(map[string]any)
	tokenID, _ := tokenData["id"].(string)
	if tokenID == "" {
		t.Fatalf("expected valid project token id, got empty")
	}

	res, body, err = h.get("/api/projects/"+projID+"/tokens", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/projects/:id/tokens to succeed, got status %d", res.StatusCode)
	}

	res, body, err = h.delete("/api/projects/"+projID+"/tokens/"+tokenID, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent) {
		t.Fatalf("expected delete project token to succeed, got status %d", res.StatusCode)
	}
}
