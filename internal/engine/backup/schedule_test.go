package backup

import (
	"codedock.run/codedock/internal/models"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManualBackupRemovesRecurringSchedule(t *testing.T) {
	manager := NewBackupManager(nil, nil, t.TempDir())
	config := &models.BackupConfig{ID: "backup", Schedule: "0 2 * * *", Status: models.BackupConfigStatusActive, BackupEnabled: true}
	if err := manager.RegisterBackup(config); err != nil {
		t.Fatal(err)
	}
	if len(manager.entries) != 1 {
		t.Fatal("recurring backup was not scheduled")
	}
	config.Schedule = "manual"
	if err := ValidateSchedule(config.Schedule); err != nil {
		t.Fatal(err)
	}
	if err := manager.RegisterBackup(config); err != nil {
		t.Fatal(err)
	}
	if len(manager.entries) != 0 {
		t.Fatal("manual backup has a recurring schedule")
	}
	if err := ValidateSchedule("invalid"); err == nil {
		t.Fatal("invalid schedule was accepted")
	}
}

func TestBucketVerificationDoesNotCreateMissingBucket(t *testing.T) {
	methods := make(chan string, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods <- r.Method
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	originalTransport := s3HTTPClient.Transport
	s3HTTPClient.Transport = server.Client().Transport
	t.Cleanup(func() { s3HTTPClient.Transport = originalTransport })
	err := CheckS3Bucket(context.Background(), &models.S3Destination{Endpoint: server.URL, Bucket: "missing", Region: "us-east-1"})
	if err == nil {
		t.Fatal("missing bucket verified successfully")
	}
	if len(methods) != 1 {
		t.Fatalf("expected one verification request, got %d", len(methods))
	}
	if method := <-methods; method != http.MethodHead {
		t.Fatalf("verification mutated bucket with %s", method)
	}
}
