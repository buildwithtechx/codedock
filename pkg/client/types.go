package client

import "time"

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProjectConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedBy   string    `json:"createdBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type EnvironmentConfig struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsDefault   bool      `json:"isDefault,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DomainConfig struct {
	ID                 string     `json:"id"`
	ServiceID          string     `json:"serviceId"`
	DomainName         string     `json:"domainName"`
	RedirectTo         string     `json:"redirectTo,omitempty"`
	SSLCertStatus      string     `json:"sslCertStatus,omitempty"`
	Primary            bool       `json:"primary,omitempty"`
	PathPrefix         string     `json:"pathPrefix,omitempty"`
	StripPrefix        bool       `json:"stripPrefix,omitempty"`
	DNSProvisionStatus string     `json:"dnsProvisionStatus,omitempty"`
	DNSProvider        string     `json:"dnsProvider,omitempty"`
	DNSProvisionedIP   string     `json:"dnsProvisionedIP,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	CertificateExpires *time.Time `json:"certificateExpires,omitempty"`
}

type AppService struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	RepositoryURL string    `json:"repositoryUrl"`
	Branch        string    `json:"branch"`
	Domain        string    `json:"domain,omitempty"`
	ImageRef      string    `json:"imageRef,omitempty"`
	ContainerID   string    `json:"containerId,omitempty"`
	DeployToken   string    `json:"deployToken,omitempty"`
	Status        string    `json:"status,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Deployment struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"serviceId"`
	EnvironmentID string    `json:"environmentId"`
	ProjectID     string    `json:"projectId"`
	ImageRef      string    `json:"imageRef"`
	Branch        string    `json:"branch"`
	CommitHash    string    `json:"commitHash"`
	Trigger       string    `json:"trigger"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type ServiceMetric struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUUsage    float64   `json:"cpuUsage"`
	MemoryUsage float64   `json:"memoryUsage"`
}

type DatabaseEngine string

const (
	DatabaseEnginePostgres DatabaseEngine = "postgres"
	DatabaseEngineMySQL    DatabaseEngine = "mysql"
	DatabaseEngineRedis    DatabaseEngine = "redis"
	DatabaseEngineMongoDB  DatabaseEngine = "mongodb"
)

type Database struct {
	ID           string         `json:"id"`
	ProjectID    string         `json:"projectId"`
	Name         string         `json:"name"`
	Engine       DatabaseEngine `json:"engine"`
	Version      string         `json:"version"`
	Status       string         `json:"status"`
	Host         string         `json:"host,omitempty"`
	Port         int            `json:"port,omitempty"`
	Username     string         `json:"username,omitempty"`
	Password     string         `json:"password,omitempty"`
	DatabaseName string         `json:"databaseName,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

type CreateDatabaseRequest struct {
	ProjectID     string         `json:"projectId"`
	EnvironmentID string         `json:"environmentId,omitempty"`
	Name          string         `json:"name"`
	Engine        DatabaseEngine `json:"engine"`
	Version       string         `json:"version,omitempty"`
	Password      string         `json:"password,omitempty"`
	DatabaseName  string         `json:"databaseName,omitempty"`
}

type ImportDatabaseRequest struct {
	SourceURL string `json:"sourceUrl"`
}

type BackupConfig struct {
	ID         string    `json:"id"`
	DatabaseID string    `json:"databaseId"`
	Name       string    `json:"name,omitempty"`
	Status     string    `json:"status,omitempty"`
	Schedule   string    `json:"schedule"`
	Retention  int       `json:"retention"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type BackupRecord struct {
	ID            string     `json:"id"`
	BackupID      string     `json:"backupId"`
	Status        string     `json:"status"`
	FilePath      string     `json:"filePath,omitempty"`
	FileSize      int64      `json:"fileSize,omitempty"`
	FileSizeBytes int64      `json:"fileSizeBytes,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type Server struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	IPAddress             string    `json:"ipAddress"`
	IsLocal               bool      `json:"isLocal,omitempty"`
	Status                string    `json:"status"`
	BuildCacheBudgetBytes *int64    `json:"buildCacheBudgetBytes,omitempty"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type CreateServerRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ipAddress"`
	IsLocal   bool   `json:"isLocal,omitempty"`
	SSHHost   string `json:"sshHost,omitempty"`
	SSHPort   int    `json:"sshPort,omitempty"`
	SSHUser   string `json:"sshUser,omitempty"`
	SSHKey    string `json:"sshKey,omitempty"`
}

type UpdateServerRequest struct {
	Name                  string `json:"name,omitempty"`
	BuildCacheBudgetBytes *int64 `json:"buildCacheBudgetBytes,omitempty"`
}

type ServerNetworkSettings struct {
	ServerID       string   `json:"serverId,omitempty"`
	InternetAccess *bool    `json:"internetAccess,omitempty"`
	IngressPorts   []int    `json:"ingressPorts,omitempty"`
	IngressAll     *bool    `json:"ingressAll,omitempty"`
	Egress         []string `json:"egress,omitempty"`
	Revision       string   `json:"revision,omitempty"`
	OutboundIP     string   `json:"outboundIp,omitempty"`
	OutboundMode   string   `json:"outboundMode,omitempty"`
	TailscaleIP    string   `json:"tailscaleIp,omitempty"`
	WireGuardIP    string   `json:"wireguardIp,omitempty"`
	DefaultGateway string   `json:"defaultGateway,omitempty"`
}

type UpdateServerNetworkSettingsRequest struct {
	InternetAccess         bool     `json:"internetAccess"`
	ExpectedInternetAccess bool     `json:"expectedInternetAccess"`
	Egress                 []string `json:"egress,omitempty"`
	ExpectedRevision       string   `json:"expectedRevision,omitempty"`
	Confirm                bool     `json:"confirm"`
	TailscaleIP            string   `json:"tailscaleIp,omitempty"`
	WireGuardIP            string   `json:"wireguardIp,omitempty"`
}

type VarsRequest map[string]string
type SetEnvVarsRequest map[string]string

type RouteProtectionStatus struct {
	Hostname string `json:"hostname"`
	Enabled  bool   `json:"enabled"`
	Status   string `json:"status,omitempty"`
	HasAuth  bool   `json:"hasAuth"`
}

type RouteProtectionInput struct {
	Hostname string `json:"hostname"`
	Action   string `json:"action,omitempty"`
	Enabled  bool   `json:"enabled"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type ActionWorkflow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Trigger   string    `json:"trigger,omitempty"`
	Enabled   bool      `json:"enabled,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type ActionRun struct {
	ID          string    `json:"id"`
	WorkflowID  string    `json:"workflowId"`
	Name        string    `json:"name,omitempty"`
	Number      int       `json:"number,omitempty"`
	Status      string    `json:"status"`
	TriggeredBy string    `json:"triggeredBy,omitempty"`
	DurationS   int       `json:"durationSeconds,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

type ActionRunner struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Mode    string `json:"mode,omitempty"`
	Status  string `json:"status"`
	Enabled bool   `json:"enabled,omitempty"`
}

type IngressConnection struct {
	ID        string    `json:"id"`
	ServerID  string    `json:"serverId"`
	Name      string    `json:"name"`
	Connector string    `json:"connector"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
