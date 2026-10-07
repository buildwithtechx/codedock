package analytics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type DomainLister interface {
	ListAll(ctx context.Context) ([]models.DomainConfig, error)
}

type Service struct {
	traffic  repositories.TrafficRepository
	deploys  repositories.DeploymentRepository
	projects repositories.ProjectRepository
	apps     repositories.AppServiceRepository
	domains  DomainLister
	users    *repositories.UserRepo
	orgs     repositories.OrganizationRepository
	issues   repositories.AttentionRepository
	usage    *UsageReader

	attrMu   sync.Mutex
	attrMap  map[string]string
	attrAt   time.Time
}

func NewService(
	traffic repositories.TrafficRepository,
	deploys repositories.DeploymentRepository,
	projects repositories.ProjectRepository,
	apps repositories.AppServiceRepository,
	domains DomainLister,
	users *repositories.UserRepo,
	orgs repositories.OrganizationRepository,
	issues repositories.AttentionRepository,
	usage *UsageReader,
) *Service {
	return &Service{traffic: traffic, deploys: deploys, projects: projects, apps: apps, domains: domains, users: users, orgs: orgs, issues: issues, usage: usage, attrMap: map[string]string{}}
}

func (s *Service) AttributeDomain(host string) string {
	canonical := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(host), "."))
	if canonical == "" {
		return ""
	}
	s.attrMu.Lock()
	projectID, ok := s.attrMap[canonical]
	fresh := time.Since(s.attrAt) < 2*time.Minute
	s.attrMu.Unlock()
	if ok && fresh {
		return projectID
	}
	if !fresh {
		s.refreshAttribution()
	}
	s.attrMu.Lock()
	defer s.attrMu.Unlock()
	return s.attrMap[canonical]
}

func (s *Service) refreshAttribution() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var domains []models.DomainConfig
	if s.domains != nil {
		listed, err := s.domains.ListAll(ctx)
		if err != nil {
			return
		}
		domains = listed
	}
	mapped := map[string]string{}
	for _, domain := range domains {
		if domain.ServiceID == "" {
			continue
		}
		service, err := s.apps.GetByID(ctx, domain.ServiceID)
		if err != nil || service == nil {
			continue
		}
		mapped[strings.ToLower(domain.DomainName)] = service.ProjectID
	}
	s.attrMu.Lock()
	s.attrMap = mapped
	s.attrAt = time.Now()
	s.attrMu.Unlock()
}

func (s *Service) Usage() *UsageReader {
	return s.usage
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

func (s *Service) requireProject(ctx context.Context, userID, orgID, projectID string) error {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return err
	}
	project, err := s.projects.GetByOrganization(ctx, projectID, orgID)
	if err != nil {
		return err
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}
	return nil
}

func (s *Service) Summary(ctx context.Context, userID, orgID, projectID, domain, from, to string) (*models.AnalyticsSummary, error) {
	if err := s.requireProject(ctx, userID, orgID, projectID); err != nil {
		return nil, err
	}
	return s.traffic.Summary(ctx, projectID, domain, from, to)
}

func (s *Service) Overview(ctx context.Context, userID, orgID, projectID, domain, from, to, step string) (*models.AnalyticsOverview, error) {
	if err := s.requireProject(ctx, userID, orgID, projectID); err != nil {
		return nil, err
	}
	return s.traffic.Overview(ctx, projectID, domain, from, to, step)
}

func (s *Service) Geo(ctx context.Context, userID, orgID, projectID, from, to string) (*models.AnalyticsGeo, error) {
	if err := s.requireProject(ctx, userID, orgID, projectID); err != nil {
		return nil, err
	}
	return s.traffic.Geo(ctx, projectID, from, to)
}

func (s *Service) PathsEnabled(ctx context.Context, userID, orgID, projectID string) (bool, error) {
	if err := s.requireProject(ctx, userID, orgID, projectID); err != nil {
		return false, err
	}
	return s.traffic.PathsEnabled(ctx, projectID)
}

