package clusters

import "codedock/internal/modules/clusters/kubernetes"

type ClusterDataSpec = kubernetes.ClusterDataSpec
type ClusterData = kubernetes.ClusterData
type OperatorManifest = kubernetes.OperatorManifest
type ClusterDataPlan = kubernetes.ClusterDataPlan

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
