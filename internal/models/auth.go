package models

import (
	"time"

	"codedock.run/codedock/pkg/types"
)

type AuthResult = types.AuthResult
type SignupRequest = types.SignupRequest
type SigninRequest = types.SigninRequest

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

type TwoFASetupResponse struct {
	QRCodeURI     string   `json:"qrCodeUri"`
	RecoveryCodes []string `json:"recoveryCodes"`
}
