package http

import (
	"context"
	"fmt"

	"codedock/internal/modules/canvas"
	modclusters "codedock/internal/modules/clusters"
	composemod "codedock/internal/modules/compose"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	projectsmod "codedock/internal/modules/projects"
	templatesmod "codedock/internal/modules/templates"
	docker "codedock/internal/platform/docker"
)

type envGitProviderAdapter struct {
	ps *gitmod.ProviderService
}

func (a *envGitProviderAdapter) GetAnyProviderByType(ctx context.Context, providerType string) (*environments.GitProviderConfig, error) {
	if a.ps == nil {
		return nil, nil
	}
	p, err := a.ps.GetAnyProviderByType(ctx, providerType)
	if err != nil || p == nil {
		return nil, err
	}
	return &environments.GitProviderConfig{
		AccessToken: p.AccessToken,
	}, nil
}

type canvasServiceAdapter struct {
	cs *canvas.CanvasService
}

func (c *canvasServiceAdapter) ListSummaries(ctx context.Context, organizationID string) (any, error) {
	return c.cs.ListSummaries(ctx, organizationID)
}

func (c *canvasServiceAdapter) GetSummary(ctx context.Context, id string) (any, error) {
	return c.cs.GetSummary(ctx, id)
}

func (c *canvasServiceAdapter) GetEnvironmentCanvas(ctx context.Context, id string) (any, error) {
	return c.cs.GetEnvironmentCanvas(ctx, id)
}

func (c *canvasServiceAdapter) ApplyTopology(ctx context.Context, envID string, req any) error {
	r, _ := req.(canvas.TopologyApplyRequest)
	return c.cs.ApplyTopology(ctx, envID, r)
}

func (c *canvasServiceAdapter) ReadObservation(ctx context.Context, canv any, node string) (any, error) {
	cnv, _ := canv.(*canvas.EnvironmentCanvas)
	return c.cs.ReadObservation(ctx, cnv, node)
}

type composeStackAdapter struct {
	cs *composemod.ComposeStackService
}

func (c *composeStackAdapter) List(ctx context.Context, projectID string) (any, error) {
	return c.cs.List(ctx, projectID)
}

func (c *composeStackAdapter) Config(ctx context.Context, projectID, stackID string) (any, error) {
	return c.cs.Get(ctx, projectID, stackID)
}

func (c *composeStackAdapter) Review(ctx context.Context, projectID string, req any) (any, error) {
	r, _ := req.(composemod.ComposeStackRequest)
	return c.cs.Review(ctx, r)
}

func (c *composeStackAdapter) Save(ctx context.Context, projectID string, req any) (any, error) {
	r, _ := req.(composemod.ComposeStackRequest)
	return c.cs.Save(ctx, projectID, r)
}

func (c *composeStackAdapter) Deploy(ctx context.Context, projectID, stackID string) (any, error) {
	err := c.cs.Deploy(ctx, projectID, stackID)
	return map[string]string{"status": "deploying"}, err
}

func (c *composeStackAdapter) Cancel(ctx context.Context, projectID, stackID string) error {
	return c.cs.Cancel(projectID, stackID)
}

func (c *composeStackAdapter) Restart(ctx context.Context, projectID, stackID string) error {
	return nil
}

func (c *composeStackAdapter) Stop(ctx context.Context, projectID, stackID string) error {
	return nil
}

type oneClickAdapter struct {
	oc *templatesmod.OneClickService
}

func (o *oneClickAdapter) ListApps() any {
	return o.oc.ListApps()
}

func (o *oneClickAdapter) GetApp(id string) (any, error) {
	return o.oc.GetApp(id)
}

func (o *oneClickAdapter) AddCustomApp(tmpl any) (any, error) {
	if t, ok := tmpl.(docker.ComposeTemplate); ok {
		return o.oc.AddCustomApp(t)
	}
	return nil, fmt.Errorf("invalid template")
}

func (o *oneClickAdapter) DeleteCustomApp(id string) error {
	return o.oc.RemoveCustomApp(id)
}

func (o *oneClickAdapter) Review(ctx context.Context, req any) (any, error) {
	if r, ok := req.(templatesmod.InstallAppInput); ok {
		return o.oc.ReviewInstall(ctx, r)
	}
	return nil, fmt.Errorf("invalid review request")
}

