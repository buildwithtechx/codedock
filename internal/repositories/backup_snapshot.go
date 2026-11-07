package repositories

import (
	"context"
	"fmt"
)

func (r *BackupRepo) SetSnapshotDumper(dumper func(ctx context.Context) ([]byte, error)) {
	r.snapshotDumper = dumper
}

func (r *BackupRepo) ControlPlaneSnapshot(ctx context.Context) (data []byte, err error) {
	if r.snapshotDumper == nil {
		return nil, fmt.Errorf("control-plane snapshotter is not configured")
	}
	return r.snapshotDumper(ctx)
}
