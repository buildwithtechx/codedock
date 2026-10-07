package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/hetzner"
	"codedock.run/codedock/internal/models"
)

func (s *managedService) ReviewProvision(ctx context.Context, userID, orgID string, req models.ReviewManagedProvisionRequest) (*models.ManagedReview, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	if err := s.checkEntitlement(ctx, userID); err != nil {
		return nil, err
	}
	tier, err := s.validateProvisionTarget(req)
	if err != nil {
		return nil, err
	}
	credential, err := s.credentialForOrg(ctx, orgID, req.CredentialID)
	if err != nil {
		return nil, err
	}
	serverID := req.ServerID
	exclude := ""
	if serverID != "" {
		server, link, err := s.managedServer(ctx, orgID, serverID)
		if err != nil {
			return nil, err
		}
		if server.Status == models.ServerStatusOnline {
			return nil, fmt.Errorf("managed server is already running")
		}
		if link.CredentialID != credential.ID {
			return nil, fmt.Errorf("retry must use the original provision credential")
		}
		exclude = serverID
	} else {
		serverID = uuid.NewString()
	}
	if _, err := s.admit(ctx, orgID, 1, tier.MemoryGB, exclude); err != nil {
		return nil, err
	}
	plan := &models.ManagedServerPlan{
		CredentialID:   credential.ID,
		OrganizationID: orgID,
		UserID:         userID,
		ServerID:       serverID,
		Name:           req.Name,
		Region:         req.Region,
		ServerType:     req.ServerType,
		Image:          req.Image,
		SSHKeyName:     sshKeyName(serverID),
		Provider:       models.ManagedProviderHetzner,
		Tier:           tier,
		MonthlyPrice:   tier.MonthlyPrice,
		Action:         models.ManagedActionProvision,
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	snapshot, err := s.provisionSnapshot(ctx, orgID, credential, serverID)
	if err != nil {
		return nil, err
	}
	effects := fmt.Sprintf("Create Hetzner server %q (%s, %s, %s) at $%.2f/mo billed by Hetzner to the stored credential. The server joins this organization's SSH targets; teardown requires a reviewed delete.", req.Name, req.ServerType, req.Region, req.Image, tier.MonthlyPrice)
	review, err := s.operations.Review(ctx, userID, orgID, ManagedOperationProvision, serverID, string(payload), snapshot, effects)
	if err != nil {
		return nil, err
	}
	return &models.ManagedReview{Review: review, Plan: plan}, nil
}

func (s *managedService) ApplyOperation(ctx context.Context, userID, operationID, confirmation string) error {
	op, err := s.operations.Get(ctx, operationID)
	if err != nil {
		return err
	}
	var plan models.ManagedServerPlan
	if err := json.Unmarshal([]byte(op.Payload), &plan); err != nil {
		return fmt.Errorf("managed operation payload is invalid: %w", err)
	}
	if err := s.requireOrgRole(ctx, userID, plan.OrganizationID, false); err != nil {
		return err
	}
	if err := s.checkEntitlement(ctx, userID); err != nil && plan.Action != models.ManagedActionDelete {
		return err
	}
	switch op.Kind {
	case ManagedOperationProvision:
		return s.applyProvision(ctx, userID, op, &plan, confirmation)
	case ManagedOperationResize:
		return s.applyResize(ctx, userID, op, &plan, confirmation)
	case ManagedOperationDelete:
		return s.applyDelete(ctx, userID, op, &plan, confirmation)
	default:
		return fmt.Errorf("unsupported managed operation")
	}
}

func (s *managedService) applyProvision(ctx context.Context, userID string, op *models.Operation, plan *models.ManagedServerPlan, confirmation string) error {
	credential, err := s.credentialForOrg(ctx, plan.OrganizationID, plan.CredentialID)
	if err != nil {
		return err
	}
	tier, err := hetzner.TierByName(plan.ServerType)
	if err != nil {
		return err
	}
	if _, err := s.admit(ctx, plan.OrganizationID, 1, tier.MemoryGB, plan.ServerID); err != nil {
		return err
	}
	snapshot, err := s.provisionSnapshot(ctx, plan.OrganizationID, credential, plan.ServerID)
	if err != nil {
		return err
	}
	provider := s.providers(credential.Token)
	return s.operations.Apply(ctx, op.ID, userID, confirmation, snapshot, func(runCtx context.Context, _ *models.Operation, progress func(string, string) error) error {
		return s.runProvision(runCtx, provider, credential, plan, progress)
	})
}

func (s *managedService) runProvision(ctx context.Context, provider ManagedProvider, _ *models.ManagedCredential, plan *models.ManagedServerPlan, progress func(string, string) error) error {
	fail := func(server *models.Server, cause error) error {
		if server != nil {
			server.Status = models.ServerStatusOffline
			server.UpdatedAt = time.Now().UTC()
			if err := s.serverRepo.Update(ctx, server); err != nil {
				slog.Error("persist managed provision failure", "server", server.ID, "error", err)
			}
		}
		return cause
	}
	server, err := s.ensureServerRow(ctx, plan)
	if err != nil {
		return err
	}
	if err := progress("PROVISIONING", "Reserving "+plan.Name+" ("+plan.ServerType+", "+plan.Region+")"); err != nil {
		return fail(server, err)
	}
	if server.SSHPrivateKey == "" {
		private, _, err := hetzner.GenerateKeyPair()
		if err != nil {
			return fail(server, fmt.Errorf("generate server ssh key: %w", err))
		}
		server.SSHPrivateKey = private
		server.SSHKey = private
		server.SSHPassword = ""
		if err := s.serverRepo.Update(ctx, server); err != nil {
			return fail(server, fmt.Errorf("persist server ssh key: %w", err))
		}
	}
	publicKey, err := deriveSSHPublicKey(server.SSHPrivateKey)
	if err != nil {
		return fail(server, fmt.Errorf("derive server ssh key: %w", err))
	}
	if err := ensureProviderKey(ctx, provider, plan.SSHKeyName, publicKey); err != nil {
		return fail(server, fmt.Errorf("register provider ssh key: %w", err))
	}
	instance, err := s.ensureProviderServer(ctx, provider, server, plan)
	if err != nil {
		return fail(server, err)
	}
	if err := progress("BOOTSTRAP", "Waiting for "+instance.ExternalID+" to boot"); err != nil {
		return fail(server, err)
	}
	instance, err = provider.WaitForStatus(ctx, instance.ExternalID, []string{hetzner.StatusRunning}, s.serverWait)
	if err != nil {
		return fail(server, fmt.Errorf("provider server did not become ready: %w", err))
	}
	if instance.PublicIP == "" {
		return fail(server, fmt.Errorf("provider server has no public ip yet"))
	}
	server.ExternalID = instance.ExternalID
	server.IPAddress = instance.PublicIP
	server.SSHHost = instance.PublicIP
	if err := s.serverRepo.Update(ctx, server); err != nil {
		return fail(server, fmt.Errorf("persist provider address: %w", err))
	}
	if err := s.bootstrapDockerHost(ctx, server, progress); err != nil {
		return fail(server, err)
	}
	server.Status = models.ServerStatusOnline
	server.UpdatedAt = time.Now().UTC()
	if err := s.serverRepo.Update(ctx, server); err != nil {
		return fail(server, fmt.Errorf("mark managed server ready: %w", err))
	}
	if err := progress("READY", "Managed server "+server.IPAddress+" is ready"); err != nil {
		return fail(server, err)
	}
	return nil
}

func (s *managedService) ensureServerRow(ctx context.Context, plan *models.ManagedServerPlan) (*models.Server, error) {
	existing, err := s.serverRepo.GetByID(ctx, plan.ServerID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if existing.Provider != string(models.ManagedProviderHetzner) {
			return nil, fmt.Errorf("server %s is not a managed server", plan.ServerID)
		}
		existing.Status = models.ServerStatusProvisioning
		existing.Name = plan.Name
		existing.Region = plan.Region
		existing.ServerType = plan.ServerType
		existing.UpdatedAt = time.Now().UTC()
		if err := s.serverRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to reopen managed server: %w", err)
		}
		return existing, nil
	}
	now := time.Now().UTC()
	server := &models.Server{
		ID:            plan.ServerID,
		UserID:        plan.UserID,
		Name:          plan.Name,
		SSHPort:       22,
		SSHUser:       "root",
		SSHAuthMethod: "key",
		SSHTransport:  "direct",
		Status:        models.ServerStatusProvisioning,
		Provider:      string(models.ManagedProviderHetzner),
		Region:        plan.Region,
		ServerType:    plan.ServerType,
		WorkerToken:   generateWorkerToken(),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.serverRepo.Create(ctx, server); err != nil {
		return nil, fmt.Errorf("failed to reserve managed server: %w", err)
	}
	if err := s.managedRepo.SaveLink(ctx, &models.ManagedServerLink{
		ServerID:       plan.ServerID,
		OrganizationID: plan.OrganizationID,
		CredentialID:   plan.CredentialID,
		SSHKeyName:     plan.SSHKeyName,
	}); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *managedService) ensureProviderServer(ctx context.Context, provider ManagedProvider, server *models.Server, plan *models.ManagedServerPlan) (models.ManagedInstance, error) {
	if server.ExternalID != "" {
		instance, err := provider.GetServer(ctx, server.ExternalID)
		if err == nil {
			return instance, nil
		}
		var apiErr *hetzner.APIError
		if errors.As(err, &apiErr) && apiErr.Status == 404 {
			server.ExternalID = ""
		} else {
			return models.ManagedInstance{}, fmt.Errorf("inspect provider server: %w", err)
		}
	}
	adopted, err := provider.FindServerByLabel(ctx, serverLabelSelector(server.ID))
	if err != nil {
		return models.ManagedInstance{}, fmt.Errorf("inspect provider servers: %w", err)
	}
	if adopted.ExternalID != "" {
		server.ExternalID = adopted.ExternalID
		if err := s.serverRepo.Update(ctx, server); err != nil {
			return models.ManagedInstance{}, fmt.Errorf("adopt provider server: %w", err)
		}
		return adopted, nil
	}
	created, err := provider.CreateServer(ctx, models.ManagedServerSpec{
		Name:       plan.Name,
		ServerType: plan.ServerType,
		Image:      plan.Image,
		Region:     plan.Region,
		SSHKeyName: plan.SSHKeyName,
		UserData:   hetzner.ProvisionUserData(plan.Name),
		Labels: map[string]string{
			"codedock-server":     server.ID,
			"codedock-region":     plan.Region,
			"codedock-credential": plan.CredentialID,
		},
	})
	if err != nil {
		return models.ManagedInstance{}, fmt.Errorf("create provider server: %w", err)
	}
	if created.ExternalID == "" {
		return models.ManagedInstance{}, fmt.Errorf("provider did not return a server id")
	}
	server.ExternalID = created.ExternalID
	if err := s.serverRepo.Update(ctx, server); err != nil {
		return models.ManagedInstance{}, fmt.Errorf("persist provider server id: %w", err)
	}
	return created, nil
}
