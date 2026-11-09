package http

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

type appLookupAdapter struct {
	appService *appsmod.Service
}

func (a *appLookupAdapter) GetAppService(ctx context.Context, id string) (*deploymentsmod.AppServiceSummary, error) {
	app, err := a.appService.GetAppService(ctx, id)
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

func (a *appLookupAdapter) UpdateAppContainerID(ctx context.Context, id, containerID string) error {
	app, err := a.appService.GetAppService(ctx, id)
	if err != nil || app == nil {
		return err
	}
	app.ContainerID = containerID
	return a.appService.UpdateAppService(ctx, app)
}

type projectLookupAdapter struct {
	projectService *projectsmod.Service
}

func (p *projectLookupAdapter) HasPermission(ctx context.Context, projectID, userID, role, action string) bool {
	return p.projectService.HasPermission(ctx, projectID, userID, users.UserRole(role), organizations.MemberPermission(action))
}

func (p *projectLookupAdapter) Get(ctx context.Context, id string) (any, error) {
	return p.projectService.GetProject(ctx, id)
}

type deployEngineAdapter struct {
	deployer   *deploy.Deployer
	appService *appsmod.Service
}

func (e *deployEngineAdapter) BeginServiceOperation(ctx context.Context, serviceID string) (context.Context, func(), error) {
	if e.deployer == nil {
		return ctx, func() {}, nil
	}
	return e.deployer.BeginServiceOperation(ctx, serviceID)
}

func (e *deployEngineAdapter) DeployApp(ctx context.Context, app *deploymentsmod.AppServiceSummary, sourceDir string, logWriter io.Writer) (string, error) {
	if e.deployer == nil {
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
	return e.deployer.DeployAppService(ctx, appService, sourceDir, logWriter)
}

func (e *deployEngineAdapter) StopApp(ctx context.Context, containerID string) error { return nil }

type autoscalerDeploymentAdapter struct {
	ds *deploymentsmod.DeploymentService
}

func (a *autoscalerDeploymentAdapter) CreateDeployment(ctx context.Context, d *deploymentsmod.Deployment) (*deploymentsmod.Deployment, error) {
	if d == nil {
		return nil, errors.New("deployment required")
	}
	return a.ds.CreateDeployment(ctx, d)
}

func (a *autoscalerDeploymentAdapter) ExecuteDeploymentAsync(d *deploymentsmod.Deployment) {
	if d != nil {
		a.ds.ExecuteDeploymentAsync(d)
	}
}

type appAccessChecker struct {
	ps *projectsmod.Service
}

func (a *appAccessChecker) HasPermission(ctx context.Context, projectID, userID, role, permission string) bool {
	return a.ps.HasPermission(ctx, projectID, userID, users.UserRole(role), organizations.MemberPermission(permission))
}

func (a *appAccessChecker) AccessibleServiceIDs(ctx context.Context, userID, role string) (map[string]bool, error) {
	return make(map[string]bool), nil
}

type envLookupAdapter struct {
	repo environments.EnvironmentRepository
}

func (a *envLookupAdapter) GetProjectID(ctx context.Context, environmentID string) (string, error) {
	env, err := a.repo.Get(ctx, environmentID)
	if err != nil || env == nil {
		return "", errors.New("environment not found")
	}
	return env.ProjectID, nil
}

type userStatusAdapter struct {
	us *users.UserService
}

func (a *userStatusAdapter) IsUserActive(ctx context.Context, userID string) (bool, error) {
	u, err := a.us.GetUserByID(ctx, userID)
	if err != nil || u == nil {
		return false, err
	}
	return u.IsActive, nil
}

type auditLoggerAdapter struct {
	audit *audit.AuditService
}

func (a *auditLoggerAdapter) LogAction(ctx any, action, resource, ip string, details map[string]string) {
	c, ok := ctx.(context.Context)
	if !ok {
		c = context.Background()
	}
	a.audit.LogAction(c, audit.AuditActionOpts{Action: action, Resource: resource, IPAddress: ip, Details: details})
}

type appLifecycleDeployerAdapter struct {
	ds         *deploymentsmod.DeploymentService
	appService *appsmod.Service
}

func (a *appLifecycleDeployerAdapter) StopAppService(ctx context.Context, appID string) error {
	return nil
}

func (a *appLifecycleDeployerAdapter) RestartAppService(ctx context.Context, appID string) error {
	app, err := a.appService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return fmt.Errorf("app not found: %w", err)
	}
	_, err = a.ds.TriggerDeployment(ctx, app.ID, app.EnvironmentID, app.ProjectID, app.Branch, "restart", "")
	return err
}

func (a *appLifecycleDeployerAdapter) RedeployAppService(ctx context.Context, appID string) (any, error) {
	app, err := a.appService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return nil, fmt.Errorf("app not found: %w", err)
	}
	depID, err := a.ds.TriggerDeployment(ctx, app.ID, app.EnvironmentID, app.ProjectID, app.Branch, "redeploy", "")
	return map[string]string{"deploymentId": depID}, err
}

func (a *appLifecycleDeployerAdapter) ApplyProjectUpdates(ctx context.Context, projectID string) (any, error) {
	return map[string]string{"status": "applied"}, nil
}

type aiExplainerAdapter struct {
	svc any
}

func (a *aiExplainerAdapter) ExplainDeploymentFailure(ctx context.Context, deploymentID string) (any, error) {
	return nil, nil
}

type prPreviewLookupAdapter struct {
	pr *gitmod.PRPreviewService
}

func (p *prPreviewLookupAdapter) ListByApp(ctx context.Context, appID string) (any, error) {
	if p.pr == nil {
		return nil, nil
	}
	return p.pr.ListByApp(ctx, appID)
}
