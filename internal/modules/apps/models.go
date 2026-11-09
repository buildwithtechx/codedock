package apps

import (
	"time"

	"codedock/internal/modules/apps/autoscaling"
	"codedock/internal/modules/apps/vars"
	"codedock/internal/modules/apps/volumes"
	"codedock/internal/modules/apps/webhooks"
)

type ServiceVolume = volumes.ServiceVolume

type ServiceRoleSpec struct {
	Name            string `json:"name,omitempty"`
	Role            string `json:"role,omitempty"`
	Kind            string `json:"kind,omitempty"`
	Command         string `json:"command"`
	StopGracePeriod string `json:"stopGracePeriod,omitempty"`
	Singleton       bool   `json:"singleton,omitempty"`
}

type AppService struct {
	GitUserID              string            `json:"-" db:"git_user_id"`
	ID                     string            `json:"id" db:"id"`
	ProjectID              string            `json:"projectId" db:"project_id"`
	AppID                  string            `json:"appId,omitempty" db:"app_id"`
	EnvironmentID          string            `json:"environmentId" db:"environment_id"`
	Name                   string            `json:"name" db:"name"`
	RepositoryURL          string            `json:"repositoryUrl" db:"repository_url"`
	ImageRef               string            `json:"imageRef,omitempty" db:"image_ref"`
	Branch                 string            `json:"branch" db:"branch"`
	RootDirectory          string            `json:"rootDirectory" db:"root_directory"`
	Icon                   string            `json:"icon" db:"icon"`
	RuntimeMode            RuntimeMode       `json:"runtimeMode" db:"runtime_mode"`
	InstallCommand         string            `json:"installCommand" db:"install_command"`
	BuildCommand           string            `json:"buildCommand" db:"build_command"`
	StartCommand           string            `json:"startCommand" db:"start_command"`
	DockerfilePath         string            `json:"dockerfilePath" db:"dockerfile_path"`
	BuildEngine            BuildEngine       `json:"buildEngine" db:"build_engine"`
	InternalPort           int               `json:"internalPort" db:"internal_port"`
	Domain                 string            `json:"domain" db:"domain"`
	StaticOutput           string            `json:"staticOutput" db:"static_output"`
	DeployToken            string            `json:"deployToken,omitempty" db:"deploy_token"`
	HealthCheckPath        string            `json:"healthCheckPath" db:"health_check_path"`
	ContainerID            string            `json:"containerId" db:"container_id"`
	Status                 AppServiceStatus  `json:"status" db:"status"`
	Replicas               int               `json:"replicas" db:"replicas"`
	CPULimit               float64           `json:"cpuLimit,omitempty" db:"cpu_limit"`
	MemoryLimit            int               `json:"memoryLimit,omitempty" db:"memory_limit"`
	SwapDisabled           bool              `json:"swapDisabled,omitempty" db:"swap_disabled"`
	CreatedAt              time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt              time.Time         `json:"updatedAt" db:"updated_at"`
	EnablePRPreviews       bool              `json:"enablePRPreviews" db:"enable_pr_previews"`
	MaintenanceMode        bool              `json:"maintenanceMode" db:"maintenance_mode"`
	Optimize               bool              `json:"optimize" db:"optimize"`
	Stack                  string            `json:"stack,omitempty" db:"stack"`
	StopGracePeriodSeconds int               `json:"stopGracePeriodSeconds,omitempty" db:"stop_grace_period_seconds"`
	RuntimeVersion         string            `json:"runtimeVersion,omitempty" db:"runtime_version"`
	BuildImage             string            `json:"buildImage,omitempty" db:"build_image"`
	RuntimeImage           string            `json:"runtimeImage,omitempty" db:"runtime_image"`
	RegistryID             *string           `json:"registryId,omitempty" db:"registry_id"`
	SPAFallback            *bool             `json:"spaFallback,omitempty" db:"spa_fallback"`
	Volumes                []ServiceVolume   `json:"volumes,omitempty" db:"-"`
	Roles                  []ServiceRoleSpec `json:"roles,omitempty" db:"-"`
}

type Variable = vars.Variable

type Webhook = webhooks.Webhook
type LogDrain = webhooks.LogDrain

type AutoscalingPolicy = autoscaling.AutoscalingPolicy
type ServerlessFunctionCode = autoscaling.ServerlessFunctionCode

var DefaultAutoscalingPolicy = autoscaling.DefaultAutoscalingPolicy
