package http

import (
	"context"

	"codedock/internal/adapters"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	"codedock/internal/modules/servers"
	"codedock/internal/modules/servers/networks"
	"codedock/internal/modules/servers/takeover"
	"codedock/internal/platform/networking"
)

func createServerHandlers(d *serverSetupDeps) (*servers.Handler, *servers.ServerMetricsWSHandler, *takeover.TakeoverHandler, *networks.ServerNetworksHandler) {
	serversRepo := servers.NewRepository(d.db, d.v)
	serversUserLookup := &adapters.ServerUserLookupAdapter{Users: d.userRepo}
	serversService := servers.NewService(serversRepo, serversUserLookup, d.sshManager)
	serverAudit := &adapters.ServerAuditLoggerAdapter{Audit: &adapters.AuditDomainAdapter{Audit: d.auditService}}
	serversHandler := servers.NewHandler(serversService, serverAudit)
	takeoverRepo := takeover.NewTakeoverRepository(d.db, d.v)
	takeoverCreator := &adapters.ProjectCreatorAdapter{Projects: d.projectRepo, Apps: d.appRepo}
	takeoverService := takeover.NewTakeoverService(takeoverRepo, takeoverCreator)
	serversTakeoverHandler := takeover.NewTakeoverHandler(takeoverService, takeoverRepo)
	serversNetworksHandler := networks.NewServerNetworksHandler()
	serversWSHandler := servers.NewServerMetricsWSHandler(&adapters.ServerTokenValidatorAdapter{Tokens: d.tokenService}, serversService, serversUserLookup)
	return serversHandler, serversWSHandler, serversTakeoverHandler, serversNetworksHandler
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
