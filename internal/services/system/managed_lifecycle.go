package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"codedock.run/codedock/internal/engine/hetzner"
	"codedock.run/codedock/internal/models"
)

func (s *managedService) managedServer(ctx context.Context, orgID, serverID string) (*models.Server, *models.ManagedServerLink, error) {
	link, err := s.managedRepo.GetLink(ctx, serverID)
	if err != nil {
		return nil, nil, err
	}
	if link == nil {
		return nil, nil, fmt.Errorf("server %s is not a managed server", serverID)
	}
	if link.OrganizationID != orgID {
		return nil, nil, fmt.Errorf("managed server belongs to another organization")
	}
	server, err := s.serverRepo.GetByID(ctx, serverID)
	if err != nil {
		return nil, nil, err
	}
	if server == nil {
		return nil, nil, fmt.Errorf("managed server not found")
	}
	return server, link, nil
}

func (s *managedService) ReviewResize(ctx context.Context, userID, orgID, serverID string, req models.ReviewManagedResizeRequest) (*models.ManagedReview, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, true); err != nil {
		return nil, err
	}
	if err := s.checkEntitlement(ctx, userID); err != nil {
		return nil, err
	}
	server, link, err := s.managedServer(ctx, orgID, serverID)
	if err != nil {
		return nil, err
	}
	if server.Status == models.ServerStatusProvisioning {
		return nil, fmt.Errorf("managed server is still provisioning")
	}
	current, err := hetzner.TierByName(server.ServerType)
	if err != nil {
		return nil, fmt.Errorf("current server type is unknown")
	}
	target, err := hetzner.TierByName(req.ServerType)
	if err != nil {
		return nil, err
	}
	if target.Name == current.Name {
		return nil, fmt.Errorf("server is already type %s", current.Name)
	}
	if target.DiskGB < current.DiskGB {
		return nil, fmt.Errorf("resize to a smaller disk is not supported")
	}
	credential, err := s.credentialForOrg(ctx, orgID, link.CredentialID)
	if err != nil {
		return nil, err
	}
	if _, err := s.admit(ctx, orgID, 0, target.MemoryGB-current.MemoryGB, serverID); err != nil {
		return nil, err
	}
	plan := &models.ManagedServerPlan{
		CredentialID:   credential.ID,
		OrganizationID: orgID,
		UserID:         userID,
		ServerID:       serverID,
		ExternalID:     server.ExternalID,
		Name:           server.Name,
		Region:         server.Region,
		ServerType:     server.ServerType,
		ResizeTo:       target.Name,
		SSHKeyName:     link.SSHKeyName,
		Provider:       models.ManagedProviderHetzner,
		Tier:           target,
		MonthlyPrice:   target.MonthlyPrice,
		Action:         models.ManagedActionResize,
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	quota, err := s.quotaWithUsage(ctx, orgID)
	if err != nil {
		return nil, err
	}
	snapshot := snapshotValue(snapshotServer(server), snapshotQuota(quota), credential.UpdatedAt)
	effects := fmt.Sprintf("Resize %q from %s to %s ($%.2f/mo). The server reboots and workloads stop during the change; data on the server disk is retained.", server.Name, current.Name, target.Name, target.MonthlyPrice)
	review, err := s.operations.Review(ctx, userID, orgID, ManagedOperationResize, serverID, string(payload), snapshot, effects)
	if err != nil {
		return nil, err
	}
	return &models.ManagedReview{Review: review, Plan: plan}, nil
}

