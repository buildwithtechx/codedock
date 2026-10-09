package migrations

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"codedock/internal/engine/dockerprobe"
	"codedock/internal/models"
)

func (s *Service) ScanSource(ctx context.Context, userID, orgID, sourceID string, progress func(step, detail string)) (*models.MaskedStack, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	source, err := s.orgSource(ctx, orgID, sourceID)
	if err != nil {
		return nil, err
	}
	containers, err := dockerprobe.ProbeAll(ctx, s.sources(source), progress)
	if err != nil {
		return nil, err
	}
	return maskStack(source.SSHHost, containers), nil
}

func (s *Service) RevealSourceEnv(ctx context.Context, userID, orgID string, req models.RevealMigrationEnvRequest) (map[string]string, error) {
	if err := s.requireOrg(ctx, userID, orgID, true); err != nil {
		return nil, err
	}
	source, err := s.orgSource(ctx, orgID, req.SourceID)
	if err != nil {
		return nil, err
	}
	if req.ContainerID == "" {
		return nil, fmt.Errorf("container id is required")
	}
	container, err := dockerprobe.InspectContainer(ctx, s.sources(source), req.ContainerID)
	if err != nil {
		return nil, err
	}
	return container.Env, nil
}

func (s *Service) RepoCompose(ctx context.Context, userID, orgID string, req models.RepoComposeMigrationRequest) (*models.MaskedStack, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	if req.RepoURL == "" {
		return nil, fmt.Errorf("repository url is required")
	}
	if s.git == nil || s.composer == nil {
		return nil, fmt.Errorf("repository inspection is not available")
	}
	branch := req.Branch
	if branch == "" {
		branch = "main"
	}
	workdir, err := os.MkdirTemp("", "codedock-migration-compose")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workdir)
	app := &models.AppService{RepositoryURL: req.RepoURL, Branch: branch, GitUserID: userID}
	if err := s.git.CloneOrPullAppRepository(ctx, app, workdir, io.Discard); err != nil {
		return nil, fmt.Errorf("clone repository: %w", err)
	}
	composePath := req.ComposePath
	if composePath == "" {
		composePath = findComposeFile(workdir)
	}
	if composePath == "" {
		return &models.MaskedStack{Host: req.RepoURL, Platform: models.TakeoverPlatformDocker}, nil
	}
	data, err := os.ReadFile(filepath.Join(workdir, composePath))
	if err != nil {
		return nil, fmt.Errorf("read compose file: %w", err)
	}
	var compose models.UserComposeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("parse compose file: %w", err)
	}
	stack := &models.MaskedStack{Host: req.RepoURL, Platform: models.TakeoverPlatformDocker}
	for name, service := range compose.Services {
		env := map[string]string{}
		for key := range service.Environment {
			env[key] = dockerprobe.MaskedSecret
		}
		stack.Containers = append(stack.Containers, models.MaskedContainer{
			Name:   name,
			Image:  service.Image,
			Ports:  service.Ports,
			Env:    env,
			EnvKeys: dockerprobe.EnvKeys(service.Environment),
			Volumes: service.Volumes,
			Labels: map[string]string{},
			Status: "compose",
		})
	}
	return stack, nil
}

func findComposeFile(root string) string {
	for _, candidate := range []string{"compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml"} {
		if _, err := os.Stat(filepath.Join(root, candidate)); err == nil {
			return candidate
		}
	}
	return ""
}

func maskStack(host string, containers []dockerprobe.Container) *models.MaskedStack {
	stack := &models.MaskedStack{Host: host, Platform: models.TakeoverPlatformDocker}
	composeSet := map[string]struct{}{}
	for _, container := range containers {
		if platform := detectSourcePlatform(container.Labels); platform != models.TakeoverPlatformDocker {
			stack.Platform = platform
		}
		if container.ComposeProject != "" {
			composeSet[container.ComposeProject] = struct{}{}
		}
		volumes := make([]string, 0, len(container.Volumes))
		for _, volume := range container.Volumes {
			volumes = append(volumes, volume.Source)
		}
		stack.Containers = append(stack.Containers, models.MaskedContainer{
			ID:             container.ID,
			Name:           container.Name,
			Image:          container.Image,
			Ports:          container.Ports,
			EnvKeys:        dockerprobe.EnvKeys(container.Env),
			Env:            dockerprobe.MaskEnv(container.Env),
			Volumes:        volumes,
			Labels:         container.Labels,
			Status:         container.State,
			ComposeProject: container.ComposeProject,
			Routes:         dockerprobe.DetectRoutes(container),
		})
	}
	for project := range composeSet {
		stack.ComposeProjects = append(stack.ComposeProjects, project)
	}
	return stack
}

func detectSourcePlatform(labels map[string]string) models.TakeoverPlatform {
	for key := range labels {
		switch {
		case strings.HasPrefix(key, "dokploy."):
			return models.TakeoverPlatformDokploy
		case strings.HasPrefix(key, "coolify."):
			return models.TakeoverPlatformCoolify
		case strings.HasPrefix(key, "com.dokku."):
			return models.TakeoverPlatformDokku
		case strings.HasPrefix(key, "com.codedock."):
			return models.TakeoverPlatformDocker
		}
	}
	return models.TakeoverPlatformDocker
}
