package git

const (
	// Git providers
	ProviderGitHub    = "github"
	ProviderGitLab    = "gitlab"
	ProviderBitbucket = "bitbucket"
	ProviderGitea     = "gitea"

	// Webhook sources
	WebhookSourceGitHub      = "github"
	WebhookSourceDeployToken = "deploy_token"

	// Webhook events
	WebhookEventPush        = "push"
	WebhookEventPullRequest = "pull_request"
	WebhookEventPing        = "ping"
	WebhookEventDeploy      = "deploy"

	// Webhook outcomes
	WebhookOutcomeSuccess    = "success"
	WebhookOutcomeSkipped    = "skipped"
	WebhookOutcomeFailed     = "failed"
	WebhookOutcomeReceived   = "received"
	WebhookOutcomeDispatched = "dispatched"

	// PR Preview status
	PRPreviewStatusPending  = "PENDING"
	PRPreviewStatusBuilding = "BUILDING"
	PRPreviewStatusReady    = "READY"
	PRPreviewStatusFailed   = "FAILED"

	// Limits
	MaxSubmoduleRepos = 25
)
