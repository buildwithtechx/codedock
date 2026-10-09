package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	handlerutils "codedock/internal/handlers/utils"
	"codedock/internal/models"
	systemservices "codedock/internal/services/system"
	"codedock/internal/utils"
)

type NotificationSettingsHandler struct {
	notifSettingsService *systemservices.NotificationSettingsService
}

func NewNotificationSettingsHandler(s *systemservices.NotificationSettingsService) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{notifSettingsService: s}
}

func maskNotificationSecrets(s *models.NotificationSettings) {
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
	s, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	masked := *s
	maskNotificationSecrets(&masked)
	return utils.Success(c, "Operation successful", masked)
}

func (h *NotificationSettingsHandler) UpdateNotificationSettings(c echo.Context) error {
	existing, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to fetch existing notification settings")
	}

	realSMTP := existing.SMTPPassword
	realResend := existing.ResendAPIKey
	realTelegram := existing.TelegramBotToken
	realPushover := existing.PushoverAPIToken
	realSlack := existing.SlackWebhookURL

	if err := c.Bind(existing); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
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

	if err := h.notifSettingsService.UpdateNotificationSettings(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	maskNotificationSecrets(existing)
	return utils.Success(c, "Notification settings updated successfully", existing)
}

func (h *NotificationSettingsHandler) ListSubscriptions(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.QueryParam("orgId")
	}
	if orgID == "" {
		orgID = "default"
	}
	subs, err := h.notifSettingsService.ListSubscriptions(c.Request().Context(), userID, orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", subs)
}

func (h *NotificationSettingsHandler) UpsertSubscription(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req models.UpsertSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.QueryParam("orgId")
	}
	if orgID == "" {
		orgID = "default"
	}
	sub := models.NotificationSubscription{
		UserID:         userID,
		OrganizationID: orgID,
		Category:       req.Category,
		Channels:       req.Channels,
		Enabled:        req.Enabled,
	}
	if err := h.notifSettingsService.UpsertSubscription(c.Request().Context(), &sub); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Subscription updated successfully", sub)
}

func (h *NotificationSettingsHandler) ListDefaults(c echo.Context) error {
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.QueryParam("orgId")
	}
	if orgID == "" {
		orgID = "default"
	}
	defs, err := h.notifSettingsService.ListDefaults(c.Request().Context(), orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", defs)
}

func (h *NotificationSettingsHandler) UpsertDefault(c echo.Context) error {
	var req models.UpsertDefaultRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.QueryParam("orgId")
	}
	if orgID == "" {
		orgID = "default"
	}
	def := models.NotificationDefault{
		OrganizationID: orgID,
		Category:       req.Category,
		Channels:       req.Channels,
		Enabled:        req.Enabled,
	}
	if err := h.notifSettingsService.UpsertDefault(c.Request().Context(), &def); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Default updated successfully", def)
}
