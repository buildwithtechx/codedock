package http

import (
	"context"
	"errors"
	"fmt"
	"io"

	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	deploymentsmod "codedock/internal/modules/deployments"
	dnsproviders "codedock/internal/modules/dns/providers"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/monitoring"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/platform/networking"
)

type dnsSettingsAdapter struct {
	repo settingsmodule.Repository
}

func (a *dnsSettingsAdapter) GetDNSConfig(ctx context.Context) (*dnsproviders.DNSConfig, error) {
	if a.repo == nil {
		return nil, nil
	}
	cfg, err := a.repo.GetServerSettings(ctx)
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

type gitServiceLookupAdapter struct {
	appRepo appsmod.AppServiceRepository
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

func (a *gitServiceLookupAdapter) GetService(ctx context.Context, serviceID string) (*gitmod.ServiceInfo, error) {
	app, err := a.appRepo.GetByID(ctx, serviceID)
	if err != nil || app == nil {
		return nil, err
	}
	return toGitServiceInfo(app), nil
}

func (a *gitServiceLookupAdapter) ListServices(ctx context.Context) ([]*gitmod.ServiceInfo, error) {
	apps, err := a.appRepo.ListAll(ctx, "")
	if err != nil {
		return nil, err
	}
	res := make([]*gitmod.ServiceInfo, len(apps))
	for i, app := range apps {
		res[i] = toGitServiceInfo(app)
	}
	return res, nil
}

type domainServiceLookupAdapter struct {
	appRepo appsmod.AppServiceRepository
}

func (d *domainServiceLookupAdapter) GetProjectID(ctx context.Context, serviceID string) (string, error) {
	app, err := d.appRepo.GetByID(ctx, serviceID)
	if err != nil || app == nil {
		return "", err
	}
	return app.ProjectID, nil
}

func (d *domainServiceLookupAdapter) ListServiceIDsByProject(ctx context.Context, projectID string) ([]string, error) {
	apps, err := d.appRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(apps))
	for i, a := range apps {
		ids[i] = a.ID
	}
	return ids, nil
}

type serverSettingsAdapter struct {
	repo settingsmodule.Repository
}

func (s *serverSettingsAdapter) GetPublicIPv4(ctx context.Context) (string, error) {
	settings, err := s.repo.GetServerSettings(ctx)
	if err != nil || settings == nil {
		return "", err
	}
	return settings.PublicIPv4, nil
}

type auditDomainAdapter struct {
	audit *audit.AuditService
}

func (a *auditDomainAdapter) Log(ctx context.Context, userID, action, resource, ip string, details map[string]string) {
	a.audit.LogAction(ctx, audit.AuditActionOpts{UserID: userID, Action: action, Resource: resource, IPAddress: ip, Details: details})
}

type serviceCheckerAdapter struct {
	appService *appsmod.Service
}

func (s *serviceCheckerAdapter) CheckServiceExists(ctx context.Context, serviceID string) bool {
	svc, err := s.appService.GetAppService(ctx, serviceID)
	return err == nil && svc != nil
}

type migrationDeployerAdapter struct {
	deploymentService *deploymentsmod.DeploymentService
}

func (m *migrationDeployerAdapter) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	return m.deploymentService.TriggerDeployment(ctx, appID, "", "", "", "migration", "")
}

func (m *migrationDeployerAdapter) StopAppService(ctx context.Context, app *appsmod.AppService) error {
	return nil
}

func (m *migrationDeployerAdapter) RemoveAppService(ctx context.Context, app *appsmod.AppService) error {
	return nil
}

type migrationGitAdapter struct {
	gitService *gitmod.ProviderService
}

func (g *migrationGitAdapter) CloneOrPullAppRepository(ctx context.Context, app *appsmod.AppService, targetDir string, logWriter io.Writer) error {
	if app == nil {
		return errors.New("app required")
	}
	return g.gitService.SyncCodebase(ctx, app.RepositoryURL, app.Branch, "", targetDir, logWriter)
}

type attentionDeployerAdapter struct {
	deploymentService *deploymentsmod.DeploymentService
	appService        *appsmod.Service
}

