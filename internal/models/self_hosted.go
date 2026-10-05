package models

type SelfHostedConfig struct {
	JWTSecret      string `json:"jwtSecret"`
	RefreshSecret  string `json:"refreshSecret"`
	TelemetrySalt  string `json:"telemetrySalt"`
	TLSEmail       string `json:"tlsEmail,omitempty"`
	WildcardDomain string `json:"wildcardDomain,omitempty"`
}
