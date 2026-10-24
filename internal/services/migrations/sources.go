package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock/internal/engine/dockerprobe"
	"codedock/internal/models"
)

func (s *Service) ListSources(ctx context.Context, userID, orgID string) ([]*models.MigrationSource, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	return s.migrations.ListSourcesByOrg(ctx, orgID)
}

func (s *Service) TestSource(ctx context.Context, userID, orgID string, req models.CreateMigrationSourceRequest) (string, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return "", err
	}
	host, port, user, err := normalizeSourceEndpoint(req)
	if err != nil {
		return "", err
	}
	check, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	fingerprint, err := dockerprobe.FingerprintFor(host, port, 15*time.Second)
	if err != nil {
		return "", fmt.Errorf("reach migration source: %w", err)
	}
	runner := dockerprobe.NewSSHRunner(dockerprobe.SSHConfig{
		Host: host, Port: port, User: user,
		PrivateKey: req.SSHKey, Password: req.SSHPassword, Fingerprint: fingerprint,
	})
	if _, err := runner.Run(check, "docker info --format '{{.ServerVersion}}'"); err != nil {
		return "", fmt.Errorf("source docker check failed: %w", err)
	}
	return fingerprint, nil
}

func (s *Service) CreateSource(ctx context.Context, userID, orgID string, req models.CreateMigrationSourceRequest) (*models.MigrationSource, error) {
	if err := s.requireOrg(ctx, userID, orgID, true); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, fmt.Errorf("source name is required")
	}
	host, port, user, err := normalizeSourceEndpoint(req)
	if err != nil {
		return nil, err
	}
	if req.SSHKey == "" && req.SSHPassword == "" {
		return nil, fmt.Errorf("ssh key or password is required")
	}
	fingerprint := req.Fingerprint
	if fingerprint == "" {
		discovered, err := dockerprobe.FingerprintFor(host, port, 15*time.Second)
		if err != nil {
			return nil, fmt.Errorf("reach migration source: %w", err)
		}
		fingerprint = discovered
	}
	runner := dockerprobe.NewSSHRunner(dockerprobe.SSHConfig{
		Host: host, Port: port, User: user,
		PrivateKey: req.SSHKey, Password: req.SSHPassword, Fingerprint: fingerprint,
	})
	check, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if _, err := runner.Run(check, "docker info --format '{{.ServerVersion}}'"); err != nil {
		return nil, fmt.Errorf("source docker check failed: %w", err)
	}
	method := req.SSHAuthMethod
	if method == "" {
		if req.SSHPassword != "" {
			method = "password"
		} else {
			method = "key"
		}
	}
	source := &models.MigrationSource{
		ID:             uuid.NewString(),
		OrganizationID: orgID,
		Name:           req.Name,
		SSHHost:        host,
		SSHPort:        port,
		SSHUser:        user,
		SSHAuthMethod:  method,
		SSHKey:         req.SSHKey,
		SSHPassword:    req.SSHPassword,
		Fingerprint:    fingerprint,
	}
	if err := s.migrations.CreateSource(ctx, source); err != nil {
		return nil, err
	}
	source.SSHKey = ""
	source.SSHPassword = ""
	return source, nil
}

func (s *Service) DeleteSource(ctx context.Context, userID, orgID, sourceID string) error {
	if err := s.requireOrg(ctx, userID, orgID, true); err != nil {
		return err
	}
	if _, err := s.orgSource(ctx, orgID, sourceID); err != nil {
		return err
	}
	return s.migrations.DeleteSource(ctx, sourceID)
}

func normalizeSourceEndpoint(req models.CreateMigrationSourceRequest) (string, int, string, error) {
	if req.SSHHost == "" {
		return "", 0, "", fmt.Errorf("ssh host is required")
	}
	port := req.SSHPort
	if port <= 0 {
		port = 22
	}
	if port > 65535 {
		return "", 0, "", fmt.Errorf("ssh port must be between 1 and 65535")
	}
	user := req.SSHUser
	if user == "" {
		user = "root"
	}
	return req.SSHHost, port, user, nil
}
