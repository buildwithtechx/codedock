package autoscaling

import "context"

// AutoscalingRepository persists autoscaling threshold metrics and evaluation results.
type AutoscalingRepository interface {
	GetByServiceID(ctx context.Context, serviceID string) (*AutoscalingPolicy, error)
	Upsert(ctx context.Context, policy *AutoscalingPolicy) error
}

// ServerlessRepository stores uploaded serverless code implementations.
type ServerlessRepository interface {
	SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*ServerlessFunctionCode, error)
	GetCodeByServiceID(ctx context.Context, serviceID string) (*ServerlessFunctionCode, error)
}

// AppServiceRepository references service app records for autoscaling.
type AppServiceRepository interface{}
