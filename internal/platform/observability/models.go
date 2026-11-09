package observability

import "time"

type TrafficSample struct {
	Time       string `json:"time"`
	ProjectID  string `json:"projectId"`
	Domain     string `json:"domain"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	Bytes      int64  `json:"bytes"`
	DurationMs int64  `json:"durationMs"`
	ClientIP   string `json:"clientIp"`
	Country    string `json:"country"`
}

type LogDrainType string

const (
	LogDrainTypeAxiom    LogDrainType = "axiom"
	LogDrainTypeNewRelic LogDrainType = "new_relic"
	LogDrainTypeDatadog  LogDrainType = "datadog"
	LogDrainTypeWebhook  LogDrainType = "webhook"
)

type LogDrain struct {
	ID          string       `json:"id" db:"id"`
	ServiceID   string       `json:"serviceId" db:"service_id"`
	ProjectID   string       `json:"projectId" db:"project_id"`
	DrainType   LogDrainType `json:"drainType" db:"drain_type"`
	EndpointURL string       `json:"endpointUrl" db:"endpoint_url"`
	AuthToken   string       `json:"authToken" db:"auth_token"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time    `json:"updatedAt" db:"updated_at"`
}

type RuntimeObservation struct {
	Logs    string         `json:"logs"`
	Metrics map[string]any `json:"metrics"`
}

type RuntimePod struct {
	Name  string `json:"name"`
	Node  string `json:"node"`
	Phase string `json:"phase"`
}

type WorkloadObservation struct {
	Kind             string       `json:"kind"`
	Status           string       `json:"status"`
	Desired          int          `json:"desired"`
	Available        int          `json:"available"`
	MetricsAvailable bool         `json:"metricsAvailable"`
	MetricsError     string       `json:"metricsError,omitempty"`
	Pods             []RuntimePod `json:"pods,omitempty"`
}

type CanvasNode struct {
	ID   string         `json:"id"`
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

type CanvasEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type CanvasApp struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	ContainerID string `json:"containerId"`
}

type CanvasDatabase struct {
	ID          string `json:"id"`
	ContainerID string `json:"containerId"`
}

type CanvasEnvironment struct {
	ProjectID string `json:"projectId"`
}

type CanvasClusterDataSpec struct {
	Name string `json:"name"`
}

type CanvasClusterData struct {
	ID     string                `json:"id"`
	Status string                `json:"status"`
	Error  string                `json:"error"`
	Spec   CanvasClusterDataSpec `json:"spec"`
}

type EnvironmentCanvas struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Environment      *CanvasEnvironment  `json:"environment"`
	Apps             []*CanvasApp        `json:"apps"`
	Databases        []*CanvasDatabase   `json:"databases"`
	ClusterDatabases []CanvasClusterData `json:"clusterDatabases"`
	Nodes            []CanvasNode        `json:"nodes"`
	Edges            []CanvasEdge        `json:"edges"`
	Revision         string              `json:"revision"`
}
