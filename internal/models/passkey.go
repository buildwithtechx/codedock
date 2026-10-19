package models

import "time"

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

type PasskeyRenameRequest struct {
	Name string `json:"name"`
}

type PasskeyRegisterVerifyRequest struct {
	Name     string `json:"name"`
	Response any    `json:"response"`
}

type PasskeyAuthVerifyRequest struct {
	TOTPCode string `json:"totpCode,omitempty"`
	Response any    `json:"response"`
}
