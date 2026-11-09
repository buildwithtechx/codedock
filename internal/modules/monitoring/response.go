package monitoring

type IngressDiscoveryAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IngressDiscoveryTunnel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Remote bool   `json:"remote"`
	Status string `json:"status"`
}

type IngressDiscoveryPolicy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Decision string `json:"decision"`
}

type IngressDiscoveryResult struct {
	Accounts      []IngressDiscoveryAccount `json:"accounts"`
	Tunnels       []IngressDiscoveryTunnel  `json:"tunnels"`
	Policies      []IngressDiscoveryPolicy  `json:"policies"`
	PoliciesError *string                   `json:"policiesError,omitempty"`
}

type SystemInfoResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}
