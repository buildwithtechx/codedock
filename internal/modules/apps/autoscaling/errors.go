package autoscaling

import "errors"

var (
	ErrPolicyNotFound     = errors.New("autoscaling policy not found")
	ErrServiceIDRequired  = errors.New("service id is required")
	ErrServerlessNotFound = errors.New("serverless code not found")
	ErrInvalidReplicas    = errors.New("minReplicas cannot exceed maxReplicas")
)
