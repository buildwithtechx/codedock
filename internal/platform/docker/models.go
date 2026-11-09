package docker

type GPURequest struct {
	Driver       string            `json:"driver"`
	Count        any               `json:"count,omitempty"`
	DeviceIDs    []string          `json:"deviceIds,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
}

type GPUConfig struct {
	Runtime string       `json:"runtime,omitempty"`
	GPUs    []GPURequest `json:"gpus,omitempty"`
}

type ComposeStack struct {
	ID            string            `json:"id" db:"id"`
	ProjectID     string            `json:"projectId" db:"project_id"`
	EnvironmentID string            `json:"environmentId" db:"environment_id"`
	Name          string            `json:"name" db:"name"`
	Revision      int               `json:"revision" db:"revision"`
	Status        string            `json:"status" db:"status"`
	Error         string            `json:"error" db:"error"`
	Results       string            `json:"results" db:"results"`
	UpdatedAt     string            `json:"updatedAt" db:"updated_at"`
	Config        string            `json:"-" db:"encrypted_config"`
	Variables     map[string]string `json:"variables,omitempty"`
}

type ComposeServiceResult struct {
	Name        string `json:"name"`
	ContainerID string `json:"containerId"`
	State       string `json:"state"`
	Health      string `json:"health"`
	ExitCode    int    `json:"exitCode"`
}
