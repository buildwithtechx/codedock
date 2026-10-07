package models

type ClusterNode struct {
	ServerID    string `json:"serverId"`
	PrivateIP   string `json:"privateIp"`
	Interface   string `json:"interface"`
	Fingerprint string `json:"fingerprint"`
}
type Cluster struct {
	ID             string        `json:"id" db:"id"`
	ProjectID      string        `json:"projectId" db:"project_id"`
	OrganizationID string        `json:"organizationId" db:"organization_id"`
	Name           string        `json:"name" db:"name"`
	Version        string        `json:"version" db:"version"`
	Controls       int           `json:"controls" db:"controls"`
	Nodes          []ClusterNode `json:"nodes" db:"-"`
	NodesJSON      string        `json:"-" db:"nodes_json"`
	JoinToken      string        `json:"-" db:"encrypted_token"`
	Revision       int           `json:"revision" db:"revision"`
	Status         string        `json:"status" db:"status"`
	Error          string        `json:"error" db:"error"`
	UpdatedAt      string        `json:"updatedAt" db:"updated_at"`
}
type ClusterReviewRequest struct {
	Cluster Cluster `json:"cluster"`
	Action  string  `json:"action"`
}
type ClusterPlan struct {
	OperationID     string  `json:"operationId"`
	PreviousVersion string  `json:"previousVersion"`
	ExistingNodes   int     `json:"existingNodes"`
	Cluster         Cluster `json:"cluster"`
	Action          string  `json:"action"`
	Installer       string  `json:"installer"`
	InstallerSHA256 string  `json:"installerSha256"`
}
