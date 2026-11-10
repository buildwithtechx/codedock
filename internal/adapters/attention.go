package adapters

import (
	"context"

	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/servers"
	"codedock/internal/modules/servers/managed"
	"codedock/internal/platform/observability"
)

type StatsHealthAdapter struct {
	Monitor *observability.StatsMonitor
}

func (a *StatsHealthAdapter) GetHealth(ctx context.Context, containerID string) (monitoring.AttentionHealth, error) {
	if a.Monitor == nil {
		return monitoring.AttentionHealth{Status: "unknown"}, nil
	}
	health, err := a.Monitor.GetHealth(ctx, containerID)
	if err != nil {
		return monitoring.AttentionHealth{Status: "unknown"}, err
	}
	return monitoring.AttentionHealth{
		Status:  string(health.Status),
		Running: health.Status == observability.ContainerHealthStatusRunning,
	}, nil
}

type AttentionManagedAdapter struct {
	Repo managed.ManagedRepository
}

func (a *AttentionManagedAdapter) GetQuota(ctx context.Context, orgID string) (*monitoring.ManagedQuota, error) {
	if a.Repo == nil {
		return nil, nil
	}
	q, err := a.Repo.GetQuota(ctx, orgID)
	if err != nil || q == nil {
		return nil, err
	}
	return &monitoring.ManagedQuota{MaxServers: q.MaxServers, MaxMemoryGB: q.MaxMemoryGB}, nil
}

func (a *AttentionManagedAdapter) ListOrgServerIDs(ctx context.Context, orgID string) ([]string, error) {
	if a.Repo == nil {
		return nil, nil
	}
	return a.Repo.ListOrgServerIDs(ctx, orgID)
}

type AttentionServerAdapter struct {
	Repo servers.Repository
}

func (a *AttentionServerAdapter) GetByID(ctx context.Context, id string) (*monitoring.Server, error) {
	if a.Repo == nil {
		return nil, nil
	}
	s, err := a.Repo.GetByID(ctx, id)
	if err != nil || s == nil {
		return nil, err
	}
	return &monitoring.Server{ID: s.ID, ServerType: s.ServerType}, nil
}

type TrafficCollectorAdapter struct {
	Repo monitoring.TrafficRepository
}

func (t *TrafficCollectorAdapter) RecordBatch(ctx context.Context, samples []observability.TrafficSample) error {
	if t.Repo == nil {
		return nil
	}
	mSamples := make([]monitoring.TrafficSample, len(samples))
	for i, s := range samples {
		mSamples[i] = monitoring.TrafficSample{
			Time:       s.Time,
			ProjectID:  s.ProjectID,
			Domain:     s.Domain,
			Path:       s.Path,
			Status:     s.Status,
			Bytes:      s.Bytes,
			DurationMs: s.DurationMs,
			ClientIP:   s.ClientIP,
			Country:    s.Country,
		}
	}
	return t.Repo.RecordBatch(ctx, mSamples)
}

func (t *TrafficCollectorAdapter) PathsEnabled(ctx context.Context, projectID string) (bool, error) {
	if t.Repo == nil {
		return false, nil
	}
	return t.Repo.PathsEnabled(ctx, projectID)
}

func (t *TrafficCollectorAdapter) DeleteBefore(ctx context.Context, cutoff string) error {
	if t.Repo == nil {
		return nil
	}
	return t.Repo.DeleteBefore(ctx, cutoff)
}

func (t *TrafficCollectorAdapter) RetentionCutoff(days int) string {
	if t.Repo == nil {
		return ""
	}
	return t.Repo.RetentionCutoff(days)
}
