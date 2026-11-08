package chatapps

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/services/chatapps/transports"
	"codedock/internal/utils"
)

// WebhookIngress handles incoming event and command webhooks from Slack, Discord, Telegram, and Teams
func (h *Handler) WebhookIngress(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	conn, err := h.service.GetConnection(ctx, id)
	if err != nil || conn == nil || !conn.Enabled {
		return utils.Error(c, http.StatusNotFound, "chat connection not found or disabled")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "cannot read request body")
	}

	creds, _ := h.service.GetDecryptedCredentials(ctx, conn)

	var extUserID, extUsername, extChannelID, commandText string

	switch conn.Provider {
	case models.ChatAppProviderSlack:
		sig := c.Request().Header.Get("X-Slack-Signature")
		ts := c.Request().Header.Get("X-Slack-Request-Timestamp")
		if secret := creds["signing_secret"]; secret != "" {
			if err := transports.VerifySlackSignature(secret, sig, ts, bodyBytes); err != nil {
				return utils.Error(c, http.StatusUnauthorized, err.Error())
			}
		}
		// Form-urlencoded for slash commands
		c.Request().ParseForm()
		extUserID = c.Request().FormValue("user_id")
		extUsername = c.Request().FormValue("user_name")
		extChannelID = c.Request().FormValue("channel_id")
		commandText = c.Request().FormValue("command") + " " + c.Request().FormValue("text")

	case models.ChatAppProviderDiscord:
		sig := c.Request().Header.Get("X-Signature-Ed25519")
		ts := c.Request().Header.Get("X-Signature-Timestamp")
		if pubKey := creds["public_key"]; pubKey != "" {
			if err := transports.VerifyDiscordSignature(pubKey, sig, ts, bodyBytes); err != nil {
				return utils.Error(c, http.StatusUnauthorized, err.Error())
			}
		}
		var payload DiscordInteractionPayload
		_ = json.Unmarshal(bodyBytes, &payload)
		// Discord Ping (type 1)
		if payload.Type == 1 {
			return c.JSON(http.StatusOK, map[string]int{"type": 1})
		}
		// Interaction user extraction
		if payload.Member != nil {
			extUserID = payload.Member.User.ID
			extUsername = payload.Member.User.Username
		} else if payload.User != nil {
			extUserID = payload.User.ID
			extUsername = payload.User.Username
		}
		extChannelID = payload.ChannelID
		if payload.Data != nil {
			commandText = "/" + payload.Data.Name
			for _, opt := range payload.Data.Options {
				if strVal, ok := opt.Value.(string); ok {
					commandText += " " + strVal
				}
			}
		}

	case models.ChatAppProviderTelegram:
		tokenHeader := c.Request().Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if secret := creds["secret_token"]; secret != "" {
			if err := transports.VerifyTelegramSecret(secret, tokenHeader); err != nil {
				return utils.Error(c, http.StatusUnauthorized, err.Error())
			}
		}
		var update TelegramWebhookUpdate
		_ = json.Unmarshal(bodyBytes, &update)
		if update.Message != nil {
			commandText = update.Message.Text
			extUserID = strconv.FormatInt(update.Message.From.ID, 10)
			extUsername = update.Message.From.Username
			extChannelID = strconv.FormatInt(update.Message.Chat.ID, 10)
		}

	case models.ChatAppProviderTeams:
		authHeader := c.Request().Header.Get("Authorization")
		if appID := creds["app_id"]; appID != "" {
			if err := transports.VerifyTeamsRequest(appID, authHeader); err != nil {
				return utils.Error(c, http.StatusUnauthorized, err.Error())
			}
		}
		var activity TeamsActivityPayload
		_ = json.Unmarshal(bodyBytes, &activity)
		commandText = activity.Text
		extUserID = activity.From.ID
		extUsername = activity.From.Name
		extChannelID = activity.Conversation.ID
	}

	if strings.TrimSpace(commandText) == "" {
		return c.JSON(http.StatusOK, IngressIgnoredResponse{Status: "ignored", Reason: "empty command"})
	}

	cmdReq := models.ChatCommandRequest{
		ConnectionID:      conn.ID,
		Provider:          conn.Provider,
		ExternalUserID:    extUserID,
		ExternalUsername:  extUsername,
		ExternalChannelID: extChannelID,
		RawCommand:        commandText,
	}

	resp, err := h.dispatcher.DispatchCommand(ctx, cmdReq)
	if err != nil {
		return c.JSON(http.StatusOK, SimpleTextResponse{Text: "Error: " + err.Error()})
	}

	if conn.Provider == models.ChatAppProviderDiscord {
		return c.JSON(http.StatusOK, DiscordInteractionResponse{
			Type: 4,
			Data: &DiscordInteractionDataOut{
				Content: resp.Text,
			},
		})
	}

	return c.JSON(http.StatusOK, SimpleTextResponse{Text: resp.Text})
}
