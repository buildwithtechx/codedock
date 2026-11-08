package deployments

type DeployResponse struct {
	DeploymentID string `json:"deploymentId"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

type RollbackResponse struct {
	DeploymentID string `json:"deploymentId"`
	PreviousID   string `json:"previousId"`
	Status       string `json:"status"`
}

type PinToggleResponse struct {
	DeploymentID string `json:"deploymentId"`
	Pinned       bool   `json:"pinned"`
}

type CancelResponse struct {
	DeploymentID string `json:"deploymentId"`
	Cancelled    bool   `json:"cancelled"`
}

type LogsResponse struct {
	DeploymentID string `json:"deploymentId"`
	BuildLogs    string `json:"buildLogs"`
}

type FailureExplanation struct {
	DeploymentID string   `json:"deploymentId"`
	RootCause    string   `json:"rootCause"`
	Explanation  string   `json:"explanation"`
	Suggestions  []string `json:"suggestions,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
