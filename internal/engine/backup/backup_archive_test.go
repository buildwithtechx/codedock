package backup

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"codedock.run/codedock/internal/models"
)

func TestArchiveDownloadUsesRecordedDestinationAndKey(t *testing.T) {
	var requested string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = r.URL.Path
		if _, err := w.Write([]byte("original archive")); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	store := newMockStore()
	store.s3dests["original"] = &models.S3Destination{ID: "original", Endpoint: server.URL, Bucket: "archives", PathPrefix: "new-prefix"}
	store.configs["config"] = &models.BackupConfig{ID: "config", S3DestinationID: "changed"}
	store.records["record"] = &models.BackupRecord{ID: "record", BackupConfigID: "config", S3DestinationID: "original", S3URL: "s3://archives/old-prefix/backup.sql", FilePath: filepath.Join(t.TempDir(), "missing.sql")}
	manager := NewBackupManager(nil, store, t.TempDir())
	reader, name, err := manager.OpenBackupArchive(context.Background(), "record")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	if string(data) != "original archive" || name != "backup.sql" || requested != "/archives/old-prefix/backup.sql" {
		t.Fatalf("wrong archive: %s %s %s", data, name, requested)
	}
	store.records["record"].S3URL = "s3://other-bucket/backup.sql"
	if _, _, err := manager.OpenBackupArchive(context.Background(), "record"); err == nil {
		t.Fatal("mismatched bucket accepted")
	}
}

func TestArchiveDownloadPrefersLocalArchive(t *testing.T) {
	store := newMockStore()
	filename := filepath.Join(t.TempDir(), "backup.sql")
	if err := os.WriteFile(filename, []byte("local archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	store.records["record"] = &models.BackupRecord{ID: "record", FilePath: filename, S3URL: "s3://unavailable/backup.sql"}
	manager := NewBackupManager(nil, store, t.TempDir())
	reader, name, err := manager.OpenBackupArchive(context.Background(), "record")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	if name != "backup.sql" || string(data) != "local archive" {
		t.Fatal("wrong local archive")
	}
	if _, _, err := manager.OpenBackupArchive(context.Background(), "missing"); err == nil {
		t.Fatal("missing archive accepted")
	}
}
