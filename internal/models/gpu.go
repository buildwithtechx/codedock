package models

type GPURequest struct {
	Driver       string            `json:"driver"`
	Count        any               `json:"count,omitempty"`
	DeviceIDs    []string          `json:"deviceIds,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Options      map[string]string `json:"options,omitempty"`
}

type GPUConfig struct {
	Runtime string       `json:"runtime,omitempty"`
	GPUs    []GPURequest `json:"gpus,omitempty"`
}
