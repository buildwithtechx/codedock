package deploy

import (
	"codedock.run/codedock/internal/engine/observability"
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"
)

type AppRepository interface {
	ListAll(context.Context) ([]*models.AppService, error)
}
type DeploymentCreator interface {
	CreateDeployment(context.Context, *models.Deployment) (*models.Deployment, error)
	ExecuteDeploymentAsync(*models.Deployment)
}
type ScalingMetrics interface {
	GetHealth(context.Context, string) (*observability.ContainerHealth, error)
}
type ScalingStore interface {
	ListEnabled(context.Context) ([]*models.AutoscalingPolicy, error)
	Observe(context.Context, *models.AutoscalingPolicy, float64, string, bool) error
	SetReplicas(context.Context, *models.AutoscalingPolicy, int, int) (bool, error)
	RollbackReplicas(context.Context, string, int, int) error
}
type AutoscalerWorker struct {
	appRepo           AppRepository
	statsMonitor      ScalingMetrics
	deploymentService DeploymentCreator
	policies          ScalingStore
	quit              chan struct{}
	startOnce         sync.Once
	stopOnce          sync.Once
}

func NewAutoscalerWorker(apps AppRepository, metrics ScalingMetrics, deployments DeploymentCreator, policies ScalingStore) *AutoscalerWorker {
	return &AutoscalerWorker{appRepo: apps, statsMonitor: metrics, deploymentService: deployments, policies: policies, quit: make(chan struct{})}
}
func (a *AutoscalerWorker) Start() {
	a.startOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(2 * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
					a.checkAndScale(ctx)
					cancel()
				case <-a.quit:
					return
				}
			}
		}()
	})
}
func (a *AutoscalerWorker) Stop() { a.stopOnce.Do(func() { close(a.quit) }) }

func scalingTarget(p *models.AutoscalingPolicy, current int, cpu float64, now time.Time) (int, string) {
	if !p.Enabled {
		return current, "disabled"
	}
	if math.IsNaN(cpu) || math.IsInf(cpu, 0) || cpu < 0 {
		return current, "invalid CPU sample"
	}
	if last, err := time.Parse(time.RFC3339Nano, p.LastScaledAt); err == nil && now.Sub(last) < time.Duration(p.CooldownSeconds)*time.Second {
		return current, "cooldown"
	}
	if current < p.MinReplicas {
		return p.MinReplicas, "below configured minimum"
	}
	if current > p.MaxReplicas {
		return p.MaxReplicas, "above configured maximum"
	}
	if cpu > p.ScaleUpCPU && current < p.MaxReplicas {
		return current + 1, "CPU above scale-up threshold"
	}
	if cpu < p.ScaleDownCPU && current > p.MinReplicas {
		return current - 1, "CPU below scale-down threshold"
	}
	return current, "within thresholds or at replica limit"
}
func (a *AutoscalerWorker) observe(ctx context.Context, p *models.AutoscalingPolicy, cpu float64, reason string, scaled bool) {
	if math.IsNaN(cpu) || math.IsInf(cpu, 0) {
		cpu = 0
	}
	if err := a.policies.Observe(ctx, p, cpu, reason, scaled); err != nil {
		slog.Error("persist autoscaling decision", "service_id", p.ServiceID, "error", err)
	}
}
func (a *AutoscalerWorker) checkAndScale(ctx context.Context) {
	policies, err := a.policies.ListEnabled(ctx)
	if err != nil {
		slog.Error("list autoscaling policies", "error", err)
		return
	}
	if len(policies) == 0 {
		return
	}
	apps, err := a.appRepo.ListAll(ctx)
	if err != nil {
		slog.Error("list autoscaling apps", "error", err)
		return
	}
	byID := make(map[string]*models.AppService, len(apps))
	for _, app := range apps {
		byID[app.ID] = app
	}
	for _, p := range policies {
		app := byID[p.ServiceID]
		if app == nil || app.Status != models.AppServiceStatusRunning || app.ContainerID == "" {
			a.observe(ctx, p, 0, "service is not running", false)
			continue
		}
		health, err := a.statsMonitor.GetHealth(ctx, app.ContainerID)
		if err != nil || health == nil {
			a.observe(ctx, p, 0, "CPU metrics unavailable", false)
			continue
		}
		current := app.Replicas
		if current < 1 {
			current = 1
		}
		target, reason := scalingTarget(p, current, health.CPUUsagePercentage, time.Now())
		if target == current {
			a.observe(ctx, p, health.CPUUsagePercentage, reason, false)
			continue
		}
		reserved, err := a.policies.SetReplicas(ctx, p, current, target)
		if err != nil {
			a.observe(ctx, p, health.CPUUsagePercentage, "replica reservation failed: "+err.Error(), false)
			continue
		}
		if !reserved {
			a.observe(ctx, p, health.CPUUsagePercentage, "skipped: policy changed, remote target, active deployment or replicas changed", false)
			continue
		}
		created, err := a.deploymentService.CreateDeployment(ctx, &models.Deployment{ServiceID: app.ID, ProjectID: app.ProjectID, EnvironmentID: app.EnvironmentID, Status: models.DeploymentStatusPending, Branch: app.Branch, Trigger: "Auto-Scaler", CommitMessage: fmt.Sprintf("Autoscale %d to %d replicas: %s (CPU %.2f%%)", current, target, reason, health.CPUUsagePercentage)})
		if err != nil || created == nil {
			recovery, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if rollbackErr := a.policies.RollbackReplicas(recovery, app.ID, target, current); rollbackErr != nil {
				slog.Error("rollback autoscaling replicas", "service_id", app.ID, "error", rollbackErr)
			}
			a.observe(recovery, p, health.CPUUsagePercentage, fmt.Sprintf("deployment creation failed: %v", err), false)
			cancel()
			continue
		}
		a.observe(ctx, p, health.CPUUsagePercentage, fmt.Sprintf("requested %d to %d replicas via deployment %s: %s", current, target, created.ID, reason), true)
		a.deploymentService.ExecuteDeploymentAsync(created)
	}
}
