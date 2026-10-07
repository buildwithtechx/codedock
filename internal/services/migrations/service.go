package migrations

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/dockerprobe"
	enginessh "codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	deploymentservices "codedock.run/codedock/internal/services/deployments"
	projectservices "codedock.run/codedock/internal/services/projects"
)

type AppDeployer interface {
	DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error)
	StopAppService(ctx context.Context, app *models.AppService) error
	RemoveAppService(ctx context.Context, app *models.AppService) error
}

type RunnerFactory func(source *models.MigrationSource) dockerprobe.Runner
type TargetFactory func(server *models.Server) dockerprobe.Runner

type Service struct {
	migrations repositories.MigrationRepository
	projects   repositories.ProjectRepository
	apps       *projectservices.AppService
	appRepo    repositories.AppServiceRepository
	volumes    repositories.ServiceVolumeRepository
	envs       repositories.EnvironmentRepository
	servers    repositories.ServerRepository
	users      *repositories.UserRepo
	orgs       repositories.OrganizationRepository
	git        *deploymentservices.GitService
	composer   *projectservices.ComposeParserService
	deployer   AppDeployer
	ssh        *enginessh.SSHManager
	runtimes   RuntimeKinds
	sources    RunnerFactory
	targets    TargetFactory
	runs       *runRegistry
}

func defaultSourceRunner(source *models.MigrationSource) dockerprobe.Runner {
	return dockerprobe.NewSSHRunner(dockerprobe.SSHConfig{
		Host:        source.SSHHost,
		Port:        source.SSHPort,
		User:        source.SSHUser,
		PrivateKey:  source.SSHKey,
		Password:    source.SSHPassword,
		Fingerprint: source.Fingerprint,
	})
}

func NewService(
	migrations repositories.MigrationRepository,
	projects repositories.ProjectRepository,
	apps *projectservices.AppService,
	appRepo repositories.AppServiceRepository,
	volumes repositories.ServiceVolumeRepository,
	envs repositories.EnvironmentRepository,
	servers repositories.ServerRepository,
	users *repositories.UserRepo,
	orgs repositories.OrganizationRepository,
	git *deploymentservices.GitService,
	composer *projectservices.ComposeParserService,
	deployer AppDeployer,
	sshManager *enginessh.SSHManager,
	runtimes RuntimeKinds,
	sources RunnerFactory,
	targets TargetFactory,
) *Service {
	if sources == nil {
		sources = defaultSourceRunner
	}
	service := &Service{
		migrations: migrations,
		projects:   projects,
		apps:       apps,
		appRepo:    appRepo,
		volumes:    volumes,
		envs:       envs,
		servers:    servers,
		users:      users,
		orgs:       orgs,
		git:        git,
		composer:   composer,
		deployer:   deployer,
		ssh:        sshManager,
		runtimes:   runtimes,
		sources:    sources,
		targets:    targets,
		runs:       newRunRegistry(),
	}
	if targets == nil {
		service.targets = service.defaultTargetRunner
	}
	return service
}

func (s *Service) defaultTargetRunner(server *models.Server) dockerprobe.Runner {
	if server == nil || server.IsLocal {
		return dockerprobe.LocalRunner{}
	}
	return &engineRunner{manager: s.ssh, server: server}
}

type engineRunner struct {
	manager *enginessh.SSHManager
	server  *models.Server
}

func (r *engineRunner) Run(ctx context.Context, cmd string) (string, error) {
	if r.manager == nil {
		return "", fmt.Errorf("ssh manager is not initialized")
	}
	client, release, err := r.manager.GetClient(r.server)
	if err != nil {
		return "", err
	}
	defer release()
	return client.RunCommand(ctx, cmd)
}

func (r *engineRunner) RunPipe(ctx context.Context, cmd string, stdin io.Reader, stdout io.Writer) error {
	if r.manager == nil {
		return fmt.Errorf("ssh manager is not initialized")
	}
	client, release, err := r.manager.GetClient(r.server)
	if err != nil {
		return err
	}
	defer release()
	return client.Stream(ctx, cmd, stdin, stdout)
}

func (s *Service) requireOrg(ctx context.Context, userID, orgID string, admin bool) error {
	if orgID == "" {
		return fmt.Errorf("organization id is required")
	}
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %w", err)
	}
	if user != nil && (user.Role == models.UserRoleOwner || user.Role == models.UserRoleAdmin) {
		return nil
	}
	member, err := s.orgs.GetMember(ctx, orgID, userID)
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

func (s *Service) orgSource(ctx context.Context, orgID, sourceID string) (*models.MigrationSource, error) {
	source, err := s.migrations.GetSource(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if source == nil || source.OrganizationID != orgID {
		return nil, fmt.Errorf("migration source not found")
	}
	return source, nil
}

func (s *Service) orgRun(ctx context.Context, userID, runID string) (*models.MigrationRun, error) {
	run, err := s.migrations.GetRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run == nil {
		return nil, fmt.Errorf("migration run not found")
	}
	if err := s.requireOrg(ctx, userID, run.OrganizationID, false); err != nil {
		return nil, err
	}
	return run, nil
}

func (s *Service) targetServer(ctx context.Context, orgID, serverID string) (*models.Server, dockerprobe.Runner, error) {
	if serverID == "" {
		return nil, s.targets(nil), nil
	}
	server, err := s.servers.GetByID(ctx, serverID)
	if err != nil {
		return nil, nil, err
	}
	if server == nil {
		return nil, nil, fmt.Errorf("target server not found")
	}
	if server.IsControlPlane {
		return nil, nil, fmt.Errorf("control plane cannot be a migration target")
	}
	if server.Provider != "" {
		return nil, nil, fmt.Errorf("managed servers join migration targets through their ssh record once provisioned")
	}
	return server, s.targets(server), nil
}

func decodeSelection(raw string) models.MigrationSelection {
	var selection models.MigrationSelection
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &selection)
	}
	if selection.Overrides == nil {
		selection.Overrides = map[string]string{}
	}
	if selection.Decisions == nil {
		selection.Decisions = map[string]string{}
	}
	return selection
}

func encodeSelection(selection models.MigrationSelection) string {
	raw, _ := json.Marshal(selection)
	return string(raw)
}

func decodeProgress(raw string) []models.MigrationServiceProgress {
	var progress []models.MigrationServiceProgress
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &progress)
	}
	return progress
}

func encodeProgress(progress []models.MigrationServiceProgress) string {
	raw, _ := json.Marshal(progress)
	return string(raw)
}

func decodePrompt(raw string) *models.MigrationPrompt {
	if raw == "" {
		return nil
	}
	var prompt models.MigrationPrompt
	if err := json.Unmarshal([]byte(raw), &prompt); err != nil {
		return nil
	}
	return &prompt
}

func encodePrompt(prompt *models.MigrationPrompt) string {
	if prompt == nil {
		return ""
	}
	raw, _ := json.Marshal(prompt)
	return string(raw)
}

func newConfirmationToken() (token, hash string, err error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", "", fmt.Errorf("generate confirmation: %w", err)
	}
	token = hex.EncodeToString(secret)
	sum := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(sum[:]), nil
}

func verifyConfirmation(hash, token string) bool {
	sum := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare([]byte(hash), []byte(hex.EncodeToString(sum[:]))) == 1
}

func newPrompt(kind models.MigrationPromptKind, subject, detail string, options []models.MigrationPromptOption) *models.MigrationPrompt {
	return &models.MigrationPrompt{
		ID:        uuid.NewString(),
		Kind:      kind,
		Subject:   subject,
		Detail:    detail,
		Options:   options,
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
	}
}
