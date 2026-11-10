package networks

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
