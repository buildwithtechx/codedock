package networks

type ServerNetworkSettings struct {
	InternetAccess *bool    `json:"internetAccess"`
	IngressPorts   []int    `json:"ingressPorts"`
	Egress         []string `json:"egress"`
	PrivateIP      string   `json:"privateIp"`
	OutboundIP     string   `json:"outboundIp"`
	OutboundMode   string   `json:"outboundMode"`
	Revision       string   `json:"revision"`
}

type ServerNetwork struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	Mode             string                   `json:"mode"`
	SubnetCIDR       string                   `json:"subnetCidr"`
	MTU              int                      `json:"mtu"`
	VerificationPort int                      `json:"verificationPort"`
	Status           string                   `json:"status"`
	Revision         int                      `json:"revision"`
	Members          any                      `json:"members"`
	AccessPolicy     any                      `json:"accessPolicy"`
	Watchdog         *NetworkRollbackWatchdog `json:"watchdog,omitempty"`
	CreatedAt        string                   `json:"createdAt"`
	UpdatedAt        string                   `json:"updatedAt"`
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

type CreateNetworkRequest struct {
	Name             string `json:"name"`
	Mode             string `json:"mode,omitempty"`
	SubnetCIDR       string `json:"subnetCidr,omitempty"`
	MTU              int    `json:"mtu,omitempty"`
	VerificationPort int    `json:"verificationPort,omitempty"`
	Members          any    `json:"members,omitempty"`
	AccessPolicy     any    `json:"accessPolicy,omitempty"`
}

type ReviseNetworkAccessRequest struct {
	AccessPolicy any `json:"accessPolicy"`
}
