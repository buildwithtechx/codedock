package settings

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/utils"
	"codedock/internal/version"
)

type SettingsHandler struct {
	settingsService *SettingsService
	mu              sync.Mutex
}

func NewSettingsHandler(s *SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s}
}

func maskSettingsSecrets(s *ServerSettings) {
	if s.CloudflareAPIToken != "" {
		s.CloudflareAPIToken = MaskedSecret
	}
	if s.NamecheapAPIKey != "" {
		s.NamecheapAPIKey = MaskedSecret
	}
	if s.SpaceshipAPIKey != "" {
		s.SpaceshipAPIKey = MaskedSecret
	}
	if s.SpaceshipAPISecret != "" {
		s.SpaceshipAPISecret = MaskedSecret
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

	notifEmailEnabled, err := h.settingsService.IsNotificationEmailEnabled(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	cfg := config.Get()
	publicSettings := PublicSettingsResponse{
		RegistrationEnabled:   s.RegistrationEnabled,
		SiteName:              s.SiteName,
		EmailEnabled:          notifEmailEnabled || (cfg.SMTP.Host != "" && cfg.SMTP.From != "") || cfg.Resend.APIKey != "",
		CloudMode:             cfg.Cloud.Enabled,
		SelfHosted:            !cfg.Cloud.Enabled,
		Version:               version.Version,
		HostDomain:            cfg.Server.Domain,
		DefaultWildcardDomain: cfg.Domains.WildcardDomain,
		ServerIP:              cfg.Server.HostIP,
	}
	return utils.Success(c, "Operation successful", publicSettings)
}

func (h *SettingsHandler) GetHealthEnv(c echo.Context) error {
	cfg := config.Get()
	return c.JSON(http.StatusOK, HealthEnvResponse{
		SelfHosted:            !cfg.Cloud.Enabled,
		CloudMode:             cfg.Cloud.Enabled,
		DeployMode:            "docker",
		IsServerHost:          !cfg.Cloud.Enabled,
		HostControlEnabled:    true,
		Version:               version.Version,
		AuthMode:              "local",
		AuthProviders:         []string{},
		ProductMode:           "platform",
		TeamMode:              "single_user",
		HostDomain:            cfg.Server.Domain,
		DefaultWildcardDomain: cfg.Domains.WildcardDomain,
		ServerIP:              cfg.Server.HostIP,
	})
}

func (h *SettingsHandler) UpdateSettings(c echo.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	existing, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, ErrFetchSettingsFailed.Error())
	}

	realCloudflare := existing.CloudflareAPIToken
	realNamecheap := existing.NamecheapAPIKey
	realSpaceship := existing.SpaceshipAPIKey
	realSpaceshipSecret := existing.SpaceshipAPISecret

	if err := c.Bind(existing); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if existing.CloudflareAPIToken == MaskedSecret {
		existing.CloudflareAPIToken = realCloudflare
	}
	if existing.NamecheapAPIKey == MaskedSecret {
		existing.NamecheapAPIKey = realNamecheap
	}
	if existing.SpaceshipAPIKey == MaskedSecret {
		existing.SpaceshipAPIKey = realSpaceship
	}
	if existing.SpaceshipAPISecret == MaskedSecret {
		existing.SpaceshipAPISecret = realSpaceshipSecret
	}

	if err := h.settingsService.UpdateSettings(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	maskSettingsSecrets(existing)
	return utils.Success(c, "Operation successful", existing)
}
