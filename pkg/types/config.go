package types

type Config struct {
	Server    ServerConfig    `json:"server"`
	Security  SecurityConfig  `json:"security"`
	Docker    DockerConfig    `json:"docker"`
	Traefik   TraefikConfig   `json:"traefik"`
	Builder   BuilderConfig   `json:"builder"`
	Defaults  DefaultsConfig  `json:"defaults"`
	Limits    LimitsConfig    `json:"limits"`
	Domains   DomainsConfig   `json:"domains"`
	Worker    WorkerConfig    `json:"worker"`
	Updates   UpdatesConfig   `json:"updates"`
	Telemetry TelemetryConfig `json:"telemetry"`
	Cloud     CloudConfig     `json:"cloud"`
	SMTP      SMTPConfig      `json:"smtp"`
	Resend    ResendConfig    `json:"resend"`
	OAuth     OAuthConfig     `json:"oauth"`
	Stripe    StripeConfig    `json:"stripe"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	Host         string `json:"host"`
	DataDir      string `json:"dataDir"`
	HostIP       string `json:"hostIp"`
	Domain       string `json:"domain"`
	DashboardURL string `json:"dashboardUrl"`
	ServerURL    string `json:"serverUrl"`
	APIHost      string `json:"apiHost"`
	StaticDir    string `json:"staticDir"`
}

type SecurityConfig struct {
	JWTSecret     string `json:"jwtSecret"`
	RefreshSecret string `json:"refreshSecret"`
	TLSEmail      string `json:"tlsEmail"`
}

type DockerConfig struct {
	SocketPath     string `json:"socketPath"`
	RuntimeNetwork string `json:"runtimeNetwork"`
	PortStart      int    `json:"portStart"`
	PortEnd        int    `json:"portEnd"`
	DryRun         bool   `json:"dryRun"`
}

type TraefikConfig struct {
	HTTPPort   int    `json:"httpPort"`
	HTTPSPort  int    `json:"httpsPort"`
	APIPort    int    `json:"apiPort"`
	Image      string `json:"image"`
	DockerHost string `json:"dockerHost"`
}

type BuilderConfig struct {
	NixpacksImage string `json:"nixpacksImage"`
	PackImage     string `json:"packImage"`
	BuilderImage  string `json:"builderImage"`
}

type DefaultsConfig struct {
	AppPort    int     `json:"appPort"`
	MemoryMB   int64   `json:"memoryMb"`
	CPU        float64 `json:"cpu"`
	DBMemoryMB int64   `json:"dbMemoryMb"`
	DBCPU      float64 `json:"dbCpu"`
}

type LimitsConfig struct {
	MaxConcurrentBuilds int `json:"maxConcurrentBuilds"`
	DeploymentTimeout   int `json:"deploymentTimeout"`
}

type DomainsConfig struct {
	WildcardDomain string `json:"wildcardDomain"`
	MagicDomain    string `json:"magicDomain"`
}

type WorkerConfig struct {
	UseWSS      bool   `json:"useWss"`
	WorkerToken string `json:"workerToken"`
}

type UpdatesConfig struct {
	UpdateURL   string `json:"updateUrl"`
	DownloadURL string `json:"downloadUrl"`
}

type TelemetryConfig struct {
	PosthogKey  string `json:"posthogKey"`
	PosthogHost string `json:"posthogHost"`
	Salt        string `json:"salt"`
}

type CloudConfig struct {
	Enabled bool `json:"enabled"`
}

type SMTPConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	From string `json:"from"`
}

type ResendConfig struct {
	APIKey string `json:"apiKey"`
}

type OAuthConfig struct {
	GitHubClientID     string `json:"githubClientId"`
	GitHubClientSecret string `json:"githubClientSecret"`
	GoogleClientID     string `json:"googleClientId"`
	GoogleClientSecret string `json:"googleClientSecret"`
}

type StripeConfig struct {
	SecretKey      string `json:"secretKey"`
	PublishableKey string `json:"publishableKey"`
	WebhookSecret  string `json:"webhookSecret"`
	PriceIDPro     string `json:"priceIdPro"`
}
