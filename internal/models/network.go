package models

type NetworkMode string

const (
	NetworkModeWireguard NetworkMode = "wireguard"
	NetworkModeNative    NetworkMode = "native"
)

type NetworkMember struct {
	ServerID      string  `json:"serverId"`
	ServerName    string  `json:"serverName,omitempty"`
	PrivateIP     string  `json:"privateIp"`
	InterfaceName string  `json:"interfaceName"`
	Endpoint      string  `json:"endpoint,omitempty"`
	PublicKey     string  `json:"publicKey,omitempty"`
	Status        string  `json:"status"`
	LatencyMS     float64 `json:"latencyMs,omitempty"`
	PacketLoss    float64 `json:"packetLoss,omitempty"`
}

type NetworkAccessPolicy struct {
	Mode         string      `json:"mode"`
	AllowedPairs [][2]string `json:"allowedPairs,omitempty"`
}

type NetworkRollbackWatchdog struct {
	OperationID      string `json:"operationId"`
	Armed            bool   `json:"armed"`
	TimeoutSeconds   int    `json:"timeoutSeconds"`
	RemainingSeconds int    `json:"remainingSeconds"`
	ExpiresAt        string `json:"expiresAt"`
	CanCommit        bool   `json:"canCommit"`
	CanRollback      bool   `json:"canRollback"`
	Status           string `json:"status"`
}

type ServerNetwork struct {
	ID               string                  `json:"id" db:"id"`
	OrganizationID   string                  `json:"organizationId" db:"organization_id"`
	Name             string                  `json:"name" db:"name"`
	Mode             string                  `json:"mode" db:"mode"`
	SubnetCIDR       string                  `json:"subnetCidr" db:"subnet_cidr"`
	MTU              int                     `json:"mtu" db:"mtu"`
	VerificationPort int                     `json:"verificationPort" db:"verification_port"`
	Status           string                  `json:"status" db:"status"`
	Revision         int                     `json:"revision" db:"revision"`
	Members          []NetworkMember         `json:"members"`
	AccessPolicy     NetworkAccessPolicy     `json:"accessPolicy"`
	Watchdog         *NetworkRollbackWatchdog `json:"watchdog,omitempty"`
	CreatedAt        string                  `json:"createdAt" db:"created_at"`
	UpdatedAt        string                  `json:"updatedAt" db:"updated_at"`
}

type CreateNetworkRequest struct {
	Name             string              `json:"name"`
	Mode             string              `json:"mode"`
	SubnetCIDR       string              `json:"subnetCidr"`
	MTU              int                 `json:"mtu,omitempty"`
	VerificationPort int                 `json:"verificationPort,omitempty"`
	Members          []NetworkMember     `json:"members"`
	AccessPolicy     NetworkAccessPolicy `json:"accessPolicy"`
}

type ReviseNetworkAccessRequest struct {
	NetworkID    string              `json:"networkId"`
	AccessPolicy NetworkAccessPolicy `json:"accessPolicy"`
}
