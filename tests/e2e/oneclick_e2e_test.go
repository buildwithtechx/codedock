package e2e_test

import (
	"net/http"
	"strings"
	"testing"
)

func oneclickTestSignup(t *testing.T, h *e2eHarness, email string) {
	res, _, err := h.post("/api/auth/signup", map[string]string{
		"email":    email,
		"password": "Password123!",
		"name":     "Catalog Owner",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected signup to succeed, got status %d, err %v", res.StatusCode, err)
	}
}

func oneclickTestProject(t *testing.T, h *e2eHarness, name string) string {
	res, body, err := h.post("/api/projects", map[string]string{"name": name}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create project to succeed, got status %d, body %v", res.StatusCode, body)
	}
	projectData, _ := body["data"].(map[string]any)
	projectID, _ := projectData["id"].(string)
	if projectID == "" {
		t.Fatal("expected valid project id, got empty")
	}
	return projectID
}

func oneclickFindApp(t *testing.T, apps []any, id string) map[string]any {
	for _, entry := range apps {
		app, _ := entry.(map[string]any)
		if app["id"] == id {
			return app
		}
	}
	t.Fatalf("expected catalogue to contain %s", id)
	return nil
}

func oneclickFindVariable(t *testing.T, app map[string]any, key string) map[string]any {
	variables, _ := app["envVariables"].([]any)
	for _, entry := range variables {
		variable, _ := entry.(map[string]any)
		if variable["key"] == key {
			return variable
		}
	}
	t.Fatalf("expected app %v to expose variable %s", app["id"], key)
	return nil
}

func TestE2EOneClickCatalogAndReview(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()
	oneclickTestSignup(t, h, "owner_catalog@codedock.local")
	projectID := oneclickTestProject(t, h, "Catalog Project")

	res, body, err := h.get("/api/one-click", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected list one-click apps to succeed, got status %d, err %v", res.StatusCode, err)
	}
	apps, _ := body["data"].([]any)
	if len(apps) == 0 {
		t.Fatal("expected embedded catalogue apps, got none")
	}
	for _, entry := range apps {
		app, _ := entry.(map[string]any)
		if app["id"] == "" || app["name"] == "" {
			t.Fatalf("expected catalogue entries to carry id and name, got %v", app)
		}
		if _, ok := app["verified"]; !ok {
			t.Fatalf("expected catalogue entry %v to carry verified flag", app["id"])
		}
	}
	plausible := oneclickFindApp(t, apps, "plausible")
	baseURL := oneclickFindVariable(t, plausible, "BASE_URL")
	if baseURL["input"] != true {
		t.Fatalf("expected BASE_URL to be flagged as human input, got %v", baseURL)
	}

	res, body, err = h.get("/api/one-click/n8n", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected get one-click app to succeed, got status %d, err %v", res.StatusCode, err)
	}
	n8n, _ := body["data"].(map[string]any)
	secret := oneclickFindVariable(t, n8n, "N8N_ENCRYPTION_KEY")
	if secret["secret"] != true {
		t.Fatalf("expected encryption key to be flagged secret, got %v", secret)
	}

	res, _, err = h.get("/api/one-click/does-not-exist", nil)
	if err != nil || res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected unknown app to return 404, got status %d, err %v", res.StatusCode, err)
	}

	res, body, err = h.post("/api/one-click/review", map[string]any{
		"appId":     "n8n",
		"projectId": projectID,
		"name":      "demo",
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected review to succeed, got status %d, body %v", res.StatusCode, body)
	}
	preview, _ := body["data"].(map[string]any)
	if preview["digest"] == "" {
		t.Fatalf("expected review digest, got %v", preview)
	}
	generated, _ := preview["generatedSecrets"].([]any)
	if len(generated) != 1 || generated[0] != "N8N_ENCRYPTION_KEY" {
		t.Fatalf("expected generated encryption key, got %v", generated)
	}
	services, _ := preview["services"].([]any)
	if len(services) == 0 {
		t.Fatalf("expected preview services, got %v", preview)
	}
	firstService, _ := services[0].(map[string]any)
	env, _ := firstService["env"].([]any)
	masked := false
	for _, entry := range env {
		if line, _ := entry.(string); strings.Contains(line, "***") {
			masked = true
		}
	}
	if !masked {
		t.Fatalf("expected masked secrets in preview env, got %v", env)
	}
	document, _ := preview["composeYaml"].(string)
	if !strings.Contains(document, "***") {
		t.Fatalf("expected masked compose document, got:\n%s", document)
	}

	res, body, err = h.post("/api/one-click/review", map[string]any{
		"appId":       "plausible",
		"projectId":   projectID,
		"name":        "stats",
		"environment": map[string]string{"BASE_URL": "https://stats.example.com"},
	}, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected input review to succeed, got status %d, body %v", res.StatusCode, body)
	}
	inputPreview, _ := body["data"].(map[string]any)
	inputDocument, _ := inputPreview["composeYaml"].(string)
	if !strings.Contains(inputDocument, "https://stats.example.com") {
		t.Fatalf("expected resolved input in document, got:\n%s", inputDocument)
	}
	if !strings.Contains(inputDocument, "healthcheck:") || !strings.Contains(inputDocument, "pg_isready") {
		t.Fatalf("expected rendered healthcheck in document, got:\n%s", inputDocument)
	}

	res, _, err = h.post("/api/one-click/review", map[string]any{
		"appId": "n8n",
		"name":  "demo",
	}, nil)
	if err != nil || res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected missing projectId to return 400, got status %d, err %v", res.StatusCode, err)
	}

	res, body, err = h.post("/api/one-click/review", map[string]any{
		"appId":     "n8n",
		"projectId": projectID,
		"name":      "demo",
		"domain":    "not a domain",
	}, nil)
	if err != nil || res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected bad domain to return 400, got status %d, body %v", res.StatusCode, body)
	}
}

func TestE2EOneClickDeployRejectsStaleDigest(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()
	oneclickTestSignup(t, h, "owner_digest@codedock.local")
	projectID := oneclickTestProject(t, h, "Digest Project")

	res, body, err := h.post("/api/one-click/deploy", map[string]any{
		"appId":     "n8n",
		"projectId": projectID,
		"name":      "demo",
		"digest":    "stale-digest",
	}, nil)
	if err != nil || res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected stale digest to return 400, got status %d, body %v", res.StatusCode, body)
	}
	message, _ := body["message"].(string)
	if !strings.Contains(message, "review again") {
		t.Fatalf("expected review-again guidance, got %q", message)
	}
}

func TestE2ECatalogRequiresAuth(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	res, _, err := h.get("/api/one-click", nil)
	if err != nil || res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected catalogue list to require auth, got status %d, err %v", res.StatusCode, err)
	}
	res, _, err = h.get("/api/one-click/n8n", nil)
	if err != nil || res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected catalogue get to require auth, got status %d, err %v", res.StatusCode, err)
	}
	res, _, err = h.post("/api/one-click/review", map[string]any{"appId": "n8n"}, nil)
	if err != nil || res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected review to require auth, got status %d, err %v", res.StatusCode, err)
	}
	res, _, err = h.get("/api/examples", nil)
	if err != nil || res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected examples list to require auth, got status %d, err %v", res.StatusCode, err)
	}
}
