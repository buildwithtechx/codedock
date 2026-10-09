package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	codedockhttp "codedock/internal/http"
	"codedock/internal/models"
	"codedock/internal/repositories"
	authservices "codedock/internal/services/auth"
	"codedock/internal/testdb"
	"codedock/internal/utils"
	"github.com/docker/docker/client"
)

func setupTestApp(t *testing.T) (*codedockhttp.Server, *sql.DB, string, string) {
	jwtSecret := "super-secure-integration-test-jwt-secret-32-chars!"
	t.Setenv("CODEDOCK_JWT_SECRET", jwtSecret)

	dataDir := t.TempDir()
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	db := testdb.Open(t)

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	server, err := codedockhttp.NewServer(db, vlt, nil, nil, dockerClient, dataDir)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	return server, db, jwtSecret, dataDir
}

func TestSecurityCSRFProtection(t *testing.T) {
	server, db, _, _ := setupTestApp(t)
	defer db.Close()

	body, _ := json.Marshal(map[string]string{
		"email":    "csrf_test@example.com",
		"password": "Password123!",
		"name":     "CSRF Tester",
	})

	for _, path := range []string{"/api/auth/signup", "/api/auth/signin", "/api/auth/refresh", "/api/v1/auth/signup"} {
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected %s without CSRF token to be rejected, got %d: %s", path, rec.Code, rec.Body.String())
		}
	}
	bootstrap := httptest.NewRecorder()
	bootstrapRequest := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	bootstrapRequest.Header.Set("Sec-Fetch-Site", "same-site")
	server.ServeHTTP(bootstrap, bootstrapRequest)
	var token struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(bootstrap.Body.Bytes(), &token); err != nil || token.Token == "" || token.Token == "_echo_csrf_using_sec_fetch_site_" {
		t.Fatalf("expected bootstrap token, got %s: %v", bootstrap.Body.String(), err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", token.Token)
	req.Header.Set("Sec-Fetch-Site", "same-site")
	for _, cookie := range bootstrap.Result().Cookies() {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Fatalf("expected signup with CSRF token to succeed, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestDeactivatedUserAuthenticationBlocking(t *testing.T) {
	server, db, _, _ := setupTestApp(t)
	defer db.Close()

	userRepo := repositories.NewUserRepo(db)
	tokenService, err := authservices.NewTokenService()
	if err != nil {
		t.Fatalf("failed to initialize token service: %v", err)
	}

	ctx := context.Background()
	user := &models.User{
		Email:        "deactivated@example.com",
		PasswordHash: "hashed",
		Name:         "Deactivated User",
		Role:         models.UserRoleMember,
		IsActive:     false,
	}

	if err := userRepo.CreateUser(ctx, user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tokenStr, err := tokenService.GenerateToken(user)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized && rec.Code != http.StatusForbidden {
		t.Fatalf("expected deactivated user token request to be rejected with 401/403, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestWorkerTokenAuthenticationAndRevocation(t *testing.T) {
	_, db, _, _ := setupTestApp(t)
	defer db.Close()

	serverRepo := repositories.NewServerRepository(db, nil)
	ctx := context.Background()

	srv := &models.Server{
		Name:        "worker-test-node",
		IPAddress:   "127.0.0.1",
		Status:      models.ServerStatusOnline,
		WorkerToken: "worker-secret-token-12345",
	}

	if err := serverRepo.Create(ctx, srv); err != nil {
		t.Fatalf("failed to create server node: %v", err)
	}

	fetched, err := serverRepo.GetByToken(ctx, srv.WorkerToken)
	if err != nil || fetched == nil {
		t.Fatalf("failed to authenticate active worker by token: %v", err)
	}

	if err := serverRepo.Delete(ctx, srv.ID); err != nil {
		t.Fatalf("failed to delete/revoke server node: %v", err)
	}

	revoked, err := serverRepo.GetByToken(ctx, srv.WorkerToken)
	if err == nil && revoked != nil {
		t.Fatalf("expected revoked worker token lookup to fail")
	}
}
