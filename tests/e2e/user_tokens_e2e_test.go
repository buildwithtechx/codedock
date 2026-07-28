package e2e_test

import (
	"net/http"
	"testing"
)

func TestE2EUserProfileAndPATs(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	res, body, err := h.post("/api/auth/signup", map[string]string{
		"email":    "pat_user@codedock.local",
		"password": "Password123!",
		"name":     "Initial Name",
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected signup to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, body, err = h.get("/api/profile", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/profile to succeed, got status %d", res.StatusCode)
	}
	profileData, _ := body["data"].(map[string]any)
	if profileData["name"] != "Initial Name" {
		t.Fatalf("expected profile name Initial Name, got %v", profileData["name"])
	}

	res, body, err = h.put("/api/profile", map[string]string{
		"name": "Updated Name",
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected PUT /api/profile to succeed, got status %d", res.StatusCode)
	}
	profileData, _ = body["data"].(map[string]any)
	if profileData["name"] != "Updated Name" {
		t.Fatalf("expected updated profile name Updated Name, got %v", profileData["name"])
	}

	res, body, err = h.post("/api/profile/tokens", map[string]string{
		"name": "Test CLI PAT",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create PAT to succeed, got status %d, body %v", res.StatusCode, body)
	}
	patData, _ := body["data"].(map[string]any)
	rawToken, _ := patData["token"].(string)
	patInner, _ := patData["pat"].(map[string]any)
	patID, _ := patInner["id"].(string)
	if patID == "" || rawToken == "" {
		t.Fatalf("expected valid PAT id and raw token, got id=%q token=%q", patID, rawToken)
	}

	res, body, err = h.get("/api/profile/tokens", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/profile/tokens to succeed, got status %d", res.StatusCode)
	}
	tokensList, _ := body["data"].([]any)
	if len(tokensList) == 0 {
		t.Fatalf("expected non-empty tokens list, got %v", tokensList)
	}

	res, body, err = h.delete("/api/profile/tokens/"+patID, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected delete PAT to succeed, got status %d", res.StatusCode)
	}
}
