package chatapps

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
	Token string `json:"token" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

type ChatCommandRequest struct {
	ConnectionID      string          `json:"connectionId"`
	Provider          ChatAppProvider `json:"provider"`
	ExternalUserID    string          `json:"externalUserId"`
	ExternalUsername  string          `json:"externalUsername"`
	ExternalChannelID string          `json:"externalChannelId"`
	RawCommand        string          `json:"rawCommand"`
}