func (s *managedService) applyResize(ctx context.Context, userID string, op *models.Operation, plan *models.ManagedServerPlan, confirmation string) error {
	server, link, err := s.managedServer(ctx, plan.OrganizationID, plan.ServerID)
	if err != nil {
		return err
	}
	if server.Status == models.ServerStatusProvisioning {
		return fmt.Errorf("managed server is still provisioning")
	}
	credential, err := s.credentialForOrg(ctx, plan.OrganizationID, plan.CredentialID)
	if err != nil {
		return err
	}
	if link.CredentialID != credential.ID {
		return fmt.Errorf("resize must use the original provision credential")
	}
	current, err := hetzner.TierByName(server.ServerType)
	if err != nil {
		return fmt.Errorf("current server type is unknown")
	}
	target, err := hetzner.TierByName(plan.ResizeTo)
	if err != nil {
		return err
	}
	if _, err := s.admit(ctx, plan.OrganizationID, 0, target.MemoryGB-current.MemoryGB, server.ID); err != nil {
		return err
	}
	quota, err := s.quotaWithUsage(ctx, plan.OrganizationID)
	if err != nil {
		return err
	}
	snapshot := snapshotValue(snapshotServer(server), snapshotQuota(quota), credential.UpdatedAt)
	provider := s.providers(credential.Token)
	return s.operations.Apply(ctx, op.ID, userID, confirmation, snapshot, func(runCtx context.Context, _ *models.Operation, progress func(string, string) error) error {
		return s.runResize(runCtx, provider, server, plan, current, target, progress)
	})
}

func (s *managedService) runResize(ctx context.Context, provider ManagedProvider, server *models.Server, plan *models.ManagedServerPlan, current, target models.ManagedTier, progress func(string, string) error) error {
	instance, err := provider.GetServer(ctx, server.ExternalID)
	if err != nil {
		return s.failResize(ctx, provider, server, fmt.Errorf("inspect provider server: %w", err))
	}
	if instance.ServerType == target.Name {
		server.ServerType = target.Name
		server.Status = models.ServerStatusOnline
		server.UpdatedAt = time.Now().UTC()
		if err := s.serverRepo.Update(ctx, server); err != nil {
			return fmt.Errorf("persist resized server type: %w", err)
		}
		return progress("READY", "Server already runs type "+target.Name)
	}
	if err := progress("RESIZING", "Stopping "+server.Name+" for resize to "+target.Name); err != nil {
		return err
	}
	if instance.Status != hetzner.StatusOff {
		if err := provider.PowerOff(ctx, server.ExternalID); err != nil {
			return s.failResize(ctx, provider, server, fmt.Errorf("stop provider server: %w", err))
		}
		if _, err := provider.WaitForStatus(ctx, server.ExternalID, []string{hetzner.StatusOff}, s.serverWait); err != nil {
			return s.failResize(ctx, provider, server, fmt.Errorf("provider server did not stop: %w", err))
		}
	}
	if err := provider.ChangeType(ctx, server.ExternalID, target.Name, target.DiskGB > current.DiskGB); err != nil {
		return s.failResize(ctx, provider, server, fmt.Errorf("change provider server type: %w", err))
	}
	if err := progress("RESIZING", "Starting "+server.Name+" on type "+target.Name); err != nil {
		return s.failResize(ctx, provider, server, err)
	}
	if err := provider.PowerOn(ctx, server.ExternalID); err != nil {
		return s.failResize(ctx, provider, server, fmt.Errorf("start provider server: %w", err))
	}
	if _, err := provider.WaitForStatus(ctx, server.ExternalID, []string{hetzner.StatusRunning}, s.serverWait); err != nil {
		return s.failResize(ctx, provider, server, fmt.Errorf("provider server did not start: %w", err))
	}
	if err := s.bootstrapDockerHost(ctx, server, progress); err != nil {
		return s.failResize(ctx, provider, server, err)
	}
	server.ServerType = target.Name
	server.Status = models.ServerStatusOnline
	server.UpdatedAt = time.Now().UTC()
	if err := s.serverRepo.Update(ctx, server); err != nil {
		return s.failResize(ctx, provider, server, fmt.Errorf("persist resized server type: %w", err))
	}
	return progress("READY", "Server resized to "+target.Name)
}

