package adapters

import (
	"context"
	"fmt"

	composemod "codedock/internal/modules/compose"
	"codedock/internal/modules/environments"
	"codedock/internal/modules/projects/canvas"
	docker "codedock/internal/platform/docker"
	"codedock/internal/platform/observability"
)

type CanvasRuntimeAdapter struct {
	Rt *observability.CanvasRuntime
}

func (a *CanvasRuntimeAdapter) Observe(ctx context.Context, c *canvas.EnvironmentCanvas) {
	if a.Rt == nil || c == nil {
		return
	}
	obsCanvas := toObservabilityCanvas(c)
	a.Rt.Observe(ctx, obsCanvas)
	for i := range c.Nodes {
		for _, on := range obsCanvas.Nodes {
			if c.Nodes[i].ID == on.ID {
				c.Nodes[i].Data = on.Data
				break
			}
		}
	}
}

func (a *CanvasRuntimeAdapter) Read(ctx context.Context, c *canvas.EnvironmentCanvas, nodeID string) (*canvas.RuntimeObservation, error) {
	if a.Rt == nil || c == nil {
		return nil, fmt.Errorf("runtime unavailable")
	}
	obsCanvas := toObservabilityCanvas(c)
	obs, err := a.Rt.Read(ctx, obsCanvas, nodeID)
	if err != nil {
		return nil, err
	}
	if obs == nil {
		return nil, nil
	}
	return &canvas.RuntimeObservation{
		Logs:    obs.Logs,
		Metrics: obs.Metrics,
	}, nil
}

func toObservabilityCanvas(c *canvas.EnvironmentCanvas) *observability.EnvironmentCanvas {
	if c == nil {
		return nil
	}
	var env *observability.CanvasEnvironment
	if c.Environment != nil {
		env = &observability.CanvasEnvironment{
			ProjectID: c.Environment.ProjectID,
		}
	}
	apps := make([]*observability.CanvasApp, len(c.Apps))
	for i, app := range c.Apps {
		if app != nil {
			apps[i] = &observability.CanvasApp{
				ID:          app.ID,
				ProjectID:   app.ProjectID,
				Name:        app.Name,
				ContainerID: app.ContainerID,
			}
		}
	}
	dbs := make([]*observability.CanvasDatabase, len(c.Databases))
	for i, db := range c.Databases {
		if db != nil {
			dbs[i] = &observability.CanvasDatabase{
				ID:          db.ID,
				ContainerID: db.ContainerID,
			}
		}
	}
	nodes := make([]observability.CanvasNode, len(c.Nodes))
	for i, n := range c.Nodes {
		dataCopy := make(map[string]any, len(n.Data))
		for k, v := range n.Data {
			dataCopy[k] = v
		}
		nodes[i] = observability.CanvasNode{
			ID:   n.ID,
			Type: n.Type,
			Data: dataCopy,
		}
	}
	cdbs := make([]observability.CanvasClusterData, len(c.ClusterDatabases))
	for i, cdb := range c.ClusterDatabases {
		cdbs[i] = observability.CanvasClusterData{
			ID:     cdb.ID,
			Status: cdb.Status,
			Error:  cdb.Error,
			Spec:   observability.CanvasClusterDataSpec{Name: cdb.Spec.Name},
		}
	}
	return &observability.EnvironmentCanvas{
		ID:               c.ID,
		Name:             c.Name,
		Environment:      env,
		Apps:             apps,
		Databases:        dbs,
		ClusterDatabases: cdbs,
		Nodes:            nodes,
		Revision:         c.Revision,
	}
}

type ComposeStackRuntimeAdapter struct {
	Rt *docker.StackRuntime
}

func (a *ComposeStackRuntimeAdapter) Canonical(ctx context.Context, name, config string, env map[string]string) (string, error) {
	if a.Rt == nil {
		return "", fmt.Errorf("runtime unavailable")
	}
	return a.Rt.Canonical(ctx, name, config, env)
}

func (a *ComposeStackRuntimeAdapter) ValidatePorts(ctx context.Context, name, config string) error {
	if a.Rt == nil {
		return fmt.Errorf("runtime unavailable")
	}
	return a.Rt.ValidatePorts(ctx, name, config)
}

func (a *ComposeStackRuntimeAdapter) Apply(ctx context.Context, stack *composemod.ComposeStack, progress func(string) error) error {
	if a.Rt == nil || stack == nil {
		return fmt.Errorf("runtime unavailable")
	}
	dStack := &docker.ComposeStack{
		ID:            stack.ID,
		ProjectID:     stack.ProjectID,
		EnvironmentID: stack.EnvironmentID,
		Name:          stack.Name,
		Revision:      stack.Revision,
		Config:        stack.Config,
		Status:        stack.Status,
		Error:         stack.Error,
		Results:       stack.Results,
		UpdatedAt:     stack.UpdatedAt,
	}
	return a.Rt.Apply(ctx, dStack, progress)
}

func (a *ComposeStackRuntimeAdapter) Results(ctx context.Context, projectID string) ([]composemod.ComposeServiceResult, error) {
	if a.Rt == nil {
		return nil, fmt.Errorf("runtime unavailable")
	}
	res, err := a.Rt.Results(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]composemod.ComposeServiceResult, len(res))
	for i, r := range res {
		out[i] = composemod.ComposeServiceResult{
			Name:        r.Name,
			ContainerID: r.ContainerID,
			State:       r.State,
			Health:      r.Health,
			ExitCode:    r.ExitCode,
		}
	}
	return out, nil
}

type CanvasEnvAdapter struct {
	Repo *environments.EnvironmentRepo
}

func (a *CanvasEnvAdapter) ListByProject(ctx context.Context, projectID string) ([]*canvas.EnvironmentConfig, error) {
	if a.Repo == nil {
		return []*canvas.EnvironmentConfig{}, nil
	}
	envs, err := a.Repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]*canvas.EnvironmentConfig, len(envs))
	for i := range envs {
		res[i] = &envs[i]
	}
	return res, nil
}

func (a *CanvasEnvAdapter) GetDefaultByProject(ctx context.Context, projectID string) (*canvas.EnvironmentConfig, error) {
	if a.Repo == nil {
		return nil, nil
	}
	envs, err := a.Repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for i := range envs {
		if envs[i].IsDefault {
			return &envs[i], nil
		}
	}
	if len(envs) > 0 {
		return &envs[0], nil
	}
	return nil, nil
}
