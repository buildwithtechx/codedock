package app

import (
	"context"

	backupmodule "codedock/internal/modules/backups"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/platform/networking"
)

func configureBindings(a *App, s *appServices) {
	s.backupManager.SetServiceRuntime(func(ctx context.Context, serviceID string) (string, error) {
		runtime, err := s.runtimeRepo.Get(ctx, serviceID)
		if err != nil {
			return "", err
		}
		return runtime.Target.Kind, nil
	})
	backupmodule.ConfigureRemoteBackups(s.backupManager, s.projectRepo, s.serverRepo, s.sshManager)

	if a.deployer == nil {
		return
	}

	a.deployer.ProtectionChecker = func(ctx context.Context, hostname string) bool {
		if s.routeProtectionRepo == nil {
			return false
		}
		rp, err := s.routeProtectionRepo.GetByHostname(ctx, hostname)
		return err == nil && rp != nil && rp.Enabled && rp.Secret != nil
	}
	a.deployer.ScopedEnvProvider = func(ctx context.Context, project, environment string) (map[string]string, error) {
		return s.serviceLinker.GetLinkedVariablesForEnvironment(ctx, project, environment)
	}
	a.deployer.ScopedEnvInterpolator = func(ctx context.Context, project, environment string) (map[string]map[string]string, error) {
		return s.serviceLinker.GetNamespacedVariablesForEnvironment(ctx, project, environment)
	}
	a.deployer.EnvProvider = func(projectID string) (map[string]string, error) {
		return s.serviceLinker.GetLinkedEnvironmentVariables(context.Background(), projectID)
	}
	a.deployer.EnvInterpolator = func(projectID string) (map[string]map[string]string, error) {
		return s.serviceLinker.GetNamespacedVariables(context.Background(), projectID)
	}
	a.deployer.RouteRuleFetcher = func(ctx context.Context, serviceID, serviceName string) (map[string]string, error) {
		rules, err := s.routeRuleRepo.ListByService(ctx, serviceID)
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
	a.deployer.DomainFetcher = func(ctx context.Context, serviceID string) ([]domainsmod.DomainConfig, error) {
		if a.appService != nil {
			if app, err := a.appService.GetAppService(ctx, serviceID); err == nil && app != nil && app.ProjectID != "" {
				if services, err := a.appService.ListByProject(ctx, app.ProjectID); err == nil && len(services) == 1 {
					_ = s.domainRepo.BindUnassignedVerifiedDomainsToSingleService(ctx, app.ProjectID, serviceID)
				}
			}
		}
		list, err := s.domainRepo.ListByService(ctx, serviceID)
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
