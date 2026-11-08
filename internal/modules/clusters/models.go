package clusters

type Vault interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) ([]byte, error)
}

type ProjectConfig struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	ServerID       string `json:"serverId"`
}

type Server struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	IP             string `json:"ip"`
}

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

type AppService struct {
	ID            string `json:"id"`
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
	Name          string `json:"name"`
}

type Variable struct {
	ID                       string `json:"id"`
	ServiceID                string `json:"serviceId"`
	EnvironmentID            string `json:"environmentId"`
	Key                      string `json:"key"`
	Value                    string `json:"value"`
	IsSecret                 bool   `json:"isSecret"`
	ExpectedTopologyRevision string `json:"expectedTopologyRevision"`
	CreatedAt                string `json:"createdAt"`
	UpdatedAt                string `json:"updatedAt"`
}

type EnvironmentConfig struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
}

type S3Destination struct {
	ID        string `json:"id"`
	Bucket    string `json:"bucket"`
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type ClusterNode struct {
	ServerID    string `json:"serverId"`
	PrivateIP   string `json:"privateIp"`
	Interface   string `json:"interface"`
	Fingerprint string `json:"fingerprint"`
}

type ClusterScaling struct {
	Status        string `json:"status"`
	K3sVersion    string `json:"k3sVersion"`
	ControlNodes  int    `json:"controlNodes"`
	WorkerNodes   int    `json:"workerNodes"`
	PodCIDR       string `json:"podCidr,omitempty"`
	ServiceCIDR   string `json:"serviceCidr,omitempty"`
	ClusterDomain string `json:"clusterDomain,omitempty"`
	Error         string `json:"error,omitempty"`
	UpdatedAt     string `json:"updatedAt"`
}

type Cluster struct {
	ID             string          `json:"id" db:"id"`
	ProjectID      string          `json:"projectId" db:"project_id"`
	OrganizationID string          `json:"organizationId" db:"organization_id"`
	Name           string          `json:"name" db:"name"`
	Version        string          `json:"version" db:"version"`
	Controls       int             `json:"controls" db:"controls"`
	Nodes          []ClusterNode   `json:"nodes" db:"-"`
	NodesJSON      string          `json:"-" db:"nodes_json"`
	JoinToken      string          `json:"-" db:"encrypted_token"`
	Revision       int             `json:"revision" db:"revision"`
	Status         string          `json:"status" db:"status"`
	Scaling        *ClusterScaling `json:"scaling,omitempty" db:"-"`
	Error          string          `json:"error" db:"error"`
	UpdatedAt      string          `json:"updatedAt" db:"updated_at"`
}

type StorageSetupPlan struct {
	ClusterID string `json:"clusterId"`
	ProjectID string `json:"projectId"`
	Version   string `json:"version"`
	Manifest  string `json:"manifest"`
	SHA256    string `json:"sha256"`
	Action    string `json:"action"`
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