func (s *managedService) failResize(ctx context.Context, provider ManagedProvider, server *models.Server, cause error) error {
	server.Status = models.ServerStatusOffline
	server.UpdatedAt = time.Now().UTC()
	if err := s.serverRepo.Update(ctx, server); err != nil {
		slog.Error("persist managed resize failure", "server", server.ID, "error", err)
	}
	if provider != nil && server.ExternalID != "" {
		if err := provider.PowerOn(ctx, server.ExternalID); err != nil {
			slog.Error("restart server after failed resize", "server", server.ID, "error", err)
		}
	}
	return cause
}

func (s *managedService) RefreshServer(ctx context.Context, userID, orgID, serverID string) (*models.Server, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	server, link, err := s.managedServer(ctx, orgID, serverID)
	if err != nil {
		return nil, err
	}
	credential, err := s.credentialForOrg(ctx, orgID, link.CredentialID)
	if err != nil {
		return nil, err
	}
	return s.refreshFromProvider(ctx, server, s.providers(credential.Token))
}

func (s *managedService) refreshFromProvider(ctx context.Context, server *models.Server, provider ManagedProvider) (*models.Server, error) {
	if server.ExternalID == "" {
		adopted, err := provider.FindServerByLabel(ctx, serverLabelSelector(server.ID))
		if err != nil {
			return nil, fmt.Errorf("inspect provider servers: %w", err)
		}
		if adopted.ExternalID == "" {
			return server, nil
		}
		server.ExternalID = adopted.ExternalID
	}
	instance, err := provider.GetServer(ctx, server.ExternalID)
	if err != nil {
		var apiErr *hetzner.APIError
		if errors.As(err, &apiErr) && apiErr.Status == 404 {
			server.Status = models.ServerStatusOffline
			server.UpdatedAt = time.Now().UTC()
			if updateErr := s.serverRepo.Update(ctx, server); updateErr != nil {
				return nil, fmt.Errorf("mark missing managed server: %w", updateErr)
			}
			return server, nil
		}
		return nil, fmt.Errorf("inspect provider server: %w", err)
	}
	if instance.PublicIP != "" {
		server.IPAddress = instance.PublicIP
		server.SSHHost = instance.PublicIP
	}
	switch instance.Status {
	case hetzner.StatusRunning:
		if s.sshAlive(ctx, server) {
			server.Status = models.ServerStatusOnline
		}
	case hetzner.StatusOff:
		server.Status = models.ServerStatusOffline
	}
	server.UpdatedAt = time.Now().UTC()
	if err := s.serverRepo.Update(ctx, server); err != nil {
		return nil, fmt.Errorf("failed to refresh managed server: %w", err)
	}
	return server, nil
}

func (s *managedService) sshAlive(ctx context.Context, server *models.Server) bool {
	if s.sshManager == nil || server.SSHHost == "" {
		return false
	}
	check, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.sshManager.TestConnection(check, sshConfigFor(server)) == nil
}

func (s *managedService) Recover(ctx context.Context) error {
	ids, err := s.managedRepo.ListProvisioningServerIDs(ctx)
	if err != nil {
		return err
	}
	for _, id := range ids {
		link, err := s.managedRepo.GetLink(ctx, id)
		if err != nil || link == nil {
			continue
		}
		credential, err := s.managedRepo.GetCredential(ctx, link.CredentialID)
		if err != nil || credential == nil || credential.Token == "" {
			slog.Error("recover managed server without credential", "server", id)
			continue
		}
		server, err := s.serverRepo.GetByID(ctx, id)
		if err != nil || server == nil {
			continue
		}
		if _, err := s.refreshFromProvider(ctx, server, s.providers(credential.Token)); err != nil {
			slog.Error("recover managed server", "server", id, "error", err)
		}
	}
	return nil
}

func deriveSSHPublicKey(privatePEM string) (string, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privatePEM))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(signer.PublicKey()))), nil
}
