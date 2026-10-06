package backup

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

type progressKey struct{}

func WithProgress(ctx context.Context, report func(string, string) error) context.Context {
	return context.WithValue(ctx, progressKey{}, report)
}
func backupProgress(ctx context.Context, phase, message string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if report, ok := ctx.Value(progressKey{}).(func(string, string) error); ok {
		return report(phase, message)
	}
	return nil
}
func (bm *BackupManager) VerifyArchive(ctx context.Context, id string) (string, error) {
	record, err := bm.store.GetBackupRecord(id)
	if err != nil {
		return "", err
	}
	if record.Status != "completed" {
		return "", fmt.Errorf("archive is not completed")
	}
	stream, _, err := bm.OpenBackupArchive(ctx, id)
	if err != nil {
		return "", err
	}
	digest := sha256.New()
	count, readErr := io.Copy(digest, io.LimitReader(stream, 1024*1024*1024+1))
	closeErr := stream.Close()
	if readErr != nil {
		return "", fmt.Errorf("verify archive: %w", readErr)
	}
	if closeErr != nil {
		return "", closeErr
	}
	if count == 0 || count > 1024*1024*1024 {
		return "", fmt.Errorf("archive must be nonempty and at most 1 GiB")
	}
	actual := hex.EncodeToString(digest.Sum(nil))
	if record.SHA256 != "" && record.SHA256 != actual {
		return "", fmt.Errorf("archive checksum mismatch")
	}
	if record.FileSizeBytes > 0 && record.FileSizeBytes != count {
		return "", fmt.Errorf("archive size changed")
	}
	if err := bm.store.UpdateBackupRecord(modelsVerification(record, actual)); err != nil {
		return "", err
	}
	return actual, nil
}

func modelsVerification(record *models.BackupRecord, actual string) models.UpdateBackupRecordOpts {
	return models.UpdateBackupRecordOpts{ID: record.ID, Status: record.Status, FilePath: record.FilePath, S3URL: record.S3URL, S3DestinationID: record.S3DestinationID, FileSizeBytes: record.FileSizeBytes, Logs: record.Logs, CompletedAt: record.CompletedAt, SHA256: actual, VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
}
