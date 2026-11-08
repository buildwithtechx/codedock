package git

type GitConnectRequest struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"accessToken"`
	AccountName string `json:"accountName"`
}

type RepoTreeRequest struct {
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch"`
	Path          string `json:"path,omitempty"`
}

type RepositoryInspectionRequest struct {
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch,omitempty"`
	RootDirectory string `json:"rootDirectory,omitempty"`
}

type GitAppsManifestRequest struct {
	Code string `json:"code"`
}

type DeployPRPreviewOpts struct {
	AppID      string `json:"appId"`
	PRNumber   int    `json:"prNumber"`
	CommitHash string `json:"commitHash"`
	Branch     string `json:"branch"`
}

// Third-party API response DTOs
type GitHubRepoDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"clone_url"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
}

type GitLabRepoDTO struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	Visibility        string `json:"visibility"`
	HTTPURL           string `json:"http_url_to_repo"`
	WebURL            string `json:"web_url"`
	DefaultBranch     string `json:"default_branch"`
}

type BitbucketCloneLinkDTO struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type BitbucketHTMLLinkDTO struct {
	Href string `json:"href"`
}

type BitbucketLinksDTO struct {
	HTML  BitbucketHTMLLinkDTO    `json:"html"`
	Clone []BitbucketCloneLinkDTO `json:"clone"`
}

type BitbucketBranchDTO struct {
	Name string `json:"name"`
}

type BitbucketRepoDTO struct {
	UUID       string             `json:"uuid"`
	Name       string             `json:"name"`
	FullName   string             `json:"full_name"`
	IsPriv     bool               `json:"is_private"`
	Links      BitbucketLinksDTO  `json:"links"`
	Mainbranch BitbucketBranchDTO `json:"mainbranch"`
}

type BitbucketRepoListDTO struct {
	Values []BitbucketRepoDTO `json:"values"`
}

type GiteaRepoDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"clone_url"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
}
