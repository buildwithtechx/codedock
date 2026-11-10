package takeover

type TakeoverScanRequest struct {
	Host           string           `json:"host"`
	SSHUser        string           `json:"sshUser"`
	SSHKey         string           `json:"sshKey"`
	SSHFingerprint string           `json:"sshFingerprint,omitempty"`
	Platform       TakeoverPlatform `json:"platform"`
}

type TakeoverAdoptRequest struct {
	RunID        string   `json:"runId"`
	ProjectName  string   `json:"projectName"`
	ServiceNames []string `json:"serviceNames"`
	ImportEnv    bool     `json:"importEnv"`
}
