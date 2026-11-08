package projects

import "time"

type CreateProjectRequest struct {
	ID              string `json:"id,omitempty"`
	AppID           string `json:"appId,omitempty"`
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	EnvironmentName string `json:"environmentName,omitempty"`
	EnvironmentType string `json:"environmentType,omitempty"`
	GitProvider     string `json:"gitProvider,omitempty"`
	GitOwner        string `json:"gitOwner,omitempty"`
	GitRepo         string `json:"gitRepo,omitempty"`
	GitBranch       string `json:"gitBranch,omitempty"`
	GitURL          string `json:"gitUrl,omitempty"`
	OrganizationID  string `json:"organizationId,omitempty"`
	ServerID        string `json:"serverId,omitempty"`
	IsApp           bool   `json:"isApp,omitempty"`
	AppTemplateID   string `json:"appTemplateId,omitempty"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectAppRequest struct {
	Name           string                  `json:"name"`
	Slug           string                  `json:"slug,omitempty"`
	GitProvider    string                  `json:"gitProvider,omitempty"`
	GitOwner       string                  `json:"gitOwner,omitempty"`
	GitRepo        string                  `json:"gitRepo,omitempty"`
	GitURL         string                  `json:"gitUrl,omitempty"`
	InstallationID int                     `json:"installationId,omitempty"`
	GitHubChecks   *GitHubDeploymentChecks `json:"githubChecks,omitempty"`
}

type UpdateProjectAppRequest struct {
	Name         string                  `json:"name,omitempty"`
	Slug         string                  `json:"slug,omitempty"`
	Favicon      string                  `json:"favicon,omitempty"`
	GitProvider  string                  `json:"gitProvider,omitempty"`
	GitOwner     string                  `json:"gitOwner,omitempty"`
	GitRepo      string                  `json:"gitRepo,omitempty"`
	GitURL       string                  `json:"gitUrl,omitempty"`
	GitHubChecks *GitHubDeploymentChecks `json:"githubChecks,omitempty"`
}

type CreateTokenRequest struct {
	Name          string     `json:"name"`
	Role          string     `json:"role"`
	EnvironmentID string     `json:"environmentId,omitempty"`
	Scopes        []string   `json:"scopes,omitempty"`
	IPAllowlist   []string   `json:"ipAllowlist,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt"`
}

type AddMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UpdateChecksConfigRequest struct {
	Enabled              bool     `json:"enabled"`
	ReportCommitStatus   bool     `json:"reportCommitStatus"`
	RequirePassingChecks bool     `json:"requirePassingChecks"`
	TargetBranches       []string `json:"targetBranches"`
	ContextName          string   `json:"contextName"`
}
