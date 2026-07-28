package e2e_test

import (
	"net/http"
	"testing"
)

func TestE2EDatabaseProvisioningAndCredentialsReveal(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	signupRes, _, signupErr := h.post("/api/auth/signup", map[string]string{
		"email":    "db_owner@codedock.local",
		"password": "Password123!",
		"name":     "Database Owner",
	}, nil)
	if signupErr != nil || (signupRes.StatusCode != http.StatusOK && signupRes.StatusCode != http.StatusCreated) {
		t.Fatalf("expected signup to succeed, got status %d, err %v", signupRes.StatusCode, signupErr)
	}

	projRes, projBody, err := h.post("/api/projects", map[string]string{
		"name": "DB E2E Project",
	}, nil)
	if err != nil || (projRes.StatusCode != http.StatusOK && projRes.StatusCode != http.StatusCreated) {
		t.Fatalf("failed to create test project for database: %v", err)
	}

	projData, _ := projBody["data"].(map[string]any)
	projectID, _ := projData["id"].(string)

	createDBReq := map[string]any{
		"projectId":    projectID,
		"name":         "postgres-prod-db",
		"engine":       "postgres",
		"version":      "16",
		"databaseName": "production",
		"username":     "pguser",
		"password":     "SecretPass123!",
		"port":         5432,
	}

	res, body, err := h.post("/api/databases", createDBReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create database to succeed, got status %d, body %v", res.StatusCode, body)
	}

	dbData, _ := body["data"].(map[string]any)
	dbID, _ := dbData["id"].(string)
	if dbID == "" {
		t.Fatalf("expected valid database id, got empty")
	}

	res, body, err = h.get("/api/databases/"+dbID, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/databases/%s to return 200, got status %d", dbID, res.StatusCode)
	}

	getDBData, _ := body["data"].(map[string]any)
	if _, exists := getDBData["password"]; exists {
		t.Fatalf("expected password to be omitted from standard GET /api/databases response")
	}

	res, body, err = h.post("/api/databases/"+dbID+"/credentials/reveal", nil, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected reveal credentials to return 200, got status %d, body %v", res.StatusCode, body)
	}

	credData, _ := body["data"].(map[string]any)
	revealedPassword, _ := credData["password"].(string)
	if revealedPassword != "SecretPass123!" {
		t.Fatalf("expected revealed password 'SecretPass123!', got '%s'", revealedPassword)
	}

	backupReq := map[string]any{
		"databaseId":    dbID,
		"name":          "Daily Database Backup",
		"schedule":      "0 2 * * *",
		"retentionDays": 7,
		"backupEnabled": true,
	}
	res, body, err = h.post("/api/backups", backupReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create backup config to succeed, got status %d, body %v", res.StatusCode, body)
	}
}
