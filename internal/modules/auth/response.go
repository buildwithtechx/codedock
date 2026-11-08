package auth

import "time"

type AuthResponse struct {
	User         *User  `json:"user"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type TokenResponse struct {
	User         *User  `json:"user"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type TwoFASetupResponse struct {
	QRCodeURI     string   `json:"qrCodeUri"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

type CreatePATResponse struct {
	Token *PersonalAccessToken `json:"token"`
	Plain string               `json:"plain"`
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
