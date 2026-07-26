package types

import "time"

type AppService struct {
	ID               string    `json:"id" db:"id"`
	ProjectID        string    `json:"projectId" db:"project_id"`
	EnvironmentID    string    `json:"environmentId" db:"environment_id"`
	Name             string    `json:"name" db:"name"`
	RepositoryURL    string    `json:"repositoryUrl" db:"repository_url"`
	ImageRef         string    `json:"imageRef,omitempty" db:"image_ref"`
	Branch           string    `json:"branch" db:"branch"`
	RootDirectory    string    `json:"rootDirectory" db:"root_directory"`
	RuntimeMode      string    `json:"runtimeMode" db:"runtime_mode"`
	InstallCommand   string    `json:"installCommand" db:"install_command"`
	BuildCommand     string    `json:"buildCommand" db:"build_command"`
	StartCommand     string    `json:"startCommand" db:"start_command"`
	DockerfilePath   string    `json:"dockerfilePath" db:"dockerfile_path"`
	BuildEngine      string    `json:"buildEngine" db:"build_engine"`
	InternalPort     int       `json:"internalPort" db:"internal_port"`
	Domain           string    `json:"domain" db:"domain"`
	StaticOutput     string    `json:"staticOutput" db:"static_output"`
	HealthCheckPath  string    `json:"healthCheckPath" db:"health_check_path"`
	ContainerID      string    `json:"containerId" db:"container_id"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
	CPULimit         float64   `json:"cpuLimit,omitempty" db:"cpu_limit"`
	MemoryLimit      int64     `json:"memoryLimit,omitempty" db:"memory_limit"`
	EnablePRPreviews bool      `json:"enablePrPreviews" db:"enable_pr_previews"`
	MaintenanceMode  bool      `json:"maintenanceMode" db:"maintenance_mode"`
}

type CreateAppServiceRequest struct {
	ProjectID       string  `json:"projectId"`
	Name            string  `json:"name"`
	RepositoryURL   string  `json:"repositoryUrl"`
	ImageRef        string  `json:"imageRef,omitempty"`
	Branch          string  `json:"branch"`
	RootDirectory   string  `json:"rootDirectory"`
	RuntimeMode     string  `json:"runtimeMode"`
	InstallCommand  string  `json:"installCommand"`
	BuildCommand    string  `json:"buildCommand"`
	StartCommand    string  `json:"startCommand"`
	DockerfilePath  string  `json:"dockerfilePath"`
	BuildEngine     string  `json:"buildEngine"`
	InternalPort    int     `json:"internalPort"`
	Domain          string  `json:"domain"`
	StaticOutput    string  `json:"staticOutput"`
	HealthCheckPath string  `json:"healthCheckPath"`
	CPULimit        float64 `json:"cpuLimit,omitempty"`
	MemoryLimit     int64   `json:"memoryLimit,omitempty"`
}

type UpdateAppServiceRequest struct {
	Name             string  `json:"name"`
	RepositoryURL    string  `json:"repositoryUrl"`
	ImageRef         string  `json:"imageRef,omitempty"`
	Branch           string  `json:"branch"`
	RootDirectory    string  `json:"rootDirectory"`
	RuntimeMode      string  `json:"runtimeMode"`
	InstallCommand   string  `json:"installCommand"`
	BuildCommand     string  `json:"buildCommand"`
	StartCommand     string  `json:"startCommand"`
	DockerfilePath   string  `json:"dockerfilePath"`
	BuildEngine      string  `json:"buildEngine"`
	InternalPort     int     `json:"internalPort"`
	Domain           string  `json:"domain"`
	StaticOutput     string  `json:"staticOutput"`
	HealthCheckPath  string  `json:"healthCheckPath"`
	ContainerID      string  `json:"containerId"`
	Status           string  `json:"status"`
	CPULimit         float64 `json:"cpuLimit,omitempty"`
	MemoryLimit      int64   `json:"memoryLimit,omitempty"`
	Icon             string  `json:"icon,omitempty"`
	DeployToken      string  `json:"deployToken,omitempty"`
	EnablePRPreviews bool    `json:"enablePrPreviews"`
	MaintenanceMode  bool    `json:"maintenanceMode"`
}
