package git

import (
	"context"
	"time"
)

// Encryptor encrypts and decrypts Git provider tokens and private keys.
type Encryptor interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}

// GitAppRepository defines persistence for GitHub App credentials and installations.
type GitAppRepository interface {
	ListGithubApps(ctx context.Context) ([]GithubApp, error)
	GetGithubApp(ctx context.Context, id string) (*GithubApp, error)
	SaveGithubApp(ctx context.Context, app *GithubApp) error
	DeleteGithubApp(ctx context.Context, id string) error
}

// PRPreviewRepository manages ephemeral pull request preview deployments.
type PRPreviewRepository interface {
	Create(ctx context.Context, p *PRPreview) error
	GetByAppAndPR(ctx context.Context, appID string, prNumber int) ([]*PRPreview, error)
	GetByApp(ctx context.Context, appID string) ([]*PRPreview, error)
	Update(ctx context.Context, p *PRPreview) error
	Delete(ctx context.Context, id string) error
}

// GitProviderRepository manages user connections to GitHub, GitLab, and Bitbucket.
type GitProviderRepository interface {
	SaveProvider(ctx context.Context, gp *GitProviderConfig) error
	GetProvider(ctx context.Context, userID, provider string) (*GitProviderConfig, error)
	GetAnyProviderByType(ctx context.Context, provider string) (*GitProviderConfig, error)
	ListProvidersByUser(ctx context.Context, userID string) ([]*GitProviderConfig, error)
	DeleteProvider(ctx context.Context, userID, provider string) error
}

// WebhookDeliveryRepository manages inbound Git webhook delivery audit and replay prevention.
type WebhookDeliveryRepository interface {
	ClaimGithub(ctx context.Context, input WebhookDeliveryInput) (ClaimResult, error)
	Record(ctx context.Context, input WebhookDeliveryInput) (string, error)
	MarkProcessed(ctx context.Context, id string, res WebhookDeliveryResult) error
	ListByProject(ctx context.Context, projectID string, limit int, cursor string) ([]WebhookDelivery, string, error)
	Prune(ctx context.Context, olderThan time.Duration) (int64, error)
}
