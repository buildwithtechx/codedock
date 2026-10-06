package backup

import (
	"codedock.run/codedock/internal/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

func (bm *BackupManager) expireRecord(cfg *models.BackupConfig, record *models.BackupRecord) (resultErr error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if store, ok := bm.store.(interface {
		ClaimRecordExpiry(context.Context, string) (bool, error)
	}); ok {
		claimed, err := store.ClaimRecordExpiry(ctx, record.ID)
		if err != nil {
			return err
		}
		if !claimed {
			return nil
		}
		defer func() {
			if resultErr != nil {
				opts := modelsVerification(record, record.SHA256)
				opts.Status = models.BackupRecordStatusFailed
				opts.Logs += "\nArchive cleanup was interrupted or failed; inspect storage before retrying."
				resultErr = errors.Join(resultErr, bm.store.UpdateBackupRecord(opts))
			}
		}()
	}
	return bm.removeArchive(ctx, cfg, record)
}

func (bm *BackupManager) removeArchive(ctx context.Context, cfg *models.BackupConfig, record *models.BackupRecord) error {
	if record.S3URL != "" {
		id := record.S3DestinationID
		if id == "" {
			id = cfg.S3DestinationID
		}
		if id == "" {
			return fmt.Errorf("retained archive destination is unavailable")
		}
		destination, err := bm.store.GetS3Destination(id)
		if err != nil {
			return err
		}
		prefix := "s3://" + destination.Bucket + "/"
		if !strings.HasPrefix(record.S3URL, prefix) {
			return fmt.Errorf("archive destination mismatch")
		}
		response, err := signedS3Request(ctx, destination, "DELETE", strings.TrimPrefix(record.S3URL, prefix), nil, "")
		if err != nil {
			return err
		}
		if err := response.Body.Close(); err != nil {
			slog.Warn("close archive deletion response", "error", err)
		}
	}
	if record.FilePath != "" {
		if err := os.Remove(record.FilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	opts := modelsVerification(record, record.SHA256)
	opts.Status = models.BackupRecordStatusExpired
	opts.Logs += "\nArchive removed from storage."
	opts.FilePath = ""
	opts.S3URL = ""
	return bm.store.UpdateBackupRecord(opts)
}
