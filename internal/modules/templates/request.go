package templates

// CustomAppRequest defines the payload for creating a custom one-click app template.
type CustomAppRequest struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Category    string             `json:"category,omitempty"`
	Logo        string             `json:"logo,omitempty"`
	Compose     string             `json:"compose"`
	EnvVars     []OneClickEnvVar   `json:"envVars,omitempty"`
	Endpoints   []OneClickEndpoint `json:"endpoints,omitempty"`
}
