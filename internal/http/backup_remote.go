package http

import (
	"codedock.run/codedock/internal/engine/backup"
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"log/slog"
)

func configureRemoteBackups(manager *backup.BackupManager, projects repositories.ProjectRepository, servers repositories.ServerRepository, ssh *ssh.SSHManager) {
	manager.SetDockerTarget(func(ctx context.Context, cfg *models.BackupConfig) (*client.Client, func(), error) {
		if cfg.ProjectID == "" {
			return nil, func() {}, nil
		}
		project, err := projects.Get(ctx, cfg.ProjectID)
		if err != nil {
			return nil, nil, err
		}
		if project.ServerID == "" {
			return nil, func() {}, nil
		}
		server, err := servers.GetByID(ctx, project.ServerID)
		if err != nil {
			return nil, nil, err
		}
		if server.IsLocal || server.OrganizationID != project.OrganizationID {
			return nil, nil, fmt.Errorf("backup target organization changed")
		}
		docker, release, err := ssh.GetDockerClient(ctx, server)
		if err != nil {
			return nil, nil, err
		}
		return docker, func() {
			if err := docker.Close(); err != nil {
				slog.Warn("close backup Docker transport", "error", err)
			}
			release()
		}, nil
	})
}
