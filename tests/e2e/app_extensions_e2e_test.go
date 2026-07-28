package e2e_test

import (
	"net/http"
	"testing"
)

func TestE2EAppWebhooksVolumesAndLogDrains(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	h.post("/api/auth/signup", map[string]string{
		"email":    "app_owner@codedock.local",
		"password": "Password123!",
		"name":     "App Owner",
	}, nil)

	res, body, err := h.post("/api/projects", map[string]string{
		"name": "App Extensions Project",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create project to succeed, got status %d", res.StatusCode)
	}

	projData, _ := body["data"].(map[string]any)
	projID, _ := projData["id"].(string)

	res, body, err = h.get("/api/projects/"+projID+"/environments", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET environments to succeed, got status %d", res.StatusCode)
	}
	envsList, _ := body["data"].([]any)
	if len(envsList) == 0 {
		t.Fatalf("expected default environment to exist")
	}
	envData, _ := envsList[0].(map[string]any)
	envID, _ := envData["id"].(string)

	res, body, err = h.post("/api/environments/"+envID+"/apps", map[string]any{
		"projectId":   projID,
		"name":        "ext-web-app",
		"runtimeMode": "docker",
		"replicas":    1,
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create app service to succeed, got status %d, body %v", res.StatusCode, body)
	}

	appData, _ := body["data"].(map[string]any)
	appID, _ := appData["id"].(string)
	if appID == "" {
		t.Fatalf("expected valid app service id, got empty")
	}

	res, body, err = h.post("/api/apps/"+appID+"/webhooks", map[string]string{
		"name":   "Deploy Alert Webhook",
		"url":    "https://example.com/hooks/deploy",
		"events": "deploy_success,deploy_failed",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create webhook to succeed, got status %d, body %v", res.StatusCode, body)
	}

	whData, _ := body["data"].(map[string]any)
	whID, _ := whData["id"].(string)

	res, body, err = h.get("/api/apps/"+appID+"/webhooks", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET webhooks to succeed, got status %d", res.StatusCode)
	}

	if whID != "" {
		res, body, err = h.delete("/api/apps/"+appID+"/webhooks/"+whID, nil)
		if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent) {
			t.Fatalf("expected delete webhook to succeed, got status %d", res.StatusCode)
		}
	}

	res, body, err = h.post("/api/apps/"+appID+"/volumes", map[string]string{
		"name":          "app-uploads-data",
		"hostPath":      "/var/data/uploads",
		"containerPath": "/app/uploads",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create volume to succeed, got status %d, body %v", res.StatusCode, body)
	}

	volData, _ := body["data"].(map[string]any)
	volID, _ := volData["id"].(string)

	res, body, err = h.get("/api/apps/"+appID+"/volumes", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET volumes to succeed, got status %d", res.StatusCode)
	}

	if volID != "" {
		res, body, err = h.delete("/api/apps/"+appID+"/volumes/"+volID, nil)
		if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent) {
			t.Fatalf("expected delete volume to succeed, got status %d", res.StatusCode)
		}
	}

	res, body, err = h.post("/api/apps/"+appID+"/log-drains", map[string]string{
		"drainType":   "syslog",
		"endpointUrl": "https://logs.example.com/syslog",
	}, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create log drain to succeed, got status %d, body %v", res.StatusCode, body)
	}

	drainData, _ := body["data"].(map[string]any)
	drainID, _ := drainData["id"].(string)

	res, body, err = h.get("/api/apps/"+appID+"/log-drains", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET log drains to succeed, got status %d", res.StatusCode)
	}

	if drainID != "" {
		res, body, err = h.delete("/api/apps/"+appID+"/log-drains/"+drainID, nil)
		if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent) {
			t.Fatalf("expected delete log drain to succeed, got status %d", res.StatusCode)
		}
	}
}
