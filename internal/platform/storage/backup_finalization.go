package storage

import (
	"fmt"
	"time"
)

type FinalizeBackupOpts = BackupFinalization

func (bm *BackupManager) finalizeBackupRecord(opts FinalizeBackupOpts) (*BackupRecord, error) {
	nowStr := time.Now().UTC().Format(time.RFC3339)
	finalLogs := opts.Record.Logs + opts.ExecLogs + "\nBackup run completed successfully."
	if err := bm.store.UpdateBackupRecord(UpdateBackupRecordOpts{
		ID:     opts.Record.ID,
		Status: BackupRecordStatusCompleted,
		SHA256: opts.Record.SHA256, VerifiedAt: opts.Record.VerifiedAt,
		FilePath:       opts.FilePath,
		S3URL:          opts.S3URL,
		SFTPURL:        opts.SFTPURL,
		ParentRecordID: opts.ParentRecordID,
		Logs:           finalLogs,
		FileSizeBytes:  opts.SizeBytes,
		CompletedAt:    nowStr,
	}); err != nil {
		return nil, fmt.Errorf("persist completed backup record: %w", err)
	}

	opts.Record.Status = BackupRecordStatusCompleted
	opts.Record.FilePath = opts.FilePath
	opts.Record.FileSizeBytes = opts.SizeBytes
	opts.Record.S3URL = opts.S3URL
	opts.Record.SFTPURL = opts.SFTPURL
	opts.Record.ParentRecordID = opts.ParentRecordID
	opts.Record.Logs = finalLogs
	opts.Record.CompletedAt = nowStr
	return opts.Record, nil
}
