package attention

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/hetzner"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/migrations"
)

type WorkloadDeployer interface {
	DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error)
	RestartAppService(ctx context.Context, app *models.AppService) error
}

type BackupRunner interface {
	ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*models.BackupRecord, error)
	TriggerBackup(ctx context.Context, configID string) (*models.BackupRecord, error)
}

type HealthReader interface {
	GetHealth(ctx context.Context, containerIDOrName string) (AttentionHealth, error)
}

type AttentionHealth struct {
	Status  string
	Running bool
}

type Service struct {
	issues     repositories.AttentionRepository
	deploys    repositories.DeploymentRepository
	apps       repositories.AppServiceRepository
	projects   repositories.ProjectRepository
	users      *repositories.UserRepo
	orgs       repositories.OrganizationRepository
	backups    BackupRunner
	migrations *migrations.Service
	managed    repositories.ManagedRepository
	servers    repositories.ServerRepository
	health     HealthReader
	deployer   WorkloadDeployer

	mu        sync.Mutex
	evaluated map[string]time.Time
}

func NewService(
	issues repositories.AttentionRepository,
	deploys repositories.DeploymentRepository,
	apps repositories.AppServiceRepository,
	projects repositories.ProjectRepository,
	users *repositories.UserRepo,
	orgs repositories.OrganizationRepository,
	backups BackupRunner,
	migrationService *migrations.Service,
	managed repositories.ManagedRepository,
	servers repositories.ServerRepository,
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

func (s *Service) requireOrg(ctx context.Context, userID, orgID string) error {
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
	return nil
}

func (s *Service) ListIssues(ctx context.Context, userID, orgID, status string) ([]*models.AttentionIssue, error) {
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

func (s *Service) stale(orgID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	last, ok := s.evaluated[orgID]
	return !ok || time.Since(last) > 5*time.Minute
}

func (s *Service) markEvaluated(orgID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evaluated[orgID] = time.Now()
}

func (s *Service) Evaluate(ctx context.Context, userID, orgID string) error {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return err
	}
	var live []string
	scanners := []func(context.Context, string, string) ([]*models.AttentionIssue, error){
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

func (s *Service) Acknowledge(ctx context.Context, userID, orgID, issueID string) error {
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
	if issue.Status != models.AttentionStatusOpen {
		return fmt.Errorf("only open issues can be acknowledged")
	}
	return s.issues.SetStatus(ctx, issueID, orgID, models.AttentionStatusAcked)
}

func (s *Service) Resolve(ctx context.Context, userID, orgID, issueID string) error {
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
	return s.issues.SetStatus(ctx, issueID, orgID, models.AttentionStatusResolved)
}

func (s *Service) scanDeployments(ctx context.Context, _ string, orgID string) ([]*models.AttentionIssue, error) {
	records, _, err := s.deploys.ListByOrganization(ctx, models.DeploymentListFilter{OrganizationID: orgID, Limit: 200})
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var issues []*models.AttentionIssue
	for _, record := range records {
		if record.Status != models.DeploymentStatusFailed {
			continue
		}
		if record.CreatedAt.Before(time.Now().UTC().Add(-24 * time.Hour)) {
			continue
		}
		if seen[record.ServiceID] {
			continue
		}
		seen[record.ServiceID] = true
		params, _ := json.Marshal(map[string]string{"serviceId": record.ServiceID, "projectId": record.ProjectID})
		issues = append(issues, &models.AttentionIssue{
			Kind: "deployment-failed", Subject: record.ServiceID,
			ProjectID: record.ProjectID, ServiceID: record.ServiceID,
			Severity: models.AttentionSeverityCritical,
			Title:    fmt.Sprintf("Deployments failing for %s", record.ServiceName),
			Detail:   fmt.Sprintf("Latest deployment %s failed at %s", record.ID, record.CreatedAt.Format(time.RFC3339)),
			Remediation: "Inspect the deployment logs, fix the build or config, then redeploy.",
			Action:      "redeploy", ActionParams: string(params),
		})
	}
	return issues, nil
}

func (s *Service) scanWorkloads(ctx context.Context, _ string, orgID string) ([]*models.AttentionIssue, error) {
	if s.health == nil {
		return nil, nil
	}
	projects, _, err := s.projects.ListByOrganization(ctx, orgID, 200, 0)
	if err != nil {
		return nil, err
	}
	var issues []*models.AttentionIssue
	for _, project := range projects {
		services, err := s.apps.ListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		for _, service := range services {
			if service.ContainerID == "" {
				continue
			}
			health, err := s.health.GetHealth(ctx, service.ContainerID)
			if err != nil || health.Running {
				continue
			}
			params, _ := json.Marshal(map[string]string{"serviceId": service.ID})
			issues = append(issues, &models.AttentionIssue{
				Kind: "service-unhealthy", Subject: service.ID,
				ProjectID: project.ID, ServiceID: service.ID,
				Severity: models.AttentionSeverityWarning,
				Title:    fmt.Sprintf("%s is not running", service.Name),
				Detail:   fmt.Sprintf("Container %s reports %s", service.ContainerID, health.Status),
				Remediation: "Check service logs for crash loops, then restart the workload.",
				Action:      "restart-service", ActionParams: string(params),
			})
		}
	}
	return issues, nil
}

func (s *Service) scanBackups(ctx context.Context, _ string, orgID string) ([]*models.AttentionIssue, error) {
	if s.backups == nil {
		return nil, nil
	}
	records, err := s.backups.ListRecordsByConfigs(ctx, nil, 100)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var issues []*models.AttentionIssue
	for _, record := range records {
		if string(record.Status) != "failed" {
			continue
		}
		if seen[record.BackupConfigID] {
			continue
		}
		seen[record.BackupConfigID] = true
		params, _ := json.Marshal(map[string]string{"configId": record.BackupConfigID})
		issues = append(issues, &models.AttentionIssue{
			Kind: "backup-failed", Subject: record.BackupConfigID,
			Severity: models.AttentionSeverityWarning,
			Title:    "Backup run failed",
			Detail:   fmt.Sprintf("Backup %s failed: %s", record.BackupConfigID, record.Logs),
			Remediation: "Verify the database and destination, then retry the backup.",
			Action:      "retry-backup", ActionParams: string(params),
		})
	}
	return issues, nil
}

func (s *Service) scanMigrations(ctx context.Context, userID, orgID string) ([]*models.AttentionIssue, error) {
	if s.migrations == nil {
		return nil, nil
	}
	runs, err := s.migrations.ListRuns(ctx, userID, orgID, "")
	if err != nil {
		return nil, err
	}
	var issues []*models.AttentionIssue
	for _, run := range runs {
		if run.Status != models.MigrationStatusFailed {
			continue
		}
		params, _ := json.Marshal(map[string]string{"runId": run.ID})
		issues = append(issues, &models.AttentionIssue{
			Kind: "migration-failed", Subject: run.ID,
			ProjectID: run.ProjectID,
			Severity: models.AttentionSeverityWarning,
			Title:    "Migration run failed",
			Detail:   fmt.Sprintf("Run %s failed in %s: %s", run.ID, run.Phase, run.Error),
			Remediation: "Review the run logs, then resume with overrides or skips.",
			Action:      "resume-migration", ActionParams: string(params),
		})
	}
	return issues, nil
}

func (s *Service) scanQuota(ctx context.Context, _ string, orgID string) ([]*models.AttentionIssue, error) {
	if s.managed == nil {
		return nil, nil
	}
	quota, err := s.managed.GetQuota(ctx, orgID)
	if err != nil {
		return nil, err
	}
	ids, err := s.managed.ListOrgServerIDs(ctx, orgID)
	if err != nil {
		return nil, err
	}
	usedServers := 0
	usedMemory := 0
	for _, id := range ids {
		server, err := s.servers.GetByID(ctx, id)
		if err != nil || server == nil {
			continue
		}
		usedServers++
		if tier, err := hetzner.TierByName(server.ServerType); err == nil {
			usedMemory += tier.MemoryGB
		}
	}
	var issues []*models.AttentionIssue
	if quota.MaxServers > 0 && float64(usedServers)/float64(quota.MaxServers) >= 0.8 {
		issues = append(issues, &models.AttentionIssue{
			Kind: "quota-pressure", Subject: "managed-servers",
			Severity: models.AttentionSeverityWarning,
			Title:    "Managed server quota nearly exhausted",
			Detail:   fmt.Sprintf("%d of %d managed servers in use", usedServers, quota.MaxServers),
			Remediation: "Raise the organization quota or delete unused managed servers.",
		})
	}
	if quota.MaxMemoryGB > 0 && float64(usedMemory)/float64(quota.MaxMemoryGB) >= 0.8 {
		issues = append(issues, &models.AttentionIssue{
			Kind: "quota-pressure", Subject: "managed-memory",
			Severity: models.AttentionSeverityWarning,
			Title:    "Managed memory quota nearly exhausted",
			Detail:   fmt.Sprintf("%d of %d GB in use", usedMemory, quota.MaxMemoryGB),
			Remediation: "Raise the organization quota or downsize managed servers.",
		})
	}
	return issues, nil
}
