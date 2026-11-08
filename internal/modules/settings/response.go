package settings

type PublicSettingsResponse struct {
	RegistrationEnabled   bool   `json:"registrationEnabled"`
	SiteName              string `json:"siteName"`
	EmailEnabled          bool   `json:"emailEnabled"`
	CloudMode             bool   `json:"cloudMode"`
	SelfHosted            bool   `json:"selfHosted"`
	Version               string `json:"version"`
	HostDomain            string `json:"hostDomain"`
	DefaultWildcardDomain string `json:"defaultWildcardDomain"`
	ServerIP              string `json:"serverIp"`
}

type HealthEnvResponse struct {
	SelfHosted            bool     `json:"selfHosted"`
	CloudMode             bool     `json:"cloudMode"`
	DeployMode            string   `json:"deployMode"`
	IsServerHost          bool     `json:"isServerHost"`
	HostControlEnabled    bool     `json:"hostControlEnabled"`
	Version               string   `json:"version"`
	AuthMode              string   `json:"authMode"`
	AuthProviders         []string `json:"authProviders"`
	ProductMode           string   `json:"productMode"`
	TeamMode              string   `json:"teamMode"`
	HostDomain            string   `json:"hostDomain"`
	DefaultWildcardDomain string   `json:"defaultWildcardDomain"`
	ServerIP              string   `json:"serverIp"`
}
