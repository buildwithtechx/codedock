package models

type ServerCapabilities struct {
	Monitor           bool `json:"monitor"`
	Terminal          bool `json:"terminal"`
	Exec              bool `json:"exec"`
	HostConfiguration bool `json:"hostConfiguration"`
	SSH               bool `json:"ssh"`
	NetworkSettings   bool `json:"networkSettings"`
}

type ServerNetworkSettings struct {
	InternetAccess *bool    `json:"internetAccess"`
	IngressPorts   []int    `json:"ingressPorts"`
	IngressAll     *bool    `json:"ingressAll,omitempty"`
	Egress         []string `json:"egress,omitempty"`
	PrivateIP      string   `json:"privateIp,omitempty"`
	OutboundIP     string   `json:"outboundIp,omitempty"`
	OutboundMode   string   `json:"outboundMode,omitempty"`
	Revision       string   `json:"revision,omitempty"`
}

type UpdateServerNetworkSettingsRequest struct {
	InternetAccess         bool     `json:"internetAccess"`
	ExpectedInternetAccess bool     `json:"expectedInternetAccess"`
	Egress                 []string `json:"egress,omitempty"`
	ExpectedRevision       string   `json:"expectedRevision,omitempty"`
	Confirm                bool     `json:"confirm"`
}

type ManagedServerResources struct {
	CPUCores float64 `json:"cpuCores"`
	MemoryMB int64   `json:"memoryMb"`
	DiskMB   int64   `json:"diskMb"`
}

type ManagedServerInfo struct {
	WorkspaceID     string                 `json:"workspaceId"`
	Image           string                 `json:"image"`
	State           string                 `json:"state"`
	Mode            string                 `json:"mode,omitempty"`
	OperatingSystem string                 `json:"operatingSystem,omitempty"`
	RestartPolicy   string                 `json:"restartPolicy,omitempty"`
	Resources       ManagedServerResources `json:"resources"`
}

type ManagedBootLogsRequest struct {
	Tail int `json:"tail"`
}

type ManagedLogsResponse struct {
	Logs      string `json:"logs"`
	Truncated bool   `json:"truncated"`
}

type ManagedSshConnection struct {
	User    string `json:"user"`
	Host    string `json:"host"`
	Bastion string `json:"bastion"`
	Command string `json:"command"`
}

type ManagedSshStatus struct {
	Enabled                *bool                 `json:"enabled"`
	KeyConfigured          *bool                 `json:"keyConfigured"`
	PasswordConfigured     *bool                 `json:"passwordConfigured"`
	RequiresIdentityAccess bool                  `json:"requiresIdentityAccess"`
	Connection             *ManagedSshConnection `json:"connection,omitempty"`
}

type ManagedSshResult struct {
	Status          ManagedSshStatus `json:"status"`
	InitialPassword string           `json:"initialPassword,omitempty"`
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

type ManagedRuntimeStatus struct {
	Enabled *bool `json:"enabled"`
	Running *bool `json:"running"`
}

type ManagedRuntimeCredential struct {
	Endpoint  string `json:"endpoint"`
	Token     string `json:"token"`
	Revision  string `json:"revision"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

type ManagedWorkload struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	State         string `json:"state"`
	RestartPolicy string `json:"restartPolicy,omitempty"`
	Source        string `json:"source"`
	ProjectID     string `json:"projectId,omitempty"`
	Manageable    bool   `json:"manageable"`
}

type ManagedWorkloadsResponse struct {
	Workloads []ManagedWorkload `json:"workloads"`
	Truncated bool              `json:"truncated"`
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

type ControlManagedWorkloadResult struct {
	Ok       bool             `json:"ok"`
	Workload *ManagedWorkload `json:"workload,omitempty"`
}

type ManagedWorkloadLogsRequest struct {
	WorkloadID string `json:"workloadId"`
	Tail       int    `json:"tail"`
}

type ManagedWorkspaceOperation struct {
	ID            string   `json:"id"`
	Kind          string   `json:"kind"`
	Status        string   `json:"status"`
	RequestedAt   string   `json:"requestedAt"`
	NextAttemptAt string   `json:"nextAttemptAt,omitempty"`
	Error         string   `json:"error,omitempty"`
	Logs          []string `json:"logs"`
}

type CloudWorkspaceSummary struct {
	ID                 string                     `json:"id"`
	ServerID           string                     `json:"serverId"`
	Name               string                     `json:"name"`
	PlanTierID         string                     `json:"planTierId"`
	SubscriptionStatus string                     `json:"subscriptionStatus"`
	ProjectCount       int                        `json:"projectCount"`
	State              string                     `json:"state"`
	Resources          *ManagedServerResources    `json:"resources,omitempty"`
	Operation          *ManagedWorkspaceOperation `json:"operation,omitempty"`
	CreatedAt          string                     `json:"createdAt"`
}
