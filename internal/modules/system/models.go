package system

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

type UpdateInfo struct {
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	HasUpdate       bool   `json:"hasUpdate"`
	ReleaseNotes    string `json:"releaseNotes"`
	DownloadURL     string `json:"downloadUrl"`
	LastChecked     string `json:"lastChecked"`
	AutoUpdate      bool   `json:"autoUpdate"`
	UpdateCheckCron string `json:"updateCheckCron"`
}
