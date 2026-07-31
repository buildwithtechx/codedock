package system

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type TakeoverAdopter struct {
	projectRepo repositories.ProjectRepository
	appRepo     repositories.AppServiceRepository
}

func NewTakeoverAdopter(
	projectRepo repositories.ProjectRepository,
	appRepo repositories.AppServiceRepository,
) *TakeoverAdopter {
	return &TakeoverAdopter{projectRepo: projectRepo, appRepo: appRepo}
}

func (a *TakeoverAdopter) Adopt(ctx context.Context, req models.TakeoverAdoptRequest, stack *models.DiscoveredStack, userID string) ([]string, error) {
	projectID := uuid.NewString()
	now := time.Now().UTC()

	project := &models.ProjectConfig{
		ID:          projectID,
		Name:        req.ProjectName,
		Description: fmt.Sprintf("Imported from %s (%s)", stack.Host, stack.Platform),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := a.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	containersByName := make(map[string]models.DiscoveredContainer, len(stack.Containers))
	for _, c := range stack.Containers {
		containersByName[c.Name] = c
	}

	for _, svcName := range req.ServiceNames {
		c, ok := containersByName[svcName]
		if !ok {
			continue
		}

		port := extractFirstPort(c.Ports)
		svc := &models.AppService{
			ID:           uuid.NewString(),
			ProjectID:    projectID,
			Name:         svcName,
			ImageRef:     c.Image,
			RuntimeMode:  "image",
			InternalPort: port,
			Status:       "imported",
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := a.appRepo.Create(ctx, svc); err != nil {
			return nil, fmt.Errorf("create app service %s: %w", svcName, err)
		}
	}

	return []string{projectID}, nil
}

func extractFirstPort(ports []string) int {
	for _, p := range ports {
		p = strings.TrimSuffix(p, "/tcp")
		p = strings.TrimSuffix(p, "/udp")
		parts := strings.Split(p, ":")
		target := parts[len(parts)-1]
		var port int
		if _, err := fmt.Sscanf(target, "%d", &port); err == nil && port > 0 {
			return port
		}
	}
	return 80
}

func serializeStack(stack *models.DiscoveredStack) (string, error) {
	b, err := json.Marshal(stack)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func DeserializeStack(raw string) (*models.DiscoveredStack, error) {
	if raw == "" {
		return nil, nil
	}
	var stack models.DiscoveredStack
	if err := json.Unmarshal([]byte(raw), &stack); err != nil {
		return nil, err
	}
	return &stack, nil
}
