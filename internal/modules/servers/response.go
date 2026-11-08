package servers

type ManagedLogsResponse struct {
	Logs      string `json:"logs"`
	Truncated bool   `json:"truncated"`
}

type ManagedWorkloadsResponse struct {
	Workloads []ManagedWorkload `json:"workloads"`
	Truncated bool              `json:"truncated"`
}

type ControlManagedWorkloadResult struct {
	Ok       bool             `json:"ok"`
	Workload *ManagedWorkload `json:"workload,omitempty"`
}

type OperationReview struct {
	OperationID string `json:"operationId"`
	Type        string `json:"type"`
	Summary     string `json:"summary"`
	CanApply    bool   `json:"canApply"`
}

type ManagedReview struct {
	Review *OperationReview   `json:"review"`
	Plan   *ManagedServerPlan `json:"plan"`
}

type TakeoverScanResponse struct {
	RunID string `json:"runId"`
}

type TakeoverAdoptResponse struct {
	ProjectIDs []string `json:"projectIds"`
}
