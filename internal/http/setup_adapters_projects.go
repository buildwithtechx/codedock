package http

import (
	"context"

	"codedock/internal/models"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	compose "codedock/internal/platform/docker"
	projectservices "codedock/internal/services/projects"
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
	cs *projectservices.CanvasService
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
	r, _ := req.(models.TopologyApplyRequest)
	return c.cs.ApplyTopology(ctx, envID, r)
}

func (c *canvasServiceAdapter) ReadObservation(ctx context.Context, canvas any, node string) (any, error) {
	canv, _ := canvas.(*models.EnvironmentCanvas)
	return c.cs.ReadObservation(ctx, canv, node)
}

type composeStackAdapter struct {
	cs *projectservices.ComposeStackService
}

func (c *composeStackAdapter) List(ctx context.Context, projectID string) (any, error) {
	return c.cs.List(ctx, projectID)
}

func (c *composeStackAdapter) Config(ctx context.Context, projectID, stackID string) (any, error) {
	return c.cs.Get(ctx, projectID, stackID)
}

func (c *composeStackAdapter) Review(ctx context.Context, projectID string, req any) (any, error) {
	r, _ := req.(models.ComposeStackRequest)
	return c.cs.Review(ctx, r)
}

func (c *composeStackAdapter) Save(ctx context.Context, projectID string, req any) (any, error) {
	r, _ := req.(models.ComposeStackRequest)
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
	oc *projectservices.OneClickService
}

func (o *oneClickAdapter) ListApps() any {
	return o.oc.ListApps()
}

func (o *oneClickAdapter) GetApp(id string) (any, error) {
	return o.oc.GetApp(id)
}

func (o *oneClickAdapter) AddCustomApp(tmpl any) (any, error) {
	t, _ := tmpl.(compose.ComposeTemplate)
	return o.oc.AddCustomApp(t)
}

func (o *oneClickAdapter) DeleteCustomApp(id string) error {
	return o.oc.RemoveCustomApp(id)
}

func (o *oneClickAdapter) Review(ctx context.Context, req any) (any, error) {
	in, _ := req.(models.InstallAppInput)
	return o.oc.ReviewInstall(ctx, in)
}

func (o *oneClickAdapter) Deploy(ctx context.Context, req any) (any, error) {
	in, _ := req.(models.InstallAppInput)
	return o.oc.InstallApp(ctx, in)
}
