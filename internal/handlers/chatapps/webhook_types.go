package chatapps

// Telegram webhook payload structures
type TelegramWebhookUpdate struct {
	UpdateID int64            `json:"update_id"`
	Message  *TelegramMessage `json:"message,omitempty"`
}

type TelegramMessage struct {
	MessageID int64        `json:"message_id"`
	Text      string       `json:"text"`
	From      TelegramUser `json:"from"`
	Chat      TelegramChat `json:"chat"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TelegramChat struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title,omitempty"`
}

// Discord interaction webhook payload structures
type DiscordInteractionPayload struct {
	ID        string                  `json:"id"`
	Type      int                     `json:"type"`
	ChannelID string                  `json:"channel_id"`
	Member    *DiscordGuildMember     `json:"member,omitempty"`
	User      *DiscordUser            `json:"user,omitempty"`
	Data      *DiscordInteractionData `json:"data,omitempty"`
}

type DiscordGuildMember struct {
	User DiscordUser `json:"user"`
}

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type DiscordInteractionData struct {
	ID      string                     `json:"id"`
	Name    string                     `json:"name"`
	Options []DiscordInteractionOption `json:"options,omitempty"`
}

type DiscordInteractionOption struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Value any    `json:"value"`
}

// Teams Bot Framework activity payload structures
type TeamsActivityPayload struct {
	Type         string            `json:"type"`
	ID           string            `json:"id"`
	Text         string            `json:"text"`
	From         TeamsAccount      `json:"from"`
	Conversation TeamsConversation `json:"conversation"`
}

type TeamsAccount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AadObjectID string `json:"aadObjectId,omitempty"`
}

type TeamsConversation struct {
	ID               string `json:"id"`
	Name             string `json:"name,omitempty"`
	ConversationType string `json:"conversationType,omitempty"`
}

// Slack slash command form payload
type SlackSlashCommandPayload struct {
	Command     string `json:"command"`
	Text        string `json:"text"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	TeamID      string `json:"team_id"`
	ResponseURL string `json:"response_url"`
	TriggerID   string `json:"trigger_id"`
}
