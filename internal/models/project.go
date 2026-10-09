package models

import (
	"time"
)

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

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectToken struct {
	ID            string     `json:"id" db:"id"`
	ProjectID     string     `json:"projectId" db:"project_id"`
	EnvironmentID string     `json:"environmentId,omitempty" db:"environment_id"`
	Name          string     `json:"name" db:"name"`
	TokenHash     string     `json:"-" db:"token_hash"`
	TokenPrefix   string     `json:"tokenPrefix,omitempty" db:"token_prefix"`
	Prefix        string     `json:"prefix" db:"prefix"`
	Role          string     `json:"role" db:"role"`
	Scopes        []string   `json:"scopes,omitempty" db:"scopes"`
	IPAllowlist   []string   `json:"ipAllowlist,omitempty" db:"ip_allowlist"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
}

type CreateTokenRequest struct {
	Name          string     `json:"name"`
	Role          string     `json:"role"`
	EnvironmentID string     `json:"environmentId,omitempty"`
	Scopes        []string   `json:"scopes,omitempty"`
	IPAllowlist   []string   `json:"ipAllowlist,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt"`
}

type CreateTokenResponse struct {
	Token        string        `json:"token"`
	ProjectToken *ProjectToken `json:"projectToken"`
}

type ProjectMember struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"projectId" db:"project_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type AddMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ServerlessFunctionCode struct {
	ID          string    `json:"id" db:"id"`
	ServiceID   string    `json:"serviceId" db:"service_id"`
	Runtime     string    `json:"runtime" db:"runtime"`
	CodeContent string    `json:"codeContent" db:"code_content"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

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
