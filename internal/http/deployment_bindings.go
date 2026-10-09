package http

import (
	"context"

	"codedock/internal/engine/networking"
	"codedock/internal/repositories"
)

func configureDeploymentBindings(srv *Server, routeRuleRepo repositories.RouteRuleRepository) {
	if srv.deployer != nil {
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
			return networking.BuildMiddlewareLabels(serviceName, rules), nil
		}
	}

}
