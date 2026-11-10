package autoscaling

// AutoscalingPolicyResponse wraps an autoscaling policy.
type AutoscalingPolicyResponse struct {
	Policy *AutoscalingPolicy `json:"policy"`
}

// ServerlessCodeResponse wraps serverless code content.
type ServerlessCodeResponse struct {
	Code *ServerlessFunctionCode `json:"code"`
}
