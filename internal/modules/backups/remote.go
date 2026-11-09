package backups

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/docker/docker/client"

	projectsmod "codedock/internal/modules/projects"
	serversmod "codedock/internal/modules/servers"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
)

type ProjectBackupTargetLookup interface {
	Get(ctx context.Context, id string) (*projectsmod.ProjectConfig, error)
}

type ServerBackupTargetLookup interface {
	GetByID(ctx context.Context, id string) (*serversmod.Server, error)
}

func ConfigureRemoteBackups(manager *backup.BackupManager, projects ProjectBackupTargetLookup, servers ServerBackupTargetLookup, sshMgr *ssh.SSHManager) {
	manager.SetDockerTarget(func(ctx context.Context, cfg *backup.BackupConfig) (*client.Client, func(), error) {
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
		sshServer := &ssh.Server{
			ID:             server.ID,
			OrganizationID: server.OrganizationID,
			UserID:         server.UserID,
			Name:           server.Name,
			IPAddress:      server.IPAddress,
			IsLocal:        server.IsLocal,
			SSHHost:        server.SSHHost,
			SSHPort:        server.SSHPort,
			SSHUser:        server.SSHUser,
			SSHAuthMethod:  server.SSHAuthMethod,
			SSHKey:         server.SSHKey,
			SSHPrivateKey:  server.SSHPrivateKey,
			SSHPassword:    server.SSHPassword,
			SSHTransport:   server.SSHTransport,
			SSHJumpHost:    server.SSHJumpHost,
			Status:         ssh.ServerStatus(server.Status),
		}
		docker, release, err := sshMgr.GetDockerClient(ctx, sshServer)
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
