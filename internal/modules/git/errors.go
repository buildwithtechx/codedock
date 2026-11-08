package git

import "errors"

var (
	ErrProviderNotFound    = errors.New("git provider not found")
	ErrUnsupportedProvider = errors.New("unsupported git provider")
	ErrAccessTokenRequired = errors.New("access token is required")
	ErrUserRequired        = errors.New("user ID is required")
	ErrRepoNotFound        = errors.New("repository not found")
	ErrInvalidRepoURL      = errors.New("use an HTTPS repository URL without credentials")
	ErrInvalidBranch       = errors.New("invalid branch name")
	ErrInvalidSignature    = errors.New("invalid webhook signature")
	ErrMissingWebhookEvent = errors.New("missing X-GitHub-Event header")
	ErrAppNotFound         = errors.New("github app not found")
	ErrPreviewNotEligible  = errors.New("pull request is not eligible for preview")
	ErrPreviewNotFound     = errors.New("pr preview not found")
	ErrTreeExceedsLimit    = errors.New("package.json exceeds inspection limits")
)
