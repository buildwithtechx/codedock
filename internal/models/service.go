package models

import "time"

type Service struct {
	ID                  string    `json:"id" db:"id"`
	ProjectID           string    `json:"projectId" db:"project_id"`
	AppID               string    `json:"appId,omitempty" db:"app_id"`
	EnvironmentID       string    `json:"environmentId,omitempty" db:"environment_id"`
	Kind                string    `json:"kind" db:"kind"`
	Name                string    `json:"name" db:"name"`
	Icon                string    `json:"icon" db:"icon"`
	Image               string    `json:"image" db:"image"`
	Build               string    `json:"build" db:"build"`
	Dockerfile          string    `json:"dockerfile" db:"dockerfile"`
	BuildArgsJSON       string    `json:"buildArgs" db:"build_args_json"`
	PortsJSON           string    `json:"ports" db:"ports_json"`
	DependsOnJSON       string    `json:"dependsOn" db:"depends_on_json"`
	EnvironmentJSON     string    `json:"environment" db:"environment_json"`
	VolumesJSON         string    `json:"volumes" db:"volumes_json"`
	NamespaceVolumes    bool      `json:"namespaceVolumes" db:"namespace_volumes"`
	Command             string    `json:"command" db:"command"`
	CommandArgvJSON     string    `json:"commandArgv" db:"command_argv_json"`
	Restart             string    `json:"restart" db:"restart"`
	AdvancedJSON        string    `json:"advanced" db:"advanced_json"`
	Exposed             bool      `json:"exposed" db:"exposed"`
	ExposedPort         string    `json:"exposedPort" db:"exposed_port"`
	Domain              string    `json:"domain" db:"domain"`
	CustomDomain        string    `json:"customDomain" db:"custom_domain"`
	DomainType          string    `json:"domainType" db:"domain_type"`
	PublicEndpointsJSON string    `json:"publicEndpoints" db:"public_endpoints_json"`
	CPURequest          float64   `json:"cpuRequest" db:"cpu_request"`
	MemoryLimitMB       int       `json:"memoryLimitMb" db:"memory_limit_mb"`
	Status              string    `json:"status" db:"status"`
	ContainerID         string    `json:"containerId" db:"container_id"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}

type ServicePublicEndpoint struct {
	Port         int    `json:"port"`
	DomainType   string `json:"domainType"`
	Domain       string `json:"domain,omitempty"`
	CustomDomain string `json:"customDomain,omitempty"`
}

type ServiceDeployment struct {
	ID           string    `json:"id" db:"id"`
	DeploymentID string    `json:"deploymentId" db:"deployment_id"`
	ServiceID    string    `json:"serviceId" db:"service_id"`
	ImageRef     string    `json:"imageRef" db:"image_ref"`
	ContainerID  string    `json:"containerId" db:"container_id"`
	Status       string    `json:"status" db:"status"`
	Error        string    `json:"error,omitempty" db:"error"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateServiceRequest struct {
	ProjectID        string   `json:"projectId"`
	AppID            string   `json:"appId,omitempty"`
	EnvironmentID    string   `json:"environmentId,omitempty"`
	Kind             string   `json:"kind,omitempty"`
	Name             string   `json:"name"`
	Icon             string   `json:"icon,omitempty"`
	Image            string   `json:"image,omitempty"`
	Build            string   `json:"build,omitempty"`
	Dockerfile       string   `json:"dockerfile,omitempty"`
	Command          string   `json:"command,omitempty"`
	Restart          string   `json:"restart,omitempty"`
	Exposed          bool     `json:"exposed,omitempty"`
	ExposedPort      string   `json:"exposedPort,omitempty"`
	Domain           string   `json:"domain,omitempty"`
	CustomDomain     string   `json:"customDomain,omitempty"`
	DomainType       string   `json:"domainType,omitempty"`
	Ports            []string `json:"ports,omitempty"`
	Volumes          []string `json:"volumes,omitempty"`
	DependsOn        []string `json:"dependsOn,omitempty"`
	CPURequest       float64  `json:"cpuRequest,omitempty"`
	MemoryLimitMB    int      `json:"memoryLimitMb,omitempty"`
	NamespaceVolumes bool     `json:"namespaceVolumes,omitempty"`
}

type UpdateServiceRequest struct {
	AppID            string   `json:"appId,omitempty"`
	EnvironmentID    string   `json:"environmentId,omitempty"`
	Name             string   `json:"name,omitempty"`
	Icon             string   `json:"icon,omitempty"`
	Image            string   `json:"image,omitempty"`
	Build            string   `json:"build,omitempty"`
	Dockerfile       string   `json:"dockerfile,omitempty"`
	Command          string   `json:"command,omitempty"`
	Restart          string   `json:"restart,omitempty"`
	Exposed          *bool    `json:"exposed,omitempty"`
	ExposedPort      string   `json:"exposedPort,omitempty"`
	Domain           string   `json:"domain,omitempty"`
	CustomDomain     string   `json:"customDomain,omitempty"`
	DomainType       string   `json:"domainType,omitempty"`
	Ports            []string `json:"ports,omitempty"`
	Volumes          []string `json:"volumes,omitempty"`
	DependsOn        []string `json:"dependsOn,omitempty"`
	CPURequest       float64  `json:"cpuRequest,omitempty"`
	MemoryLimitMB    int      `json:"memoryLimitMb,omitempty"`
	Status           string   `json:"status,omitempty"`
	ContainerID      string   `json:"containerId,omitempty"`
	NamespaceVolumes *bool    `json:"namespaceVolumes,omitempty"`
}
