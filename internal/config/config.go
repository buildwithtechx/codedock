package config

import (
	"os"
	"strconv"
	"sync"

	"codedock.run/codedock/pkg/types"
)

var (
	globalConfig *types.Config
	once         sync.Once
)

func Get() *types.Config {
	once.Do(func() {
		globalConfig = Load()
	})
	return globalConfig
}

func Load() *types.Config {
	return &types.Config{
		Server: types.ServerConfig{
			Port:         getEnvInt("PORT", 8080),
			Host:         getEnv("HOST", "0.0.0.0"),
			DataDir:      getEnv("CODEDOCK_DATA_DIR", "data"),
			HostIP:       getEnv("CODEDOCK_HOST_IP", ""),
			Domain:       getEnv("CODEDOCK_DOMAIN", ""),
			DashboardURL: getEnv("CODEDOCK_DASHBOARD_URL", "http://localhost:3000"),
			ServerURL:    getEnv("CODEDOCK_SERVER_URL", "http://localhost:8080"),
			ApiHost:      getEnv("CODEDOCK_API_HOST", "http://localhost:8080"),
			StaticDir:    getEnv("CODEDOCK_STATIC_DIR", ""),
		},
		Security: types.SecurityConfig{
			JWTSecret:     getEnv("CODEDOCK_JWT_SECRET", ""),
			RefreshSecret: getEnv("CODEDOCK_REFRESH_SECRET", ""),
			TLSEmail:      getEnv("CODEDOCK_TLS_EMAIL", ""),
		},
		Docker: types.DockerConfig{
			SocketPath:     getEnv("DOCKER_SOCKET_PATH", "/var/run/docker.sock"),
			RuntimeNetwork: getEnv("CODEDOCK_RUNTIME_NETWORK", "codedock-network"),
			PortStart:      getEnvInt("DEPLOY_HOST_PORT_START", 4100),
			PortEnd:        getEnvInt("DEPLOY_HOST_PORT_END", 4999),
			DryRun:         os.Getenv("DEPLOY_DRY_RUN") == "true",
		},
		Traefik: types.TraefikConfig{
			HTTPPort:   getEnvInt("CODEDOCK_TRAEFIK_HTTP_PORT", 80),
			HTTPSPort:  getEnvInt("CODEDOCK_TRAEFIK_HTTPS_PORT", 443),
			APIPort:    getEnvInt("CODEDOCK_TRAEFIK_API_PORT", 8082),
			Image:      getEnv("CODEDOCK_TRAEFIK_IMAGE", "traefik:v3.3"),
			DockerHost: getEnv("CODEDOCK_TRAEFIK_DOCKER_HOST", ""),
		},
		Builder: types.BuilderConfig{
			NixpacksImage: getEnv("CODEDOCK_NIXPACKS_IMAGE", "ghcr.io/railwayapp/nixpacks:latest"),
			PackImage:     getEnv("CODEDOCK_PACK_IMAGE", "buildpacksio/pack:latest"),
			BuilderImage:  getEnv("CODEDOCK_BUILDER_IMAGE", ""),
		},
		Defaults: types.DefaultsConfig{
			AppPort:    getEnvInt("CODEDOCK_DEFAULT_APP_PORT", 3000),
			MemoryMB:   int64(getEnvInt("CODEDOCK_DEFAULT_MEMORY_MB", 512)),
			CPU:        getEnvFloat("CODEDOCK_DEFAULT_CPU", 0.5),
			DBMemoryMB: int64(getEnvInt("CODEDOCK_DEFAULT_DB_MEMORY_MB", 1024)),
			DBCPU:      getEnvFloat("CODEDOCK_DEFAULT_DB_CPU", 1.0),
		},
		Limits: types.LimitsConfig{
			MaxConcurrentBuilds: getEnvInt("CODEDOCK_MAX_CONCURRENT_BUILDS", 2),
			DeploymentTimeout:   getEnvInt("CODEDOCK_DEPLOYMENT_TIMEOUT", 1800),
		},
		Domains: types.DomainsConfig{
			WildcardDomain: getEnv("CODEDOCK_WILDCARD_DOMAIN", ""),
			MagicDomain:    getEnv("CODEDOCK_MAGIC_DOMAIN", ""),
		},
		Worker: types.WorkerConfig{
			UseWSS:      os.Getenv("CODEDOCK_USE_WSS") == "true",
			WorkerToken: getEnv("CODEDOCK_WORKER_TOKEN", ""),
		},
		Updates: types.UpdatesConfig{
			UpdateURL:   getEnv("CODEDOCK_UPDATE_URL", ""),
			DownloadURL: getEnv("CODEDOCK_DOWNLOAD_URL", ""),
		},
		Telemetry: types.TelemetryConfig{
			PosthogKey:  getEnv("POSTHOG_API_KEY", ""),
			PosthogHost: getEnv("POSTHOG_HOST", "https://us.i.posthog.com"),
			Salt:        getEnv("POSTHOG_DISTINCT_ID_SALT", ""),
		},
		Cloud: types.CloudConfig{
			Enabled: os.Getenv("CODEDOCK_CLOUD_MODE") == "true",
		},
		SMTP: types.SMTPConfig{
			Host: getEnv("SMTP_HOST", ""),
			Port: getEnvInt("SMTP_PORT", 587),
			User: getEnv("SMTP_USER", ""),
			Pass: getEnv("SMTP_PASS", ""),
			From: getEnv("SMTP_FROM", ""),
		},
		Resend: types.ResendConfig{
			APIKey: getEnv("RESEND_API_KEY", ""),
		},
		OAuth: types.OAuthConfig{
			GitHubClientID:     firstEnv("GITHUB_CLIENT_ID", "CODEDOCK_GITHUB_CLIENT_ID"),
			GitHubClientSecret: firstEnv("GITHUB_CLIENT_SECRET", "CODEDOCK_GITHUB_CLIENT_SECRET"),
			GoogleClientID:     firstEnv("GOOGLE_CLIENT_ID", "CODEDOCK_GOOGLE_CLIENT_ID"),
			GoogleClientSecret: firstEnv("GOOGLE_CLIENT_SECRET", "CODEDOCK_GOOGLE_CLIENT_SECRET"),
		},
		Stripe: types.StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
			PriceIDPro:     getEnv("STRIPE_PRICE_ID_PRO", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if val := os.Getenv(key); val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return fallback
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if val := os.Getenv(k); val != "" {
			return val
		}
	}
	return ""
}
