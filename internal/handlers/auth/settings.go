package auth

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/models"
	"codedock/internal/utils"
	"codedock/internal/version"

	authservices "codedock/internal/services/auth"
	systemservices "codedock/internal/services/system"
)

type SettingsHandler struct {
	settingsService      *authservices.SettingsService
	notifSettingsService *systemservices.NotificationSettingsService
	mu                   sync.Mutex
}

func NewSettingsHandler(s *authservices.SettingsService, ns *systemservices.NotificationSettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s, notifSettingsService: ns}
}

func maskSettingsSecrets(s *models.ServerSettings) {
	if s.CloudflareAPIToken != "" {
		s.CloudflareAPIToken = "********"
	}
	if s.NamecheapAPIKey != "" {
		s.NamecheapAPIKey = "********"
	}
	if s.SpaceshipAPIKey != "" {
		s.SpaceshipAPIKey = "********"
	}
	if s.SpaceshipAPISecret != "" {
		s.SpaceshipAPISecret = "********"
	}
}

func (h *SettingsHandler) GetSettings(c echo.Context) error {
	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	masked := *s
	maskSettingsSecrets(&masked)
	return utils.Success(c, "Operation successful", masked)
}

func (h *SettingsHandler) GetPublicSettings(c echo.Context) error {
	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	notif, err := h.notifSettingsService.GetNotificationSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	cfg := config.Get()
	publicSettings := map[string]any{
		"registrationEnabled":   s.RegistrationEnabled,
		"siteName":              s.SiteName,
		"emailEnabled":          notif.SMTPEnabled || notif.ResendEnabled || (cfg.SMTP.Host != "" && cfg.SMTP.From != "") || cfg.Resend.APIKey != "",
		"cloudMode":             cfg.Cloud.Enabled,
		"selfHosted":            !cfg.Cloud.Enabled,
		"version":               version.Version,
		"hostDomain":            cfg.Server.Domain,
		"defaultWildcardDomain": cfg.Domains.WildcardDomain,
		"serverIp":              cfg.Server.HostIP,
	}
	return utils.Success(c, "Operation successful", publicSettings)
}

func (h *SettingsHandler) GetHealthEnv(c echo.Context) error {
	cfg := config.Get()
	return c.JSON(http.StatusOK, map[string]any{
		"selfHosted":            !cfg.Cloud.Enabled,
		"cloudMode":             cfg.Cloud.Enabled,
		"deployMode":            "docker",
		"isServerHost":          !cfg.Cloud.Enabled,
		"hostControlEnabled":    true,
		"version":               version.Version,
		"authMode":              "local",
		"authProviders":         []string{},
		"productMode":           "platform",
		"teamMode":              "single_user",
		"hostDomain":            cfg.Server.Domain,
		"defaultWildcardDomain": cfg.Domains.WildcardDomain,
		"serverIp":              cfg.Server.HostIP,
	})
}

func (h *SettingsHandler) UpdateSettings(c echo.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	existing, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to fetch existing settings")
	}

	realCloudflare := existing.CloudflareAPIToken
	realNamecheap := existing.NamecheapAPIKey
	realSpaceship := existing.SpaceshipAPIKey
	realSpaceshipSecret := existing.SpaceshipAPISecret

	if err := c.Bind(existing); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if existing.CloudflareAPIToken == "********" {
		existing.CloudflareAPIToken = realCloudflare
	}
	if existing.NamecheapAPIKey == "********" {
		existing.NamecheapAPIKey = realNamecheap
	}
	if existing.SpaceshipAPIKey == "********" {
		existing.SpaceshipAPIKey = realSpaceship
	}
	if existing.SpaceshipAPISecret == "********" {
		existing.SpaceshipAPISecret = realSpaceshipSecret
	}

	if err := h.settingsService.UpdateSettings(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	maskSettingsSecrets(existing)
	return utils.Success(c, "Operation successful", existing)
}
