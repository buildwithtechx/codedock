package backup

import (
	"context"
	"fmt"
	"time"
)

func (bm *BackupManager) waitExecSuccess(ctx context.Context, execID string) error {
	for {
		result, err := bm.dockerClient.ContainerExecInspect(ctx, execID)
		if err != nil {
			return fmt.Errorf("inspect backup command result: %w", err)
		}
		if !result.Running {
			if result.ExitCode != 0 {
				return fmt.Errorf("backup command exited with status %d", result.ExitCode)
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for backup command: %w", ctx.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}
