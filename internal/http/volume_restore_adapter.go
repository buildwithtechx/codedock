package http

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"errors"
	"fmt"
)

func (a *engineAdapter) VolumeRestoreOwner(ctx context.Context, cfg *models.BackupConfig) (string, error) {
	var projectID, containerID string
	if cfg.ServiceID != "" && cfg.DatabaseID != "" {
		return "", errors.New("volume backup must have exactly one owner")
	}
	if cfg.ServiceID != "" {
		app, err := a.appRepo.GetByID(ctx, cfg.ServiceID)
		if err != nil {
			return "", fmt.Errorf("load volume service: %w", err)
		}
		projectID = app.ProjectID
		containerID = utils.NormalizeContainerName(app.ID)
		if app.ContainerID != "" {
			containerID = app.ContainerID
		}
	} else if cfg.DatabaseID != "" && cfg.DatabaseID != "global" {
		db, err := a.dbRepo.GetByID(ctx, cfg.DatabaseID)
		if err != nil {
			return "", fmt.Errorf("load volume database: %w", err)
		}

		projectID = db.ProjectID
		containerID = utils.NormalizeContainerName(db.ID)
		if cfg.VolumeName != "codedock-db-data-"+db.ID {
			return "", errors.New("volume is not the database data volume")
		}
	} else {
		return "", errors.New("volume backup requires a registered service or database owner")
	}
	project, err := a.projectRepo.Get(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("load volume project: %w", err)
	}
	if project.ServerID != "" {
		return "", errors.New("volume restore supports local Docker targets only")
	}
	return containerID, nil
}
