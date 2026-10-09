package projects

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"

	"codedock/internal/engine/compose"
	"codedock/internal/models"
	"codedock/internal/repositories"
)

type OneClickService struct {
	tmplManager *compose.TemplateManager
	stacks      *ComposeStackService
	envRepo     repositories.EnvironmentRepository
}

func NewOneClickService(
	tm *compose.TemplateManager,
	stacks *ComposeStackService,
	er repositories.EnvironmentRepository,
) *OneClickService {
	return &OneClickService{
		tmplManager: tm,
		stacks:      stacks,
		envRepo:     er,
	}
}

func (s *OneClickService) ListApps() []models.OneClickApp {
	apps := []models.OneClickApp{}
	for _, id := range s.tmplManager.ListTemplates() {
		tmpl, err := s.tmplManager.GetTemplate(id)
		if err != nil {
			continue
		}
		if app := extractOneClickApp(id, &tmpl); app != nil {
			apps = append(apps, *app)
		}
	}
	sort.Slice(apps, func(i, j int) bool { return apps[i].Name < apps[j].Name })
	return apps
}

func (s *OneClickService) GetApp(appID string) (*models.OneClickApp, error) {
	tmpl, err := s.tmplManager.GetTemplate(appID)
	if err != nil {
		return nil, errors.New("unknown app: " + appID)
	}
	app := extractOneClickApp(appID, &tmpl)
	if app == nil {
		return nil, errors.New("app has no one-click metadata")
	}
	return app, nil
}

func (s *OneClickService) ReviewInstall(_ context.Context, input models.InstallAppInput) (*models.InstallPreview, error) {
	tmpl, err := s.tmplManager.GetTemplate(input.AppID)
	if err != nil {
		return nil, errors.New("unknown app: " + input.AppID)
	}
	if findOneClickMetadata(&tmpl) == nil {
		return nil, errors.New("app has no one-click metadata")
	}
	plan, err := compose.ResolveInstallPlan(input.AppID, tmpl, installRequest(input), nil)
	if err != nil {
		return nil, err
	}
	return previewInstall(input.AppID, plan.Masked()), nil
}

func (s *OneClickService) InstallApp(ctx context.Context, input models.InstallAppInput) (*models.AppInstallResult, error) {
	tmpl, err := s.tmplManager.GetTemplate(input.AppID)
	if err != nil {
		return nil, errors.New("unknown app: " + input.AppID)
	}
	if findOneClickMetadata(&tmpl) == nil {
		return nil, errors.New("app has no one-click metadata")
	}
	if input.ProjectID == "" {
		return nil, errors.New("projectId is required")
	}
	envID, err := s.resolveEnvironment(ctx, input.ProjectID, input.EnvironmentID)
	if err != nil {
		return nil, err
	}
	plan, err := compose.ResolveInstallPlan(input.AppID, tmpl, installRequest(input), nil)
	if err != nil {
		return nil, err
	}
	if input.Digest != "" && input.Digest != plan.Digest {
		return nil, errors.New("install inputs changed since review; review again before installing")
	}
	return s.installStack(ctx, input, envID, plan)
}

func (s *OneClickService) resolveEnvironment(ctx context.Context, projectID, environmentID string) (string, error) {
	envs, err := s.envRepo.ListByProject(ctx, projectID)
	if err != nil || len(envs) == 0 {
		return "", errors.New("project has no environments")
	}
	if environmentID == "" {
		return envs[0].ID, nil
	}
	for _, env := range envs {
		if env.ID == environmentID {
			return env.ID, nil
		}
	}
	return "", errors.New("environment does not belong to project")
}

func (s *OneClickService) installStack(ctx context.Context, input models.InstallAppInput, envID string, plan *compose.InstallPlan) (*models.AppInstallResult, error) {
	if s.stacks == nil {
		return nil, errors.New("stack installer is unavailable")
	}
	appName := input.Name
	if appName == "" {
		appName = input.AppID
	}
	request := models.ComposeStackRequest{
		ID:            uuid.New().String(),
		EnvironmentID: envID,
		Name:          appName,
		Content:       plan.ComposeYAML,
		Revision:      0,
	}
	review, err := s.stacks.Review(ctx, request)
	if err != nil {
		return nil, err
	}
	request.Digest = review.Digest
	stack, err := s.stacks.Save(ctx, input.ProjectID, request)
	if err != nil {
		return nil, err
	}
	if err := s.stacks.Deploy(ctx, input.ProjectID, stack.ID); err != nil {
		return nil, err
	}
	return &models.AppInstallResult{Kind: "stack", Stack: stack}, nil
}

func installRequest(input models.InstallAppInput) compose.InstallRequest {
	name := input.Name
	if name == "" {
		name = input.AppID
	}
	return compose.InstallRequest{
		Name:        name,
		Secrets:     input.Secrets,
		Environment: input.Environment,
		HostPort:    input.HostPort,
		Domain:      input.Domain,
	}
}

func previewInstall(appID string, plan *compose.InstallPlan) *models.InstallPreview {
	services := make([]models.InstallPreviewService, 0, len(plan.Services))
	for _, service := range plan.Services {
		services = append(services, models.InstallPreviewService{
			Service: service.Service,
			Image:   service.Image,
			Env:     service.Env,
			Ports:   service.Ports,
		})
	}
	volumes := make([]models.InstallPreviewVolume, 0, len(plan.Volumes))
	for _, volume := range plan.Volumes {
		volumes = append(volumes, models.InstallPreviewVolume{
			Service: volume.Service,
			Name:    volume.Name,
			Target:  volume.Target,
		})
	}
	return &models.InstallPreview{
		AppID:            appID,
		Name:             plan.Name,
		Services:         services,
		Volumes:          volumes,
		GeneratedSecrets: plan.GeneratedSecrets,
		ComposeYAML:      plan.ComposeYAML,
		Digest:           plan.Digest,
		Kind:             "stack",
	}
}
