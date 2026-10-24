package models

type ComposeStackRequest struct {
	RepositoryURL string            `json:"repositoryUrl,omitempty"`
	Branch        string            `json:"branch,omitempty"`
	RootDirectory string            `json:"rootDirectory,omitempty"`
	ID            string            `json:"id"`
	EnvironmentID string            `json:"environmentId"`
	Name          string            `json:"name"`
	Content       string            `json:"content"`
	Variables     map[string]string `json:"variables"`
	Revision      int               `json:"revision"`
	Digest        string            `json:"digest"`
}

type ComposeStack struct {
	ID            string `json:"id" db:"id"`
	ProjectID     string `json:"projectId" db:"project_id"`
	EnvironmentID string `json:"environmentId" db:"environment_id"`
	Name          string `json:"name" db:"name"`
	Revision      int    `json:"revision" db:"revision"`
	Status        string `json:"status" db:"status"`
	Error         string `json:"error" db:"error"`
	Results       string `json:"results" db:"results"`
	UpdatedAt     string `json:"updatedAt" db:"updated_at"`
	Config        string `json:"-" db:"encrypted_config"`
}

type ComposeReview struct {
	Config   string   `json:"config"`
	Digest   string   `json:"digest"`
	Services []string `json:"services"`
	Effects  []string `json:"effects"`
}

type ComposeServiceResult struct {
	Name        string `json:"name"`
	ContainerID string `json:"containerId"`
	State       string `json:"state"`
	Health      string `json:"health"`
	ExitCode    int    `json:"exitCode"`
}

type TopologyApplyRequest struct {
	Revision     string       `json:"revision"`
	Dependencies []CanvasEdge `json:"dependencies"`
}
