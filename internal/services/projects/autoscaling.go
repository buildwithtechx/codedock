package projects

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"fmt"
	"math"
)

type AutoscalingPolicies interface {
	Get(context.Context, string) (*models.AutoscalingPolicy, error)
	Save(context.Context, *models.AutoscalingPolicy) error
}
type AutoscalingApps interface {
	GetByID(context.Context, string) (*models.AppService, error)
}
type AutoscalingProjects interface {
	Get(context.Context, string) (*models.ProjectConfig, error)
}
type AutoscalingRuntimes interface {
	Get(context.Context, string) (*models.ServiceRuntime, error)
}
type AutoscalingService struct {
	runtimes AutoscalingRuntimes
	policies AutoscalingPolicies
	apps     AutoscalingApps
	projects AutoscalingProjects
}

func NewAutoscalingService(p AutoscalingPolicies, a AutoscalingApps, projects AutoscalingProjects) *AutoscalingService {
	return &AutoscalingService{policies: p, apps: a, projects: projects}
}
func (s *AutoscalingService) Get(ctx context.Context, id string) (*models.AutoscalingPolicy, error) {
	p, err := s.policies.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load autoscaling policy: %w", err)
	}
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("load autoscaling service: %w", err)
	}
	project, err := s.projects.Get(ctx, app.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("load autoscaling target: %w", err)
	}
	p.Supported = project.ServerID == ""
	if s.runtimes != nil {
		runtime, err := s.runtimes.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		p.Supported = p.Supported && runtime.Target.Kind == "docker" && runtime.Journal == ""
	}
	return p, nil
}
func ValidateAutoscaling(p *models.AutoscalingPolicy) error {
	if p == nil || p.MinReplicas < 1 || p.MaxReplicas < p.MinReplicas || p.MaxReplicas > 10 {
		return utils.NewValidationError("replica limits must satisfy 1 <= minimum <= maximum <= 10")
	}
	if math.IsNaN(p.ScaleUpCPU) || math.IsNaN(p.ScaleDownCPU) || math.IsInf(p.ScaleUpCPU, 0) || math.IsInf(p.ScaleDownCPU, 0) || p.ScaleDownCPU < 0 || p.ScaleUpCPU > 100 || p.ScaleUpCPU <= p.ScaleDownCPU {
		return utils.NewValidationError("CPU thresholds must satisfy 0 <= down < up <= 100")
	}
	if p.CooldownSeconds < 120 || p.CooldownSeconds > 86400 {
		return utils.NewValidationError("cooldown must be between 120 and 86400 seconds")
	}
	return nil
}
func (s *AutoscalingService) Save(ctx context.Context, p *models.AutoscalingPolicy) error {
	if err := ValidateAutoscaling(p); err != nil {
		return err
	}
	app, err := s.apps.GetByID(ctx, p.ServiceID)
	if err != nil {
		return fmt.Errorf("load autoscaling service: %w", err)
	}
	project, err := s.projects.Get(ctx, app.ProjectID)
	if err != nil {
		return fmt.Errorf("load autoscaling target: %w", err)
	}
	supported := project.ServerID == ""
	if s.runtimes != nil {
		runtime, err := s.runtimes.Get(ctx, p.ServiceID)
		if err != nil {
			return err
		}
		supported = supported && runtime.Target.Kind == "docker" && runtime.Journal == ""
	}
	if p.Enabled && !supported {
		return utils.NewValidationError("autoscaling currently supports local Docker targets only")
	}
	return s.policies.Save(ctx, p)
}

func (s *AutoscalingService) SetRuntimes(r AutoscalingRuntimes) { s.runtimes = r }
