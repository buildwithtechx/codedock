package models

type ClusterDataSpec struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	EnvironmentID   string `json:"environmentId"`
	Engine          string `json:"engine"`
	Image           string `json:"image"`
	Instances       int    `json:"instances"`
	Shards          int    `json:"shards"`
	StorageGiB      int    `json:"storageGiB"`
	StorageClass    string `json:"storageClass"`
	S3DestinationID string `json:"s3DestinationId"`
	Synchronous     bool   `json:"synchronous"`
}
type ClusterData struct {
	ID              string          `json:"id" db:"id"`
	ClusterID       string          `json:"clusterId" db:"cluster_id"`
	ProjectID       string          `json:"projectId" db:"project_id"`
	Spec            ClusterDataSpec `json:"spec" db:"-"`
	Status          string          `json:"status" db:"status"`
	Error           string          `json:"error" db:"error"`
	ObservedRoles   []string        `json:"observedRoles"`
	ObservedVolumes []string        `json:"observedVolumes"`
	Config          string          `json:"-" db:"encrypted_config"`
	UpdatedAt       string          `json:"updatedAt" db:"updated_at"`
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

type RedisSnapshot struct {
	ID              string `json:"id" db:"id"`
	DatabaseID      string `json:"databaseId" db:"database_id"`
	ProjectID       string `json:"projectId" db:"project_id"`
	ClusterID       string `json:"clusterId" db:"cluster_id"`
	S3DestinationID string `json:"s3DestinationId" db:"s3_destination_id"`
	S3Key           string `json:"s3Key" db:"s3_key"`
	SizeBytes       int64  `json:"sizeBytes" db:"size_bytes"`
	Status          string `json:"status" db:"status"`
	Error           string `json:"error" db:"error"`
	CreatedAt       string `json:"createdAt" db:"created_at"`
}
