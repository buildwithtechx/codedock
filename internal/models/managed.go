package models

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

type ManagedServerRequest struct {
	CredentialID string `json:"credentialId"`
	Name         string `json:"name"`
	Region       string `json:"region"`
	ServerType   string `json:"serverType"`
	Image        string `json:"image"`
	SSHKeyName   string `json:"sshKeyName"`
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

type ManagedCatalog struct {
	Provider ManagedProvider `json:"provider"`
	Tiers    []ManagedTier   `json:"tiers"`
	Regions  []string        `json:"regions"`
	Images   []string        `json:"images"`
}

type ManagedReview struct {
	Review *OperationReview   `json:"review"`
	Plan   *ManagedServerPlan `json:"plan"`
}
