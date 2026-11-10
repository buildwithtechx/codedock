package settings

import "context"

// Vault encrypts and decrypts sensitive server-wide configuration parameters.
type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// SettingsRepository defines database persistence for global instance settings and defaults.
type SettingsRepository interface {
	GetServerSettings(ctx context.Context) (*ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *ServerSettings) error
	ListProjects(ctx context.Context) ([]map[string]any, error)
	IsNotificationEmailEnabled(ctx context.Context) (bool, error)
}

type Repository = SettingsRepository
