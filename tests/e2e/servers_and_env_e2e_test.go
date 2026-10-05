package e2e_test

import (
	"net/http"
	"testing"
)

func TestE2EServerLifecycle(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	signupRes, body, signupErr := h.post("/api/auth/signup", map[string]string{
		"email":    "server_admin@codedock.local",
		"password": "Password123!",
		"name":     "Server Admin",
	}, nil)
	if signupErr != nil || (signupRes.StatusCode != http.StatusOK && signupRes.StatusCode != http.StatusCreated) {
		t.Fatalf("expected signup to succeed, got status %d, body %v", signupRes.StatusCode, body)
	}

	localServerReq := map[string]any{
		"name":    "Local Daemon Engine",
		"isLocal": true,
	}

	res, body, err := h.post("/api/servers", localServerReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create local server to succeed, got status %d, body %v", res.StatusCode, body)
	}

	localData, _ := body["data"].(map[string]any)
	localID, _ := localData["id"].(string)
	if localID == "" {
		t.Fatalf("expected valid local server id, got empty")
	}
	if isLocal, _ := localData["isLocal"].(bool); !isLocal {
		t.Fatalf("expected isLocal to be true, got %v", localData["isLocal"])
	}

	if _, err := h.db.Exec("UPDATE users SET plan_type = 'pro' WHERE email = ?", "server_admin@codedock.local"); err != nil {
		t.Fatalf("failed to update user plan: %v", err)
	}

	remoteServerReq := map[string]any{
		"name":          "Remote Worker Node",
		"ipAddress":     "198.51.100.25",
		"sshHost":       "198.51.100.25",
		"sshPort":       2222,
		"sshUser":       "deploy",
		"sshAuthMethod": "key",
		"sshKey":        testSSHPrivateKey(t),
	}

	res, body, err = h.post("/api/servers", remoteServerReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create remote server to succeed, got status %d, body %v", res.StatusCode, body)
	}

	remoteData, _ := body["data"].(map[string]any)
	remoteID, _ := remoteData["id"].(string)
	if remoteID == "" {
		t.Fatalf("expected valid remote server id, got empty")
	}
	if token, _ := remoteData["workerToken"].(string); token != "********" {
		t.Fatalf("expected worker token to be redacted, got %s", token)
	}

	res, body, err = h.get("/api/servers/"+remoteID, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/servers/%s to return 200, got status %d", remoteID, res.StatusCode)
	}
	fetchedData, _ := body["data"].(map[string]any)
	if name, _ := fetchedData["name"].(string); name != "Remote Worker Node" {
		t.Fatalf("expected name 'Remote Worker Node', got %s", name)
	}

	updateReq := map[string]any{
		"name":    "Renamed Worker Node",
		"sshPort": 2200,
	}
	res, body, err = h.patch("/api/servers/"+remoteID, updateReq, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected PATCH /api/servers/%s to return 200, got status %d, body %v", remoteID, res.StatusCode, body)
	}
	updatedData, _ := body["data"].(map[string]any)
	if name, _ := updatedData["name"].(string); name != "Renamed Worker Node" {
		t.Fatalf("expected name to be updated, got %s", name)
	}

	res, body, err = h.get("/api/servers", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/servers to return 200, got status %d", res.StatusCode)
	}
	serversList, _ := body["data"].([]any)
	if len(serversList) < 2 {
		t.Fatalf("expected at least 2 servers in list, got %d", len(serversList))
	}

	res, _, err = h.delete("/api/servers/"+remoteID, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent) {
		t.Fatalf("expected DELETE /api/servers/%s to succeed, got status %d", remoteID, res.StatusCode)
	}

	res, _, err = h.get("/api/servers/"+remoteID, nil)
	if err != nil || res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected GET /api/servers/%s to return 404 after deletion, got status %d", remoteID, res.StatusCode)
	}
}

func TestE2EProjectEnvVarsLifecycle(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	signupRes, _, signupErr := h.post("/api/auth/signup", map[string]string{
		"email":    "env_admin@codedock.local",
		"password": "Password123!",
		"name":     "Env Admin",
	}, nil)
	if signupErr != nil || (signupRes.StatusCode != http.StatusOK && signupRes.StatusCode != http.StatusCreated) {
		t.Fatalf("expected signup to succeed, got status %d", signupRes.StatusCode)
	}

	projectReq := map[string]string{
		"name":        "Env Test Project",
		"description": "Testing environment variables sync",
	}

	res, body, err := h.post("/api/projects", projectReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create project to succeed, got status %d, body %v", res.StatusCode, body)
	}

	projectData, _ := body["data"].(map[string]any)
	projectID, _ := projectData["id"].(string)
	if projectID == "" {
		t.Fatalf("expected valid project id, got empty")
	}

	envPayload := map[string]any{
		"variables": map[string]string{
			"APP_ENV":   "production",
			"PORT":      "8080",
			"LOG_LEVEL": "info",
		},
	}
	res, body, err = h.post("/api/projects/"+projectID+"/env", envPayload, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected set env vars to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, body, err = h.get("/api/projects/"+projectID+"/env", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/projects/%s/env to return 200, got status %d", projectID, res.StatusCode)
	}

	varsData, _ := body["data"].(map[string]any)
	if varsData["APP_ENV"] != "production" || varsData["PORT"] != "8080" {
		t.Fatalf("expected variables to match saved values, got %v", varsData)
	}
}
