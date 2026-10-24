package backup

import (
	"codedock/internal/models"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/mount"
)

func (bm *BackupManager) validateVolumeProducer(ctx context.Context, cfg *models.BackupConfig) error {
	var id string
	var err error
	if owner, ok := bm.store.(interface {
		VolumeBackupOwner(context.Context, *models.BackupConfig) (string, error)
	}); ok {
		id, err = owner.VolumeBackupOwner(ctx, cfg)
	} else if owner, ok := bm.store.(VolumeRestoreOwner); ok {
		id, err = owner.VolumeRestoreOwner(ctx, cfg)
	} else {
		return fmt.Errorf("volume producer ownership validation unavailable")
	}
	if err != nil {
		return err
	}

	if bm.dockerClient == nil {
		return fmt.Errorf("volume producer requires Docker")
	}
	volume, err := bm.dockerClient.VolumeInspect(ctx, cfg.VolumeName)
	if err != nil {
		return err
	}
	if volume.Name != cfg.VolumeName {
		return fmt.Errorf("volume producer identity changed")
	}
	inspected, err := bm.dockerClient.ContainerInspect(ctx, id)
	if err != nil {
		return err
	}
	for _, binding := range inspected.Mounts {
		if binding.Type == mount.TypeVolume && binding.Name == cfg.VolumeName {
			return nil
		}
	}
	return fmt.Errorf("named volume is not mounted by its registered owner")
}
