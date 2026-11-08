package clusters

type ClusterReviewRequest struct {
	Cluster Cluster `json:"cluster"`
	Action  string  `json:"action"`
}

type StorageSetupRequest struct {
	Version string `json:"version"`
}

type VolumeSnapshotRequest struct {
	Claim       string `json:"claim"`
	Namespace   string `json:"namespace"`
	Snapshot    string `json:"snapshot"`
	RestoreName string `json:"restoreName"`
}

type ClusterDataRequest struct {
	Action             string          `json:"action"`
	Spec               ClusterDataSpec `json:"spec"`
	SourceID           string          `json:"sourceId"`
	RestoreTime        string          `json:"restoreTime"`
	OperatorVersion    string          `json:"operatorVersion"`
	BarmanVersion      string          `json:"barmanVersion"`
	CertManagerVersion string          `json:"certManagerVersion"`
}

type OperationApply struct {
	Confirmation string `json:"confirmation"`
}

type BindingReviewRequest struct {
	AppID      string `json:"appId"`
	DatabaseID string `json:"databaseId"`
}

type BindingApplyRequest struct {
	OperationID  string `json:"operationId"`
	Confirmation string `json:"confirmation"`
}

type LifecycleReviewRequest struct {
	DatabaseID string `json:"databaseId"`
	Action     string `json:"action"`
}

type LifecycleApplyRequest struct {
	OperationID  string `json:"operationId"`
	Confirmation string `json:"confirmation"`
}

type RedisSnapshotApplyRequest struct {
	OperationID  string `json:"operationId"`
	Confirmation string `json:"confirmation"`
}

type RedisRestoreReviewRequest struct {
	SnapshotID string `json:"snapshotId"`
	TargetID   string `json:"targetId"`
}

type RedisRestoreApplyRequest struct {
	OperationID  string `json:"operationId"`
	Confirmation string `json:"confirmation"`
}

type RedisFailoverApplyRequest struct {
	OperationID  string `json:"operationId"`
	Confirmation string `json:"confirmation"`
}
