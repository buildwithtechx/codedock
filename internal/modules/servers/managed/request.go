package managed

type CreateManagedCredentialRequest struct {
	Provider ManagedProvider `json:"provider"`
	Label    string          `json:"label"`
	Token    string          `json:"token"`
}

type ReviewManagedProvisionRequest struct {
	CredentialID string `json:"credentialId"`
	Region       string `json:"region"`
	ServerType   string `json:"serverType"`
	Image        string `json:"image"`
	Name         string `json:"name"`
	SSHKeyName   string `json:"sshKeyName"`
	UserData     string `json:"userData"`
}

type ReviewManagedResizeRequest struct {
	ServerType string `json:"serverType"`
}

type ApplyManagedOperationRequest struct {
	Confirmation string `json:"confirmation"`
}

type SetManagedSshRequest struct {
	Enabled bool `json:"enabled"`
}

type SetManagedSshKeyRequest struct {
	PublicKey string `json:"publicKey"`
}

type SetManagedSshPasswordRequest struct {
	Password string `json:"password"`
}

type CreateManagedWorkloadRequest struct {
	Name          string `json:"name"`
	Image         string `json:"image"`
	RestartPolicy string `json:"restartPolicy,omitempty"`
	Port          int    `json:"port,omitempty"`
}

type ControlManagedWorkloadRequest struct {
	Action string `json:"action"`
}

type ManagedQuotaRequest struct {
	MaxServers  int `json:"maxServers"`
	MaxMemoryGB int `json:"maxMemoryGB"`
}

type UpdateServerNetworkSettingsRequest struct {
	InternetAccess bool     `json:"internetAccess"`
	Egress         []string `json:"egress"`
}
