package chatapps

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(nil, nil)
	if h == nil {
		t.Fatal("expected NewHandler to return a valid instance")
	}
}

func TestTelegramWebhookUpdateDeserialization(t *testing.T) {
	payload := `{
		"update_id": 998877,
		"message": {
			"message_id": 42,
			"text": "/codedock status",
			"from": {
				"id": 1234567890,
				"username": "solomon",
				"first_name": "Solomon",
				"last_name": "Olatunji"
			},
			"chat": {
				"id": -9876543210,
				"type": "group",
				"title": "Ops Team"
			}
		}
	}`

	var update TelegramWebhookUpdate
	if err := json.Unmarshal([]byte(payload), &update); err != nil {
		t.Fatalf("failed to unmarshal Telegram update: %v", err)
	}

	if update.Message == nil {
		t.Fatal("expected message to be present")
	}
	if update.Message.Text != "/codedock status" {
		t.Errorf("unexpected command text: %s", update.Message.Text)
	}

	extUserID := strconv.FormatInt(update.Message.From.ID, 10)
	if extUserID != "1234567890" {
		t.Errorf("expected numeric string ID '1234567890', got %q", extUserID)
	}

	extChatID := strconv.FormatInt(update.Message.Chat.ID, 10)
	if extChatID != "-9876543210" {
		t.Errorf("expected numeric string Chat ID '-9876543210', got %q", extChatID)
	}
}

func TestDiscordInteractionPayloadDeserialization(t *testing.T) {
	payload := `{
		"id": "interaction-123",
		"type": 2,
		"channel_id": "channel-456",
		"member": {
			"user": {
				"id": "discord-user-789",
				"username": "devops"
			}
		},
		"data": {
			"id": "cmd-1",
			"name": "codedock",
			"options": [
				{
					"name": "action",
					"type": 3,
					"value": "deploy"
				}
			]
		}
	}`

	var interaction DiscordInteractionPayload
	if err := json.Unmarshal([]byte(payload), &interaction); err != nil {
		t.Fatalf("failed to unmarshal Discord interaction: %v", err)
	}

	if interaction.Type != 2 {
		t.Errorf("expected interaction type 2, got %d", interaction.Type)
	}
	if interaction.Member == nil || interaction.Member.User.ID != "discord-user-789" {
		t.Fatalf("unexpected member user ID: %+v", interaction.Member)
	}
	if interaction.Data == nil || interaction.Data.Name != "codedock" {
		t.Fatalf("unexpected data name: %+v", interaction.Data)
	}
	if len(interaction.Data.Options) != 1 || interaction.Data.Options[0].Value != "deploy" {
		t.Errorf("unexpected options: %+v", interaction.Data.Options)
	}
}

func TestTeamsActivityPayloadDeserialization(t *testing.T) {
	payload := `{
		"type": "message",
		"id": "activity-999",
		"text": "/codedock logs api",
		"from": {
			"id": "teams-user-1",
			"name": "Solomon"
		},
		"conversation": {
			"id": "teams-channel-primary",
			"name": "General"
		}
	}`

	var act TeamsActivityPayload
	if err := json.Unmarshal([]byte(payload), &act); err != nil {
		t.Fatalf("failed to unmarshal Teams activity: %v", err)
	}

	if act.Type != "message" {
		t.Errorf("expected activity type message, got %s", act.Type)
	}
	if act.From.ID != "teams-user-1" {
		t.Errorf("expected From ID teams-user-1, got %s", act.From.ID)
	}
	if act.Conversation.ID != "teams-channel-primary" {
		t.Errorf("expected Conversation ID teams-channel-primary, got %s", act.Conversation.ID)
	}
}
