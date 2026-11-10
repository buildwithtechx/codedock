package actions

// WorkflowResponse wraps a single workflow payload.
type WorkflowResponse struct {
	Workflow *ActionWorkflow `json:"workflow"`
}

// WorkflowListResponse wraps a list of workflows.
type WorkflowListResponse struct {
	Workflows []*ActionWorkflow `json:"workflows"`
	Total     int               `json:"total"`
}

// RunResponse wraps a single run execution payload.
type RunResponse struct {
	Run *ActionRun `json:"run"`
}

// RunListResponse wraps a list of runs.
type RunListResponse struct {
	Runs  []*ActionRun `json:"runs"`
	Total int          `json:"total"`
}

// RunnerResponse wraps a runner payload.
type RunnerResponse struct {
	Runner *ActionRunner `json:"runner"`
}

// RunnerListResponse wraps a list of runners.
type RunnerListResponse struct {
	Runners []*ActionRunner `json:"runners"`
	Total   int             `json:"total"`
}

// DestinationInspectionResponse wraps runner host capability detection.
type DestinationInspectionResponse struct {
	ServerID     string              `json:"serverId"`
	Capabilities *ActionCapabilities `json:"capabilities"`
}
