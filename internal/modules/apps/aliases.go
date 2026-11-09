package apps

import (
	"codedock/internal/modules/apps/autoscaling"
)

// Autoscaling aliases
type AutoscalingHandler = autoscaling.AutoscalingHandler
type AutoscalingService = autoscaling.AutoscalingService
type AutoscalingRepository = autoscaling.AutoscalingRepository
type ServerlessRepository = autoscaling.ServerlessRepository
type WorkerAutoscalingRepo = autoscaling.WorkerAutoscalingRepo

var (
	NewAutoscalingRepository = autoscaling.NewAutoscalingRepository
	NewServerlessRepository  = autoscaling.NewServerlessRepository
	NewWorkerAutoscalingRepo = autoscaling.NewWorkerAutoscalingRepo
	NewAutoscalingService    = autoscaling.NewAutoscalingService
	NewAutoscalingHandler    = autoscaling.NewAutoscalingHandler
)
