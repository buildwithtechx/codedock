package http

import (
	"context"
	"fmt"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/engine/systemdb"
)

func newSystemSnapshotDumper(dataDir string) func(ctx context.Context) ([]byte, error) {
	return func(ctx context.Context) ([]byte, error) {
		if databaseURL := config.Get().Database.URL; databaseURL != "" {
			return systemdb.DumpURL(ctx, databaseURL)
		}
		password, err := systemdb.ReadPassword(dataDir)
		if err != nil {
			return nil, err
		}
		return systemdb.DumpContainer(ctx, systemdb.ContainerName, systemdb.DBUser, systemdb.DBName, password)
	}
}

func (a *engineAdapter) ControlPlaneSnapshot(ctx context.Context) ([]byte, error) {
	store, ok := a.backupRepo.(interface {
		ControlPlaneSnapshot(context.Context) ([]byte, error)
	})
	if !ok {
		return nil, fmt.Errorf("consistent control-plane snapshots unavailable")
	}
	return store.ControlPlaneSnapshot(ctx)
}
