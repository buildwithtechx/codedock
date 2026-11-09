package adapters

import (
	"context"
	"errors"
	"fmt"
	"io"

	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	deploymentsmod "codedock/internal/modules/deployments"
	dnsproviders "codedock/internal/modules/dns/providers"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/monitoring"
	settingsmodule "codedock/internal/modules/settings"
)

type DNSSettingsAdapter struct {
	Repo settingsmodule.Repository
}

func (a *DNSSettingsAdapter) GetDNSConfig(ctx context.Context) (*dnsproviders.DNSConfig, error) {
	if a.Repo == nil {
		return nil, nil
	}
	cfg, err := a.Repo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return nil, err
	}
	return &dnsproviders.DNSConfig{
		PublicIPv4:         cfg.PublicIPv4,
		CloudflareAPIToken: cfg.CloudflareAPIToken,
		NamecheapAPIKey:    cfg.NamecheapAPIKey,
		NamecheapAPIUser:   cfg.NamecheapAPIUser,
		NamecheapClientIP:  cfg.NamecheapClientIP,
		SpaceshipAPIKey:    cfg.SpaceshipAPIKey,
		SpaceshipAPISecret: cfg.SpaceshipAPISecret,
	}, nil
}

type GitServiceLookupAdapter struct {
	AppRepo appsmod.AppServiceRepository
}

func toGitServiceInfo(app *appsmod.AppService) *gitmod.ServiceInfo {
	if app == nil {
		return nil
	}
	return &gitmod.ServiceInfo{
		ID:               app.ID,
		ProjectID:        app.ProjectID,
		EnvironmentID:    app.EnvironmentID,
		Name:             app.Name,
		RepositoryURL:    app.RepositoryURL,
		Branch:           app.Branch,
		RootDirectory:    app.RootDirectory,
		BuildCommand:     app.BuildCommand,
		RuntimeMode:      string(app.RuntimeMode),
		EnablePRPreviews: app.EnablePRPreviews,
		DeployToken:      app.DeployToken,
	}
}

func (a *GitServiceLookupAdapter) GetService(ctx context.Context, serviceID string) (*gitmod.ServiceInfo, error) {
	app, err := a.AppRepo.GetByID(ctx, serviceID)
	if err != nil || app == nil {
		return nil, err
	}
	return toGitServiceInfo(app), nil
}

func (a *GitServiceLookupAdapter) ListServices(ctx context.Context) ([]*gitmod.ServiceInfo, error) {
	apps, err := a.AppRepo.ListAll(ctx, "")
	if err != nil {
		return nil, err
	}
	res := make([]*gitmod.ServiceInfo, len(apps))
	for i, app := range apps {
		res[i] = toGitServiceInfo(app)
	}
	return res, nil
}

type DomainServiceLookupAdapter struct {
	AppRepo appsmod.AppServiceRepository
}

func (d *DomainServiceLookupAdapter) GetProjectID(ctx context.Context, serviceID string) (string, error) {
	app, err := d.AppRepo.GetByID(ctx, serviceID)
	if err != nil || app == nil {
		return "", err
	}
	return app.ProjectID, nil
}

func (d *DomainServiceLookupAdapter) ListServiceIDsByProject(ctx context.Context, projectID string) ([]string, error) {
	apps, err := d.AppRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(apps))
	for i, a := range apps {
		ids[i] = a.ID
	}
	return ids, nil
}

type ServerSettingsAdapter struct {
	Repo settingsmodule.Repository
}

func (s *ServerSettingsAdapter) GetPublicIPv4(ctx context.Context) (string, error) {
	settings, err := s.Repo.GetServerSettings(ctx)
	if err != nil || settings == nil {
		return "", err
	}
	return settings.PublicIPv4, nil
}

type AuditDomainAdapter struct {
	Audit *audit.AuditService
}

func (a *AuditDomainAdapter) Log(ctx context.Context, userID, action, resource, ip string, details map[string]string) {
	if a.Audit != nil {
		a.Audit.LogAction(ctx, audit.AuditActionOpts{UserID: userID, Action: action, Resource: resource, IPAddress: ip, Details: details})
	}
}

type ServiceCheckerAdapter struct {
	AppService *appsmod.Service
}

func (s *ServiceCheckerAdapter) CheckServiceExists(ctx context.Context, serviceID string) bool {
	svc, err := s.AppService.GetAppService(ctx, serviceID)
	return err == nil && svc != nil
}

type MigrationDeployerAdapter struct {
	DeploymentService *deploymentsmod.DeploymentService
}

func (m *MigrationDeployerAdapter) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	return m.DeploymentService.TriggerDeployment(ctx, appID, "", "", "", "migration", "")
}

func (m *MigrationDeployerAdapter) StopAppService(ctx context.Context, app *appsmod.AppService) error {
	return nil
}

func (m *MigrationDeployerAdapter) RemoveAppService(ctx context.Context, app *appsmod.AppService) error {
	return nil
}

type MigrationGitAdapter struct {
	GitService *gitmod.ProviderService
}

func (g *MigrationGitAdapter) CloneOrPullAppRepository(ctx context.Context, app *appsmod.AppService, targetDir string, logWriter io.Writer) error {
	if app == nil {
		return errors.New("app required")
	}
	return g.GitService.SyncCodebase(ctx, app.RepositoryURL, app.Branch, "", targetDir, logWriter)
}

type AttentionDeployerAdapter struct {
	DeploymentService *deploymentsmod.DeploymentService
	AppService        *appsmod.Service
}

func (a *AttentionDeployerAdapter) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	app, err := a.AppService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return "", fmt.Errorf("app not found: %w", err)
	}
	return a.DeploymentService.TriggerDeployment(ctx, app.ID, app.EnvironmentID, app.ProjectID, app.Branch, "attention", "")
}

func (a *AttentionDeployerAdapter) RestartAppService(ctx context.Context, app *monitoring.AppService) error {
	if app == nil {
		return errors.New("app required")
	}
	appSvc, err := a.AppService.GetAppService(ctx, app.ID)
	if err != nil || appSvc == nil {
		return fmt.Errorf("app not found")
	}
	_, err = a.DeploymentService.TriggerDeployment(ctx, appSvc.ID, appSvc.EnvironmentID, appSvc.ProjectID, appSvc.Branch, "restart", "")
	return err
}
