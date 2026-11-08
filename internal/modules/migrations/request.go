package migrations

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

type TestMigrationSourceRequest struct {
	SSHHost     string `json:"sshHost"`
	SSHPort     int    `json:"sshPort"`
	SSHUser     string `json:"sshUser"`
	Fingerprint string `json:"fingerprint"`
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
	RepoURL     string `json:"repoUrl"`
	Branch      string `json:"branch"`
	ComposePath string `json:"composePath"`
}

type PreviewMigrationRequest struct {
	SourceID       string        `json:"sourceId"`
	ProjectID      string        `json:"projectId"`
	ContainerIDs   []string      `json:"containerIds"`
	TargetServerID string        `json:"targetServerId"`
	Mode           MigrationMode `json:"mode"`
}

type StartMigrationRequest struct {
	SourceID       string        `json:"sourceId"`
	ContainerIDs   []string      `json:"containerIds"`
	ProjectName    string        `json:"projectName"`
	TargetServerID string        `json:"targetServerId"`
	Mode           MigrationMode `json:"mode"`
	ImportEnv      bool          `json:"importEnv"`
	KillOriginals  bool          `json:"killOriginals"`
}

type StartProjectMoveRequest struct {
	ProjectID      string        `json:"projectId"`
	ServiceIDs     []string      `json:"serviceIds"`
	TargetServerID string        `json:"targetServerId"`
	Mode           MigrationMode `json:"mode"`
	KillOriginals  bool          `json:"killOriginals"`
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

type ExportRequest struct {
	Passphrase string `json:"passphrase"`
}
