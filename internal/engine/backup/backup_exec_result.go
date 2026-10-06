package backup

import (
	"context"
	"fmt"
	"time"
)

func (bm *BackupManager) waitExecSuccess(ctx context.Context, execID, operation string) error {
	for {
		result, err := bm.dockerClient.ContainerExecInspect(ctx, execID)
		if err != nil {
			return fmt.Errorf("inspect %s command result: %w", operation, err)
		}
		if !result.Running {
			if result.ExitCode != 0 {
				return fmt.Errorf("%s command exited with status %d", operation, result.ExitCode)
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for %s command: %w", operation, ctx.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
}
