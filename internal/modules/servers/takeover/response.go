package takeover

// TakeoverScanResponse wraps a scan result.
type TakeoverScanResponse struct {
	RunID string           `json:"runId"`
	Stack *DiscoveredStack `json:"stack"`
}

// TakeoverAdoptResponse wraps adopted project details.
type TakeoverAdoptResponse struct {
	RunID     string   `json:"runId"`
	ProjectID string   `json:"projectId"`
	Services  []string `json:"services"`
}
