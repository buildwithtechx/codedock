package models

type RepositoryInspectionRequest struct {
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch"`
	RootDirectory string `json:"rootDirectory"`
}

type ServiceRoleSpec struct {
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	Command         string `json:"command"`
	StopGracePeriod string `json:"stopGracePeriod,omitempty"`
	Singleton       bool   `json:"singleton,omitempty"`
}

type MonorepoCandidateApp struct {
	Name           string `json:"name"`
	RootDirectory  string `json:"rootDirectory"`
	Framework      string `json:"framework"`
	InternalPort   int    `json:"internalPort"`
	BuildEngine    string `json:"buildEngine"`
	PackageManager string `json:"packageManager,omitempty"`
	InstallCommand string `json:"installCommand,omitempty"`
	BuildCommand   string `json:"buildCommand,omitempty"`
	StartCommand   string `json:"startCommand,omitempty"`
	DockerfilePath string `json:"dockerfilePath,omitempty"`
}

type RepoTreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "tree" (dir) or "blob" (file)
}

type RepoTreeRequest struct {
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch,omitempty"`
}

type RepositoryInspection struct {
	Framework      string                 `json:"framework"`
	PackageManager string                 `json:"packageManager"`
	InstallCommand string                 `json:"installCommand"`
	BuildCommand   string                 `json:"buildCommand"`
	StartCommand   string                 `json:"startCommand"`
	InternalPort   int                    `json:"internalPort"`
	StaticOutput   string                 `json:"staticOutput"`
	BuildEngine    string                 `json:"buildEngine"`
	DockerfilePath string                 `json:"dockerfilePath"`
	Warnings       []string               `json:"warnings"`
	Stack          string                 `json:"stack,omitempty"`
	Roles          []ServiceRoleSpec      `json:"roles,omitempty"`
	CandidateApps  []MonorepoCandidateApp `json:"candidateApps,omitempty"`
	IsMonorepo     bool                   `json:"isMonorepo,omitempty"`
}
