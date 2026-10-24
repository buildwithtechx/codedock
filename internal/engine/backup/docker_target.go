package backup

import (
	"codedock/internal/models"
	"context"
	"github.com/docker/docker/client"
)

type DockerTarget func(context.Context, *models.BackupConfig) (*client.Client, func(), error)

func (bm *BackupManager) SetDockerTarget(selector DockerTarget) { bm.dockerTarget = selector }
func (bm *BackupManager) producer(ctx context.Context, cfg *models.BackupConfig) (*BackupManager, func(), error) {
	if bm.dockerTarget == nil {
		return bm, func() {}, nil
	}
	client, release, err := bm.dockerTarget(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}
	if client == nil {
		return bm, release, nil
	}
	manager := NewBackupManager(client, bm.store, bm.backupDir)
	manager.volumeOperations = bm.volumeOperations
	return manager, release, nil
}
