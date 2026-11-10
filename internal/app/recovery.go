package app

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"codedock/internal/deploy"
)

func recoverRollouts(deployer *deploy.Deployer, dataDir string) error {
	if deployer == nil || dataDir == "" {
		return nil
	}
	if err := deployer.SetRolloutDirectory(filepath.Join(dataDir, "rollouts")); err != nil {
		return err
	}
	recoveryCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := deployer.RecoverRollouts(recoveryCtx); err != nil {
		return fmt.Errorf("recover interrupted rollout: %w", err)
	}
	return nil
}

func runRecovery(ctx context.Context, s *appServices) error {
	if err := s.deployRepo.RecoverInterrupted(ctx); err != nil {
		return fmt.Errorf("recover deployments: %w", err)
	}
	if err := s.runtimeRepo.SyncKinds(ctx); err != nil {
		return fmt.Errorf("recover runtime kinds: %w", err)
	}
	if err := s.operationService.Recover(ctx); err != nil {
		return fmt.Errorf("recover operations: %w", err)
	}
	if err := s.backupRepo.RecoverRecords(ctx); err != nil {
		return fmt.Errorf("recover backup records: %w", err)
	}
	return nil
}
