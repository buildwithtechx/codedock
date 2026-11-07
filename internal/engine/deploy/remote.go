package deploy

import (
	"codedock/internal/engine/build"
	"context"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"path/filepath"
)

func (d *Deployer) ForDockerHost(ctx context.Context, client *client.Client, serverID string) (*Deployer, error) {
	if _, err := uuid.Parse(serverID); err != nil {
		return nil, err
	}
	remote := *d
	remote.builder = build.NewBuilder(client)
	remote.containerManager = NewContainerManager(client, d.store)
	remote.containerManager.volumes = d.containerManager.volumes
	if d.rolloutDirectory != "" {
		if err := remote.SetRolloutDirectory(filepath.Join(d.rolloutDirectory, serverID)); err != nil {
			return nil, err
		}
	}
	if err := remote.RecoverRollouts(ctx); err != nil {
		return nil, err
	}
	return &remote, nil
}
