package apps

import (
	"codedock/internal/modules/apps/vars"
	"codedock/internal/modules/apps/volumes"
)

type CreateAppServiceRequest struct {
	ProjectID              string      `json:"projectId"`
	AppID                  string      `json:"appId,omitempty"`
	EnvironmentID          string      `json:"environmentId,omitempty"`
	Name                   string      `json:"name"`
	RepositoryURL          string      `json:"repositoryUrl"`
	ImageRef               string      `json:"imageRef,omitempty"`
	Branch                 string      `json:"branch"`
	RootDirectory          string      `json:"rootDirectory"`
	RuntimeMode            RuntimeMode `json:"runtimeMode"`
	InstallCommand         string      `json:"installCommand"`
	BuildCommand           string      `json:"buildCommand"`
	StartCommand           string      `json:"startCommand"`
	DockerfilePath         string      `json:"dockerfilePath"`
	BuildEngine            BuildEngine `json:"buildEngine"`
	InternalPort           int         `json:"internalPort"`
	Domain                 string      `json:"domain"`
	StaticOutput           string      `json:"staticOutput"`
	HealthCheckPath        string      `json:"healthCheckPath"`
	CPULimit               float64     `json:"cpuLimit,omitempty"`
	MemoryLimit            int         `json:"memoryLimit,omitempty"`
	SwapDisabled           bool        `json:"swapDisabled,omitempty"`
	EnablePRPreviews       bool        `json:"enablePRPreviews"`
	Optimize               bool        `json:"optimize,omitempty"`
	Stack                  string      `json:"stack,omitempty"`
	StopGracePeriodSeconds int         `json:"stopGracePeriodSeconds,omitempty"`
	RuntimeVersion         string      `json:"runtimeVersion,omitempty"`
	BuildImage             string      `json:"buildImage,omitempty"`
	RuntimeImage           string      `json:"runtimeImage,omitempty"`
	RegistryID             *string     `json:"registryId,omitempty"`
	SPAFallback            *bool       `json:"spaFallback,omitempty"`
}

type UpdateAppServiceRequest struct {
	Name            *string           `json:"name,omitempty"`
	Branch          *string           `json:"branch,omitempty"`
	RootDirectory   *string           `json:"rootDirectory,omitempty"`
	InstallCommand  *string           `json:"installCommand,omitempty"`
	BuildCommand    *string           `json:"buildCommand,omitempty"`
	StartCommand    *string           `json:"startCommand,omitempty"`
	DockerfilePath  *string           `json:"dockerfilePath,omitempty"`
	BuildEngine     *BuildEngine      `json:"buildEngine,omitempty"`
	InternalPort    *int              `json:"internalPort,omitempty"`
	Domain          *string           `json:"domain,omitempty"`
	StaticOutput    *string           `json:"staticOutput,omitempty"`
	HealthCheckPath *string           `json:"healthCheckPath,omitempty"`
	CPULimit        *float64          `json:"cpuLimit,omitempty"`
	MemoryLimit     *int              `json:"memoryLimit,omitempty"`
	Status          *AppServiceStatus `json:"status,omitempty"`
	RuntimeVersion  *string           `json:"runtimeVersion,omitempty"`
	BuildImage      *string           `json:"buildImage,omitempty"`
	RuntimeImage    *string           `json:"runtimeImage,omitempty"`
	RegistryID      *string           `json:"registryId,omitempty"`
	SPAFallback     *bool             `json:"spaFallback,omitempty"`
	MaintenanceMode *bool             `json:"maintenanceMode,omitempty"`
	Optimize        *bool             `json:"optimize,omitempty"`
}

type CreateVariableRequest = vars.CreateVariableRequest
type UpdateVariableRequest = vars.UpdateVariableRequest

type CreateVolumeRequest = volumes.CreateVolumeRequest

type CreateWebhookRequest struct {
	URL    string `json:"url"`
	Secret string `json:"secret,omitempty"`
	Events string `json:"events,omitempty"`
}

type CreateLogDrainRequest struct {
	ProjectID   string    `json:"projectId"`
	DrainType   DrainType `json:"drainType"`
	EndpointURL string    `json:"endpointUrl"`
	APIKey      string    `json:"apiKey,omitempty"`
}

type SaveServerlessCodeRequest struct {
	Runtime     string `json:"runtime"`
	CodeContent string `json:"codeContent"`
}
