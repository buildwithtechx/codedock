package backups

import (
	"context"
	"errors"
	"io"
)

func (s *BackupService) OpenArchive(ctx context.Context, recordID string) (io.ReadCloser, string, error) {
	if s.manager == nil {
		return nil, "", errors.New("backup runtime is unavailable")
	}
	return s.manager.OpenBackupArchive(ctx, recordID)
}
