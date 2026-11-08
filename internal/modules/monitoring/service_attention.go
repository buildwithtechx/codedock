package monitoring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"

	"codedock/internal/engine/hetzner"
	"codedock/internal/modules/backups"
)

var ErrClearHistoryForbidden = errors.New("an organization administrator is required to clear monitoring history")

type WorkloadDeployer interface {
	DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error)
	RestartAppService(ctx context.Context, app *AppService) error
}

type BackupRunner interface {
	ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*backups.BackupRecord, error)
	TriggerBackup(ctx context.Context, configID string) (*backups.BackupRecord, error)
}

type HealthReader interface {
	GetHealth(ctx context.Context, containerIDOrName string) (AttentionHealth, error)
}

type AttentionHealth struct {
	Status  string
	Running bool
}

type MigrationRunner interface {
	ListRuns(ctx context.Context, userID, orgID, status string) (any, error)
	ResumeRun(ctx context.Context, userID, runID string, req any) (any, error)
}

type AttentionManager struct {
	issues     AttentionRepository
	deploys    DeploymentRepository
	apps       AppServiceRepository
	projects   ProjectRepository
	users      *UserRepo
	orgs       OrganizationRepository
	backups    BackupRunner
	migrations MigrationRunner
	managed    ManagedRepository
	servers    ServerRepository
	health     HealthReader
	deployer   WorkloadDeployer

	mu        sync.Mutex
	evaluated map[string]time.Time
}

func NewAttentionManager(
	issues AttentionRepository,
	deploys DeploymentRepository,
	apps AppServiceRepository,
	projects ProjectRepository,
	users *UserRepo,
	orgs OrganizationRepository,
	backups BackupRunner,
	migrationService MigrationRunner,
	managed ManagedRepository,
	servers ServerRepository,
	health HealthReader,
	deployer WorkloadDeployer,
) *Service {
	return &Service{
		issues: issues, deploys: deploys, apps: apps, projects: projects,
		users: users, orgs: orgs, backups: backups, migrations: migrationService,
		managed: managed, servers: servers, health: health, deployer: deployer,
		evaluated: map[string]time.Time{},
	}
}

func (s *AttentionManager) requireOrg(ctx context.Context, userID, orgID string) error {
	if orgID == "" {
		return fmt.Errorf("organization id is required")
	}
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %w", err)
	}
	if user != nil && (user.Role == UserRoleOwner || user.Role == UserRoleAdmin) {
		return nil
	}
	member, err := s.orgs.GetMember(ctx, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to load organization membership: %w", err)
	}
	if member == nil {
		return fmt.Errorf("organization membership is required")
	}
	return nil
}

func (s *AttentionManager) ListIssues(ctx context.Context, userID, orgID, status string) ([]*AttentionIssue, error) {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return nil, err
	}
	if s.stale(orgID) {
		if err := s.Evaluate(ctx, userID, orgID); err != nil {
			return nil, err
		}
	}
	return s.issues.List(ctx, orgID, status, 100)
}

func (s *AttentionManager) stale(orgID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	last, ok := s.evaluated[orgID]
	return !ok || time.Since(last) > 5*time.Minute
}

func (s *AttentionManager) markEvaluated(orgID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evaluated[orgID] = time.Now()
}

func (s *AttentionManager) Evaluate(ctx context.Context, userID, orgID string) error {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return err
	}
	var live []string
	scanners := []func(context.Context, string, string) ([]*AttentionIssue, error){
		s.scanDeployments, s.scanWorkloads, s.scanBackups, s.scanMigrations, s.scanQuota,
	}
	for _, scanner := range scanners {
		issues, err := scanner(ctx, userID, orgID)
		if err != nil {
			return err
		}
		for _, issue := range issues {
			issue.ID = uuid.NewString()
			issue.OrganizationID = orgID
			if err := s.issues.Upsert(ctx, issue); err != nil {
				return err
			}
			live = append(live, issue.Kind+":"+issue.Subject)
		}
	}
	if err := s.issues.ResolveMissing(ctx, orgID, live); err != nil {
		return err
	}
	s.markEvaluated(orgID)
	return nil
}

func (s *AttentionManager) Acknowledge(ctx context.Context, userID, orgID, issueID string) error {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return err
	}
	issue, err := s.issues.Get(ctx, issueID)
	if err != nil {
		return err
	}
	if issue == nil || issue.OrganizationID != orgID {
		return fmt.Errorf("attention issue not found")
	}
	if issue.Status != AttentionStatusOpen {
		return fmt.Errorf("only open issues can be acknowledged")
	}
	return s.issues.SetStatus(ctx, issueID, orgID, AttentionStatusAcked)
}

func (s *AttentionManager) Resolve(ctx context.Context, userID, orgID, issueID string) error {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return err
	}
	issue, err := s.issues.Get(ctx, issueID)
	if err != nil {
		return err
	}
	if issue == nil || issue.OrganizationID != orgID {
		return fmt.Errorf("attention issue not found")
	}
	return s.issues.SetStatus(ctx, issueID, orgID, AttentionStatusResolved)
}

func (s *AttentionManager) CanClearHistory(ctx context.Context, userID, orgID string) (bool, error) {
	if orgID == "" {
		return false, nil
	}
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return false, nil
	}
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to load user: %w", err)
	}
	if user != nil && (user.Role == UserRoleOwner || user.Role == UserRoleAdmin) {
		return true, nil
	}
	member, err := s.orgs.GetMember(ctx, orgID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to load organization membership: %w", err)
	}
	if member == nil {
		return false, nil
	}
	return member.Permission == MemberPermissionAdmin || member.Permission == MemberPermissionOwner, nil
}

func (s *AttentionManager) ClearHistory(ctx context.Context, userID, orgID string) (int64, error) {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return 0, err
	}
	allowed, err := s.CanClearHistory(ctx, userID, orgID)
	if err != nil {
		return 0, err
	}
	if !allowed {
		return 0, ErrClearHistoryForbidden
	}
	return s.issues.ClearResolved(ctx, orgID)
}
