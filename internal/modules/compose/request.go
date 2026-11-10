package compose

// ReviewStackInput contains the parameters to review a compose stack before deployment.
type ReviewStackInput struct {
	EnvironmentID string            `json:"environmentId"`
	Name          string            `json:"name"`
	Content       string            `json:"content"`
	Variables     map[string]string `json:"variables"`
}

// SaveStackInput contains parameters to persist or update a compose stack definition.
type SaveStackInput struct {
	ID            string            `json:"id,omitempty"`
	EnvironmentID string            `json:"environmentId"`
	Name          string            `json:"name"`
	Content       string            `json:"content"`
	Variables     map[string]string `json:"variables"`
}
