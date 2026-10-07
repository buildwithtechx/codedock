package system

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/engine/hetzner"
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/operations"
)

const (
	ManagedOperationProvision = "managed-provision"
	ManagedOperationResize    = "managed-resize"
	ManagedOperationDelete    = "managed-delete"
)

var managedNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.-]{0,62}$`)

type ManagedProvider interface {
	CreateServer(ctx context.Context, spec models.ManagedServerSpec) (models.ManagedInstance, error)
	GetServer(ctx context.Context, externalID string) (models.ManagedInstance, error)
	FindServerByLabel(ctx context.Context, selector string) (models.ManagedInstance, error)
	DeleteServer(ctx context.Context, externalID string) error
	PowerOn(ctx context.Context, externalID string) error
	PowerOff(ctx context.Context, externalID string) error
	ChangeType(ctx context.Context, externalID, serverType string, upgradeDisk bool) error
	WaitForStatus(ctx context.Context, externalID string, wanted []string, timeout time.Duration) (models.ManagedInstance, error)
	CreateSSHKey(ctx context.Context, name, publicKey string) error
	DeleteSSHKey(ctx context.Context, name string) error
}

type ProviderFactory func(token string) ManagedProvider

func DefaultProviderFactory(token string) ManagedProvider {
	return hetzner.NewProvider(token)
}

type ManagedService interface {
	Catalog() models.ManagedCatalog
	CreateCredential(ctx context.Context, userID, orgID string, req models.CreateManagedCredentialRequest) (*models.ManagedCredential, error)
	ListCredentials(ctx context.Context, userID, orgID string) ([]*models.ManagedCredential, error)
	DeleteCredential(ctx context.Context, userID, orgID, credentialID string) error
	GetQuota(ctx context.Context, userID, orgID string) (*models.ManagedQuota, error)
	SetQuota(ctx context.Context, userID, orgID string, req models.ManagedQuotaRequest) (*models.ManagedQuota, error)
	ReviewProvision(ctx context.Context, userID, orgID string, req models.ReviewManagedProvisionRequest) (*models.ManagedReview, error)
	ReviewResize(ctx context.Context, userID, orgID, serverID string, req models.ReviewManagedResizeRequest) (*models.ManagedReview, error)
	ReviewDelete(ctx context.Context, userID, orgID, serverID string) (*models.ManagedReview, error)
	ApplyOperation(ctx context.Context, userID, operationID, confirmation string) error
	RefreshServer(ctx context.Context, userID, orgID, serverID string) (*models.Server, error)
	Recover(ctx context.Context) error
}

type managedService struct {
	managedRepo repositories.ManagedRepository
	serverRepo  repositories.ServerRepository
	orgRepo     repositories.OrganizationRepository
	userRepo    *repositories.UserRepo
	sshManager  *ssh.SSHManager
	operations  *operations.Service
	providers   ProviderFactory

	pollInterval time.Duration
	serverWait   time.Duration
	sshWait      time.Duration
}

func NewManagedService(managedRepo repositories.ManagedRepository, serverRepo repositories.ServerRepository, orgRepo repositories.OrganizationRepository, userRepo *repositories.UserRepo, sshManager *ssh.SSHManager, operations *operations.Service, factory ProviderFactory) ManagedService {
	if factory == nil {
		factory = DefaultProviderFactory
	}
	return &managedService{
		managedRepo:  managedRepo,
		serverRepo:   serverRepo,
		orgRepo:      orgRepo,
		userRepo:     userRepo,
		sshManager:   sshManager,
		operations:   operations,
		providers:    factory,
		pollInterval: 5 * time.Second,
		serverWait:   10 * time.Minute,
		sshWait:      10 * time.Minute,
	}
}

func (s *managedService) SetProvisionTimeouts(pollInterval, serverWait, sshWait time.Duration) {
	s.pollInterval = pollInterval
	s.serverWait = serverWait
	s.sshWait = sshWait
}

func (s *managedService) Catalog() models.ManagedCatalog {
	return models.ManagedCatalog{
		Provider: models.ManagedProviderHetzner,
		Tiers:    hetzner.Tiers(),
		Regions:  hetzner.Regions(),
		Images:   hetzner.Images(),
	}
}

func (s *managedService) CreateCredential(ctx context.Context, userID, orgID string, req models.CreateManagedCredentialRequest) (*models.ManagedCredential, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, true); err != nil {
		return nil, err
	}
	if req.Provider != models.ManagedProviderHetzner {
		return nil, fmt.Errorf("unsupported managed provider")
	}
	if req.Label == "" {
		return nil, fmt.Errorf("credential label is required")
	}
	if req.Token == "" {
		return nil, fmt.Errorf("provider token is required")
	}
	credential := &models.ManagedCredential{
		ID:             uuid.NewString(),
		OrganizationID: orgID,
		Provider:       req.Provider,
		Label:          req.Label,
		Token:          req.Token,
	}
	if err := s.managedRepo.CreateCredential(ctx, credential); err != nil {
		return nil, err
	}
	credential.Token = ""
	return credential, nil
}

func (s *managedService) ListCredentials(ctx context.Context, userID, orgID string) ([]*models.ManagedCredential, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	return s.managedRepo.ListCredentialsByOrg(ctx, orgID)
}

func (s *managedService) DeleteCredential(ctx context.Context, userID, orgID, credentialID string) error {
	if err := s.requireOrgRole(ctx, userID, orgID, true); err != nil {
		return err
	}
	credential, err := s.managedRepo.GetCredential(ctx, credentialID)
	if err != nil {
		return err
	}
	if credential == nil || credential.OrganizationID != orgID {
		return fmt.Errorf("managed credential not found")
	}
	return s.managedRepo.DeleteCredential(ctx, credentialID)
}

func (s *managedService) GetQuota(ctx context.Context, userID, orgID string) (*models.ManagedQuota, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	return s.quotaWithUsage(ctx, orgID)
}

func (s *managedService) SetQuota(ctx context.Context, userID, orgID string, req models.ManagedQuotaRequest) (*models.ManagedQuota, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, true); err != nil {
		return nil, err
	}
	if req.MaxServers < 0 || req.MaxMemoryGB < 0 {
		return nil, fmt.Errorf("quota limits must not be negative")
	}
	quota, err := s.quotaWithUsage(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if req.MaxServers < quota.UsedServers {
		return nil, fmt.Errorf("quota below current usage of %d server(s)", quota.UsedServers)
	}
	if req.MaxMemoryGB < quota.UsedMemoryGB {
		return nil, fmt.Errorf("quota below current usage of %d GB", quota.UsedMemoryGB)
	}
	quota.MaxServers = req.MaxServers
	quota.MaxMemoryGB = req.MaxMemoryGB
	if err := s.managedRepo.SetQuota(ctx, quota); err != nil {
		return nil, err
	}
	return quota, nil
}

func (s *managedService) quotaWithUsage(ctx context.Context, orgID string) (*models.ManagedQuota, error) {
	quota, err := s.managedRepo.GetQuota(ctx, orgID)
	if err != nil {
		return nil, err
	}
	ids, err := s.managedRepo.ListOrgServerIDs(ctx, orgID)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		server, err := s.serverRepo.GetByID(ctx, id)
		if err != nil || server == nil {
			continue
		}
		quota.UsedServers++
		if tier, err := hetzner.TierByName(server.ServerType); err == nil {
			quota.UsedMemoryGB += tier.MemoryGB
		}
	}
	return quota, nil
}

func (s *managedService) admit(ctx context.Context, orgID string, extraServers, extraMemoryGB int, excludeServerID string) (*models.ManagedQuota, error) {
	quota, err := s.quotaWithUsage(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if excludeServerID != "" {
		if server, err := s.serverRepo.GetByID(ctx, excludeServerID); err == nil && server != nil {
			quota.UsedServers--
			if tier, err := hetzner.TierByName(server.ServerType); err == nil {
				quota.UsedMemoryGB -= tier.MemoryGB
			}
		}
	}
	if quota.UsedServers+extraServers > quota.MaxServers {
		return nil, fmt.Errorf("managed server quota exceeded: %d of %d in use", quota.UsedServers, quota.MaxServers)
	}
	if quota.UsedMemoryGB+extraMemoryGB > quota.MaxMemoryGB {
		return nil, fmt.Errorf("managed memory quota exceeded: %d of %d GB in use", quota.UsedMemoryGB, quota.MaxMemoryGB)
	}
	return quota, nil
}

func (s *managedService) checkEntitlement(ctx context.Context, userID string) error {
	if !config.Get().Cloud.Enabled {
		return nil
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load billing entitlement: %w", err)
	}
	if user == nil || user.PlanType != "pro" {
		return fmt.Errorf("managed servers require a pro subscription")
	}
	return nil
}

func (s *managedService) requireOrgRole(ctx context.Context, userID, orgID string, admin bool) error {
	if orgID == "" {
		return fmt.Errorf("organization id is required")
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %w", err)
	}
	if user != nil && (user.Role == models.UserRoleOwner || user.Role == models.UserRoleAdmin) {
		return nil
	}
	member, err := s.orgRepo.GetMember(ctx, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to load organization membership: %w", err)
	}
	if member == nil {
		return fmt.Errorf("organization membership is required")
	}
	if admin && member.Permission != models.MemberPermissionAdmin && member.Permission != models.MemberPermissionOwner {
		return fmt.Errorf("organization admin permission is required")
	}
	return nil
}

func (s *managedService) validateProvisionTarget(req models.ReviewManagedProvisionRequest) (models.ManagedTier, error) {
	if !managedNamePattern.MatchString(req.Name) {
		return models.ManagedTier{}, fmt.Errorf("server name must be 1-63 characters of letters, digits, dots or dashes")
	}
	regionOK := false
	for _, region := range hetzner.Regions() {
		if region == req.Region {
			regionOK = true
			break
		}
	}
	if !regionOK {
		return models.ManagedTier{}, fmt.Errorf("unsupported region %q", req.Region)
	}
	imageOK := false
	for _, image := range hetzner.Images() {
		if image == req.Image {
			imageOK = true
			break
		}
	}
	if !imageOK {
		return models.ManagedTier{}, fmt.Errorf("unsupported image %q", req.Image)
	}
	tier, err := hetzner.TierByName(req.ServerType)
	if err != nil {
		return models.ManagedTier{}, err
	}
	return tier, nil
}

func (s *managedService) credentialForOrg(ctx context.Context, orgID, credentialID string) (*models.ManagedCredential, error) {
	if credentialID == "" {
		return nil, fmt.Errorf("provider credential is required")
	}
	credential, err := s.managedRepo.GetCredential(ctx, credentialID)
	if err != nil {
		return nil, err
	}
	if credential == nil || credential.OrganizationID != orgID {
		return nil, fmt.Errorf("managed credential not found")
	}
	if credential.Provider != models.ManagedProviderHetzner || credential.Token == "" {
		return nil, fmt.Errorf("managed credential is not usable")
	}
	return credential, nil
}

func snapshotValue(parts ...string) string {
	raw, _ := json.Marshal(parts)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
