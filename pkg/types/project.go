package types

import "time"

type ProjectConfig struct {
	ID                  string     `json:"id" db:"id"`
	AppID               string     `json:"appId,omitempty" db:"app_id"`
	OrganizationID      string     `json:"organizationId,omitempty" db:"organization_id"`
	ServerID            string     `json:"serverId,omitempty" db:"server_id"`
	Name                string     `json:"name" db:"name"`
	Slug                string     `json:"slug" db:"slug"`
	Description         string     `json:"description" db:"description"`
	EnvironmentName     string     `json:"environmentName" db:"environment_name"`
	EnvironmentSlug     string     `json:"environmentSlug" db:"environment_slug"`
	EnvironmentType     string     `json:"environmentType" db:"environment_type"`
	IsApp               bool       `json:"isApp" db:"is_app"`
	AppTemplateID       string     `json:"appTemplateId,omitempty" db:"app_template_id"`
	LocalPath           string     `json:"localPath,omitempty" db:"local_path"`
	GitProvider         string     `json:"gitProvider,omitempty" db:"git_provider"`
	GitOwner            string     `json:"gitOwner,omitempty" db:"git_owner"`
	GitRepo             string     `json:"gitRepo,omitempty" db:"git_repo"`
	GitBranch           string     `json:"gitBranch,omitempty" db:"git_branch"`
	GitURL              string     `json:"gitUrl,omitempty" db:"git_url"`
	CloneTokenEncrypted string     `json:"-" db:"clone_token_encrypted"`
	CloneTokenSetAt     *time.Time `json:"cloneTokenSetAt,omitempty" db:"clone_token_set_at"`
	RoutingConfigJSON   string     `json:"routingConfig,omitempty" db:"routing_config_json"`
	ReadinessJSON       string     `json:"readiness,omitempty" db:"readiness_json"`
	Status              string     `json:"status" db:"status"`
	CreatedAt           time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
}

type EnvironmentConfig struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"projectId" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	IsDefault bool      `json:"isDefault" db:"is_default"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

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
