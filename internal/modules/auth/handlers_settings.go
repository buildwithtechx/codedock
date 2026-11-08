package auth

import (
	"net/http"
	"os"
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

	notifEmailEnabled, err := h.settingsService.IsNotificationEmailEnabled(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	cfg := config.Get()
	publicSettings := map[string]any{
		"registrationEnabled":   s.RegistrationEnabled,
		"siteName":              s.SiteName,
		"emailEnabled":          notifEmailEnabled || (cfg.SMTP.Host != "" && cfg.SMTP.From != "") || cfg.Resend.APIKey != "",
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

type AppearanceHandler struct {
	userRepo UserRepository
}

func NewAppearanceHandler(userRepo UserRepository) *AppearanceHandler {
	return &AppearanceHandler{userRepo: userRepo}
}

func (h *AppearanceHandler) GetAppearance(c echo.Context) error {
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	pref, err := h.userRepo.GetAppearance(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if pref == nil {
		pref = &AppearancePreferences{
			ActiveTheme: "default",
			Mode:        "dark",
		}
	}
	return c.JSON(http.StatusOK, map[string]any{"appearance": pref})
}

func (h *AppearanceHandler) SetAppearance(c echo.Context) error {
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req struct {
		Appearance *AppearancePreferences `json:"appearance"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid appearance payload")
	}

	if req.Appearance != nil {
		if err := ValidateAppearance(req.Appearance); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
	}

	if err := h.userRepo.SetAppearance(c.Request().Context(), userID, req.Appearance); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"appearance": req.Appearance})
}

type OnboardingHandler struct {
	userService *UserService
}

func NewOnboardingHandler(userService *UserService) *OnboardingHandler {
	return &OnboardingHandler{
		userService: userService,
	}
}

func (h *OnboardingHandler) SetupStatus(c echo.Context) error {
	count, err := h.userService.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	cwd, _ := os.Getwd()
	return utils.Success(c, "Setup status", map[string]any{
		"setupRequired": count == 0,
		"cwd":           cwd,
	})
}
