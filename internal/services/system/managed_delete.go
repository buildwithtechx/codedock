package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"codedock.run/codedock/internal/engine/hetzner"
	"codedock.run/codedock/internal/models"
)

func (s *managedService) ReviewDelete(ctx context.Context, userID, orgID, serverID string) (*models.ManagedReview, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, true); err != nil {
		return nil, err
	}
	server, link, err := s.managedServer(ctx, orgID, serverID)
	if err != nil {
		return nil, err
	}
	if server.Status == models.ServerStatusProvisioning {
		return nil, fmt.Errorf("managed server is still provisioning")
	}
	tier, err := hetzner.TierByName(server.ServerType)
	if err != nil {
		tier = models.ManagedTier{Name: server.ServerType}
	}
	plan := &models.ManagedServerPlan{
		CredentialID:   link.CredentialID,
		OrganizationID: orgID,
		UserID:         userID,
		ServerID:       serverID,
		ExternalID:     server.ExternalID,
		Name:           server.Name,
		Region:         server.Region,
		ServerType:     server.ServerType,
		SSHKeyName:     link.SSHKeyName,
		Provider:       models.ManagedProviderHetzner,
		Tier:           tier,
		MonthlyPrice:   tier.MonthlyPrice,
		Action:         models.ManagedActionDelete,
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	snapshot := snapshotValue(snapshotServer(server), link.CredentialID)
	effects := fmt.Sprintf("Permanently destroy Hetzner server %q (%s) and all data on its disk. Billing stops after Hetzner removes the server.", server.Name, server.ExternalID)
	review, err := s.operations.Review(ctx, userID, orgID, ManagedOperationDelete, serverID, string(payload), snapshot, effects)
	if err != nil {
		return nil, err
	}
	return &models.ManagedReview{Review: review, Plan: plan}, nil
}

func (s *managedService) applyDelete(ctx context.Context, userID string, op *models.Operation, plan *models.ManagedServerPlan, confirmation string) error {
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
		return fmt.Errorf("delete must use the original provision credential")
	}
	snapshot := snapshotValue(snapshotServer(server), link.CredentialID)
	provider := s.providers(credential.Token)
	return s.operations.Apply(ctx, op.ID, userID, confirmation, snapshot, func(runCtx context.Context, _ *models.Operation, progress func(string, string) error) error {
		return s.runDelete(runCtx, provider, server, plan, progress)
	})
}

func (s *managedService) runDelete(ctx context.Context, provider ManagedProvider, server *models.Server, plan *models.ManagedServerPlan, progress func(string, string) error) error {
	if err := progress("DELETING", "Destroying provider server "+server.ExternalID); err != nil {
		return err
	}
	if server.ExternalID != "" {
		if err := provider.DeleteServer(ctx, server.ExternalID); err != nil {
			var apiErr *hetzner.APIError
			if !errors.As(err, &apiErr) || apiErr.Status != 404 {
				return fmt.Errorf("delete provider server: %w", err)
			}
		}
	}
	if err := provider.DeleteSSHKey(ctx, plan.SSHKeyName); err != nil {
		slog.Error("delete provider ssh key", "server", server.ID, "key", plan.SSHKeyName, "error", err)
	}
	if s.sshManager != nil {
		s.sshManager.RemoveClient(server.ID)
	}
	if err := s.managedRepo.DeleteLink(ctx, server.ID); err != nil {
		return err
	}
	if err := s.serverRepo.Delete(ctx, server.ID); err != nil {
		return fmt.Errorf("failed to remove managed server: %w", err)
	}
	return progress("READY", "Managed server destroyed")
}
