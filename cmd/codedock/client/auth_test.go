package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"
)

func TestVerifyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-pat-token" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "success",
			"message": "Profile retrieved",
			"data": &models.User{
				ID:    "user-123",
				Email: "admin@codedock.local",
				Role:  models.UserRoleAdmin,
			},
		})
	}))
	defer ts.Close()

	cli := NewClient(ts.URL, "")
	user, err := cli.VerifyToken("test-pat-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "admin@codedock.local" {
		t.Errorf("expected admin@codedock.local, got %s", user.Email)
	}
}

func TestExchangeLocalCLISession(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/system/local-cli-session":
			if r.Header.Get("X-Internal-Token") != "secret-operator-key" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"token":     "session-jwt-xyz",
				"expiresAt": "2026-10-17T12:00:00Z",
			})
		case "/api/v1/auth/me":
			if r.Header.Get("Authorization") != "Bearer session-jwt-xyz" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": &models.User{
					ID:    "admin-root",
					Email: "root@codedock.internal",
					Role:  models.UserRoleAdmin,
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	cli := NewClient(ts.URL, "")
	resp, err := cli.ExchangeLocalCLISession("secret-operator-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "session-jwt-xyz" {
		t.Errorf("expected session-jwt-xyz, got %s", resp.Token)
	}
	if resp.User.Email != "root@codedock.internal" {
		t.Errorf("expected root@codedock.internal, got %s", resp.User.Email)
	}
}