func (a *attentionDeployerAdapter) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	app, err := a.appService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return "", fmt.Errorf("app not found: %w", err)
	}
	return a.deploymentService.TriggerDeployment(ctx, app.ID, app.EnvironmentID, app.ProjectID, app.Branch, "attention", "")
}

func (a *attentionDeployerAdapter) RestartAppService(ctx context.Context, app *monitoring.AppService) error {
	if app == nil {
		return errors.New("app required")
	}
	appSvc, err := a.appService.GetAppService(ctx, app.ID)
	if err != nil || appSvc == nil {
		return fmt.Errorf("app not found")
	}
	_, err = a.deploymentService.TriggerDeployment(ctx, appSvc.ID, appSvc.EnvironmentID, appSvc.ProjectID, appSvc.Branch, "restart", "")
	return err
}

func configureDeploymentBindings(srv *Server, routeRuleRepo environments.RouteRuleRepository, domainRepo domainsmod.DomainRepository, routeProtectionRepo environments.RouteProtectionRepository) {
	if srv.deployer == nil {
		return
	}
	srv.deployer.ProtectionChecker = func(ctx context.Context, hostname string) bool {
		if routeProtectionRepo == nil {
			return false
		}
		rp, err := routeProtectionRepo.GetByHostname(ctx, hostname)
		return err == nil && rp != nil && rp.Enabled && rp.Secret != nil
	}
	srv.deployer.ScopedEnvProvider = func(ctx context.Context, project, environment string) (map[string]string, error) {
		return srv.serviceLinker.GetLinkedVariablesForEnvironment(ctx, project, environment)
	}
	srv.deployer.ScopedEnvInterpolator = func(ctx context.Context, project, environment string) (map[string]map[string]string, error) {
		return srv.serviceLinker.GetNamespacedVariablesForEnvironment(ctx, project, environment)
	}
	srv.deployer.EnvProvider = func(projectID string) (map[string]string, error) {
		return srv.serviceLinker.GetLinkedEnvironmentVariables(context.Background(), projectID)
	}
	srv.deployer.EnvInterpolator = func(projectID string) (map[string]map[string]string, error) {
		return srv.serviceLinker.GetNamespacedVariables(context.Background(), projectID)
	}
	srv.deployer.RouteRuleFetcher = func(ctx context.Context, serviceID, serviceName string) (map[string]string, error) {
		rules, err := routeRuleRepo.ListByService(ctx, serviceID)
		if err != nil {
			return nil, err
		}
		mRules := make([]*networking.RouteRule, len(rules))
		for i, r := range rules {
			mRules[i] = &networking.RouteRule{
				ID:        r.ID,
				ServiceID: r.ServiceID,
				Name:      r.Name,
				Enabled:   r.Enabled,
				RuleType:  networking.RouteRuleType(r.RuleType),
				SpecJSON:  r.SpecJSON,
			}
		}
		return networking.BuildMiddlewareLabels(serviceName, mRules), nil
	}
	srv.deployer.DomainFetcher = func(ctx context.Context, serviceID string) ([]domainsmod.DomainConfig, error) {
		if srv.appService != nil {
			if app, err := srv.appService.GetAppService(ctx, serviceID); err == nil && app != nil && app.ProjectID != "" {
				if services, err := srv.appService.ListByProject(ctx, app.ProjectID); err == nil && len(services) == 1 {
					_ = domainRepo.BindUnassignedVerifiedDomainsToSingleService(ctx, app.ProjectID, serviceID)
				}
			}
		}
		list, err := domainRepo.ListByService(ctx, serviceID)
		if err != nil {
			return nil, err
		}
		res := make([]domainsmod.DomainConfig, len(list))
		for i, d := range list {
			res[i] = domainsmod.DomainConfig{
				ID:                 d.ID,
				ServiceID:          d.ServiceID,
				DomainName:         d.DomainName,
				RedirectTo:         d.RedirectTo,
				SSLCertStatus:      d.SSLCertStatus,
				PathPrefix:         d.PathPrefix,
				DNSProvisionStatus: d.DNSProvisionStatus,
				DNSProvider:        d.DNSProvider,
				DNSProvisionedIP:   d.DNSProvisionedIP,
				CreatedAt:          d.CreatedAt,
				UpdatedAt:          d.UpdatedAt,
			}
		}
		return res, nil
	}
}
