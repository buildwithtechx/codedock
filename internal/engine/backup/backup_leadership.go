package backup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"codedock/internal/models"
)

func (bm *BackupManager) Reconcile() error {
	cfgs, err := bm.store.ListAllActiveBackupConfigs()
	if err != nil {
		return fmt.Errorf("failed to load active backup configs during reconcile: %w", err)
	}
	bm.mu.Lock()
	defer bm.mu.Unlock()
	active := make(map[string]*models.BackupConfig, len(cfgs))
	for _, cfg := range cfgs {
		active[cfg.ID] = cfg
		if err := bm.registerBackupLocked(cfg); err != nil {
			slog.Warn("failed to schedule backup", "name", cfg.Name, "id", cfg.ID, "err", err)
		}
	}
	for id, entryID := range bm.entries {
		if _, ok := active[id]; !ok {
			bm.cronEngine.Remove(entryID)
			delete(bm.entries, id)
		}
	}
	return nil
}

func (bm *BackupManager) ServeElected(ctx context.Context) {
	if err := bm.Reconcile(); err != nil {
		slog.Warn("backup reconcile failed", "error", err)
	}
	bm.cronEngine.Start()
	defer bm.cronEngine.Stop()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := bm.Reconcile(); err != nil {
				slog.Warn("backup reconcile failed", "error", err)
			}
		}
	}
}
