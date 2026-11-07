package backup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func (bm *BackupManager) OpenBackupArchive(ctx context.Context, recordID string) (io.ReadCloser, string, error) {
	rec, err := bm.store.GetBackupRecord(recordID)
	if err != nil {
		return nil, "", fmt.Errorf("load backup record: %w", err)
	}
	if rec == nil {
		return nil, "", utils.NewNotFoundError("Backup record", recordID)
	}
	if rec.FilePath != "" {
		file, err := os.Open(rec.FilePath)
		if err == nil {
			return file, filepath.Base(rec.FilePath), nil
		}
		if !os.IsNotExist(err) {
			return nil, "", fmt.Errorf("open local backup archive: %w", err)
		}
	}
	if rec.S3URL == "" && rec.SFTPURL == "" {
		return nil, "", utils.NewNotFoundError("Backup archive", recordID)
	}
	if rec.S3URL != "" {
		destinationID := rec.S3DestinationID
		if destinationID == "" {
			cfg, err := bm.store.GetBackupConfig(rec.BackupConfigID)
			if err != nil {
				return nil, "", fmt.Errorf("load backup configuration: %w", err)
			}
			if cfg == nil {
				return nil, "", utils.NewNotFoundError("Backup configuration", rec.BackupConfigID)
			}
			destinationID = cfg.S3DestinationID
		}
		dest, err := bm.store.GetS3Destination(destinationID)
		if err != nil {
			return nil, "", fmt.Errorf("load backup destination: %w", err)
		}
		if dest == nil {
			return nil, "", utils.NewNotFoundError("Backup destination", destinationID)
		}
		prefix := "s3://" + dest.Bucket + "/"
		if !strings.HasPrefix(rec.S3URL, prefix) {
			return nil, "", fmt.Errorf("backup archive does not belong to its destination bucket")
		}
		key := strings.TrimPrefix(rec.S3URL, prefix)
		if key == "" {
			return nil, "", fmt.Errorf("backup archive object key is empty")
		}
		response, err := signedS3Request(ctx, dest, "GET", key, nil, "")
		if err != nil {
			return nil, "", fmt.Errorf("download backup archive: %w", err)
		}
		return response.Body, path.Base(key), nil
	}
	getter, ok := bm.store.(interface {
		GetSFTPDestination(string) (*models.SFTPDestination, error)
	})
	if !ok {
		return nil, "", fmt.Errorf("SFTP destination lookup unavailable")
	}
	destinationID := rec.SFTPDestinationID
	if destinationID == "" {
		cfg, err := bm.store.GetBackupConfig(rec.BackupConfigID)
		if err != nil {
			return nil, "", fmt.Errorf("load backup configuration: %w", err)
		}
		if cfg == nil {
			return nil, "", utils.NewNotFoundError("Backup configuration", rec.BackupConfigID)
		}
		destinationID = cfg.SFTPDestinationID
	}
	dest, err := getter.GetSFTPDestination(destinationID)
	if err != nil || dest == nil {
		return nil, "", fmt.Errorf("load backup destination: %w", err)
	}
	prefix := "sftp://" + dest.Host + "/"
	if !strings.HasPrefix(rec.SFTPURL, prefix) {
		return nil, "", fmt.Errorf("backup archive does not belong to its destination host")
	}
	key := strings.TrimPrefix(rec.SFTPURL, prefix)
	if key == "" {
		return nil, "", fmt.Errorf("backup archive object key is empty")
	}
	stream, err := sftpGet(ctx, dest, key)
	if err != nil {
		return nil, "", fmt.Errorf("download backup archive: %w", err)
	}
	return stream, path.Base(key), nil
}
