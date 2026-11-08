package auth

import (
	"time"

	"codedock/internal/modules/users"
)

type AuthResponse struct {
	User         *users.User `json:"user"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken,omitempty"`
}

type TokenResponse struct {
	User         *users.User `json:"user"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken,omitempty"`
}

type TwoFASetupResponse struct {
	QRCodeURI     string   `json:"qrCodeUri"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

type PasskeyResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	DeviceType string     `json:"deviceType"`
	BackedUp   bool       `json:"backedUp"`
	Transports []string   `json:"transports"`
	AAGUID     string     `json:"aaguid,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
}
