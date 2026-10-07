package migrations

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/dockerprobe"
	"codedock.run/codedock/internal/models"
)

type AdoptedService struct {
	Name      string `json:"name"`
	ServiceID string `json:"serviceId"`
}

type AdoptFailure struct {
	Service string `json:"service"`
	Error   string `json:"error"`
}

type AdoptResult struct {
	ProjectID string           `json:"projectId"`
	Services  []AdoptedService `json:"services"`
	Failures  []AdoptFailure   `json:"failures"`
}

func (s *Service) AdoptSource(ctx context.Context, userID, orgID string, req models.AdoptMigrationRequest) (*AdoptResult, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	if req.ProjectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if len(req.ContainerIDs) == 0 {
		return nil, fmt.Errorf("select at least one service")
	}
	source, err := s.orgSource(ctx, orgID, req.SourceID)
	if err != nil {
		return nil, err
	}
	containers, err := s.selectContainers(ctx, source, req.ContainerIDs)
	if err != nil {
		return nil, err
	}
	projectID := uuid.NewString()
	if err := s.createAdoptProject(ctx, orgID, projectID, req.ProjectName, source); err != nil {
		return nil, err
	}
	return s.adoptContainers(ctx, orgID, projectID, containers, req.ImportEnv), nil
}

func (s *Service) ReimportSource(ctx context.Context, userID, orgID string, req models.ReimportMigrationRequest) (*AdoptResult, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	if len(req.ContainerIDs) == 0 {
		return nil, fmt.Errorf("select at least one service")
	}
	source, err := s.orgSource(ctx, orgID, req.SourceID)
	if err != nil {
		return nil, err
	}
	containers, err := s.selectContainers(ctx, source, req.ContainerIDs)
	if err != nil {
		return nil, err
	}
	projectID := ""
	for _, container := range containers {
		if id := container.Labels["com.codedock.project"]; id != "" {
			projectID = id
			break
		}
	}
	if projectID == "" {
		return nil, fmt.Errorf("selected containers carry no codedock project identity")
	}
	if existing, err := s.projects.GetByOrganization(ctx, projectID, orgID); err == nil && existing != nil {
		return nil, fmt.Errorf("project %s is already registered", projectID)
	}
	name := "reimported-" + projectID
	if len(name) > 64 {
		name = name[:64]
	}
	if err := s.createAdoptProject(ctx, orgID, projectID, name, source); err != nil {
		return nil, err
	}
	return s.adoptContainers(ctx, orgID, projectID, containers, true), nil
}

func (s *Service) selectContainers(ctx context.Context, source *models.MigrationSource, ids []string) ([]dockerprobe.Container, error) {
	containers, err := dockerprobe.ProbeAll(ctx, s.sources(source), nil)
	if err != nil {
		return nil, err
	}
	byID := map[string]dockerprobe.Container{}
	byName := map[string]dockerprobe.Container{}
	for _, container := range containers {
		byID[container.ID] = container
		byName[container.Name] = container
		if len(container.ID) >= 12 {
			byID[container.ID[:12]] = container
		}
	}
	var selected []dockerprobe.Container
	var missing []string
	for _, id := range ids {
		container, ok := byID[id]
		if !ok {
			container, ok = byName[id]
		}
		if !ok {
			missing = append(missing, id)
			continue
		}
		selected = append(selected, container)
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("containers no longer present on source: %s", strings.Join(missing, ", "))
	}
	return selected, nil
}

func (s *Service) createAdoptProject(ctx context.Context, orgID, projectID, name string, source *models.MigrationSource) error {
	now := time.Now().UTC()
	project := &models.ProjectConfig{
		ID:             projectID,
		OrganizationID: orgID,
		Name:           name,
		Description:    fmt.Sprintf("Migrated from %s (%s)", source.Name, source.SSHHost),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.projects.Create(ctx, project); err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (s *Service) adoptContainers(ctx context.Context, orgID, projectID string, containers []dockerprobe.Container, importEnv bool) *AdoptResult {
	result := &AdoptResult{ProjectID: projectID}
	envs, err := s.envs.ListByProject(ctx, projectID)
	if err != nil || len(envs) == 0 {
		for _, container := range containers {
			result.Failures = append(result.Failures, AdoptFailure{Service: container.Name, Error: "project has no environment"})
		}
		return result
	}
	environmentID := envs[0].ID
	for _, container := range containers {
		serviceID, err := s.adoptContainer(ctx, orgID, projectID, environmentID, container, importEnv)
		if err != nil {
			result.Failures = append(result.Failures, AdoptFailure{Service: container.Name, Error: err.Error()})
			continue
		}
		result.Services = append(result.Services, AdoptedService{Name: container.Name, ServiceID: serviceID})
	}
	return result
}

func (s *Service) adoptContainer(ctx context.Context, orgID, projectID, environmentID string, container dockerprobe.Container, importEnv bool) (string, error) {
	service, err := s.apps.CreateAppService(ctx, &models.AppService{
		ProjectID:     projectID,
		EnvironmentID: environmentID,
		Name:          container.Name,
		ImageRef:      container.Image,
		RuntimeMode:   models.RuntimeModeWeb,
		InternalPort:  firstContainerPort(container.Ports),
		Status:        models.AppServiceStatusCreated,
	})
	if err != nil {
		return "", fmt.Errorf("create service record: %w", err)
	}
	if importEnv {
		for key, value := range container.Env {
			if _, err := s.apps.CreateVariable(ctx, &models.Variable{
				ServiceID:     service.ID,
				ProjectID:     projectID,
				EnvironmentID: environmentID,
				Key:           key,
				Value:         value,
				IsSecret:      true,
			}); err != nil {
				return "", fmt.Errorf("import variable %s: %w", key, err)
			}
		}
	}
	for _, volume := range container.Volumes {
		if volume.Source == "" || volume.Destination == "" {
			continue
		}
		if _, err := s.apps.CreateVolume(ctx, &models.ServiceVolume{
			ServiceID:     service.ID,
			HostPath:      volume.Source,
			ContainerPath: volume.Destination,
		}); err != nil {
			return "", fmt.Errorf("import volume %s: %w", volume.Source, err)
		}
	}
	return service.ID, nil
}

func firstContainerPort(ports []string) int {
	for _, port := range ports {
		binding := port
		if index := strings.Index(binding, "/"); index >= 0 {
			binding = binding[:index]
		}
		value, err := strconv.Atoi(binding)
		if err == nil && value > 0 {
			return value
		}
	}
	return 3000
}
