package backup

import (
	"codedock/internal/models"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type protectedStore struct {
	Store
	record *models.BackupRecord
	claims int
}

func (s *protectedStore) ListBackupRecords(string) ([]*models.BackupRecord, error) {
	return []*models.BackupRecord{s.record}, nil
}
func (s *protectedStore) ClaimRecordExpiry(context.Context, string) (bool, error) {
	s.claims++
	return false, nil
}
func TestRetentionDoesNotDeleteWhenProtectionWinsClaim(t *testing.T) {
	path := filepath.Join(t.TempDir(), "backup")
	if err := os.WriteFile(path, []byte("archive"), 0600); err != nil {
		t.Fatal(err)
	}
	store := &protectedStore{record: &models.BackupRecord{ID: "record", Status: models.BackupRecordStatusCompleted, FilePath: path, StartedAt: time.Now().Add(-72 * time.Hour).Format(time.RFC3339)}}
	manager := NewBackupManager(nil, store, t.TempDir())
	manager.enforceRetentionPolicy(&models.BackupConfig{ID: "policy", RetentionDays: 1})
	if store.claims != 1 {
		t.Fatal("retention did not atomically claim the record")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("protected archive was removed", err)
	}
}
