package models

type StorageSetupRequest struct {
	Version string `json:"version"`
}

type StorageSetupPlan struct {
	ClusterID string `json:"clusterId"`
	ProjectID string `json:"projectId"`
	Version   string `json:"version"`
	Manifest  string `json:"manifest"`
	SHA256    string `json:"sha256"`
	Action    string `json:"action"`
}

type VolumeSnapshotRequest struct {
	Claim       string `json:"claim"`
	Namespace   string `json:"namespace"`
	Snapshot    string `json:"snapshot"`
	RestoreName string `json:"restoreName"`
}
