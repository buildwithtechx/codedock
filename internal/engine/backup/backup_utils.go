package backup

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"codedock.run/codedock/internal/models"
)

func (bm *BackupManager) uploadToS3(ctx context.Context, dest *models.S3Destination, fileName string, data []byte) (string, error) {
	resp, err := signedS3Request(ctx, dest, "PUT", fileName, data, "application/octet-stream")
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return fmt.Sprintf("s3://%s/%s", dest.Bucket, fileName), nil
}

func (bm *BackupManager) enforceRetentionPolicy(cfg *models.BackupConfig) {
	records, err := bm.store.ListBackupRecords(cfg.ID)
	if err != nil || len(records) == 0 {
		return
	}

	var activeRecords []*models.BackupRecord
	for _, rec := range records {
		if rec.Status == models.BackupRecordStatusCompleted {
			activeRecords = append(activeRecords, rec)
		}
	}
	if len(activeRecords) == 0 {
		return
	}

	toExpire := make(map[string]*models.BackupRecord)

	if cfg.RetentionDays > 0 {
		cutoff := time.Now().Add(-time.Duration(cfg.RetentionDays) * 24 * time.Hour)
		for _, rec := range activeRecords {
			started, err := time.Parse(time.RFC3339, rec.StartedAt)
			if err == nil && started.Before(cutoff) {
				toExpire[rec.ID] = rec
			}
		}
	}

	var validRecords []*models.BackupRecord
	for _, rec := range activeRecords {
		if _, expired := toExpire[rec.ID]; !expired {
			validRecords = append(validRecords, rec)
		}
	}

	sort.Slice(validRecords, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, validRecords[i].StartedAt)
		t2, _ := time.Parse(time.RFC3339, validRecords[j].StartedAt)
		return t1.After(t2)
	})

	if cfg.MaxBackups > 0 && len(validRecords) > cfg.MaxBackups {
		for i := cfg.MaxBackups; i < len(validRecords); i++ {
			toExpire[validRecords[i].ID] = validRecords[i]
		}
	}

	if cfg.MaxStorageGB > 0 {
		maxGB := cfg.MaxStorageGB
		if maxGB > 8500000000 {
			maxGB = 8500000000
		}
		maxBytes := int64(maxGB) * 1024 * 1024 * 1024
		var currentBytes int64
		for _, rec := range validRecords {
			if _, expired := toExpire[rec.ID]; expired {
				continue
			}
			if currentBytes+rec.FileSizeBytes > maxBytes {
				toExpire[rec.ID] = rec
			} else {
				currentBytes += rec.FileSizeBytes
			}
		}
	}

	for _, rec := range toExpire {
		if rec.FilePath != "" {
			_ = os.Remove(rec.FilePath)
		}
		if rec.S3URL != "" && cfg.S3DestinationID != "" {
			if dest, err := bm.store.GetS3Destination(cfg.S3DestinationID); err == nil && dest != nil {
				prefix := fmt.Sprintf("s3://%s/", dest.Bucket)
				key := strings.TrimPrefix(rec.S3URL, prefix)
				if resp, err := signedS3Request(context.Background(), dest, "DELETE", key, nil, ""); err == nil {
					resp.Body.Close()
				}
			}
		}
		_ = bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
			ID:          rec.ID,
			Status:      models.BackupRecordStatusExpired,
			FilePath:    "",
			S3URL:       rec.S3URL,
			Logs:        rec.Logs + "\nFile pruned by retention policy.",
			CompletedAt: time.Now().UTC().Format(time.RFC3339),
		})
	}
}
