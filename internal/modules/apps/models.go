package apps

import "time"

type ServiceVolume struct {
	ID            string    `json:"id" db:"id"`
	ServiceID     string    `json:"serviceId" db:"service_id"`
	HostPath      string    `json:"hostPath" db:"host_path"`
	ContainerPath string    `json:"containerPath" db:"container_path"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

type ServiceRoleSpec struct {
	Role    string `json:"role"`
	Command string `json:"command"`
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

type Variable struct {
	ID                       string    `json:"id" db:"id"`
	ServiceID                string    `json:"serviceId" db:"service_id"`
	EnvironmentID            string    `json:"environmentId,omitempty" db:"environment_id"`
	Key                      string    `json:"key" db:"key"`
	Value                    string    `json:"value" db:"value"`
	IsSecret                 bool      `json:"isSecret" db:"is_secret"`
	ExpectedTopologyRevision string    `json:"-" db:"-"`
	CreatedAt                time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt                time.Time `json:"updatedAt" db:"updated_at"`
}

type Webhook struct {
	ID        string    `json:"id" db:"id"`
	ServiceID string    `json:"serviceId" db:"service_id"`
	URL       string    `json:"url" db:"url"`
	Secret    string    `json:"secret,omitempty" db:"secret"`
	Events    string    `json:"events" db:"events"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type LogDrain struct {
	ID          string    `json:"id" db:"id"`
	ProjectID   string    `json:"projectId" db:"project_id"`
	ServiceID   string    `json:"serviceId" db:"service_id"`
	DrainType   DrainType `json:"drainType" db:"drain_type"`
	EndpointURL string    `json:"endpointUrl" db:"endpoint_url"`
	APIKey      string    `json:"apiKey,omitempty" db:"api_key"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type AutoscalingPolicy struct {
	Supported       bool    `json:"supported"`
	ServiceID       string  `json:"serviceId"`
	Enabled         bool    `json:"enabled"`
	MinReplicas     int     `json:"minReplicas"`
	MaxReplicas     int     `json:"maxReplicas"`
	ScaleUpCPU      float64 `json:"scaleUpCpu"`
	ScaleDownCPU    float64 `json:"scaleDownCpu"`
	CooldownSeconds int     `json:"cooldownSeconds"`
	LastCPU         float64 `json:"lastCpu"`
	LastDecision    string  `json:"lastDecision"`
	LastEvaluatedAt string  `json:"lastEvaluatedAt"`
	LastScaledAt    string  `json:"lastScaledAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

func DefaultAutoscalingPolicy(serviceID string) *AutoscalingPolicy {
	return &AutoscalingPolicy{
		ServiceID:       serviceID,
		MinReplicas:     1,
		MaxReplicas:     5,
		ScaleUpCPU:      80,
		ScaleDownCPU:    20,
		CooldownSeconds: 300,
	}
}

type ServerlessFunctionCode struct {
	ID          string    `json:"id" db:"id"`
	ServiceID   string    `json:"serviceId" db:"service_id"`
	Runtime     string    `json:"runtime" db:"runtime"`
	CodeContent string    `json:"codeContent" db:"code_content"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
