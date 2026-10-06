package repositories

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func (r *BackupRepo) ControlPlaneSnapshot(ctx context.Context) (data []byte, err error) {
	directory, err := os.MkdirTemp("", "codedock-snapshot-")
	if err != nil {
		return nil, fmt.Errorf("create snapshot directory: %w", err)
	}
	defer func() {
		if cleanup := os.RemoveAll(directory); cleanup != nil && err == nil {
			err = fmt.Errorf("remove snapshot directory: %w", cleanup)
		}
	}()
	target := filepath.Join(directory, "control-plane.db")
	if _, err := r.db.ExecContext(ctx, "VACUUM INTO ?", target); err != nil {
		return nil, fmt.Errorf("snapshot control-plane database: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err = os.ReadFile(target)
	if err != nil {
		return nil, fmt.Errorf("read control-plane snapshot: %w", err)
	}
	return data, nil
}
