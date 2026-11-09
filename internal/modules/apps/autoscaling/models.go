package autoscaling

import "time"

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
	return &AutoscalingPolicy{
		ServiceID:       serviceID,
		MinReplicas:     1,
		MaxReplicas:     5,
		ScaleUpCPU:      80,
		ScaleDownCPU:    20,
		CooldownSeconds: 300,
	}
}

type ServerlessFunctionCode struct {
	ID          string    `json:"id" db:"id"`
	ServiceID   string    `json:"serviceId" db:"service_id"`
	Runtime     string    `json:"runtime" db:"runtime"`
	CodeContent string    `json:"codeContent" db:"code_content"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
