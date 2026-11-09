package compose

type ComposeStackRequest struct {
	ProjectID     string            `json:"projectId,omitempty"`
	RepositoryURL string            `json:"repositoryUrl,omitempty"`
	Branch        string            `json:"branch,omitempty"`
	RootDirectory string            `json:"rootDirectory,omitempty"`
	ID            string            `json:"id"`
	EnvironmentID string            `json:"environmentId"`
	Name          string            `json:"name"`
	Content       string            `json:"content"`
	Variables     map[string]string `json:"variables"`
	Revision      int               `json:"revision"`
	Digest        string            `json:"digest"`
}

type ComposeStack struct {
	ID            string `json:"id" db:"id"`
	ProjectID     string `json:"projectId" db:"project_id"`
	EnvironmentID string `json:"environmentId" db:"environment_id"`
	Name          string `json:"name" db:"name"`
	Revision      int    `json:"revision" db:"revision"`
	Status        string `json:"status" db:"status"`
	Error         string `json:"error" db:"error"`
	Results       string `json:"results" db:"results"`
	UpdatedAt     string `json:"updatedAt" db:"updated_at"`
	Config        string `json:"-" db:"encrypted_config"`
}

type ComposeSourceChanges struct {
	Added   []string `json:"added"`
	Changed []string `json:"changed"`
	Removed []string `json:"removed"`
}

type ComposeReview struct {
	Config   string                `json:"config"`
	Digest   string                `json:"digest"`
	Services []string              `json:"services"`
	Effects  []string              `json:"effects"`
	Changes  *ComposeSourceChanges `json:"changes,omitempty"`
}

type ComposeServiceResult struct {
	Name        string `json:"name"`
	ContainerID string `json:"containerId"`
	State       string `json:"state"`
	Health      string `json:"health"`
	ExitCode    int    `json:"exitCode"`
}

type DatabaseEngine = string

const (
	DatabaseEnginePostgres   DatabaseEngine = "postgres"
	DatabaseEngineMySQL      DatabaseEngine = "mysql"
	DatabaseEngineRedis      DatabaseEngine = "redis"
	DatabaseEngineMongoDB    DatabaseEngine = "mongodb"
	DatabaseEngineMariaDB    DatabaseEngine = "mariadb"
	DatabaseEngineClickhouse DatabaseEngine = "clickhouse"

	RuntimeModeWeb        = "web"
	BuildEngineDockerfile = "dockerfile"
)

type CreateAppServiceRequest struct {
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	ImageRef    string `json:"imageRef,omitempty"`
	RuntimeMode string `json:"runtimeMode,omitempty"`
	BuildEngine string `json:"buildEngine,omitempty"`
}

type CreateDatabaseRequest struct {
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Engine    string `json:"engine"`
	Version   string `json:"version"`
	Port      int    `json:"port"`
}

type UserComposeFile struct {
	Services map[string]UserComposeService `yaml:"services" json:"services"`
	Networks map[string]struct {
		External bool `yaml:"external" json:"external"`
	} `yaml:"networks" json:"networks"`
}

type UserComposeService struct {
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Build       any               `yaml:"build,omitempty" json:"build,omitempty"`
	Ports       []string          `yaml:"ports,omitempty" json:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	EnvFile     string            `yaml:"env_file,omitempty" json:"env_file,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	Restart     string            `yaml:"restart,omitempty" json:"restart,omitempty"`
}
