package backup

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"codedock/internal/models"
)

func (bm *BackupManager) uploadToS3(ctx context.Context, dest *models.S3Destination, fileName string, data []byte) (string, error) {
	key := strings.Trim(dest.PathPrefix, "/")
	if key != "" {
		key += "/"
	}
	key += fileName
	resp, err := signedS3Request(ctx, dest, "PUT", key, data, "application/octet-stream")
	if err != nil {
		return "", err
	}
	if err := resp.Body.Close(); err != nil {
		slog.Warn("close upload response", "error", err)
	}
	return fmt.Sprintf("s3://%s/%s", dest.Bucket, key), nil
}

func (bm *BackupManager) enforceRetentionPolicy(cfg *models.BackupConfig) {
	records, err := bm.store.ListBackupRecords(cfg.ID)
	if err != nil || len(records) == 0 {
		return
	}

	var activeRecords []*models.BackupRecord
	for _, rec := range records {
		if rec.Status == models.BackupRecordStatusCompleted && rec.ProtectedUntil <= time.Now().Unix() {
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
		if err := bm.expireRecord(cfg, rec); err != nil {
			slog.Error("expire backup record", "record", rec.ID, "error", err)
		}
	}
}
