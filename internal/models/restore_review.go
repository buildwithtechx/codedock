package models

type RestoreReviewRequest struct {
	TargetDatabaseID string `json:"targetDatabaseId"`
}
type RestoreTarget struct {
	RecordID       string `json:"recordId"`
	DatabaseID     string `json:"databaseId"`
	VolumeName     string `json:"volumeName"`
	Engine         string `json:"engine"`
	ContainerID    string `json:"-"`
	Snapshot       string `json:"-"`
	SHA256         string `json:"sha256"`
	Mode           string `json:"mode"`
	ProtectedUntil int64  `json:"protectedUntil"`
}

type BackupProtectionRequest struct {
	ProtectedUntil int64 `json:"protectedUntil"`
}
