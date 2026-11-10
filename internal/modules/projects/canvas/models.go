package canvas

import (
	"time"

	"codedock/internal/modules/environments"
)

type EnvironmentConfig = environments.EnvironmentConfig

type CanvasSummary struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	EnvironmentsCount  int                `json:"environmentsCount"`
	AppsCount          int                `json:"appsCount"`
	DatabasesCount     int                `json:"databasesCount"`
	TotalServices      int                `json:"totalServices"`
	OnlineServices     int                `json:"onlineServices"`
	DefaultEnvironment *EnvironmentConfig `json:"defaultEnvironment,omitempty"`
	ServiceIcons       []string           `json:"serviceIcons"`
	Stack              string             `json:"stack,omitempty"`
	DeployTarget       string             `json:"deployTarget"`
	ServerName         string             `json:"serverName,omitempty"`
	Nodes              []string           `json:"nodes,omitempty"`
	NodeCount          int                `json:"nodeCount,omitempty"`
	EdgeCount          int                `json:"edgeCount,omitempty"`
}

type CanvasPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type CanvasNode struct {
	ID   string         `json:"id"`
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
	Pos  CanvasPosition `json:"position"`
}

type CanvasEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
	Label  string `json:"label"`
}

const (
	AppServiceStatusRunning = "running"
	DatabaseStatusRunning   = "running"
)

type CanvasAppService struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	EnvironmentID string `json:"environmentId"`
}

type AppService struct {
	ID              string    `json:"id" db:"id"`
	ProjectID       string    `json:"projectId" db:"project_id"`
	EnvironmentID   string    `json:"environmentId" db:"environment_id"`
	Name            string    `json:"name" db:"name"`
	Stack           string    `json:"stack,omitempty" db:"stack"`
	RepositoryURL   string    `json:"repositoryUrl" db:"repository_url"`
	ImageRef        string    `json:"imageRef" db:"image_ref"`
	Branch          string    `json:"branch" db:"branch"`
	InternalPort    int       `json:"internalPort" db:"internal_port"`
	Domain          string    `json:"domain" db:"domain"`
	ContainerID     string    `json:"containerId" db:"container_id"`
	Status          string    `json:"status" db:"status"`
	Replicas        int       `json:"replicas" db:"replicas"`
	MaintenanceMode bool      `json:"maintenanceMode" db:"maintenance_mode"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type Database struct {
	ID            string    `json:"id" db:"id"`
	ProjectID     string    `json:"projectId" db:"project_id"`
	EnvironmentID string    `json:"environmentId" db:"environment_id"`
	Name          string    `json:"name" db:"name"`
	Engine        string    `json:"engine" db:"engine"`
	Version       string    `json:"version" db:"version"`
	Port          int       `json:"port" db:"port"`
	Username      string    `json:"username" db:"username"`
	DatabaseName  string    `json:"databaseName" db:"database_name"`
	VolumePath    string    `json:"volumePath" db:"volume_path"`
	ContainerID   string    `json:"containerId" db:"container_id"`
	Status        string    `json:"status" db:"status"`
	InternalDNS   string    `json:"internalDns" db:"internal_dns"`
	ExternalDNS   string    `json:"externalDns" db:"external_dns"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type CanvasDatabase struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	EnvironmentID string `json:"environmentId"`
}

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

type CanvasClusterData struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type EnvironmentCanvas struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Environment      *EnvironmentConfig `json:"environment"`
	Apps             []*AppService      `json:"apps"`
	Databases        []*Database        `json:"databases"`
	ClusterDatabases []ClusterData      `json:"clusterDatabases"`
	Nodes            []CanvasNode       `json:"nodes"`
	Edges            []CanvasEdge       `json:"edges"`
	Revision         string             `json:"revision"`
}

type TopologyApplyRequest struct {
	Revision     string       `json:"revision"`
	Dependencies []CanvasEdge `json:"dependencies"`
}

type RuntimeObservation struct {
	Logs    string         `json:"logs"`
	Metrics map[string]any `json:"metrics"`
}

type ComposeStack struct {
	ID            string `json:"id" db:"id"`
	ProjectID     string `json:"projectId" db:"project_id"`
	EnvironmentID string `json:"environmentId" db:"environment_id"`
	Name          string `json:"name" db:"name"`
	Revision      int    `json:"revision" db:"revision"`
	Status        string `json:"status" db:"status"`
	Error         string `json:"error" db:"error"`
	Results       string `json:"results" db:"results"`
	UpdatedAt     string `json:"updatedAt" db:"updated_at"`
	Config        string `json:"-" db:"encrypted_config"`
}

type ComposeServiceResult struct {
	Name        string `json:"name"`
	ContainerID string `json:"containerId,omitempty"`
	State       string `json:"state,omitempty"`
}
