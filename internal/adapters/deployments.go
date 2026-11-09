package adapters

import (
	"context"
	"errors"
	"fmt"
	"io"

	"codedock/internal/engine/deploy"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/users"
)

type AppLookupAdapter struct {
	AppService *appsmod.Service
}

func (a *AppLookupAdapter) GetAppService(ctx context.Context, id string) (*deploymentsmod.AppServiceSummary, error) {
	app, err := a.AppService.GetAppService(ctx, id)
	if err != nil || app == nil {
		return nil, err
	}
	return &deploymentsmod.AppServiceSummary{
		ID:            app.ID,
		ProjectID:     app.ProjectID,
		EnvironmentID: app.EnvironmentID,
		Name:          app.Name,
		RepositoryURL: app.RepositoryURL,
		Branch:        app.Branch,
		ImageRef:      app.ImageRef,
		ContainerID:   app.ContainerID,
		Status:        string(app.Status),
		DeployToken:   app.DeployToken,
	}, nil
}

func (a *AppLookupAdapter) UpdateAppContainerID(ctx context.Context, id, containerID string) error {
	app, err := a.AppService.GetAppService(ctx, id)
	if err != nil || app == nil {
		return err
	}
	app.ContainerID = containerID
	return a.AppService.UpdateAppService(ctx, app)
}

type ProjectLookupAdapter struct {
	ProjectService *projectsmod.Service
}

func (p *ProjectLookupAdapter) HasPermission(ctx context.Context, projectID, userID, role, action string) bool {
	return p.ProjectService.HasPermission(ctx, projectID, userID, users.UserRole(role), organizations.MemberPermission(action))
}

func (p *ProjectLookupAdapter) Get(ctx context.Context, id string) (any, error) {
	return p.ProjectService.GetProject(ctx, id)
}

type DeployEngineAdapter struct {
	Deployer   *deploy.Deployer
	AppService *appsmod.Service
}

func (e *DeployEngineAdapter) BeginServiceOperation(ctx context.Context, serviceID string) (context.Context, func(), error) {
	if e.Deployer == nil {
		return ctx, func() {}, nil
	}
	return e.Deployer.BeginServiceOperation(ctx, serviceID)
}

func (e *DeployEngineAdapter) DeployApp(ctx context.Context, app *deploymentsmod.AppServiceSummary, sourceDir string, logWriter io.Writer) (string, error) {
	if e.Deployer == nil {
		return "", errors.New("deployer unavailable")
	}
	var appService *appsmod.AppService
	if app != nil {
		appService = &appsmod.AppService{
			ID:            app.ID,
			ProjectID:     app.ProjectID,
			EnvironmentID: app.EnvironmentID,
			Name:          app.Name,
			RepositoryURL: app.RepositoryURL,
			Branch:        app.Branch,
			ImageRef:      app.ImageRef,
			ContainerID:   app.ContainerID,
			DeployToken:   app.DeployToken,
		}
	}
	return e.Deployer.DeployAppService(ctx, appService, sourceDir, logWriter)
}

func (e *DeployEngineAdapter) StopApp(ctx context.Context, containerID string) error { return nil }

type AutoscalerDeploymentAdapter struct {
	Ds *deploymentsmod.DeploymentService
}

func (a *AutoscalerDeploymentAdapter) CreateDeployment(ctx context.Context, d *deploymentsmod.Deployment) (*deploymentsmod.Deployment, error) {
	if d == nil {
		return nil, errors.New("deployment required")
	}
	return a.Ds.CreateDeployment(ctx, d)
}

func (a *AutoscalerDeploymentAdapter) ExecuteDeploymentAsync(d *deploymentsmod.Deployment) {
	if d != nil {
		a.Ds.ExecuteDeploymentAsync(d)
	}
}

type AppAccessChecker struct {
	Ps *projectsmod.Service
}

func (a *AppAccessChecker) HasPermission(ctx context.Context, projectID, userID, role, permission string) bool {
	return a.Ps.HasPermission(ctx, projectID, userID, users.UserRole(role), organizations.MemberPermission(permission))
}

func (a *AppAccessChecker) AccessibleServiceIDs(ctx context.Context, userID, role string) (map[string]bool, error) {
	return make(map[string]bool), nil
}

type EnvLookupAdapter struct {
	Repo environments.EnvironmentRepository
}

func (a *EnvLookupAdapter) GetProjectID(ctx context.Context, environmentID string) (string, error) {
	env, err := a.Repo.Get(ctx, environmentID)
	if err != nil || env == nil {
		return "", errors.New("environment not found")
	}
	return env.ProjectID, nil
}

type UserStatusAdapter struct {
	Us *users.UserService
}

func (a *UserStatusAdapter) IsUserActive(ctx context.Context, userID string) (bool, error) {
	u, err := a.Us.GetUserByID(ctx, userID)
	if err != nil || u == nil {
		return false, err
	}
	return u.IsActive, nil
}

type AuditLoggerAdapter struct {
	Audit *audit.AuditService
}

func (a *AuditLoggerAdapter) LogAction(ctx any, action, resource, ip string, details map[string]string) {
	c, ok := ctx.(context.Context)
	if !ok {
		c = context.Background()
	}
	a.Audit.LogAction(c, audit.AuditActionOpts{Action: action, Resource: resource, IPAddress: ip, Details: details})
}

type AppLifecycleDeployerAdapter struct {
	Ds         *deploymentsmod.DeploymentService
	AppService *appsmod.Service
}

func (a *AppLifecycleDeployerAdapter) StopAppService(ctx context.Context, appID string) error {
	return nil
}

func (a *AppLifecycleDeployerAdapter) RestartAppService(ctx context.Context, appID string) error {
	app, err := a.AppService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return fmt.Errorf("app not found: %w", err)
	}
	_, err = a.Ds.TriggerDeployment(ctx, app.ID, app.EnvironmentID, app.ProjectID, app.Branch, "restart", "")
	return err
}

func (a *AppLifecycleDeployerAdapter) RedeployAppService(ctx context.Context, appID string) (any, error) {
	app, err := a.AppService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return nil, fmt.Errorf("app not found: %w", err)
	}
	depID, err := a.Ds.TriggerDeployment(ctx, app.ID, app.EnvironmentID, app.ProjectID, app.Branch, "redeploy", "")
	return map[string]string{"deploymentId": depID}, err
}

func (a *AppLifecycleDeployerAdapter) ApplyProjectUpdates(ctx context.Context, projectID string) (any, error) {
	return map[string]string{"status": "applied"}, nil
}

type AIExplainerAdapter struct {
	Svc any
}

func (a *AIExplainerAdapter) ExplainDeploymentFailure(ctx context.Context, deploymentID string) (any, error) {
	return nil, nil
}

type PRPreviewLookupAdapter struct {
	Pr *gitmod.PRPreviewService
}

func (p *PRPreviewLookupAdapter) ListByApp(ctx context.Context, appID string) (any, error) {
	if p.Pr == nil {
		return nil, nil
	}
	return p.Pr.ListByApp(ctx, appID)
}
