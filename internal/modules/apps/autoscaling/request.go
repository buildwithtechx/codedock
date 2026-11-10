package autoscaling

// UpdateAutoscalingPolicyRequest defines payload for updating autoscaling boundaries.
type UpdateAutoscalingPolicyRequest struct {
	Enabled         bool    `json:"enabled"`
	MinReplicas     int     `json:"minReplicas"`
	MaxReplicas     int     `json:"maxReplicas"`
	ScaleUpCPU      float64 `json:"scaleUpCpu"`
	ScaleDownCPU    float64 `json:"scaleDownCpu"`
	CooldownSeconds int     `json:"cooldownSeconds"`
}

// SaveServerlessCodeRequest defines payload for uploading serverless function code.
type SaveServerlessCodeRequest struct {
	Runtime     string `json:"runtime"`
	CodeContent string `json:"codeContent"`
}
