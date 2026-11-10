package apps

import (
	"time"
)

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
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type OperationReview struct {
	Operation    *Operation `json:"operation"`
	Confirmation string     `json:"confirmation"`
}

type OperationApply struct {
	Confirmation string `json:"confirmation"`
}

type ClusterNode = RuntimeClusterNode

type ClusterDataSpec struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	EnvironmentID string `json:"environmentId"`
	Engine        string `json:"engine"`
	Instances     int    `json:"instances"`
	Shards        int    `json:"shards"`
}

type ClusterData struct {
	ID        string          `json:"id"`
	ProjectID string          `json:"projectId"`
	Status    string          `json:"status"`
	Spec      ClusterDataSpec `json:"spec"`
}

type ClusterDataPlan struct {
	Record   ClusterData `json:"record"`
	Password string      `json:"password"`
}

type RuntimeCluster struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type RuntimeServer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type ProjectConfig struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ServerID string `json:"serverId"`
}

type RuntimeDatabase struct {
	ID            string `json:"id" db:"id"`
	ProjectID     string `json:"projectId" db:"project_id"`
	EnvironmentID string `json:"environmentId" db:"environment_id"`
	Name          string `json:"name" db:"name"`
	Engine        string `json:"engine" db:"engine"`
	Username      string `json:"username" db:"username"`
	Password      string `json:"password" db:"password"`
	DatabaseName  string `json:"databaseName" db:"database_name"`
	InternalDNS   string `json:"internalDns" db:"internal_dns"`
	Port          int    `json:"port" db:"port"`
}

type Database = RuntimeDatabase

type RuntimeVolume struct {
	Name         string `json:"name"`
	MountPath    string `json:"mountPath"`
	SizeGiB      int    `json:"sizeGiB"`
	StorageClass string `json:"storageClass"`
	Shared       bool   `json:"shared"`
}

type RuntimeClusterNode struct {
	ServerID    string `json:"serverId"`
	PrivateIP   string `json:"privateIp"`
	Interface   string `json:"interface"`
	Fingerprint string `json:"fingerprint"`
}

type RuntimeRegistry struct {
	ID            string `json:"id"`
	RegistryURL   string `json:"registryUrl,omitempty"`
	ServerURL     string `json:"serverUrl,omitempty"`
	Username      string `json:"username"`
	Password      string `json:"password,omitempty"`
	PasswordToken string `json:"passwordToken,omitempty"`
}

type RuntimeTarget struct {
	BareReleaseURL     string             `json:"bareReleaseUrl,omitempty"`
	BareSHA256         string             `json:"bareSha256,omitempty"`
	BareCommand        []string           `json:"bareCommand"`
	BareRepoURL        string             `json:"bareRepoUrl,omitempty"`
	BareBranch         string             `json:"bareBranch,omitempty"`
	BareToolchain      string             `json:"bareToolchain,omitempty"`
	BareInstallCommand string             `json:"bareInstallCommand,omitempty"`
	BareBuildCommand   string             `json:"bareBuildCommand,omitempty"`
	BareOutput         string             `json:"bareOutput,omitempty"`
	BareStaticDir      string             `json:"bareStaticDir,omitempty"`
	Kind               string             `json:"kind"`
	ClusterID          string             `json:"clusterId,omitempty"`
	NodeIDs            []string           `json:"nodeIds"`
	ImageRepository    string             `json:"imageRepository,omitempty"`
	RegistryID         string             `json:"registryId,omitempty"`
	Volumes            []RuntimeVolume    `json:"volumes"`
	BareNode           RuntimeClusterNode `json:"bareNode,omitempty"`
}

type ServiceRuntime struct {
	RuntimeKind string        `json:"-" db:"runtime_kind"`
	ServiceID   string        `json:"serviceId" db:"service_id"`
	ProjectID   string        `json:"projectId" db:"project_id"`
	Target      RuntimeTarget `json:"target" db:"-"`
	Revision    int           `json:"revision" db:"revision"`
	Status      string        `json:"status" db:"status"`
	Error       string        `json:"error" db:"error"`
	Journal     string        `json:"-" db:"encrypted_journal"`
	Config      string        `json:"-" db:"encrypted_config"`
	UpdatedAt   string        `json:"updatedAt" db:"updated_at"`
}

type RuntimeReviewRequest struct {
	Target   RuntimeTarget `json:"target"`
	Revision int           `json:"revision"`
}

type RuntimeOperationApply struct {
	OperationID  string `json:"operationId"`
	Confirmation string `json:"confirmation"`
}

type RuntimeReviewPlan struct {
	ClusterRevision int           `json:"clusterRevision"`
	ServiceID       string        `json:"serviceId"`
	Target          RuntimeTarget `json:"target"`
	Revision        int           `json:"revision"`
}

type RuntimeExecRequest struct {
	Command []string `json:"command"`
	Pod     string   `json:"pod,omitempty"`
}

type RuntimePod struct {
	Name        string  `json:"name"`
	Node        string  `json:"node"`
	Phase       string  `json:"phase"`
	Ready       bool    `json:"ready"`
	Restarts    int     `json:"restarts"`
	CPU         float64 `json:"cpu"`
	MemoryBytes int64   `json:"memoryBytes"`
}

type WorkloadObservation struct {
	Kind             string       `json:"kind"`
	Status           string       `json:"status"`
	Desired          int          `json:"desired"`
	Available        int          `json:"available"`
	MetricsAvailable bool         `json:"metricsAvailable"`
	MetricsError     string       `json:"metricsError,omitempty"`
	Pods             []RuntimePod `json:"pods"`
	Error            string       `json:"error,omitempty"`
}

type KubernetesWorkload struct {
	App       AppService        `json:"app"`
	Target    RuntimeTarget     `json:"target"`
	Image     string            `json:"image"`
	Variables map[string]string `json:"variables"`
	Registry  *RuntimeRegistry  `json:"registry,omitempty"`
}

type DesiredRuntime struct {
	Revision int                `json:"revision"`
	Workload KubernetesWorkload `json:"workload"`
	Manifest string             `json:"manifest"`
}

type NativeJournal struct {
	ReleaseID   string     `json:"releaseId"`
	PreviousApp AppService `json:"previousApp"`
}

type RuntimeImageStatus string

const (
	RuntimeImageStatusCurrent         RuntimeImageStatus = "current"
	RuntimeImageStatusUpdateAvailable RuntimeImageStatus = "update_available"
	RuntimeImageStatusPinned          RuntimeImageStatus = "pinned"
	RuntimeImageStatusUnknown         RuntimeImageStatus = "unknown"
)

type RuntimeImageInfo struct {
	ServiceID     *string            `json:"serviceId"`
	ServiceName   *string            `json:"serviceName"`
	Role          string             `json:"role"`
	Ref           string             `json:"ref"`
	CurrentDigest *string            `json:"currentDigest"`
	LatestDigest  *string            `json:"latestDigest"`
	Status        RuntimeImageStatus `json:"status"`
}

type RuntimeImageCheckResult struct {
	CheckedAt time.Time          `json:"checkedAt"`
	Images    []RuntimeImageInfo `json:"images"`
}

type UpdateRuntimeRequest struct {
	ServiceIDs []string `json:"serviceIds,omitempty"`
}
