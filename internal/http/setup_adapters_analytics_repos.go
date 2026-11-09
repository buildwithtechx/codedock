package http

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

type analyticsDeployRepoAdapter struct {
	repo deploymentsmod.DeploymentRepository
}

func (a *analyticsDeployRepoAdapter) ListByOrganization(ctx context.Context, filter monitoring.DeploymentListFilter) ([]monitoring.DeploymentRecord, int, error) {
	if a.repo == nil {
		return nil, 0, nil
	}
	f := deploymentsmod.DeploymentListFilter{
		OrganizationID: filter.OrganizationID,
		ProjectID:      filter.ProjectID,
		ServiceID:      filter.ServiceID,
		Limit:          filter.Limit,
		Offset:         filter.Offset,
	}
	items, count, err := a.repo.ListByOrganization(ctx, f)
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

type analyticsProjectRepoAdapter struct {
	repo projectsmod.ProjectRepository
}

func (a *analyticsProjectRepoAdapter) GetByOrganization(ctx context.Context, projectID, orgID string) (*monitoring.Project, error) {
	if a.repo == nil {
		return nil, nil
	}
	p, err := a.repo.GetByOrganization(ctx, projectID, orgID)
	if err != nil || p == nil {
		return nil, err
	}
	return &monitoring.Project{
		ID:             p.ID,
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
	}, nil
}

func (a *analyticsProjectRepoAdapter) ListByOrganization(ctx context.Context, orgID, search string, limit, offset int) ([]monitoring.Project, int, error) {
	if a.repo == nil {
		return nil, 0, nil
	}
	list, total, err := a.repo.ListByOrganization(ctx, orgID, search, limit, offset)
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

type analyticsAppRepoAdapter struct {
	repo appsmod.AppServiceRepository
}

func (a *analyticsAppRepoAdapter) GetByID(ctx context.Context, id string) (*monitoring.AppService, error) {
	if a.repo == nil {
		return nil, nil
	}
	app, err := a.repo.GetByID(ctx, id)
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

func (a *analyticsAppRepoAdapter) ListByProject(ctx context.Context, projectID string) ([]monitoring.AppService, error) {
	if a.repo == nil {
		return nil, nil
	}
	apps, err := a.repo.ListByProject(ctx, projectID)
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

type analyticsDomainListerAdapter struct {
	repo domainsmod.DomainRepository
}

func (a *analyticsDomainListerAdapter) ListAll(ctx context.Context) ([]monitoring.DomainConfig, error) {
	if a.repo == nil {
		return nil, nil
	}
	list, err := a.repo.ListAll(ctx)
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

type analyticsUserRepoAdapter struct {
	repo users.UserRepository
}

func (a *analyticsUserRepoAdapter) GetUserByID(ctx context.Context, id string) (*monitoring.User, error) {
	if a.repo == nil {
		return nil, nil
	}
	u, err := a.repo.GetUserByID(ctx, id)
	if err != nil || u == nil {
		return nil, err
	}
	return toMonitoringUser(u), nil
}

type analyticsOrgRepoAdapter struct {
	repo organizations.OrganizationRepository
}

func (a *analyticsOrgRepoAdapter) GetMember(ctx context.Context, orgID, userID string) (*monitoring.OrgMember, error) {
	if a.repo == nil {
		return nil, nil
	}
	m, err := a.repo.GetMember(ctx, orgID, userID)
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
