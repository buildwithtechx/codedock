package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	codedockhttp "codedock.run/codedock/internal/http"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	authservices "codedock.run/codedock/internal/services/auth"
	"codedock.run/codedock/internal/utils"
	"github.com/docker/docker/client"
)

func setupTestApp(t *testing.T) (*codedockhttp.Server, *sql.DB, string, string) {
	dataDir := t.TempDir()
	vlt, err := utils.NewVault(dataDir)
	if err != nil {
		t.Fatalf("failed to create vault: %v", err)
	}

	dbPath := filepath.Join(dataDir, "codedock.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	jwtSecret := "super-secure-integration-test-jwt-secret-32-chars!"
	t.Setenv("CODEDOCK_JWT_SECRET", jwtSecret)

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

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Fatalf("expected initial signup to succeed under CSRF skipper, got %d: %s", rec.Code, rec.Body.String())
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

	serverRepo := repositories.NewServerRepository(db)
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