func (o *oneClickAdapter) Deploy(ctx context.Context, req any) (any, error) {
	if r, ok := req.(templatesmod.InstallAppInput); ok {
		return o.oc.InstallApp(ctx, r)
	}
	return nil, fmt.Errorf("invalid deploy request")
}

type projectEnvAdapter struct {
	repo *environments.EnvironmentRepo
}

func (a *projectEnvAdapter) ListByProject(ctx context.Context, projectID string) ([]*projectsmod.EnvironmentConfig, error) {
	if a.repo == nil {
		return []*projectsmod.EnvironmentConfig{}, nil
	}
	envs, err := a.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]*projectsmod.EnvironmentConfig, len(envs))
	for i := range envs {
		res[i] = &projectsmod.EnvironmentConfig{
			ID:        envs[i].ID,
			ProjectID: envs[i].ProjectID,
			Name:      envs[i].Name,
			IsDefault: envs[i].IsDefault,
			CreatedAt: envs[i].CreatedAt,
			UpdatedAt: envs[i].UpdatedAt,
		}
	}
	return res, nil
}

func (a *projectEnvAdapter) GetDefaultByProject(ctx context.Context, projectID string) (*projectsmod.EnvironmentConfig, error) {
	if a.repo == nil {
		return nil, nil
	}
	envs, err := a.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for i := range envs {
		if envs[i].IsDefault {
			return &projectsmod.EnvironmentConfig{
				ID:        envs[i].ID,
				ProjectID: envs[i].ProjectID,
				Name:      envs[i].Name,
				IsDefault: envs[i].IsDefault,
				CreatedAt: envs[i].CreatedAt,
				UpdatedAt: envs[i].UpdatedAt,
			}, nil
		}
	}
	if len(envs) > 0 {
		return &projectsmod.EnvironmentConfig{
			ID:        envs[0].ID,
			ProjectID: envs[0].ProjectID,
			Name:      envs[0].Name,
			IsDefault: envs[0].IsDefault,
			CreatedAt: envs[0].CreatedAt,
			UpdatedAt: envs[0].UpdatedAt,
		}, nil
	}
	return nil, nil
}

type templateEnvAdapter struct {
	repo *environments.EnvironmentRepo
}

func (a *templateEnvAdapter) ListByProject(ctx context.Context, projectID string) ([]*templatesmod.EnvironmentConfig, error) {
	if a.repo == nil {
		return []*templatesmod.EnvironmentConfig{}, nil
	}
	envs, err := a.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]*templatesmod.EnvironmentConfig, len(envs))
	for i := range envs {
		res[i] = &templatesmod.EnvironmentConfig{
			ID:        envs[i].ID,
			ProjectID: envs[i].ProjectID,
			Name:      envs[i].Name,
			IsDefault: envs[i].IsDefault,
		}
	}
	return res, nil
}

type canvasClusterDataAdapter struct {
	repo *modclusters.ClusterDataRepo
}

func (a *canvasClusterDataAdapter) ListByProject(ctx context.Context, project string) ([]canvas.ClusterData, error) {
	if a.repo == nil {
		return nil, nil
	}
	items, err := a.repo.ListByProject(ctx, project)
	if err != nil {
		return nil, err
	}
	res := make([]canvas.ClusterData, len(items))
	for i, it := range items {
		res[i] = canvas.ClusterData{
			ID:        it.ID,
			ClusterID: it.ClusterID,
			ProjectID: it.ProjectID,
			Spec: canvas.ClusterDataSpec{
				ID:              it.Spec.ID,
				Name:            it.Spec.Name,
				EnvironmentID:   it.Spec.EnvironmentID,
				Engine:          it.Spec.Engine,
				Image:           it.Spec.Image,
				Instances:       it.Spec.Instances,
				Shards:          it.Spec.Shards,
				StorageGiB:      it.Spec.StorageGiB,
				StorageClass:    it.Spec.StorageClass,
				S3DestinationID: it.Spec.S3DestinationID,
				Synchronous:     it.Spec.Synchronous,
			},
			Status:          it.Status,
			Error:           it.Error,
			ObservedRoles:   it.ObservedRoles,
			ObservedVolumes: it.ObservedVolumes,
			Config:          it.Config,
			UpdatedAt:       it.UpdatedAt,
		}
	}
	return res, nil
}
