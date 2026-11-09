package apps

import (
	"codedock/internal/modules/apps/autoscaling"
	"codedock/internal/modules/apps/vars"
	"codedock/internal/modules/apps/volumes"
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

// Vars aliases
type VariableHandler = vars.VariableHandler
type ServiceVarRepository = vars.ServiceVarRepository
type ServiceVarRepo = vars.ServiceVarRepo

var (
	NewServiceVarRepo  = vars.NewServiceVarRepo
	NewVariableHandler = vars.NewVariableHandler
)

// Volumes aliases
type VolumeHandler = volumes.VolumeHandler
type ServiceVolumeRepository = volumes.ServiceVolumeRepository

var (
	NewServiceVolumeRepo = volumes.NewServiceVolumeRepo
	NewVolumeHandler     = volumes.NewVolumeHandler
)
