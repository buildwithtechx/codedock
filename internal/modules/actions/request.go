package actions

type CreateWorkflowRequest struct {
	ProjectID   string `json:"projectId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Trigger     string `json:"trigger,omitempty"`
	YAML        string `json:"yaml,omitempty"`
}

type InspectDestinationRequest struct {
	ServerID string `json:"serverId"`
}

type CreateRunnerRequest struct {
	ServerID    string             `json:"serverId,omitempty"`
	CloudPoolID string             `json:"cloudPoolId,omitempty"`
	Kind        string             `json:"kind,omitempty"`
	Name        string             `json:"name"`
	Mode        string             `json:"mode,omitempty"`
	Image       string             `json:"image,omitempty"`
	Labels      []string           `json:"labels,omitempty"`
	Enabled     bool               `json:"enabled"`
	Config      ActionRunnerConfig `json:"config,omitempty"`
}
