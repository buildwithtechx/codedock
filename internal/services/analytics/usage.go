package analytics

import (
	"context"
	"fmt"
	"time"

	"codedock/internal/engine/observability"
	"codedock/internal/repositories"
	systemservices "codedock/internal/services/system"
)

type HealthChecker interface {
	GetHealth(ctx context.Context, containerIDOrName string) (*observability.ContainerHealth, error)
}

type UsageReader struct {
	apps    repositories.AppServiceRepository
	metrics *systemservices.MetricsService
	health  HealthChecker
}

func NewUsageReader(apps repositories.AppServiceRepository, metrics *systemservices.MetricsService, health HealthChecker) *UsageReader {
	return &UsageReader{apps: apps, metrics: metrics, health: health}
}

type ServiceUsage struct {
	ServiceID   string  `json:"serviceId"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	CPUPercent  float64 `json:"cpuPercent"`
	MemoryBytes int64   `json:"memoryBytes"`
	MemoryLimit int64   `json:"memoryLimit"`
	Uptime      int64   `json:"uptimeSeconds"`
}

type ResourcesSummary struct {
	ProjectID    string         `json:"projectId"`
	Services     int            `json:"services"`
	Running      int            `json:"running"`
	CPUPercent   float64        `json:"cpuPercent"`
	MemoryBytes  int64          `json:"memoryBytes"`
	Breakdown    []ServiceUsage `json:"breakdown"`
}

func (u *UsageReader) Current(ctx context.Context, projectID string) ([]ServiceUsage, error) {
	services, err := u.apps.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var usage []ServiceUsage
	for _, service := range services {
		entry := ServiceUsage{ServiceID: service.ID, Name: service.Name, Status: string(service.Status)}
		if service.ContainerID != "" && u.health != nil {
			health, err := u.health.GetHealth(ctx, service.ContainerID)
			if err == nil && health != nil {
				entry.Status = string(health.Status)
				entry.CPUPercent = health.CPUUsagePercentage
				entry.MemoryBytes = health.MemoryUsageBytes
				entry.MemoryLimit = health.MemoryLimitBytes
				entry.Uptime = health.UptimeSeconds
			}
		}
		usage = append(usage, entry)
	}
	return usage, nil
}

func (u *UsageReader) Resources(ctx context.Context, projectID string) (*ResourcesSummary, error) {
	usage, err := u.Current(ctx, projectID)
	if err != nil {
		return nil, err
	}
	summary := &ResourcesSummary{ProjectID: projectID, Services: len(usage), Breakdown: usage}
	for _, entry := range usage {
		if entry.Status == "running" {
			summary.Running++
		}
		summary.CPUPercent += entry.CPUPercent
		summary.MemoryBytes += entry.MemoryBytes
	}
	return summary, nil
}

func (u *UsageReader) History(ctx context.Context, serviceID string, hours int) (map[string]any, error) {
	if u.metrics == nil {
		return nil, fmt.Errorf("metrics are not available")
	}
	if hours <= 0 || hours > 168 {
		hours = 24
	}
	end := time.Now().UTC()
	return u.metrics.GetServiceMetrics(ctx, systemservices.ServiceMetricsOpts{
		ServiceID: serviceID, Start: end.Add(-time.Duration(hours) * time.Hour), End: end, Step: "5m",
	})
}

func (u *UsageReader) Container(ctx context.Context, projectID, serviceID string) (*ServiceUsage, error) {
	usage, err := u.Current(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, entry := range usage {
		if entry.ServiceID == serviceID {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("service not found")
}
