package e2e_test

import (
	"net/http"
	"testing"
)

func TestE2EProjectAndServiceLifecycle(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	signupRes, _, signupErr := h.post("/api/auth/signup", map[string]string{
		"email":    "owner_project@codedock.local",
		"password": "Password123!",
		"name":     "Project Owner",
	}, nil)
	if signupErr != nil || (signupRes.StatusCode != http.StatusOK && signupRes.StatusCode != http.StatusCreated) {
		t.Fatalf("expected signup to succeed, got status %d, err %v", signupRes.StatusCode, signupErr)
	}

	projectReq := map[string]string{
		"name":        "E2E Test Project",
		"description": "Integration test project",
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

	res, body, err = h.get("/api/projects/"+projectID, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/projects/%s to return 200, got status %d", projectID, res.StatusCode)
	}

	res, body, err = h.get("/api/projects/"+projectID+"/environments", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/projects/%s/environments to return 200, got status %d", projectID, res.StatusCode)
	}

	envsData, _ := body["data"].([]any)
	if len(envsData) == 0 {
		t.Fatalf("expected default environment to exist for created project")
	}
	envObj, _ := envsData[0].(map[string]any)
	envID, _ := envObj["id"].(string)

	serviceReq := map[string]any{
		"projectId":   projectID,
		"name":        "e2e-api-service",
		"runtimeMode": "docker",
		"replicas":    1,
	}

	res, body, err = h.post("/api/environments/"+envID+"/apps", serviceReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create service to succeed, got status %d, body %v", res.StatusCode, body)
	}

	serviceData, _ := body["data"].(map[string]any)
	serviceID, _ := serviceData["id"].(string)
	if serviceID == "" {
		t.Fatalf("expected valid service id, got empty")
	}

	varReq := map[string]string{
		"key":   "DATABASE_URL",
		"value": "postgres://localhost:5432/testdb",
	}

	res, body, err = h.post("/api/services/"+serviceID+"/variables", varReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create service var to succeed, got status %d, body %v", res.StatusCode, body)
	}

	varData, _ := body["data"].(map[string]any)
	varID, _ := varData["id"].(string)
	if varID == "" {
		t.Fatalf("expected valid var id, got empty")
	}

	varUpdate := map[string]string{
		"key":   "DATABASE_URL",
		"value": "postgres://prodhost:5432/testdb",
	}
	res, body, err = h.put("/api/services/"+serviceID+"/variables/"+varID, varUpdate, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected update service var to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, body, err = h.put("/api/services/invalid-service-id/variables/"+varID, varUpdate, nil)
	if err != nil || (res.StatusCode != http.StatusBadRequest && res.StatusCode != http.StatusNotFound) {
		t.Fatalf("expected mismatched serviceId path validation to reject with 400/404, got status %d, err %v", res.StatusCode, err)
	}

	res, body, err = h.delete("/api/services/"+serviceID+"/variables/"+varID, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent) {
		t.Fatalf("expected delete service var to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, body, err = h.delete("/api/projects/"+projectID, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected delete project to succeed, got status %d, body %v", res.StatusCode, body)
	}
}
