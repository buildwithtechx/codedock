package system

import (
	"context"
	"fmt"

	"codedock/internal/engine/dockerprobe"
	"codedock/internal/models"
)

func (s *serverService) InspectSwarm(ctx context.Context, id, userID string) (*models.SwarmInspection, error) {
	server, err := s.GetServer(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, fmt.Errorf("server not found")
	}

	if server.IsLocal {
		return dockerprobe.InspectSwarm(ctx, dockerprobe.LocalRunner{})
	}

	if s.sshManager == nil {
		return nil, fmt.Errorf("ssh manager is not initialized")
	}

	client, release, err := s.sshManager.GetClient(server)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server via SSH: %w", err)
	}
	defer release()

	return dockerprobe.InspectSwarm(ctx, client)
}
