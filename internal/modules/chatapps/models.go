package chatapps

import (
	"time"
)

type ChatAppConnection struct {
	ID                   string           `json:"id"`
	UserID               string           `json:"userId"`
	Provider             ChatAppProvider  `json:"provider"`
	Transport            ChatAppTransport `json:"transport"`
	Purpose              ChatAppPurpose   `json:"purpose"`
	Name                 string           `json:"name"`
	BotUsername          string           `json:"botUsername"`
	CredentialsEncrypted string           `json:"-"`
	Config               map[string]any   `json:"config"`
	Enabled              bool             `json:"enabled"`
	Status               ChatAppStatus    `json:"status"`
	LastError            *string          `json:"lastError,omitempty"`
	LastEventAt          *time.Time       `json:"lastEventAt,omitempty"`
	CreatedAt            time.Time        `json:"createdAt"`
	UpdatedAt            time.Time        `json:"updatedAt"`
}

type ChatAppIdentity struct {
	ID               string          `json:"id"`
	ConnectionID     string          `json:"connectionId"`
	UserID           string          `json:"userId"`
	Provider         ChatAppProvider `json:"provider"`
	ExternalUserID   string          `json:"externalUserId"`
	ExternalUsername string          `json:"externalUsername"`
	Role             string          `json:"role"`
	Permissions      []string        `json:"permissions"`
	Active           bool            `json:"active"`
	PairedAt         time.Time       `json:"pairedAt"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

type ChatAppChallenge struct {
	ID                string          `json:"id"`
	ConnectionID      string          `json:"connectionId"`
	Provider          ChatAppProvider `json:"provider"`
	ExternalUserID    string          `json:"externalUserId"`
	ExternalChannelID string          `json:"externalChannelId"`
	Token             string          `json:"token"`
	Code              string          `json:"code"`
	ExpiresAt         time.Time       `json:"expiresAt"`
	VerifiedAt        *time.Time      `json:"verifiedAt,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
}

type ChatAppCommand struct {
	ID                string          `json:"id"`
	ConnectionID      string          `json:"connectionId"`
	IdentityID        *string         `json:"identityId,omitempty"`
	Provider          ChatAppProvider `json:"provider"`
	Command           string          `json:"command"`
	Args              string          `json:"args"`
	Status            string          `json:"status"`
	Response          string          `json:"response"`
	ConfirmationToken *string         `json:"confirmationToken,omitempty"`
	ConfirmedAt       *time.Time      `json:"confirmedAt,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
}
