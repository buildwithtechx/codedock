package backup

import (
	"codedock.run/codedock/internal/models"
	"fmt"
	"time"
)

type FinalizeBackupOpts = models.BackupFinalization

func (bm *BackupManager) finalizeBackupRecord(opts FinalizeBackupOpts) (*models.BackupRecord, error) {
	nowStr := time.Now().UTC().Format(time.RFC3339)
	finalLogs := opts.Record.Logs + opts.ExecLogs + "\nBackup run completed successfully."
	if err := bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
		ID:     opts.Record.ID,
		Status: models.BackupRecordStatusCompleted,
		SHA256: opts.Record.SHA256, VerifiedAt: opts.Record.VerifiedAt,
		FilePath:      opts.FilePath,
		S3URL:         opts.S3URL,
		Logs:          finalLogs,
		FileSizeBytes: opts.SizeBytes,
		CompletedAt:   nowStr,
	}); err != nil {
		return nil, fmt.Errorf("persist completed backup record: %w", err)
	}

	opts.Record.Status = models.BackupRecordStatusCompleted
	opts.Record.FilePath = opts.FilePath
	opts.Record.FileSizeBytes = opts.SizeBytes
	opts.Record.S3URL = opts.S3URL
	opts.Record.Logs = finalLogs
	opts.Record.CompletedAt = nowStr
	return opts.Record, nil
}
