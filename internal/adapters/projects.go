package adapters

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

type EnvGitProviderAdapter struct {
	Ps *gitmod.ProviderService
}

func (a *EnvGitProviderAdapter) GetAnyProviderByType(ctx context.Context, providerType string) (*environments.GitProviderConfig, error) {
	if a.Ps == nil {
		return nil, nil
	}
	p, err := a.Ps.GetAnyProviderByType(ctx, providerType)
	if err != nil || p == nil {
		return nil, err
	}
	return &environments.GitProviderConfig{
		AccessToken: p.AccessToken,
	}, nil
}

type CanvasServiceAdapter struct {
	Cs *canvas.CanvasService
}

func (c *CanvasServiceAdapter) ListSummaries(ctx context.Context, organizationID string) (any, error) {
	return c.Cs.ListSummaries(ctx, organizationID)
}

func (c *CanvasServiceAdapter) GetSummary(ctx context.Context, id string) (any, error) {
	return c.Cs.GetSummary(ctx, id)
}

func (c *CanvasServiceAdapter) GetEnvironmentCanvas(ctx context.Context, id string) (any, error) {
	return c.Cs.GetEnvironmentCanvas(ctx, id)
}

func (c *CanvasServiceAdapter) ApplyTopology(ctx context.Context, envID string, req any) error {
	r, _ := req.(canvas.TopologyApplyRequest)
	return c.Cs.ApplyTopology(ctx, envID, r)
}

func (c *CanvasServiceAdapter) ReadObservation(ctx context.Context, canv any, node string) (any, error) {
	cnv, _ := canv.(*canvas.EnvironmentCanvas)
	return c.Cs.ReadObservation(ctx, cnv, node)
}

type ComposeStackAdapter struct {
	Cs *composemod.ComposeStackService
}

func (c *ComposeStackAdapter) List(ctx context.Context, projectID string) (any, error) {
	return c.Cs.List(ctx, projectID)
}

func (c *ComposeStackAdapter) Config(ctx context.Context, projectID, stackID string) (any, error) {
	return c.Cs.Get(ctx, projectID, stackID)
}

func (c *ComposeStackAdapter) Review(ctx context.Context, projectID string, req any) (any, error) {
	r, _ := req.(composemod.ComposeStackRequest)
	return c.Cs.Review(ctx, r)
}

func (c *ComposeStackAdapter) Save(ctx context.Context, projectID string, req any) (any, error) {
	r, _ := req.(composemod.ComposeStackRequest)
	return c.Cs.Save(ctx, projectID, r)
}

func (c *ComposeStackAdapter) Deploy(ctx context.Context, projectID, stackID string) (any, error) {
	err := c.Cs.Deploy(ctx, projectID, stackID)
	return map[string]string{"status": "deploying"}, err
}

func (c *ComposeStackAdapter) Cancel(ctx context.Context, projectID, stackID string) error {
	return c.Cs.Cancel(projectID, stackID)
}

func (c *ComposeStackAdapter) Restart(ctx context.Context, projectID, stackID string) error {
	return nil
}

func (c *ComposeStackAdapter) Stop(ctx context.Context, projectID, stackID string) error {
	return nil
}

type OneClickAdapter struct {
	Oc *templatesmod.OneClickService
}

func (o *OneClickAdapter) ListApps() any {
	return o.Oc.ListApps()
}

func (o *OneClickAdapter) GetApp(id string) (any, error) {
	return o.Oc.GetApp(id)
}

func (o *OneClickAdapter) AddCustomApp(tmpl any) (any, error) {
	if t, ok := tmpl.(docker.ComposeTemplate); ok {
		return o.Oc.AddCustomApp(t)
	}
	return nil, fmt.Errorf("invalid template")
}

func (o *OneClickAdapter) DeleteCustomApp(id string) error {
	return o.Oc.RemoveCustomApp(id)
}

func (o *OneClickAdapter) Review(ctx context.Context, req any) (any, error) {
	if r, ok := req.(templatesmod.InstallAppInput); ok {
		return o.Oc.ReviewInstall(ctx, r)
	}
	return nil, fmt.Errorf("invalid review request")
}

func (o *OneClickAdapter) Deploy(ctx context.Context, req any) (any, error) {
	if r, ok := req.(templatesmod.InstallAppInput); ok {
		return o.Oc.InstallApp(ctx, r)
	}
	return nil, fmt.Errorf("invalid deploy request")
}

type ProjectEnvAdapter struct {
	Repo *environments.EnvironmentRepo
}

func (a *ProjectEnvAdapter) ListByProject(ctx context.Context, projectID string) ([]*projectsmod.EnvironmentConfig, error) {
	if a.Repo == nil {
		return []*projectsmod.EnvironmentConfig{}, nil
	}
	envs, err := a.Repo.ListByProject(ctx, projectID)
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

func (a *ProjectEnvAdapter) GetDefaultByProject(ctx context.Context, projectID string) (*projectsmod.EnvironmentConfig, error) {
	if a.Repo == nil {
		return nil, nil
	}
	envs, err := a.Repo.ListByProject(ctx, projectID)
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

type TemplateEnvAdapter struct {
	Repo *environments.EnvironmentRepo
}

func (a *TemplateEnvAdapter) ListByProject(ctx context.Context, projectID string) ([]*templatesmod.EnvironmentConfig, error) {
	if a.Repo == nil {
		return []*templatesmod.EnvironmentConfig{}, nil
	}
	envs, err := a.Repo.ListByProject(ctx, projectID)
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

type CanvasClusterDataAdapter struct {
	Repo *modclusters.ClusterDataRepo
}

func (a *CanvasClusterDataAdapter) ListByProject(ctx context.Context, project string) ([]canvas.ClusterData, error) {
	if a.Repo == nil {
		return nil, nil
	}
	items, err := a.Repo.ListByProject(ctx, project)
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
