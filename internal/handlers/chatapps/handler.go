package chatapps

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	handlerutils "codedock/internal/handlers/utils"
	"codedock/internal/models"
	chatappservices "codedock/internal/services/chatapps"
	"codedock/internal/services/chatapps/transports"
	"codedock/internal/utils"
)

type Handler struct {
	service    *chatappservices.Service
	dispatcher *chatappservices.Dispatcher
}

func NewHandler(service *chatappservices.Service, dispatcher *chatappservices.Dispatcher) *Handler {
	return &Handler{
		service:    service,
		dispatcher: dispatcher,
	}
}

func (h *Handler) List(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	conns, err := h.service.ListConnections(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Chat connections retrieved", conns)
}

func (h *Handler) Create(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req models.CreateChatAppRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	conn, err := h.service.CreateConnection(c.Request().Context(), userID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Chat connection created", conn)
}

func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	conn, err := h.service.GetConnection(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Chat connection retrieved", conn)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req models.UpdateChatAppRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	conn, err := h.service.UpdateConnection(c.Request().Context(), id, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Chat connection updated", conn)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteConnection(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Chat connection deleted", nil)
}

func (h *Handler) Test(c echo.Context) error {
	id := c.Param("id")
	res, err := h.service.TestConnection(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Connection test completed", res)
}

func (h *Handler) ListIdentities(c echo.Context) error {
	id := c.Param("id")
	identities, err := h.service.ListIdentities(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Paired identities retrieved", identities)
}

func (h *Handler) DeleteIdentity(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteIdentity(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Paired identity removed", nil)
}

func (h *Handler) GetChallenge(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		token = c.QueryParam("token")
	}
	ch, err := h.service.GetChallenge(c.Request().Context(), token)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Pairing challenge found", ch)
}

func (h *Handler) VerifyPairing(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req models.VerifyPairingRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	identity, err := h.service.VerifyPairing(c.Request().Context(), userID, req.Token, req.Code)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Account pairing verified successfully", identity)
}

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
		var payload map[string]any
		_ = json.Unmarshal(bodyBytes, &payload)
		// Discord Ping (type 1)
		if typ, ok := payload["type"].(float64); ok && int(typ) == 1 {
			return c.JSON(http.StatusOK, map[string]any{"type": 1})
		}
		// Interaction user
		if member, ok := payload["member"].(map[string]any); ok {
			if user, ok := member["user"].(map[string]any); ok {
				extUserID, _ = user["id"].(string)
				extUsername, _ = user["username"].(string)
			}
		}
		extChannelID, _ = payload["channel_id"].(string)
		if data, ok := payload["data"].(map[string]any); ok {
			cmdName, _ := data["name"].(string)
			commandText = "/" + cmdName
			if opts, ok := data["options"].([]any); ok {
				for _, o := range opts {
					if optMap, ok := o.(map[string]any); ok {
						val, _ := optMap["value"].(string)
						commandText += " " + val
					}
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
		var update struct {
			Message struct {
				Text string `json:"text"`
				From struct {
					ID       int64  `json:"id"`
					Username string `json:"username"`
				} `json:"from"`
				Chat struct {
					ID int64 `json:"id"`
				} `json:"chat"`
			} `json:"message"`
		}
		_ = json.Unmarshal(bodyBytes, &update)
		commandText = update.Message.Text
		extUserID = strings.TrimSpace(string(rune(update.Message.From.ID)))
		extUsername = update.Message.From.Username
		extChannelID = strings.TrimSpace(string(rune(update.Message.Chat.ID)))

	case models.ChatAppProviderTeams:
		authHeader := c.Request().Header.Get("Authorization")
		if appID := creds["app_id"]; appID != "" {
			if err := transports.VerifyTeamsRequest(appID, authHeader); err != nil {
				return utils.Error(c, http.StatusUnauthorized, err.Error())
			}
		}
		var activity struct {
			Type string `json:"type"`
			Text string `json:"text"`
			From struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"from"`
			Conversation struct {
				ID string `json:"id"`
			} `json:"conversation"`
		}
		_ = json.Unmarshal(bodyBytes, &activity)
		commandText = activity.Text
		extUserID = activity.From.ID
		extUsername = activity.From.Name
		extChannelID = activity.Conversation.ID
	}

	if strings.TrimSpace(commandText) == "" {
		return c.JSON(http.StatusOK, map[string]string{"status": "ignored", "reason": "empty command"})
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
		return c.JSON(http.StatusOK, map[string]string{"text": "Error: " + err.Error()})
	}

	if conn.Provider == models.ChatAppProviderDiscord {
		return c.JSON(http.StatusOK, map[string]any{
			"type": 4,
			"data": map[string]any{
				"content": resp.Text,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"text": resp.Text})
}
