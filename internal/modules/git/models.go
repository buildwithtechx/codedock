package git

import (
	"encoding/json"
	"time"
)

type GitProviderConfig struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userId" db:"user_id"`
	Provider    string    `json:"provider" db:"provider"`
	AccessToken string    `json:"accessToken,omitempty" db:"encrypted_access_token"`
	AccountName string    `json:"accountName" db:"account_name"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type GitRepository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"fullName"`
	Private       bool   `json:"private"`
	CloneURL      string `json:"cloneUrl"`
	HTMLURL       string `json:"htmlUrl"`
	DefaultBranch string `json:"defaultBranch"`
}

type TreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" or "tree"
}

type GithubApp struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	AppID          string    `json:"appId" db:"app_id"`
	InstallationID string    `json:"installationId" db:"installation_id"`
	ClientID       string    `json:"clientId" db:"client_id"`
	ClientSecret   string    `json:"clientSecret" db:"client_secret"`
	WebhookSecret  string    `json:"webhookSecret" db:"webhook_secret"`
	PrivateKey     string    `json:"privateKey" db:"private_key"`
	IsPublic       bool      `json:"isPublic" db:"is_public"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type WebhookDelivery struct {
	ID             string          `json:"id" db:"id"`
	OrganizationID *string         `json:"organizationId,omitempty" db:"organization_id"`
	ProjectID      *string         `json:"projectId,omitempty" db:"project_id"`
	Source         string          `json:"source" db:"source"`
	HookID         *string         `json:"hookId,omitempty" db:"hook_id"`
	DeliveryID     *string         `json:"deliveryId,omitempty" db:"delivery_id"`
	Event          string          `json:"event" db:"event"`
	AuthResult     *string         `json:"authResult,omitempty" db:"auth_result"`
	Outcome        string          `json:"outcome" db:"outcome"`
	ActionRef      *string         `json:"actionRef,omitempty" db:"action_ref"`
	StatusCode     *int            `json:"statusCode,omitempty" db:"status_code"`
	Error          *string         `json:"error,omitempty" db:"error"`
	Summary        json.RawMessage `json:"summary,omitempty" db:"summary"`
	ClientIP       *string         `json:"clientIp,omitempty" db:"client_ip"`
	UserAgent      *string         `json:"userAgent,omitempty" db:"user_agent"`
	DurationMs     *int            `json:"durationMs,omitempty" db:"duration_ms"`
	ReceivedAt     time.Time       `json:"receivedAt" db:"received_at"`
	ProcessedAt    *time.Time      `json:"processedAt,omitempty" db:"processed_at"`
}

type WebhookDeliveryInput struct {
	OrganizationID *string
	ProjectID      *string
	Source         string
	HookID         *string
	DeliveryID     *string
	Event          string
	AuthResult     *string
	Outcome        string
	ActionRef      *string
	StatusCode     *int
	Error          *string
	Summary        any
	ClientIP       *string
	UserAgent      *string
	DurationMs     *int
}

type WebhookDeliveryResult struct {
	Outcome    string
	ActionRef  *string
	StatusCode *int
	Error      *string
	Summary    any
	DurationMs *int
}

type ClaimResult struct {
	Claimed           bool
	ID                string
	HandledProjectIDs []string
}

type GithubWebhookPayload struct {
	Action      string                `json:"action"`
	Number      int                   `json:"number"`
	PullRequest GithubPullRequest     `json:"pull_request"`
	Ref         string                `json:"ref"`
	After       string                `json:"after"`
	Deleted     bool                  `json:"deleted"`
	Repository  GithubRepoInfo        `json:"repository"`
	HeadCommit  *GithubCommitInfo     `json:"head_commit"`
	Commits     []GithubWebhookCommit `json:"commits"`
}

type GithubRepoInfo struct {
	Name          string          `json:"name"`
	FullName      string          `json:"full_name"`
	DefaultBranch string          `json:"default_branch"`
	Owner         GithubRepoOwner `json:"owner"`
}

type GithubRepoOwner struct {
	Login string `json:"login"`
}

type GithubCommitInfo struct {
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

type GithubWebhookCommit struct {
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

type GithubPullRequest struct {
	State string                `json:"state"`
	Draft bool                  `json:"draft"`
	Head  GithubPullRequestHead `json:"head"`
	Base  GithubPullRequestBase `json:"base"`
}

type GithubPullRequestHead struct {
	Ref  string                `json:"ref"`
	Sha  string                `json:"sha"`
	Repo GithubPullRequestRepo `json:"repo"`
}

type GithubPullRequestBase struct {
	Ref  string                `json:"ref"`
	Repo GithubPullRequestRepo `json:"repo"`
}

type GithubPullRequestRepo struct {
	FullName string `json:"full_name"`
}

type PRPreview struct {
	ID            string    `json:"id" db:"id"`
	ServiceID     string    `json:"serviceId" db:"service_id"`
	ProjectID     string    `json:"projectId" db:"project_id"`
	PRNumber      int       `json:"prNumber" db:"pr_number"`
	Branch        string    `json:"branch" db:"branch"`
	CommitHash    string    `json:"commitHash" db:"commit_hash"`
	Status        string    `json:"status" db:"status"`
	PreviewDomain string    `json:"previewDomain" db:"preview_domain"`
	ContainerID   string    `json:"containerId" db:"container_id"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type PullRequestInfo struct {
	Number   int    `json:"number"`
	State    string `json:"state"`
	Draft    bool   `json:"draft"`
	HeadSHA  string `json:"headSha"`
	HeadRef  string `json:"headRef"`
	HeadRepo string `json:"headRepo"`
	BaseRef  string `json:"baseRef"`
	BaseRepo string `json:"baseRepo"`
}

type ServiceRoleSpec struct {
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	Command         string `json:"command"`
	StopGracePeriod string `json:"stopGracePeriod"`
	Singleton       bool   `json:"singleton"`
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

type InspectionResult struct {
	Framework      string                 `json:"framework"`
	InternalPort   int                    `json:"internalPort"`
	BuildEngine    string                 `json:"buildEngine"`
	DockerfilePath string                 `json:"dockerfilePath"`
	PackageManager string                 `json:"packageManager,omitempty"`
	InstallCommand string                 `json:"installCommand,omitempty"`
	BuildCommand   string                 `json:"buildCommand,omitempty"`
	StartCommand   string                 `json:"startCommand,omitempty"`
	StaticOutput   string                 `json:"staticOutput,omitempty"`
	Stack          string                 `json:"stack,omitempty"`
	Warnings       []string               `json:"warnings"`
	Roles          []ServiceRoleSpec      `json:"roles,omitempty"`
	IsMonorepo     bool                   `json:"isMonorepo,omitempty"`
	CandidateApps  []MonorepoCandidateApp `json:"candidateApps,omitempty"`
}
