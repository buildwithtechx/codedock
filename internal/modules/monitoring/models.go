package monitoring

type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type Operation struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"-" db:"user_id"`
	ProjectID string `json:"projectId" db:"project_id"`
	Kind      string `json:"kind" db:"kind"`
	Target    string `json:"target" db:"target"`
	Status    string `json:"status" db:"status"`
	Phase     string `json:"phase" db:"phase"`
	Effects   string `json:"effects" db:"effects"`
	Payload   string `json:"-" db:"encrypted_payload"`
	Snapshot  string `json:"-" db:"snapshot"`
	TokenHash string `json:"-" db:"token_hash"`
	Error     string `json:"error" db:"error"`
	Logs      string `json:"logs" db:"logs"`
	ExpiresAt int64  `json:"expiresAt" db:"expires_at"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type OperationReview struct {
	Operation    *Operation `json:"operation"`
	Confirmation string     `json:"confirmation"`
}

type SystemStats struct {
	CPU       CPUStats    `json:"cpu"`
	Memory    MemoryStats `json:"memory"`
	Disk      DiskStats   `json:"disk"`
	Uptime    int64       `json:"uptimeSeconds"`
	LoadAvg   [3]float64  `json:"loadAvg"`
	Processes int         `json:"processes"`
	Docker    DockerStats `json:"docker"`
}

type DockerStats struct {
	Images        DockerLayerStat `json:"images"`
	Containers    DockerLayerStat `json:"containers"`
	Volumes       DockerLayerStat `json:"volumes"`
	BuildCache    DockerLayerStat `json:"buildCache"`
	ReclaimableGB float64         `json:"reclaimableGb"`
}

type DockerLayerStat struct {
	Active      string `json:"active"`
	TotalCount  string `json:"totalCount"`
	Size        string `json:"size"`
	Reclaimable string `json:"reclaimable"`
}

type CPUStats struct {
	Percent float64 `json:"percent"`
	Cores   int     `json:"cores"`
}

type MemoryStats struct {
	TotalMB int64   `json:"totalMb"`
	UsedMB  int64   `json:"usedMb"`
	FreeMB  int64   `json:"freeMb"`
	Percent float64 `json:"percent"`
}

type DiskStats struct {
	TotalGB int64   `json:"totalGb"`
	UsedGB  int64   `json:"usedGb"`
	FreeGB  int64   `json:"freeGb"`
	Percent float64 `json:"percent"`
}

type UpdateStatus struct {
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	UpdateAvail    bool   `json:"updateAvailable"`
	Updating       bool   `json:"updating"`
}
