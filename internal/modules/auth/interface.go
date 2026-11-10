package auth

import (
	"context"
	"time"
)

// OAuthRepository defines persistence for third-party OAuth providers and TOTP secrets.
type OAuthRepository interface {
	ListProviders(ctx context.Context) ([]OAuthProviderConfig, error)
	GetProvider(ctx context.Context, idOrName string) (*OAuthProviderConfig, error)
	SaveProvider(ctx context.Context, p *OAuthProviderConfig) error
	GetUserTOTPSecret(ctx context.Context, userID string) (secret string, recoveryCodes []string, err error)
	UpdateUserTOTP(ctx context.Context, userID string, enabled bool, secret string, recoveryCodes []string) error
}

// PasskeyRepository defines persistence for WebAuthn FIDO2 passkeys and authentication challenges.
type PasskeyRepository interface {
	CreateCredential(ctx context.Context, cred *PasskeyCredential) error
	GetCredentialByID(ctx context.Context, id string) (*PasskeyCredential, error)
	GetCredentialByCredentialID(ctx context.Context, credentialID string) (*PasskeyCredential, error)
	ListCredentialsByUserID(ctx context.Context, userID string) ([]*PasskeyCredential, error)
	UpdateCredentialCounter(ctx context.Context, id string, counter int64, lastUsedAt time.Time) error
	RenameCredential(ctx context.Context, id, userID, name string) error
	DeleteCredential(ctx context.Context, id, userID string) error
	SaveChallenge(ctx context.Context, challenge *PasskeyChallenge) error
	GetAndConsumeChallenge(ctx context.Context, challenge, ceremonyType string) (*PasskeyChallenge, error)
}

// RefreshTokenRepository defines persistence and revocation tracking for session refresh tokens.
type RefreshTokenRepository interface {
	StoreToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, tokenHash string) (bool, error)
	RevokeToken(ctx context.Context, tokenHash string) error
	RevokeAllForUser(ctx context.Context, userID string) error
	PruneExpired(ctx context.Context) error
}
