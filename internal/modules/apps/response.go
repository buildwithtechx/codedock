package apps

type AppServiceResponse struct {
	Service *AppService `json:"service"`
}

type AppServiceListResponse struct {
	Services []*AppService `json:"services"`
}

type VariableResponse struct {
	Variable *Variable `json:"variable"`
}

type VariableListResponse struct {
	Variables []*Variable `json:"variables"`
}

type VolumeResponse struct {
	Volume *ServiceVolume `json:"volume"`
}

type VolumeListResponse struct {
	Volumes []ServiceVolume `json:"volumes"`
}

type AutoscalingResponse struct {
	Policy *AutoscalingPolicy `json:"policy"`
}

type ServerlessCodeResponse struct {
	Code *ServerlessFunctionCode `json:"code"`
}
