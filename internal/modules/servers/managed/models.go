package managed

type ManagedProvider string

const (
	ManagedProviderHetzner ManagedProvider = "hetzner"
)

const (
	ManagedActionProvision = "provision"
	ManagedActionResize    = "resize"
	ManagedActionDelete    = "delete"
)

type ManagedCredential struct {
	ID             string          `json:"id" db:"id"`
	OrganizationID string          `json:"organizationId" db:"organization_id"`
	Provider       ManagedProvider `json:"provider" db:"provider"`
	Label          string          `json:"label" db:"label"`
	Token          string          `json:"-" db:"encrypted_token"`
	CreatedAt      string          `json:"createdAt" db:"created_at"`
	UpdatedAt      string          `json:"updatedAt" db:"updated_at"`
}

type ManagedTier struct {
	Name         string  `json:"name"`
	Cores        int     `json:"cores"`
	MemoryGB     int     `json:"memoryGB"`
	DiskGB       int     `json:"diskGB"`
	MonthlyPrice float64 `json:"monthlyPrice"`
}

type ManagedQuota struct {
	OrganizationID string `json:"organizationId" db:"organization_id"`
	MaxServers     int    `json:"maxServers" db:"max_servers"`
	MaxMemoryGB    int    `json:"maxMemoryGB" db:"max_memory_gb"`
	UsedServers    int    `json:"usedServers" db:"-"`
	UsedMemoryGB   int    `json:"usedMemoryGB" db:"-"`
	UpdatedAt      string `json:"updatedAt" db:"updated_at"`
}

type ManagedServerPlan struct {
	CredentialID   string          `json:"credentialId"`
	OrganizationID string          `json:"organizationId"`
	UserID         string          `json:"userId"`
	ServerID       string          `json:"serverId"`
	ExternalID     string          `json:"externalId"`
	Name           string          `json:"name"`
	Region         string          `json:"region"`
	ServerType     string          `json:"serverType"`
	ResizeTo       string          `json:"resizeTo,omitempty"`
	Image          string          `json:"image"`
	SSHKeyName     string          `json:"sshKeyName"`
	Provider       ManagedProvider `json:"provider"`
	Tier           ManagedTier     `json:"tier"`
	MonthlyPrice   float64         `json:"monthlyPrice"`
	Action         string          `json:"action"`
}

type ManagedServerSpec struct {
	Name       string            `json:"name"`
	ServerType string            `json:"serverType"`
	Image      string            `json:"image"`
	Region     string            `json:"region"`
	SSHKeyName string            `json:"sshKeyName"`
	UserData   string            `json:"userData"`
	Labels     map[string]string `json:"labels"`
}

type ManagedInstance struct {
	ExternalID string            `json:"externalId"`
	Name       string            `json:"name"`
	Status     string            `json:"status"`
	PublicIP   string            `json:"publicIp"`
	ServerType string            `json:"serverType"`
	Region     string            `json:"region"`
	Labels     map[string]string `json:"labels"`
}

type ManagedServerLink struct {
	ServerID       string `json:"serverId" db:"server_id"`
	OrganizationID string `json:"organizationId" db:"organization_id"`
	CredentialID   string `json:"credentialId" db:"credential_id"`
	SSHKeyName     string `json:"sshKeyName" db:"ssh_key_name"`
	CreatedAt      string `json:"createdAt" db:"created_at"`
	UpdatedAt      string `json:"updatedAt" db:"updated_at"`
}

type ManagedCatalog struct {
	Provider ManagedProvider `json:"provider"`
	Tiers    []ManagedTier   `json:"tiers"`
	Regions  []string        `json:"regions"`
	Images   []string        `json:"images"`
}

type ManagedLocation struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type ManagedOffer struct {
	Tier   ManagedTier        `json:"tier"`
	Prices map[string]float64 `json:"prices"`
}

type ManagedOffers struct {
	Provider  ManagedProvider   `json:"provider"`
	Locations []ManagedLocation `json:"locations"`
	Offers    []ManagedOffer    `json:"offers"`
	Images    []string          `json:"images"`
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

type ServerNetworkSettings struct {
	InternetAccess *bool    `json:"internetAccess"`
	IngressPorts   []int    `json:"ingressPorts"`
	Egress         []string `json:"egress"`
	PrivateIP      string   `json:"privateIp"`
	OutboundIP     string   `json:"outboundIp"`
	OutboundMode   string   `json:"outboundMode"`
	Revision       string   `json:"revision"`
}
