package models

type RepositoryInspectionRequest struct {
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch"`
	RootDirectory string `json:"rootDirectory"`
}

type RepositoryInspection struct {
	Framework      string   `json:"framework"`
	PackageManager string   `json:"packageManager"`
	InstallCommand string   `json:"installCommand"`
	BuildCommand   string   `json:"buildCommand"`
	StartCommand   string   `json:"startCommand"`
	InternalPort   int      `json:"internalPort"`
	StaticOutput   string   `json:"staticOutput"`
	BuildEngine    string   `json:"buildEngine"`
	DockerfilePath string   `json:"dockerfilePath"`
	Warnings       []string `json:"warnings"`
}
