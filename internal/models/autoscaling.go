package models

type AutoscalingPolicy struct {
	Supported       bool    `json:"supported"`
	ServiceID       string  `json:"serviceId"`
	Enabled         bool    `json:"enabled"`
	MinReplicas     int     `json:"minReplicas"`
	MaxReplicas     int     `json:"maxReplicas"`
	ScaleUpCPU      float64 `json:"scaleUpCpu"`
	ScaleDownCPU    float64 `json:"scaleDownCpu"`
	CooldownSeconds int     `json:"cooldownSeconds"`
	LastCPU         float64 `json:"lastCpu"`
	LastDecision    string  `json:"lastDecision"`
	LastEvaluatedAt string  `json:"lastEvaluatedAt"`
	LastScaledAt    string  `json:"lastScaledAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

func DefaultAutoscalingPolicy(serviceID string) *AutoscalingPolicy {
	return &AutoscalingPolicy{ServiceID: serviceID, MinReplicas: 1, MaxReplicas: 5, ScaleUpCPU: 80, ScaleDownCPU: 20, CooldownSeconds: 300}
}
