package notifications

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type NotificationSettingsServiceFull interface {
	GetNotificationSettings(ctx context.Context) (*NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, settings *NotificationSettings) error
	ListSubscriptions(ctx context.Context, userID, orgID string) (any, error)
	UpsertSubscription(ctx context.Context, userID, orgID, category string, channels []string, enabled bool) (any, error)
	ListDefaults(ctx context.Context, orgID string) (any, error)
	UpsertDefault(ctx context.Context, orgID, category string, channels []string, enabled bool) (any, error)
	ListDeliveries(ctx context.Context, limit int) (any, error)
	RegisterDeviceToken(ctx context.Context, userID, token, platform string) (any, error)
	DeleteDeviceToken(ctx context.Context, userID, token string) error
	ListDeviceTokens(ctx context.Context, userID string) (any, error)
	ResolveOrgID(ctx context.Context, userID string) string
}

type NotificationSettingsHandler struct {
	service NotificationSettingsServiceFull
}

func NewNotificationSettingsHandler(s NotificationSettingsServiceFull) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{service: s}
}

func maskNotificationSecrets(s *NotificationSettings) {
	if s.SMTPPassword != "" {
		s.SMTPPassword = "********"
	}
	if s.ResendAPIKey != "" {
		s.ResendAPIKey = "********"
	}
	if s.TelegramBotToken != "" {
		s.TelegramBotToken = "********"
	}
	if s.PushoverAPIToken != "" {
		s.PushoverAPIToken = "********"
	}
	if s.SlackWebhookURL != "" {
		s.SlackWebhookURL = "********"
	}
}

func (h *NotificationSettingsHandler) GetNotificationSettings(c echo.Context) error {
	s, err := h.service.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	masked := *s
	maskNotificationSecrets(&masked)
	return utils.Success(c, "Operation successful", masked)
}

func (h *NotificationSettingsHandler) UpdateNotificationSettings(c echo.Context) error {
	existing, err := h.service.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to fetch existing notification settings")
	}

	realSMTP := existing.SMTPPassword
	realResend := existing.ResendAPIKey
	realTelegram := existing.TelegramBotToken
	realPushover := existing.PushoverAPIToken
	realSlack := existing.SlackWebhookURL

	if err := c.Bind(existing); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	if existing.SMTPPassword == "********" {
		existing.SMTPPassword = realSMTP
	}
	if existing.ResendAPIKey == "********" {
		existing.ResendAPIKey = realResend
	}
	if existing.TelegramBotToken == "********" {
		existing.TelegramBotToken = realTelegram
	}
	if existing.PushoverAPIToken == "********" {
		existing.PushoverAPIToken = realPushover
	}
	if existing.SlackWebhookURL == "********" {
		existing.SlackWebhookURL = realSlack
	}

	if err := h.service.UpdateNotificationSettings(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	maskNotificationSecrets(existing)
	return utils.Success(c, "Notification settings updated successfully", existing)
}

func (h *NotificationSettingsHandler) resolveOrgID(ctx context.Context, c echo.Context, userID string) string {
	orgID := c.QueryParam("orgId")
	if orgID == "" {
		orgID = c.QueryParam("organizationId")
	}
	if orgID == "default" {
		orgID = ""
	}
	if orgID == "" && h.service != nil {
		orgID = h.service.ResolveOrgID(ctx, userID)
	}
	return orgID
}

func (h *NotificationSettingsHandler) ListSubscriptions(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	orgID := h.resolveOrgID(c.Request().Context(), c, userID)
	if orgID == "" || orgID == "default" {
		return utils.Success(c, "Operation successful", []any{})
	}
	subs, err := h.service.ListSubscriptions(c.Request().Context(), userID, orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", subs)
}

func (h *NotificationSettingsHandler) UpsertSubscription(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req UpsertSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	orgID := h.resolveOrgID(c.Request().Context(), c, userID)
	if orgID == "" || orgID == "default" {
		return utils.Error(c, http.StatusBadRequest, "organization not found")
	}
	sub, err := h.service.UpsertSubscription(c.Request().Context(), userID, orgID, req.Category, req.Channels, req.Enabled)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Subscription updated successfully", sub)
}

func (h *NotificationSettingsHandler) ListDefaults(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	orgID := h.resolveOrgID(c.Request().Context(), c, userID)
	if orgID == "" || orgID == "default" {
		return utils.Success(c, "Operation successful", []any{})
	}
	defs, err := h.service.ListDefaults(c.Request().Context(), orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", defs)
}

func (h *NotificationSettingsHandler) UpsertDefault(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	var req UpsertDefaultRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	orgID := h.resolveOrgID(c.Request().Context(), c, userID)
	if orgID == "" || orgID == "default" {
		return utils.Error(c, http.StatusBadRequest, "organization not found")
	}
	def, err := h.service.UpsertDefault(c.Request().Context(), orgID, req.Category, req.Channels, req.Enabled)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Default updated successfully", def)
}

func (h *NotificationSettingsHandler) ListDeliveries(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	deliveries, err := h.service.ListDeliveries(c.Request().Context(), limit)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", deliveries)
}

func (h *NotificationSettingsHandler) RegisterDeviceToken(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req RegisterDeviceRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	if req.Token == "" {
		return utils.Error(c, http.StatusBadRequest, "token is required")
	}
	token, err := h.service.RegisterDeviceToken(c.Request().Context(), userID, req.Token, req.Platform)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Device token registered successfully", token)
}

func (h *NotificationSettingsHandler) DeleteDeviceToken(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	token := c.Param("token")
	if token == "" {
		token = c.QueryParam("token")
	}
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "token parameter is required")
	}
	if err := h.service.DeleteDeviceToken(c.Request().Context(), userID, token); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Device token deleted successfully", nil)
}

func (h *NotificationSettingsHandler) ListDeviceTokens(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	tokens, err := h.service.ListDeviceTokens(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", tokens)
}
