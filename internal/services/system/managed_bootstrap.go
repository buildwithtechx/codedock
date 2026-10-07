package system

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"codedock.run/codedock/internal/engine/hetzner"
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
)

func serverLabelSelector(serverID string) string {
	return "codedock-server==" + serverID
}

func sshConfigFor(server *models.Server) ssh.Config {
	return ssh.Config{
		Host:     server.SSHHost,
		Port:     server.SSHPort,
		User:     server.SSHUser,
		Key:      server.SSHPrivateKey,
		Password: server.SSHPassword,
	}
}

func sshKeyName(serverID string) string {
	compact := strings.ReplaceAll(serverID, "-", "")
	if len(compact) > 12 {
		compact = compact[:12]
	}
	return "codedock-m-" + compact
}

func (s *managedService) bootstrapDockerHost(ctx context.Context, server *models.Server, progress func(string, string) error) error {
	if s.sshManager == nil {
		return nil
	}
	if err := progress("BOOTSTRAP", "Waiting for ssh on "+server.SSHHost); err != nil {
		return err
	}
	deadline := time.Now().Add(s.sshWait)
	for {
		err := s.sshManager.TestConnection(ctx, sshConfigFor(server))
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("ssh did not become reachable: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.pollInterval):
		}
	}
	if err := progress("BOOTSTRAP", "Waiting for docker on "+server.SSHHost); err != nil {
		return err
	}
	deadline = time.Now().Add(s.sshWait)
	for {
		ready, err := s.dockerReady(ctx, server)
		if err != nil {
			return err
		}
		if ready {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("docker did not become ready; inspect cloud-init on %s before retrying", server.SSHHost)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.pollInterval):
		}
	}
}

func (s *managedService) dockerReady(ctx context.Context, server *models.Server) (bool, error) {
	client, release, err := s.sshManager.GetClient(server)
	if err != nil {
		return false, nil
	}
	defer release()
	if _, err := client.RunCommand(ctx, "docker info --format '{{.ServerVersion}}'"); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *managedService) provisionSnapshot(ctx context.Context, orgID string, credential *models.ManagedCredential, serverID string) (string, error) {
	quota, err := s.quotaWithUsage(ctx, orgID)
	if err != nil {
		return "", err
	}
	parts := []string{credential.ID, credential.UpdatedAt, snapshotQuota(quota)}
	if serverID != "" {
		server, link, err := s.managedServer(ctx, orgID, serverID)
		if err == nil {
			parts = append(parts, snapshotServer(server), link.CredentialID)
		}
	}
	return snapshotValue(parts...), nil
}

func snapshotQuota(quota *models.ManagedQuota) string {
	return fmt.Sprintf("%d/%d/%d/%d", quota.MaxServers, quota.MaxMemoryGB, quota.UsedServers, quota.UsedMemoryGB)
}

func snapshotServer(server *models.Server) string {
	return fmt.Sprintf("%s/%s/%s/%s/%d", server.Status, server.ExternalID, server.ServerType, server.IPAddress, server.UpdatedAt.Unix())
}

func ensureProviderKey(ctx context.Context, provider ManagedProvider, name, publicKey string) error {
	if err := provider.CreateSSHKey(ctx, name, publicKey); err != nil {
		var apiErr *hetzner.APIError
		if errors.As(err, &apiErr) && apiErr.Code == "uniqueness_error" {
			return nil
		}
		return err
	}
	return nil
}
