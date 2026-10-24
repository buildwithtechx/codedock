package ssh

import (
	"context"
	"fmt"
	"sync"
	"time"

	dockerclient "github.com/docker/docker/client"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type clientEntry struct {
	client *Client
	refs   int
	closed bool
}

type SSHManager struct {
	mu                  sync.Mutex
	clients             map[string]*clientEntry
	serverRepo          repositories.ServerRepository
	OnServerUnreachable func(server *models.Server, err error)
}

func NewSSHManager(repo repositories.ServerRepository) *SSHManager {
	return &SSHManager{
		clients:    make(map[string]*clientEntry),
		serverRepo: repo,
	}
}

func (m *SSHManager) GetClient(server *models.Server) (*Client, func(), error) {
	if server == nil {
		return nil, nil, fmt.Errorf("server cannot be nil")
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, ok := m.clients[server.ID]; ok && !entry.closed {
		entry.refs++
		return entry.client, func() { m.releaseClient(entry) }, nil
	}

	host := server.SSHHost
	if host == "" {
		host = server.IPAddress
	}

	key := server.SSHPrivateKey
	if key == "" {
		key = server.SSHKey
	}
	client, err := NewClient(Config{
		Host:     host,
		Port:     server.SSHPort,
		User:     server.SSHUser,
		Key:      key,
		Password: server.SSHPassword,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed creating ssh client for server %s: %w", server.ID, err)
	}

	entry := &clientEntry{
		client: client,
		refs:   1,
	}
	m.clients[server.ID] = entry
	return client, func() { m.releaseClient(entry) }, nil
}

func (m *SSHManager) releaseClient(entry *clientEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	entry.refs--
	if entry.closed && entry.refs == 0 {
		_ = entry.client.Close()
	}
}

func (m *SSHManager) GetDockerClient(ctx context.Context, server *models.Server) (*dockerclient.Client, func(), error) {
	client, release, err := m.GetClient(server)
	if err != nil {
		return nil, nil, err
	}

	dc, err := client.DockerClient(ctx)
	if err != nil {
		release()
		return nil, nil, err
	}

	return dc, release, nil
}

func (m *SSHManager) TestConnection(ctx context.Context, cfg Config) error {
	client, err := NewClient(cfg)
	if err != nil {
		return fmt.Errorf("ssh connection test failed: %w", err)
	}
	defer client.Close()

	out, err := client.RunCommand(ctx, "docker --version")
	if err != nil {
		return fmt.Errorf("remote docker check failed: %w", err)
	}

	if out == "" {
		return fmt.Errorf("docker not responding on remote host")
	}

	return nil
}

func (m *SSHManager) RefreshServerStatus(ctx context.Context, server *models.Server) error {
	client, release, err := m.GetClient(server)
	if err != nil {
		if server.Status != models.ServerStatusOffline && m.OnServerUnreachable != nil {
			m.OnServerUnreachable(server, err)
		}
		_ = m.serverRepo.UpdateStatus(ctx, server.ID, models.ServerStatusOffline)
		return err
	}
	defer release()

	metrics, err := client.CollectMetrics(ctx)
	if err != nil {
		m.RemoveClient(server.ID)
		if server.Status != models.ServerStatusOffline && m.OnServerUnreachable != nil {
			m.OnServerUnreachable(server, err)
		}
		_ = m.serverRepo.UpdateStatus(ctx, server.ID, models.ServerStatusOffline)
		return fmt.Errorf("failed to collect remote metrics: %w", err)
	}

	if err := m.serverRepo.UpdateMetrics(ctx, server.ID, metrics); err != nil {
		return fmt.Errorf("failed to update metrics: %w", err)
	}

	return m.serverRepo.UpdateStatus(ctx, server.ID, models.ServerStatusOnline)
}

func (m *SSHManager) RemoveClient(serverID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, ok := m.clients[serverID]; ok {
		entry.closed = true
		delete(m.clients, serverID)
		if entry.refs == 0 {
			_ = entry.client.Close()
		}
	}
}

func (m *SSHManager) StartMetricsPoller(ctx context.Context, userID string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			servers, err := m.serverRepo.ListByUser(ctx, userID, "")
			if err != nil {
				continue
			}
			for _, s := range servers {
				if s.SSHHost != "" || s.IPAddress != "" {
					_ = m.RefreshServerStatus(ctx, s)
				}
			}
		}
	}
}

// AssertOperatorSSHAccess ensures tenant boundaries are enforced on SSH connections.
func AssertOperatorSSHAccess(server *models.Server, expectedOrgID string) error {
	if server == nil {
		return fmt.Errorf("server cannot be nil")
	}
	if server.OrganizationID != "" && expectedOrgID != "" && server.OrganizationID != expectedOrgID {
		return fmt.Errorf("server %s belongs to organization %s, cannot be accessed by %s", server.ID, server.OrganizationID, expectedOrgID)
	}
	return nil
}
