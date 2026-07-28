package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"codedock.run/codedock/internal/repositories"
)

func TestE2ESetupStatusAndFirstUserSignup(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	res, body, err := h.get("/api/system/setup-status", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected setup-status to return 200, got status %d, err %v", res.StatusCode, err)
	}

	data, _ := body["data"].(map[string]any)
	if data["setupRequired"] != true {
		t.Fatalf("expected setupRequired to be true for new instance, got %v", data["setupRequired"])
	}

	signupPayload := map[string]string{
		"email":    "owner@codedock.local",
		"password": "Password123!",
		"name":     "Admin Owner",
	}

	res, body, err = h.post("/api/auth/signup", signupPayload, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected signup to succeed, got status %d, body %v, err %v", res.StatusCode, body, err)
	}

	res, body, err = h.get("/api/system/setup-status", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected setup-status after signup to return 200, got status %d", res.StatusCode)
	}
	data, _ = body["data"].(map[string]any)
	if data["setupRequired"] != false {
		t.Fatalf("expected setupRequired to be false after first user signup, got %v", data["setupRequired"])
	}

	res, body, err = h.get("/api/auth/me", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected /api/auth/me to return 200, got status %d, body %v", res.StatusCode, body)
	}

	userData, _ := body["data"].(map[string]any)
	if userData["role"] != "owner" {
		t.Fatalf("expected first registered user to have role 'owner', got %v", userData["role"])
	}
}

func TestE2ESignInAndDeactivatedUserBlocking(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	signupPayload := map[string]string{
		"email":    "user@codedock.local",
		"password": "Password123!",
		"name":     "Test User",
	}
	h.post("/api/auth/signup", signupPayload, nil)

	signinPayload := map[string]string{
		"email":    "user@codedock.local",
		"password": "Password123!",
	}
	res, body, err := h.post("/api/auth/signin", signinPayload, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected signin to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, _, err = h.get("/api/projects?organizationId=none", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected authenticated GET /api/projects?organizationId=none to return 200, got status %d", res.StatusCode)
	}

	userRepo := repositories.NewUserRepo(h.db)
	user, err := userRepo.GetUserByEmail(context.Background(), "user@codedock.local")
	if err != nil || user == nil {
		t.Fatalf("failed to query created user: %v", err)
	}

	user.IsActive = false
	if err := userRepo.UpdateUser(context.Background(), user); err != nil {
		t.Fatalf("failed to deactivate user: %v", err)
	}

	res, _, err = h.get("/api/projects?organizationId=none", nil)
	if res.StatusCode != http.StatusUnauthorized && res.StatusCode != http.StatusForbidden {
		t.Fatalf("expected deactivated user request to be blocked with 401/403, got %d", res.StatusCode)
	}
}
