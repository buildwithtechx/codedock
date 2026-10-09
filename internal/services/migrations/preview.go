package migrations

import (
	"context"
	"fmt"
	"strings"

	"codedock/internal/engine/dockerprobe"
	"codedock/internal/models"
)

func (s *Service) PreviewMigration(ctx context.Context, userID, orgID string, req models.PreviewMigrationRequest) (*models.MigrationPreview, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	_, targetRunner, err := s.targetServer(ctx, orgID, req.TargetServerID)
	if err != nil {
		return nil, err
	}
	preview := &models.MigrationPreview{}
	if req.ProjectID != "" {
		return s.previewProject(ctx, orgID, req, targetRunner, preview)
	}
	source, err := s.orgSource(ctx, orgID, req.SourceID)
	if err != nil {
		return nil, err
	}
	containers, err := s.selectContainers(ctx, source, req.ContainerIDs)
	if err != nil {
		return nil, err
	}
	targetPorts, err := usedHostPorts(ctx, targetRunner)
	if err != nil {
		return nil, err
	}
	targetVolumes, err := existingVolumes(ctx, targetRunner)
	if err != nil {
		return nil, err
	}
	seenImages := map[string]struct{}{}
	for _, container := range containers {
		preview.Services = append(preview.Services, models.PreviewService{
			Name: container.Name, Image: container.Image,
			Ports: container.HostPorts, Volumes: volumeSources(container), Status: container.State,
		})
		if _, seen := seenImages[container.Image]; !seen {
			seenImages[container.Image] = struct{}{}
			preview.Images = append(preview.Images, container.Image)
		}
		for _, volume := range container.Volumes {
			if volume.Type == "volume" && volume.Name != "" {
				preview.Volumes = append(preview.Volumes, volume.Name)
			}
		}
		for _, port := range container.HostPorts {
			if targetPorts[port] {
				preview.Conflicts = append(preview.Conflicts, models.PreviewConflict{
					Kind: "port", Subject: port,
					Detail: fmt.Sprintf("host port %s on the target is already published", port),
				})
			}
		}
		for _, volume := range container.Volumes {
			if volume.Type == "volume" && targetVolumes[volume.Name] {
				preview.Conflicts = append(preview.Conflicts, models.PreviewConflict{
					Kind: "volume", Subject: volume.Name,
					Detail: fmt.Sprintf("target already has volume %s; migration will prompt for overwrite, reuse or rename", volume.Name),
				})
			}
		}
		if container.Running {
			preview.Warnings = append(preview.Warnings, fmt.Sprintf("%s is running on the source and restarts on the target; cutover stops the original", container.Name))
		}
		if len(container.Env) == 0 {
			preview.Warnings = append(preview.Warnings, fmt.Sprintf("%s exposes no environment to import", container.Name))
		}
	}
	preview.Downtime = downtimeText(req.Mode, len(containers))
	return preview, nil
}

func (s *Service) previewProject(ctx context.Context, orgID string, req models.PreviewMigrationRequest, targetRunner dockerprobe.Runner, preview *models.MigrationPreview) (*models.MigrationPreview, error) {
	project, err := s.projects.GetByOrganization(ctx, req.ProjectID, orgID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}
	services, err := s.appRepo.ListByProject(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	wanted := map[string]bool{}
	for _, id := range req.ContainerIDs {
		wanted[id] = true
	}
	targetPorts, err := usedHostPorts(ctx, targetRunner)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if len(wanted) > 0 && !wanted[service.ID] {
			continue
		}
		preview.Services = append(preview.Services, models.PreviewService{
			Name: service.Name, Image: service.ImageRef,
			Ports: []string{fmt.Sprintf("%d", service.InternalPort)}, Status: string(service.Status),
		})
		if service.ImageRef != "" {
			preview.Images = append(preview.Images, service.ImageRef)
		}
		if service.InternalPort > 0 && targetPorts[fmt.Sprintf("%d", service.InternalPort)] {
			preview.Conflicts = append(preview.Conflicts, models.PreviewConflict{
				Kind: "port", Subject: fmt.Sprintf("%d", service.InternalPort),
				Detail: fmt.Sprintf("host port %d on the target is already published", service.InternalPort),
			})
		}
		records, err := s.volumes.ListByService(ctx, service.ID)
		if err != nil {
			return nil, err
		}
		for _, volume := range records {
			preview.Volumes = append(preview.Volumes, volume.HostPath)
		}
	}
	preview.Downtime = downtimeText(req.Mode, len(preview.Services))
	if project.ServerID == req.TargetServerID {
		preview.Warnings = append(preview.Warnings, "target matches the current project server; migration redeploys in place")
	}
	return preview, nil
}

func volumeSources(container dockerprobe.Container) []string {
	var sources []string
	for _, volume := range container.Volumes {
		sources = append(sources, volume.Source)
	}
	return sources
}

func downtimeText(mode models.MigrationMode, services int) string {
	if services == 0 {
		return "no services selected"
	}
	if mode == models.MigrationModeCopy {
		return "copy keeps originals running; target services start independently"
	}
	return fmt.Sprintf("move restarts %d service(s) on the target; originals stop at cutover", services)
}

func usedHostPorts(ctx context.Context, runner dockerprobe.Runner) (map[string]bool, error) {
	used := map[string]bool{}
	rows, err := dockerprobe.ListContainers(ctx, runner)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		for _, field := range strings.Split(row.Ports, ",") {
			field = strings.TrimSpace(field)
			arrow := strings.Index(field, "->")
			host := field
			if arrow >= 0 {
				host = field[:arrow]
			}
			if index := strings.LastIndex(host, ":"); index >= 0 {
				host = host[index+1:]
			}
			host = strings.TrimSuffix(host, "/tcp")
			host = strings.TrimSuffix(host, "/udp")
			if host != "" {
				used[host] = true
			}
		}
	}
	return used, nil
}

func existingVolumes(ctx context.Context, runner dockerprobe.Runner) (map[string]bool, error) {
	existing := map[string]bool{}
	out, err := runner.Run(ctx, "docker volume ls -q")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(out, "\n") {
		if name := strings.TrimSpace(line); name != "" {
			existing[name] = true
		}
	}
	return existing, nil
}
