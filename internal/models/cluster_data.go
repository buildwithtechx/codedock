package models

type ClusterDataSpec struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	EnvironmentID   string `json:"environmentId"`
	Engine          string `json:"engine"`
	Image           string `json:"image"`
	Instances       int    `json:"instances"`
	StorageGiB      int    `json:"storageGiB"`
	StorageClass    string `json:"storageClass"`
	S3DestinationID string `json:"s3DestinationId"`
}
type ClusterData struct {
	ID        string          `json:"id" db:"id"`
	ClusterID string          `json:"clusterId" db:"cluster_id"`
	ProjectID string          `json:"projectId" db:"project_id"`
	Spec      ClusterDataSpec `json:"spec" db:"-"`
	Status    string          `json:"status" db:"status"`
	Error     string          `json:"error" db:"error"`
	Config    string          `json:"-" db:"encrypted_config"`
	UpdatedAt string          `json:"updatedAt" db:"updated_at"`
}
type OperatorManifest struct {
	Name     string `json:"name"`
	Manifest string `json:"manifest"`
	SHA256   string `json:"sha256"`
}
type ClusterDataPlan struct {
	Record      ClusterData        `json:"record"`
	Password    string             `json:"password"`
	Manifest    string             `json:"manifest"`
	Action      string             `json:"action"`
	SourceID    string             `json:"sourceId"`
	RestoreTime string             `json:"restoreTime"`
	Operators   []OperatorManifest `json:"operators"`
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
type ClusterDataCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}
