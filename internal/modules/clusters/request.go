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

type RedisSnapshotApplyRequest struct {
	S3DestinationID string `json:"s3DestinationId"`
}

type RedisRestoreApplyRequest struct {
	SnapshotID string `json:"snapshotId"`
}

type RedisFailoverApplyRequest struct {
	NodeID string `json:"nodeId"`
}

type BindingApplyRequest struct {
	ServiceID string `json:"serviceId"`
}

type LifecycleApplyRequest struct {
	Action string `json:"action"`
}
