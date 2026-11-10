package canvas

// CanvasSummaryResponse wraps a single canvas summary.
type CanvasSummaryResponse struct {
	Summary *CanvasSummary `json:"summary"`
}

// CanvasSummaryListResponse wraps a slice of canvas summaries.
type CanvasSummaryListResponse struct {
	Summaries []CanvasSummary `json:"summaries"`
	Total     int             `json:"total"`
}

// EnvironmentCanvasResponse wraps an entire environment graph.
type EnvironmentCanvasResponse struct {
	Canvas *EnvironmentCanvas `json:"canvas"`
}
