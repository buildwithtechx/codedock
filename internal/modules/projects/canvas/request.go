package canvas

// CanvasQueryRequest filters canvas nodes and connections.
type CanvasQueryRequest struct {
	EnvironmentID string `json:"environmentId,omitempty"`
	ProjectID     string `json:"projectId,omitempty"`
}
