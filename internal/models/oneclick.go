package models

type OneClickEnvVar struct {
	Key          string `json:"key"`
	Label        string `json:"label"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Secret       bool   `json:"secret"`
}

type OneClickApp struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Icon         string           `json:"icon"`
	Category     string           `json:"category"`
	DockerImage  string           `json:"dockerImage"`
	DefaultPort  int              `json:"defaultPort"`
	Services     []string         `json:"services"`
	Volumes      []string         `json:"volumes"`
	EnvVariables []OneClickEnvVar `json:"envVariables"`
}

type InstallAppInput struct {
	AppID         string            `json:"appId"`
	ProjectID     string            `json:"projectId"`
	EnvironmentID string            `json:"environmentId,omitempty"`
	Name          string            `json:"name"`
	Secrets       map[string]string `json:"secrets,omitempty"`
	Environment   map[string]string `json:"environment,omitempty"`
	HostPort      int               `json:"hostPort,omitempty"`
	Domain        string            `json:"domain,omitempty"`
	Digest        string            `json:"digest,omitempty"`
}

type InstallPreviewService struct {
	Service string   `json:"service"`
	Image   string   `json:"image"`
	Env     []string `json:"env"`
	Ports   []string `json:"ports"`
}

type InstallPreviewVolume struct {
	Service string `json:"service"`
	Name    string `json:"name"`
	Target  string `json:"target"`
}

type InstallPreview struct {
	AppID            string                  `json:"appId"`
	Name             string                  `json:"name"`
	Services         []InstallPreviewService `json:"services"`
	Volumes          []InstallPreviewVolume  `json:"volumes"`
	GeneratedSecrets []string                `json:"generatedSecrets"`
	ComposeYAML      string                  `json:"composeYaml"`
	Digest           string                  `json:"digest"`
	Kind             string                  `json:"kind"`
}

type AppInstallResult struct {
	Kind  string        `json:"kind"`
	Stack *ComposeStack `json:"stack,omitempty"`
}
