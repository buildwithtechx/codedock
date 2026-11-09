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

type DockerPsRow struct {
	ID    string `json:"ID"`
	Names string `json:"Names"`
	Image string `json:"Image"`
	State string `json:"State"`
	Ports string `json:"Ports"`
}

type DockerInspectMount struct {
	Source string `json:"Source"`
	Type   string `json:"Type"`
}

type DockerInspectResult struct {
	Name   string `json:"Name"`
	Config struct {
		Image  string            `json:"Image"`
		Env    []string          `json:"Env"`
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
	HostConfig struct {
		Binds        []string               `json:"Binds"`
		PortBindings map[string]interface{} `json:"PortBindings"`
	} `json:"HostConfig"`
	Mounts []DockerInspectMount `json:"Mounts"`
}

type dockerPsRow = DockerPsRow
type dockerInspectMount = DockerInspectMount
type dockerInspectResult = DockerInspectResult
