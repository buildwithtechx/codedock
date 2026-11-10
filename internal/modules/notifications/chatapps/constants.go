package chatapps

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