func (s *Service) SetPathsEnabled(ctx context.Context, userID, orgID, projectID string, enabled bool) error {
	if err := s.requireProject(ctx, userID, orgID, projectID); err != nil {
		return err
	}
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || (user.Role != models.UserRoleOwner && user.Role != models.UserRoleAdmin) {
		member, err := s.orgs.GetMember(ctx, orgID, userID)
		if err != nil {
			return err
		}
		if member.Permission != models.MemberPermissionAdmin && member.Permission != models.MemberPermissionOwner {
			return fmt.Errorf("organization admin permission is required")
		}
	}
	return s.traffic.SetPathsEnabled(ctx, projectID, enabled)
}

func (s *Service) Periods() []string {
	return []string{"1h", "24h", "7d", "30d"}
}

func (s *Service) DeploymentStats(ctx context.Context, userID, orgID, projectID string, days int) (*models.AnalyticsDeploymentStats, error) {
	if err := s.requireProject(ctx, userID, orgID, projectID); err != nil {
		return nil, err
	}
	if days <= 0 {
		days = 30
	}
	since := time.Now().UTC().AddDate(0, 0, -days)
	stats := &models.AnalyticsDeploymentStats{ProjectID: projectID}
	records, _, err := s.deploys.ListByOrganization(ctx, models.DeploymentListFilter{OrganizationID: orgID, ProjectID: projectID, Limit: 1000})
	if err != nil {
		return nil, err
	}
	var durations float64
	var finished int
	for _, record := range records {
		if record.CreatedAt.Before(since) {
			continue
		}
		stats.Total++
		switch record.Status {
		case models.DeploymentStatusActive, models.DeploymentStatusReady:
			stats.Succeeded++
		case models.DeploymentStatusFailed:
			stats.Failed++
		}
		if record.FinishedAt != nil {
			durations += record.FinishedAt.Sub(record.CreatedAt).Seconds()
			finished++
		}
	}
	if stats.Total > 0 {
		stats.SuccessRate = float64(stats.Succeeded) / float64(stats.Total)
	}
	if finished > 0 {
		stats.AvgDurationSeconds = durations / float64(finished)
	}
	return stats, nil
}

func (s *Service) Dashboard(ctx context.Context, userID, orgID string) (*models.AnalyticsDashboard, error) {
	if err := s.requireOrg(ctx, userID, orgID); err != nil {
		return nil, err
	}
	dashboard := &models.AnalyticsDashboard{OrganizationID: orgID}
	projects, _, err := s.projects.ListByOrganization(ctx, orgID, 200, 0)
	if err != nil {
		return nil, err
	}
	dayStart := time.Now().UTC().Add(-24 * time.Hour).Format("2006-01-02T15:04")
	now := time.Now().UTC().Format("2006-01-02T15:04")
	for _, project := range projects {
		summary, err := s.traffic.Summary(ctx, project.ID, "", dayStart, now)
		if err != nil {
			continue
		}
		dashboard.Requests24h += summary.Requests
	}
	if dashboard.Requests24h > 0 {
		var errors float64
		for _, project := range projects {
			summary, err := s.traffic.Summary(ctx, project.ID, "", dayStart, now)
			if err != nil {
				continue
			}
			errors += summary.ErrorRate * float64(summary.Requests)
		}
		dashboard.ErrorRate24h = errors / float64(dashboard.Requests24h)
	}
	stats, _, err := s.deploys.ListByOrganization(ctx, models.DeploymentListFilter{OrganizationID: orgID, Limit: 1000})
	if err == nil {
		week := time.Now().UTC().AddDate(0, 0, -7)
		for _, record := range stats {
			if record.CreatedAt.Before(week) {
				continue
			}
			dashboard.Deployments7d++
			if record.Status == models.DeploymentStatusActive || record.Status == models.DeploymentStatusReady {
				dashboard.DeploySuccess++
			}
		}
		if dashboard.Deployments7d > 0 {
			dashboard.DeploySuccess /= float64(dashboard.Deployments7d)
		}
	}
	if s.issues != nil {
		if open, err := s.issues.CountOpen(ctx, orgID); err == nil {
			dashboard.OpenIssues = open
		}
	}
	return dashboard, nil
}
