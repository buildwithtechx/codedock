package adapters

import (
	"context"

	appsmod "codedock/internal/modules/apps"
	deploymentsmod "codedock/internal/modules/deployments"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/users"
)

type AnalyticsDeployRepoAdapter struct {
	Repo deploymentsmod.DeploymentRepository
}

func (a *AnalyticsDeployRepoAdapter) ListByOrganization(ctx context.Context, filter monitoring.DeploymentListFilter) ([]monitoring.DeploymentRecord, int, error) {
	if a.Repo == nil {
		return nil, 0, nil
	}
	f := deploymentsmod.DeploymentListFilter{
		OrganizationID: filter.OrganizationID,
		ProjectID:      filter.ProjectID,
		ServiceID:      filter.ServiceID,
		Limit:          filter.Limit,
		Offset:         filter.Offset,
	}
	items, count, err := a.Repo.ListByOrganization(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	res := make([]monitoring.DeploymentRecord, len(items))
	for i, it := range items {
		res[i] = monitoring.DeploymentRecord{
			ID:          it.ID,
			ProjectID:   it.ProjectID,
			ServiceID:   it.ServiceID,
			ServiceName: it.ServiceName,
			Status:      string(it.Status),
			CreatedAt:   it.CreatedAt,
			FinishedAt:  it.FinishedAt,
		}
	}
	return res, count, nil
}

type AnalyticsProjectRepoAdapter struct {
	Repo projectsmod.ProjectRepository
}

func (a *AnalyticsProjectRepoAdapter) GetByOrganization(ctx context.Context, projectID, orgID string) (*monitoring.Project, error) {
	if a.Repo == nil {
		return nil, nil
	}
	p, err := a.Repo.GetByOrganization(ctx, projectID, orgID)
	if err != nil || p == nil {
		return nil, err
	}
	return &monitoring.Project{
		ID:             p.ID,
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
	}, nil
}

func (a *AnalyticsProjectRepoAdapter) ListByOrganization(ctx context.Context, orgID, search string, limit, offset int) ([]monitoring.Project, int, error) {
	if a.Repo == nil {
		return nil, 0, nil
	}
	list, total, err := a.Repo.ListByOrganization(ctx, orgID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	res := make([]monitoring.Project, len(list))
	for i, p := range list {
		res[i] = monitoring.Project{
			ID:             p.ID,
			OrganizationID: p.OrganizationID,
			Name:           p.Name,
		}
	}
	return res, total, nil
}

type AnalyticsAppRepoAdapter struct {
	Repo appsmod.AppServiceRepository
}

func (a *AnalyticsAppRepoAdapter) GetByID(ctx context.Context, id string) (*monitoring.AppService, error) {
	if a.Repo == nil {
		return nil, nil
	}
	app, err := a.Repo.GetByID(ctx, id)
	if err != nil || app == nil {
		return nil, err
	}
	return &monitoring.AppService{
		ID:          app.ID,
		ProjectID:   app.ProjectID,
		Name:        app.Name,
		Status:      string(app.Status),
		Replicas:    app.Replicas,
		ContainerID: app.ContainerID,
	}, nil
}

func (a *AnalyticsAppRepoAdapter) ListByProject(ctx context.Context, projectID string) ([]monitoring.AppService, error) {
	if a.Repo == nil {
		return nil, nil
	}
	apps, err := a.Repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]monitoring.AppService, len(apps))
	for i, app := range apps {
		res[i] = monitoring.AppService{
			ID:          app.ID,
			ProjectID:   app.ProjectID,
			Name:        app.Name,
			Status:      string(app.Status),
			Replicas:    app.Replicas,
			ContainerID: app.ContainerID,
		}
	}
	return res, nil
}

type AnalyticsDomainListerAdapter struct {
	Repo domainsmod.DomainRepository
}

func (a *AnalyticsDomainListerAdapter) ListAll(ctx context.Context) ([]monitoring.DomainConfig, error) {
	if a.Repo == nil {
		return nil, nil
	}
	list, err := a.Repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]monitoring.DomainConfig, len(list))
	for i, d := range list {
		res[i] = monitoring.DomainConfig{
			ID:            d.ID,
			ServiceID:     d.ServiceID,
			DomainName:    d.DomainName,
			SSLCertStatus: d.SSLCertStatus,
			PathPrefix:    d.PathPrefix,
		}
	}
	return res, nil
}

type AnalyticsUserRepoAdapter struct {
	Repo users.UserRepository
}

func (a *AnalyticsUserRepoAdapter) GetUserByID(ctx context.Context, id string) (*monitoring.User, error) {
	if a.Repo == nil {
		return nil, nil
	}
	u, err := a.Repo.GetUserByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	return toMonitoringUser(u), nil
}

type AnalyticsOrgRepoAdapter struct {
	Repo organizations.OrganizationRepository
}

func (a *AnalyticsOrgRepoAdapter) GetMember(ctx context.Context, orgID, userID string) (*monitoring.OrgMember, error) {
	if a.Repo == nil {
		return nil, nil
	}
	m, err := a.Repo.GetMember(ctx, orgID, userID)
	if err != nil || m == nil {
		return nil, err
	}
	return &monitoring.OrgMember{
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		Role:           string(m.Permission),
		Permission:     monitoring.MemberPermission(m.Permission),
	}, nil
}
