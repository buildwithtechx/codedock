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

type ManagedReview struct {
	Review any `json:"review"`
	Plan   any `json:"plan"`
}

type TakeoverScanResponse struct {
	RunID string `json:"runId"`
}

type TakeoverAdoptResponse struct {
	ProjectIDs []string `json:"projectIds"`
}
