package projects

import (
	"codedock.run/codedock/internal/models"
	"context"
	"errors"
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
type AutoscalingService struct {
	policies AutoscalingPolicies
	apps     AutoscalingApps
	projects AutoscalingProjects
}

func NewAutoscalingService(p AutoscalingPolicies, a AutoscalingApps, projects AutoscalingProjects) *AutoscalingService {
	return &AutoscalingService{policies: p, apps: a, projects: projects}
}
func (s *AutoscalingService) Get(ctx context.Context, id string) (*models.AutoscalingPolicy, error) {
	return s.policies.Get(ctx, id)
}
func ValidateAutoscaling(p *models.AutoscalingPolicy) error {
	if p == nil || p.MinReplicas < 1 || p.MaxReplicas < p.MinReplicas || p.MaxReplicas > 10 {
		return errors.New("replica limits must satisfy 1 <= minimum <= maximum <= 10")
	}
	if math.IsNaN(p.ScaleUpCPU) || math.IsNaN(p.ScaleDownCPU) || math.IsInf(p.ScaleUpCPU, 0) || math.IsInf(p.ScaleDownCPU, 0) || p.ScaleDownCPU < 0 || p.ScaleUpCPU > 100 || p.ScaleUpCPU <= p.ScaleDownCPU {
		return errors.New("CPU thresholds must satisfy 0 <= down < up <= 100")
	}
	if p.CooldownSeconds < 120 || p.CooldownSeconds > 86400 {
		return errors.New("cooldown must be between 120 and 86400 seconds")
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
	if p.Enabled && project.ServerID != "" {
		return errors.New("autoscaling currently supports local Docker targets only")
	}
	return s.policies.Save(ctx, p)
}
