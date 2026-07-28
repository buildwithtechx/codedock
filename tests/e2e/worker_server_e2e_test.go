package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestE2EWorkerServerAndS3DestinationManagement(t *testing.T) {
	h := newE2EHarness(t)
	defer h.Close()

	h.post("/api/auth/signup", map[string]string{
		"email":    "sys_admin@codedock.local",
		"password": "Password123!",
		"name":     "Sys Admin",
	}, nil)

	serverReq := map[string]string{
		"name":      "worker-us-east-1",
		"ipAddress": "192.168.1.100",
	}

	res, body, err := h.post("/api/servers", serverReq, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create server node to succeed, got status %d, body %v", res.StatusCode, body)
	}

	srvData, _ := body["data"].(map[string]any)
	srvID, _ := srvData["id"].(string)
	if srvID == "" {
		t.Fatalf("expected valid server id, got empty")
	}

	res, body, err = h.get("/api/servers", nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /api/servers to return 200, got status %d", res.StatusCode)
	}

	mockS3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockS3Server.Close()

	s3Req := map[string]string{
		"name":            "Production MinIO Bucket",
		"provider":        "minio",
		"endpoint":        mockS3Server.URL,
		"bucket":          "codedock-backups",
		"region":          "us-east-1",
		"accessKeyId":     "minioadmin",
		"secretAccessKey": "minioSecret123!",
	}

	res, body, err = h.post("/api/s3-destinations", s3Req, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated) {
		t.Fatalf("expected create s3 destination to succeed, got status %d, body %v", res.StatusCode, body)
	}

	s3Data, _ := body["data"].(map[string]any)
	s3ID, _ := s3Data["id"].(string)
	if s3ID == "" {
		t.Fatalf("expected valid s3 destination id, got empty")
	}

	s3Update := map[string]string{
		"name":            "Updated MinIO Bucket",
		"provider":        "minio",
		"endpoint":        mockS3Server.URL,
		"bucket":          "codedock-backups-v2",
		"region":          "us-east-1",
		"accessKeyId":     "minioadmin",
		"secretAccessKey": "newMinioSecret456!",
	}

	res, body, err = h.put("/api/s3-destinations/"+s3ID, s3Update, nil)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected update s3 destination to succeed, got status %d, body %v", res.StatusCode, body)
	}

	res, body, err = h.delete("/api/s3-destinations/"+s3ID, nil)
	if err != nil || (res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent) {
		t.Fatalf("expected delete s3 destination to succeed, got status %d, body %v", res.StatusCode, body)
	}
}
