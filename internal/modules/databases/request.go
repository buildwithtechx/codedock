package databases

type CreateDatabaseRequest struct {
	ProjectID          string         `json:"projectId"`
	EnvironmentID      string         `json:"environmentId,omitempty"`
	Name               string         `json:"name"`
	Engine             DatabaseEngine `json:"engine"`
	Version            string         `json:"version"`
	DatabaseName       string         `json:"databaseName"`
	Username           string         `json:"username"`
	Password           string         `json:"password"`
	Port               int            `json:"port,omitempty"`
	VolumePath         string         `json:"volumePath,omitempty"`
	CustomArgs         string         `json:"customArgs,omitempty"`
	LogicalReplication bool           `json:"logicalReplication,omitempty"`
}

type UpdateDatabaseRequest struct {
	Name               string  `json:"name,omitempty"`
	Version            string  `json:"version,omitempty"`
	ExternalDNS        string  `json:"externalDns,omitempty"`
	CPULimit           float64 `json:"cpuLimit,omitempty"`
	MemoryLimit        int     `json:"memoryLimit,omitempty"`
	LogicalReplication bool    `json:"logicalReplication,omitempty"`
	CustomArgs         string  `json:"customArgs,omitempty"`
}

type DatabaseQueryRequest struct {
	Query string `json:"query"`
}

type ImportDatabaseRequest struct {
	SQL       string `json:"sql,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
}

type UpdateTableRowOpts struct {
	ID    string
	Table string
	Keys  map[string]any
	Data  map[string]any
}
