package apps

import (
	"context"
	"errors"
)

type AutoscalingService struct {
	autoRepo       AutoscalingRepository
	serverlessRepo ServerlessRepository
	appRepo        AppServiceRepository
}

func NewAutoscalingService(
	autoRepo AutoscalingRepository,
	serverlessRepo ServerlessRepository,
	appRepo AppServiceRepository,
) *AutoscalingService {
	return &AutoscalingService{
		autoRepo:       autoRepo,
		serverlessRepo: serverlessRepo,
		appRepo:        appRepo,
	}
}

func (s *AutoscalingService) GetPolicy(ctx context.Context, serviceID string) (*AutoscalingPolicy, error) {
	if serviceID == "" {
		return nil, errors.New("service id is required")
	}
	policy, err := s.autoRepo.GetByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		policy = DefaultAutoscalingPolicy(serviceID)
		policy.Supported = true
		policy.Enabled = false
	}
	return policy, nil
}

func (s *AutoscalingService) UpdatePolicy(ctx context.Context, policy *AutoscalingPolicy) error {
	if policy == nil || policy.ServiceID == "" {
		return errors.New("valid policy with service id is required")
	}
	if policy.MinReplicas < 1 {
		policy.MinReplicas = 1
	}
	if policy.MaxReplicas < policy.MinReplicas {
		policy.MaxReplicas = policy.MinReplicas
	}
	if policy.ScaleUpCPU <= 0 {
		policy.ScaleUpCPU = 80
	}
	if policy.ScaleDownCPU <= 0 {
		policy.ScaleDownCPU = 20
	}
	if policy.CooldownSeconds <= 0 {
		policy.CooldownSeconds = 300
	}
	return s.autoRepo.Upsert(ctx, policy)
}

func (s *AutoscalingService) SaveServerlessCode(ctx context.Context, serviceID, runtime, codeContent string) (*ServerlessFunctionCode, error) {
	if serviceID == "" || runtime == "" {
		return nil, errors.New("service id and runtime are required")
	}
	if s.serverlessRepo == nil {
		return nil, errors.New("serverless repository not configured")
	}
	return s.serverlessRepo.SaveCode(ctx, serviceID, runtime, codeContent)
}

func (s *AutoscalingService) GetServerlessCode(ctx context.Context, serviceID string) (*ServerlessFunctionCode, error) {
	if serviceID == "" {
		return nil, errors.New("service id is required")
	}
	if s.serverlessRepo == nil {
		return nil, errors.New("serverless repository not configured")
	}
	return s.serverlessRepo.GetCodeByServiceID(ctx, serviceID)
}
