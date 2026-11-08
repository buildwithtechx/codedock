package models

import (
	"time"
)

type ChatAppProvider string

const (
	ChatAppProviderSlack    ChatAppProvider = "slack"
	ChatAppProviderDiscord  ChatAppProvider = "discord"
	ChatAppProviderTelegram ChatAppProvider = "telegram"
	ChatAppProviderTeams    ChatAppProvider = "teams"
)

type ChatAppTransport string

const (
	ChatAppTransportWebhook  ChatAppTransport = "webhook"
	ChatAppTransportOutbound ChatAppTransport = "outbound"
)

type ChatAppPurpose string

const (
	ChatAppPurposeCommands  ChatAppPurpose = "commands"
	ChatAppPurposeAssistant ChatAppPurpose = "assistant"
	ChatAppPurposeBoth      ChatAppPurpose = "both"
)

type ChatAppStatus string

const (
	ChatAppStatusConnected       ChatAppStatus = "connected"
	ChatAppStatusDisconnected    ChatAppStatus = "disconnected"
	ChatAppStatusError           ChatAppStatus = "error"
	ChatAppStatusPairingRequired ChatAppStatus = "pairing_required"
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

type CreateChatAppRequest struct {
	Provider    ChatAppProvider   `json:"provider"`
	Transport   ChatAppTransport  `json:"transport"`
	Purpose     ChatAppPurpose    `json:"purpose"`
	Name        string            `json:"name"`
	BotUsername string            `json:"botUsername"`
	Credentials map[string]string `json:"credentials"`
	Config      map[string]any    `json:"config"`
}

type UpdateChatAppRequest struct {
	Name        *string           `json:"name,omitempty"`
	Transport   *ChatAppTransport `json:"transport,omitempty"`
	Purpose     *ChatAppPurpose   `json:"purpose,omitempty"`
	BotUsername *string           `json:"botUsername,omitempty"`
	Enabled     *bool             `json:"enabled,omitempty"`
	Credentials map[string]string `json:"credentials,omitempty"`
	Config      map[string]any    `json:"config,omitempty"`
}

type VerifyPairingRequest struct {
	Token string `json:"token"`
	Code  string `json:"code"`
}

type ChatCommandRequest struct {
	ConnectionID      string          `json:"connectionId"`
	Provider          ChatAppProvider `json:"provider"`
	ExternalUserID    string          `json:"externalUserId"`
	ExternalUsername  string          `json:"externalUsername"`
	ExternalChannelID string          `json:"externalChannelId"`
	RawCommand        string          `json:"rawCommand"`
}

type ChatCommandResponse struct {
	Text                 string  `json:"text"`
	RequiresConfirmation bool    `json:"requiresConfirmation,omitempty"`
	ConfirmationToken    *string `json:"confirmationToken,omitempty"`
	DeepLink             *string `json:"deepLink,omitempty"`
}
