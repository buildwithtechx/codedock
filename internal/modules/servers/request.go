package servers

type CreateServerRequest struct {
	Name          string `json:"name"`
	IPAddress     string `json:"ipAddress,omitempty"`
	IsLocal       bool   `json:"isLocal,omitempty"`
	SSHHost       string `json:"sshHost,omitempty"`
	SSHPort       int    `json:"sshPort,omitempty"`
	SSHUser       string `json:"sshUser,omitempty"`
	SSHAuthMethod string `json:"sshAuthMethod,omitempty"`
	SSHKey        string `json:"sshKey,omitempty"`
	SSHPrivateKey string `json:"sshPrivateKey,omitempty"`
	SSHPassword   string `json:"sshPassword,omitempty"`
	SSHTransport  string `json:"sshTransport,omitempty"`
	SSHJumpHost   string `json:"sshJumpHost,omitempty"`
}

type UpdateServerRequest struct {
	Name                  *string `json:"name,omitempty"`
	IPAddress             *string `json:"ipAddress,omitempty"`
	IsLocal               *bool   `json:"isLocal,omitempty"`
	SSHHost               *string `json:"sshHost,omitempty"`
	SSHPort               *int    `json:"sshPort,omitempty"`
	SSHUser               *string `json:"sshUser,omitempty"`
	SSHAuthMethod         *string `json:"sshAuthMethod,omitempty"`
	SSHKey                *string `json:"sshKey,omitempty"`
	SSHPrivateKey         *string `json:"sshPrivateKey,omitempty"`
	SSHPassword           *string `json:"sshPassword,omitempty"`
	SSHTransport          *string `json:"sshTransport,omitempty"`
	SSHJumpHost           *string `json:"sshJumpHost,omitempty"`
	BuildCacheBudgetBytes *int64  `json:"buildCacheBudgetBytes,omitempty"`
}

type TestSSHRequest struct {
	SSHHost       string `json:"sshHost"`
	SSHPort       int    `json:"sshPort"`
	SSHUser       string `json:"sshUser"`
	SSHKey        string `json:"sshKey"`
	SSHPrivateKey string `json:"sshPrivateKey,omitempty"`
	SSHPassword   string `json:"sshPassword"`
}

type UpdateServerNetworkSettingsRequest struct {
	InternetAccess         bool     `json:"internetAccess"`
	ExpectedInternetAccess bool     `json:"expectedInternetAccess"`
	Egress                 []string `json:"egress,omitempty"`
	ExpectedRevision       string   `json:"expectedRevision,omitempty"`
	Confirm                bool     `json:"confirm"`
}

type CreateNetworkRequest struct {
	Name     string `json:"name"`
	Driver   string `json:"driver,omitempty"`
	Subnet   string `json:"subnet,omitempty"`
	Gateway  string `json:"gateway,omitempty"`
	Internal bool   `json:"internal,omitempty"`
}

type ManagedBootLogsRequest struct {
	Tail int `json:"tail"`
}

type SetManagedSshRequest struct {
	Enabled         bool `json:"enabled"`
	ExpectedEnabled bool `json:"expectedEnabled"`
	Confirm         bool `json:"confirm"`
}

type SetManagedSshKeyRequest struct {
	PublicKey string `json:"publicKey"`
	Confirm   bool   `json:"confirm"`
}

type SetManagedSshPasswordRequest struct {
	Password string `json:"password"`
	Confirm  bool   `json:"confirm"`
}

type CreateManagedWorkloadRequest struct {
	Name             string   `json:"name"`
	Command          string   `json:"command"`
	WorkingDirectory string   `json:"workingDirectory"`
	Environment      []string `json:"environment"`
	RestartPolicy    string   `json:"restartPolicy"`
	Confirm          bool     `json:"confirm"`
	IdempotencyKey   string   `json:"idempotencyKey"`
}

type ControlManagedWorkloadRequest struct {
	WorkloadID string `json:"workloadId"`
	Action     string `json:"action"`
	Confirm    bool   `json:"confirm"`
}

type ManagedWorkloadLogsRequest struct {
	WorkloadID string `json:"workloadId"`
	Tail       int    `json:"tail"`
}

type CreateManagedCredentialRequest struct {
	Provider ManagedProvider `json:"provider"`
	Label    string          `json:"label"`
	Token    string          `json:"token"`
}

type ReviewManagedProvisionRequest struct {
	CredentialID string `json:"credentialId"`
	ServerID     string `json:"serverId,omitempty"`
	Name         string `json:"name"`
	Region       string `json:"region"`
	ServerType   string `json:"serverType"`
	Image        string `json:"image"`
}

type ReviewManagedResizeRequest struct {
	ServerType string `json:"serverType"`
}

type ManagedQuotaRequest struct {
	MaxServers  int `json:"maxServers"`
	MaxMemoryGB int `json:"maxMemoryGB"`
}

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
