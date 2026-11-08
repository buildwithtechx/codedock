package migrations

type Vault interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) ([]byte, error)
}

type MigrationSource struct {
	ID             string `json:"id" db:"id"`
	OrganizationID string `json:"organizationId" db:"organization_id"`
	Name           string `json:"name" db:"name"`
	SSHHost        string `json:"sshHost" db:"ssh_host"`
	SSHPort        int    `json:"sshPort" db:"ssh_port"`
	SSHUser        string `json:"sshUser" db:"ssh_user"`
	SSHAuthMethod  string `json:"sshAuthMethod" db:"ssh_auth_method"`
	SSHKey         string `json:"-" db:"ssh_key"`
	SSHPassword    string `json:"-" db:"ssh_password"`
	Fingerprint    string `json:"fingerprint" db:"fingerprint"`
	CreatedAt      string `json:"createdAt" db:"created_at"`
	UpdatedAt      string `json:"updatedAt" db:"updated_at"`
}

type MigrationRun struct {
	ID              string          `json:"id" db:"id"`
	OrganizationID  string          `json:"organizationId" db:"organization_id"`
	UserID          string          `json:"userId" db:"user_id"`
	SourceID        string          `json:"sourceId,omitempty" db:"source_id"`
	SourceKind      string          `json:"sourceKind" db:"source_kind"`
	ProjectID       string          `json:"projectId,omitempty" db:"project_id"`
	TargetServerID  string          `json:"targetServerId,omitempty" db:"target_server_id"`
	Mode            MigrationMode   `json:"mode" db:"mode"`
	Status          MigrationStatus `json:"status" db:"status"`
	Phase           string          `json:"phase" db:"phase"`
	Selection       string          `json:"-" db:"selection"`
	Progress        string          `json:"-" db:"progress"`
	Logs            string          `json:"logs" db:"logs"`
	Prompt          string          `json:"-" db:"prompt"`
	TokenHash       string          `json:"-" db:"token_hash"`
	Error           string          `json:"error,omitempty" db:"error"`
	CancelRequested bool            `json:"cancelRequested" db:"cancel_requested"`
	CreatedAt       string          `json:"createdAt" db:"created_at"`
	UpdatedAt       string          `json:"updatedAt" db:"updated_at"`
}

type MigrationSelection struct {
	ContainerIDs   []string          `json:"containerIds"`
	Names          []string          `json:"names"`
	Overrides      map[string]string `json:"overrides"`
	Skips          []string          `json:"skips"`
	Decisions      map[string]string `json:"decisions"`
	KillOriginals  bool              `json:"killOriginals"`
	ProjectName    string            `json:"projectName"`
	ImportEnv      bool              `json:"importEnv"`
	SourceServerID string            `json:"sourceServerId"`
}

type MigrationPromptOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type MigrationPrompt struct {
	ID        string                  `json:"id"`
	Kind      MigrationPromptKind     `json:"kind"`
	Subject   string                  `json:"subject"`
	Detail    string                  `json:"detail"`
	Options   []MigrationPromptOption `json:"options"`
	ExpiresAt int64                   `json:"expiresAt"`
}

type MigrationServiceProgress struct {
	Service       string   `json:"service"`
	ServiceID     string   `json:"serviceId"`
	ContainerID   string   `json:"containerId"`
	SourceState   string   `json:"sourceState"`
	CopiedVolumes []string `json:"copiedVolumes"`
	Adopted       bool     `json:"adopted"`
	Transferred   bool     `json:"transferred"`
	Deployed      bool     `json:"deployed"`
	Verified      bool     `json:"verified"`
	Skipped       bool     `json:"skipped"`
	Error         string   `json:"error,omitempty"`
}

type MigrationPreview struct {
	Images    []string          `json:"images"`
	Volumes   []string          `json:"volumes"`
	Services  []PreviewService  `json:"services"`
	Conflicts []PreviewConflict `json:"conflicts"`
	Warnings  []string          `json:"warnings"`
	Downtime  string            `json:"downtime"`
}

type PreviewService struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Ports   []string `json:"ports"`
	Volumes []string `json:"volumes"`
	Status  string   `json:"status"`
}

type PreviewConflict struct {
	Kind    string `json:"kind"`
	Subject string `json:"subject"`
	Detail  string `json:"detail"`
}

type MaskedContainer struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Image          string            `json:"image"`
	Ports          []string          `json:"ports"`
	EnvKeys        []string          `json:"envKeys"`
	Env            map[string]string `json:"env"`
	Volumes        []string          `json:"volumes"`
	Labels         map[string]string `json:"labels"`
	Status         string            `json:"status"`
	ComposeProject string            `json:"composeProject,omitempty"`
	Routes         []string          `json:"routes"`
}

type MaskedStack struct {
	Containers      []MaskedContainer `json:"containers"`
	ComposeProjects []string          `json:"composeProjects"`
	Platform        string            `json:"platform"`
	Host            string            `json:"host"`
}
