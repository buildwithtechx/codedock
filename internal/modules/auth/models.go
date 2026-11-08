package auth

import (
	"time"
)

type OAuthProviderConfig struct {
	ID           string    `json:"id" db:"id"`
	ProviderName string    `json:"providerName" db:"provider_name"`
	Enabled      bool      `json:"enabled" db:"enabled"`
	ClientID     string    `json:"clientId" db:"client_id"`
	ClientSecret string    `json:"clientSecret,omitempty" db:"client_secret"`
	RedirectURI  string    `json:"redirectUri" db:"redirect_uri"`
	BaseURL      string    `json:"baseUrl,omitempty" db:"base_url"`
	Tenant       string    `json:"tenant,omitempty" db:"tenant"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type PasskeyCredential struct {
	ID           string     `json:"id" db:"id"`
	UserID       string     `json:"userId" db:"user_id"`
	Name         string     `json:"name" db:"name"`
	CredentialID string     `json:"credentialId" db:"credential_id"`
	PublicKey    string     `json:"publicKey" db:"public_key"`
	Counter      int64      `json:"counter" db:"counter"`
	DeviceType   string     `json:"deviceType" db:"device_type"`
	BackedUp     bool       `json:"backedUp" db:"backed_up"`
	Transports   string     `json:"transports" db:"transports"`
	AAGUID       string     `json:"aaguid" db:"aaguid"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty" db:"last_used_at"`
}

type PasskeyChallenge struct {
	Challenge    string    `json:"challenge" db:"challenge"`
	UserID       string    `json:"userId" db:"user_id"`
	CeremonyType string    `json:"ceremonyType" db:"ceremony_type"`
	SessionData  string    `json:"sessionData" db:"session_data"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
}

type DeviceSession struct {
	DeviceCode string              `json:"deviceCode"`
	UserCode   string              `json:"userCode"`
	ClientName string              `json:"clientName"`
	Status     DeviceSessionStatus `json:"status"`
	UserID     string              `json:"userId,omitempty"`
	Token      string              `json:"token,omitempty"`
	UserEmail  string              `json:"userEmail,omitempty"`
	UserName   string              `json:"userName,omitempty"`
	ExpiresAt  time.Time           `json:"expiresAt"`
	CreatedAt  time.Time           `json:"createdAt"`
}
