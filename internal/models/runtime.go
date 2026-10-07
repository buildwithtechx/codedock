package models

type RuntimeVolume struct {
	Name         string `json:"name"`
	MountPath    string `json:"mountPath"`
	SizeGiB      int    `json:"sizeGiB"`
	StorageClass string `json:"storageClass"`
	Shared       bool   `json:"shared"`
}
type RuntimeTarget struct {
	BareReleaseURL  string          `json:"bareReleaseUrl,omitempty"`
	BareSHA256      string          `json:"bareSha256,omitempty"`
	BareCommand     []string        `json:"bareCommand"`
	Kind            string          `json:"kind"`
	ClusterID       string          `json:"clusterId,omitempty"`
	NodeIDs         []string        `json:"nodeIds"`
	ImageRepository string          `json:"imageRepository,omitempty"`
	RegistryID      string          `json:"registryId,omitempty"`
	Volumes         []RuntimeVolume `json:"volumes"`
	BareNode        ClusterNode     `json:"bareNode,omitempty"`
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
	Registry  *Registry         `json:"registry,omitempty"`
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
