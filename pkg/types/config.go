package types

type Config struct {
	Server    ServerConfig
	Security  SecurityConfig
	Docker    DockerConfig
	Traefik   TraefikConfig
	Builder   BuilderConfig
	Defaults  DefaultsConfig
	Limits    LimitsConfig
	Domains   DomainsConfig
	Worker    WorkerConfig
	Updates   UpdatesConfig
	Telemetry TelemetryConfig
	Cloud     CloudConfig
	SMTP      SMTPConfig
	Resend    ResendConfig
	OAuth     OAuthConfig
	Stripe    StripeConfig
}

type ServerConfig struct {
	Port         int
	Host         string
	DataDir      string
	HostIP       string
	Domain       string
	DashboardURL string
	ServerURL    string
	ApiHost      string
	StaticDir    string
}

type SecurityConfig struct {
	JWTSecret     string
	RefreshSecret string
	TLSEmail      string
}

type DockerConfig struct {
	SocketPath     string
	RuntimeNetwork string
	PortStart      int
	PortEnd        int
	DryRun         bool
}

type TraefikConfig struct {
	HTTPPort   int
	HTTPSPort  int
	APIPort    int
	Image      string
	DockerHost string
}

type BuilderConfig struct {
	NixpacksImage string
	PackImage     string
	BuilderImage  string
}

type DefaultsConfig struct {
	AppPort    int
	MemoryMB   int64
	CPU        float64
	DBMemoryMB int64
	DBCPU      float64
}

type LimitsConfig struct {
	MaxConcurrentBuilds int
	DeploymentTimeout   int
}

type DomainsConfig struct {
	WildcardDomain string
	MagicDomain    string
}

type WorkerConfig struct {
	UseWSS      bool
	WorkerToken string
}

type UpdatesConfig struct {
	UpdateURL   string
	DownloadURL string
}

type TelemetryConfig struct {
	PosthogKey  string
	PosthogHost string
	Salt        string
}

type CloudConfig struct {
	Enabled bool
}

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

type ResendConfig struct {
	APIKey string
}

type OAuthConfig struct {
	GitHubClientID     string
	GitHubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
}

type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	PriceIDPro     string
}
