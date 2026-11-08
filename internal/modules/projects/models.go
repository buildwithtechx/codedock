package projects

import (
	"time"

	"codedock/internal/modules/organizations"
	"codedock/internal/modules/users"
)

type UserRole = users.UserRole
type MemberPermission = organizations.MemberPermission
type MemberStatus = organizations.MemberStatus
type Organization = organizations.Organization
type OrganizationMember = organizations.OrganizationMember

const (
	UserRoleOwner  = users.UserRoleOwner
	UserRoleAdmin  = users.UserRoleAdmin
	UserRoleMember = users.UserRoleMember

	MemberPermissionOwner  = organizations.MemberPermissionOwner
	MemberPermissionAdmin  = organizations.MemberPermissionAdmin
	MemberPermissionMember = organizations.MemberPermissionMember

	MemberStatusAccepted = organizations.MemberStatusAccepted
)

type ProjectConfig struct {
	ID                      string     `json:"id" db:"id"`
	AppID                   string     `json:"appId,omitempty" db:"app_id"`
	OrganizationID          string     `json:"organizationId,omitempty" db:"organization_id"`
	ServerID                string     `json:"serverId,omitempty" db:"server_id"`
	Name                    string     `json:"name" db:"name"`
	Slug                    string     `json:"slug" db:"slug"`
	Description             string     `json:"description" db:"description"`
	EnvironmentName         string     `json:"environmentName" db:"environment_name"`
	EnvironmentSlug         string     `json:"environmentSlug" db:"environment_slug"`
	EnvironmentType         string     `json:"environmentType" db:"environment_type"`
	IsApp                   bool       `json:"isApp" db:"is_app"`
	AppTemplateID           string     `json:"appTemplateId,omitempty" db:"app_template_id"`
	LocalPath               string     `json:"localPath,omitempty" db:"local_path"`
	GitProvider             string     `json:"gitProvider,omitempty" db:"git_provider"`
	GitOwner                string     `json:"gitOwner,omitempty" db:"git_owner"`
	GitRepo                 string     `json:"gitRepo,omitempty" db:"git_repo"`
	GitBranch               string     `json:"gitBranch,omitempty" db:"git_branch"`
	GitURL                  string     `json:"gitUrl,omitempty" db:"git_url"`
	CloneTokenEncrypted     string     `json:"-" db:"clone_token_encrypted"`
	CloneTokenSetAt         *time.Time `json:"cloneTokenSetAt,omitempty" db:"clone_token_set_at"`
	RoutingConfigJSON       string     `json:"routingConfig,omitempty" db:"routing_config_json"`
	ReadinessJSON           string     `json:"readiness,omitempty" db:"readiness_json"`
	Status                  string     `json:"status" db:"status"`
	DefaultRollbackStrategy string     `json:"defaultRollbackStrategy,omitempty" db:"default_rollback_strategy"`
	RollbackWindow          *int       `json:"rollbackWindow,omitempty" db:"rollback_window"`
	SnapshotSizeBytes       *int64     `json:"snapshotSizeBytes,omitempty" db:"snapshot_size_bytes"`
	CapacityMeasuredAt      *time.Time `json:"capacityMeasuredAt,omitempty" db:"capacity_measured_at"`
	CreatedAt               time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt               time.Time  `json:"updatedAt" db:"updated_at"`
}

type EnvironmentConfig struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"projectId" db:"project_id"`
	Name      string    `json:"name" db:"name"`
	IsDefault bool      `json:"isDefault" db:"is_default"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type GitHubDeploymentChecks struct {
	Enabled        bool     `json:"enabled"`
	RequiredChecks []string `json:"requiredChecks,omitempty"`
}

type ProjectApp struct {
	ID               string                  `json:"id" db:"id"`
	OrganizationID   string                  `json:"organizationId" db:"organization_id"`
	Name             string                  `json:"name" db:"name"`
	Slug             string                  `json:"slug" db:"slug"`
	GitProvider      string                  `json:"gitProvider" db:"git_provider"`
	GitOwner         string                  `json:"gitOwner" db:"git_owner"`
	GitRepo          string                  `json:"gitRepo" db:"git_repo"`
	GitURL           string                  `json:"gitUrl" db:"git_url"`
	InstallationID   int                     `json:"installationId" db:"installation_id"`
	Favicon          string                  `json:"favicon" db:"favicon"`
	FaviconCheckedAt *time.Time              `json:"faviconCheckedAt,omitempty" db:"favicon_checked_at"`
	DeletedAt        *time.Time              `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt        time.Time               `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time               `json:"updatedAt" db:"updated_at"`
	GitHubChecks     *GitHubDeploymentChecks `json:"githubChecks,omitempty" db:"-"`
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

type ProjectMember struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"projectId" db:"project_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type BuildAdmissionStatus struct {
	CanAdmit             bool   `json:"canAdmit"`
	ActiveBuilds         int    `json:"activeBuilds"`
	MaxConcurrentBuilds  int    `json:"maxConcurrentBuilds"`
	QueuePosition        int    `json:"queuePosition"`
	EstimatedWaitSeconds int    `json:"estimatedWaitSeconds"`
	Reason               string `json:"reason"`
}

type ProjectChecksConfig struct {
	ProjectID            string   `json:"projectId"`
	Enabled              bool     `json:"enabled"`
	ReportCommitStatus   bool     `json:"reportCommitStatus"`
	RequirePassingChecks bool     `json:"requirePassingChecks"`
	TargetBranches       []string `json:"targetBranches"`
	ContextName          string   `json:"contextName"`
	UpdatedAt            string   `json:"updatedAt"`
}
