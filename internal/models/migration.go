package models

type MigrationMode string

const (
	MigrationModeMove MigrationMode = "move"
	MigrationModeCopy MigrationMode = "copy"
)

type MigrationStatus string

const (
	MigrationStatusPending        MigrationStatus = "pending"
	MigrationStatusScanning       MigrationStatus = "scanning"
	MigrationStatusPreviewed      MigrationStatus = "previewed"
	MigrationStatusTransferring   MigrationStatus = "transferring"
	MigrationStatusVerifying      MigrationStatus = "verifying"
	MigrationStatusAwaitingCutover MigrationStatus = "awaiting_cutover"
	MigrationStatusCompleted      MigrationStatus = "completed"
	MigrationStatusFailed         MigrationStatus = "failed"
	MigrationStatusCancelled      MigrationStatus = "cancelled"
)

type MigrationPromptKind string

const (
	MigrationPromptVolumeConflict MigrationPromptKind = "volume-conflict"
	MigrationPromptPortConflict   MigrationPromptKind = "port-conflict"
)

type MigrationSource struct {
	ID            string `json:"id" db:"id"`
	OrganizationID string `json:"organizationId" db:"organization_id"`
	Name          string `json:"name" db:"name"`
	SSHHost       string `json:"sshHost" db:"ssh_host"`
	SSHPort       int    `json:"sshPort" db:"ssh_port"`
	SSHUser       string `json:"sshUser" db:"ssh_user"`
	SSHAuthMethod string `json:"sshAuthMethod" db:"ssh_auth_method"`
	SSHKey        string `json:"-" db:"ssh_key"`
	SSHPassword   string `json:"-" db:"ssh_password"`
	Fingerprint   string `json:"fingerprint" db:"fingerprint"`
	CreatedAt     string `json:"createdAt" db:"created_at"`
	UpdatedAt     string `json:"updatedAt" db:"updated_at"`
}

type MigrationRun struct {
	ID             string          `json:"id" db:"id"`
	OrganizationID string          `json:"organizationId" db:"organization_id"`
	UserID         string          `json:"userId" db:"user_id"`
	SourceID       string          `json:"sourceId,omitempty" db:"source_id"`
	SourceKind     string          `json:"sourceKind" db:"source_kind"`
	ProjectID      string          `json:"projectId,omitempty" db:"project_id"`
	TargetServerID string          `json:"targetServerId,omitempty" db:"target_server_id"`
	Mode           MigrationMode   `json:"mode" db:"mode"`
	Status         MigrationStatus `json:"status" db:"status"`
	Phase          string          `json:"phase" db:"phase"`
	Selection      string          `json:"-" db:"selection"`
	Progress       string          `json:"-" db:"progress"`
	Logs           string          `json:"logs" db:"logs"`
	Prompt         string          `json:"-" db:"prompt"`
	TokenHash      string          `json:"-" db:"token_hash"`
	Error          string          `json:"error,omitempty" db:"error"`
	CancelRequested bool           `json:"cancelRequested" db:"cancel_requested"`
	CreatedAt      string          `json:"createdAt" db:"created_at"`
	UpdatedAt      string          `json:"updatedAt" db:"updated_at"`
}

type MigrationSelection struct {
	ContainerIDs  []string          `json:"containerIds"`
	Names         []string          `json:"names"`
	Overrides     map[string]string `json:"overrides"`
	Skips         []string          `json:"skips"`
	Decisions     map[string]string `json:"decisions"`
	KillOriginals bool              `json:"killOriginals"`
	ProjectName   string            `json:"projectName"`
	ImportEnv     bool              `json:"importEnv"`
	SourceServerID string           `json:"sourceServerId"`
}

type MigrationPromptOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type MigrationPrompt struct {
	ID        string               `json:"id"`
	Kind      MigrationPromptKind  `json:"kind"`
	Subject   string               `json:"subject"`
	Detail    string               `json:"detail"`
	Options   []MigrationPromptOption `json:"options"`
	ExpiresAt int64                `json:"expiresAt"`
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

type CreateMigrationSourceRequest struct {
	Name          string `json:"name"`
	SSHHost       string `json:"sshHost"`
	SSHPort       int    `json:"sshPort"`
	SSHUser       string `json:"sshUser"`
	SSHAuthMethod string `json:"sshAuthMethod"`
	SSHKey        string `json:"sshKey"`
	SSHPassword   string `json:"sshPassword"`
	Fingerprint   string `json:"fingerprint"`
}

type ScanMigrationRequest struct {
	SourceID string `json:"sourceId"`
}

type RevealMigrationEnvRequest struct {
	SourceID    string `json:"sourceId"`
	ContainerID string `json:"containerId"`
}

type AdoptMigrationRequest struct {
	SourceID     string   `json:"sourceId"`
	ContainerIDs []string `json:"containerIds"`
	ProjectName  string   `json:"projectName"`
	ImportEnv    bool     `json:"importEnv"`
}

type ReimportMigrationRequest struct {
	SourceID     string   `json:"sourceId"`
	ContainerIDs []string `json:"containerIds"`
}

type RepoComposeMigrationRequest struct {
	RepoURL   string `json:"repoUrl"`
	Branch    string `json:"branch"`
	ComposePath string `json:"composePath"`
}

type PreviewMigrationRequest struct {
	SourceID       string                 `json:"sourceId"`
	ProjectID      string                 `json:"projectId"`
	ContainerIDs   []string               `json:"containerIds"`
	TargetServerID string                 `json:"targetServerId"`
	Mode           MigrationMode          `json:"mode"`
}

type StartMigrationRequest struct {
	SourceID       string   `json:"sourceId"`
	ContainerIDs   []string `json:"containerIds"`
	ProjectName    string   `json:"projectName"`
	TargetServerID string   `json:"targetServerId"`
	Mode           MigrationMode `json:"mode"`
	ImportEnv      bool     `json:"importEnv"`
	KillOriginals  bool     `json:"killOriginals"`
}

type StartProjectMoveRequest struct {
	ProjectID      string   `json:"projectId"`
	ServiceIDs     []string `json:"serviceIds"`
	TargetServerID string   `json:"targetServerId"`
	Mode           MigrationMode `json:"mode"`
	KillOriginals  bool     `json:"killOriginals"`
}

type CutoverMigrationRequest struct {
	Confirmation string `json:"confirmation"`
	Kill         bool   `json:"kill"`
}

type RespondMigrationRequest struct {
	PromptID string `json:"promptId"`
	OptionID string `json:"optionId"`
	Value    string `json:"value"`
}

type ResumeMigrationRequest struct {
	Overrides map[string]string `json:"overrides"`
	Skips     []string          `json:"skips"`
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
	Platform        TakeoverPlatform  `json:"platform"`
	Host            string            `json:"host"`
}
