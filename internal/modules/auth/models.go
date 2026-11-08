package auth

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type User struct {
	ID                   string     `json:"id" db:"id"`
	Email                string     `json:"email" db:"email"`
	Name                 string     `json:"name" db:"name"`
	PasswordHash         string     `json:"-" db:"password_hash"`
	Role                 UserRole   `json:"role" db:"role"`
	IsActive             bool       `json:"isActive" db:"is_active"`
	EmailVerified        bool       `json:"emailVerified" db:"email_verified"`
	TOTPEnabled          bool       `json:"totpEnabled" db:"totp_enabled"`
	OAuthProvider        string     `json:"oauthProvider,omitempty" db:"oauth_provider"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time  `json:"updatedAt" db:"updated_at"`
	LastLogin            *time.Time `json:"lastLogin,omitempty" db:"last_login"`
	ProjectsCount        int        `json:"projectsCount" db:"projects_count"`
	ServicesCount        int        `json:"servicesCount" db:"services_count"`
	APIKeysCount         int        `json:"apiKeysCount" db:"api_keys_count"`
	PlanType             string     `json:"planType" db:"plan_type"`
	StripeCustomerID     *string    `json:"stripeCustomerId,omitempty" db:"stripe_customer_id"`
	StripeSubscriptionID *string    `json:"stripeSubscriptionId,omitempty" db:"stripe_subscription_id"`
	StripePriceID        *string    `json:"stripePriceId,omitempty" db:"stripe_price_id"`
}

type UserClaims struct {
	UserID      string   `json:"sub"`
	Email       string   `json:"email"`
	Role        UserRole `json:"role"`
	TOTPEnabled bool     `json:"totpEnabled"`
	PlanType    string   `json:"planType"`
}

type PersonalAccessToken struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"userId" db:"user_id"`
	Name            string     `json:"name" db:"name"`
	TokenHash       string     `json:"-" db:"token_hash"`
	Prefix          string     `json:"prefix" db:"prefix"`
	AccessLevel     string     `json:"accessLevel" db:"access_level"`
	ProjectScope    string     `json:"projectScope" db:"project_scope"`
	AllowedProjects *string    `json:"allowedProjects" db:"allowed_projects"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
}

func (p PersonalAccessToken) MarshalJSON() ([]byte, error) {
	type tokenAlias PersonalAccessToken
	var projects []string
	if p.AllowedProjects != nil {
		if err := json.Unmarshal([]byte(*p.AllowedProjects), &projects); err != nil {
			return nil, err
		}
	}
	return json.Marshal(struct {
		tokenAlias
		AllowedProjects []string `json:"allowedProjects"`
	}{tokenAlias: tokenAlias(p), AllowedProjects: projects})
}

type OAuthProviderConfig struct {
	ID           string    `json:"id" db:"id"`
	ProviderName string    `json:"providerName" db:"provider_name"`
	Enabled      bool      `json:"enabled" db:"enabled"`
	ClientID     string    `json:"clientId" db:"client_id"`
	ClientSecret string    `json:"clientSecret,omitempty" db:"client_secret"`
	RedirectURI  string    `json:"redirectUri" db:"redirect_uri"`
	BaseURL      string    `json:"baseUrl,omitempty" db:"base_url"`
	Tenant       string    `json:"tenant,omitempty" db:"tenant"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type PasskeyCredential struct {
	ID           string     `json:"id" db:"id"`
	UserID       string     `json:"userId" db:"user_id"`
	Name         string     `json:"name" db:"name"`
	CredentialID string     `json:"credentialId" db:"credential_id"`
	PublicKey    string     `json:"publicKey" db:"public_key"`
	Counter      int64      `json:"counter" db:"counter"`
	DeviceType   string     `json:"deviceType" db:"device_type"`
	BackedUp     bool       `json:"backedUp" db:"backed_up"`
	Transports   string     `json:"transports" db:"transports"`
	AAGUID       string     `json:"aaguid" db:"aaguid"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty" db:"last_used_at"`
}

type PasskeyChallenge struct {
	Challenge    string    `json:"challenge" db:"challenge"`
	UserID       string    `json:"userId" db:"user_id"`
	CeremonyType string    `json:"ceremonyType" db:"ceremony_type"`
	SessionData  string    `json:"sessionData" db:"session_data"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
}

type Organization struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type OrganizationMember struct {
	ID             string           `json:"id" db:"id"`
	OrganizationID string           `json:"organizationId" db:"organization_id"`
	UserID         string           `json:"userId,omitempty" db:"user_id"`
	Email          string           `json:"email" db:"email"`
	Permission     MemberPermission `json:"permission" db:"permission"`
	Status         MemberStatus     `json:"status" db:"status"`
	InvitedAt      time.Time        `json:"invitedAt" db:"invited_at"`
	AcceptedAt     time.Time        `json:"acceptedAt" db:"accepted_at"`
	InviteLink     string           `json:"inviteLink,omitempty" db:"-"`
}

type DeviceSession struct {
	DeviceCode string              `json:"deviceCode"`
	UserCode   string              `json:"userCode"`
	ClientName string              `json:"clientName"`
	Status     DeviceSessionStatus `json:"status"`
	UserID     string              `json:"userId,omitempty"`
	Token      string              `json:"token,omitempty"`
	UserEmail  string              `json:"userEmail,omitempty"`
	UserName   string              `json:"userName,omitempty"`
	ExpiresAt  time.Time           `json:"expiresAt"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type CustomTheme struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Mode   string            `json:"mode"`
	Tokens map[string]string `json:"tokens"`
}

type AppearancePreferences struct {
	ActiveTheme  string            `json:"activeTheme"`
	Mode         string            `json:"mode"`
	FontFamily   string            `json:"fontFamily,omitempty"`
	Radius       string            `json:"radius,omitempty"`
	Tokens       map[string]string `json:"tokens,omitempty"`
	CustomThemes []CustomTheme     `json:"customThemes,omitempty"`
}

var illegalTokenChars = []string{";", "{", "}", "<", ">", "\\", "url("}

func ValidateAppearance(pref *AppearancePreferences) error {
	if pref == nil {
		return nil
	}
	if pref.Mode != "" {
		mode := strings.ToLower(pref.Mode)
		if mode != "light" && mode != "dim" && mode != "dark" && mode != "system" {
			return errors.New("mode must be light, dim, dark, or system")
		}
		pref.Mode = mode
	}
	if len(pref.ActiveTheme) > 64 {
		return errors.New("activeTheme cannot exceed 64 characters")
	}
	if len(pref.FontFamily) > 128 {
		return errors.New("fontFamily cannot exceed 128 characters")
	}
	if len(pref.Radius) > 32 {
		return errors.New("radius cannot exceed 32 characters")
	}
	if err := validateTokens(pref.Tokens); err != nil {
		return err
	}
	if len(pref.CustomThemes) > 32 {
		return errors.New("at most 32 custom themes are supported")
	}
	for _, ct := range pref.CustomThemes {
		if strings.TrimSpace(ct.ID) == "" || len(ct.ID) > 64 {
			return errors.New("custom theme ID must be non-empty and at most 64 characters")
		}
		if len(ct.Name) > 64 {
			return errors.New("custom theme name cannot exceed 64 characters")
		}
		if ct.Mode != "" {
			mode := strings.ToLower(ct.Mode)
			if mode != "light" && mode != "dim" && mode != "dark" && mode != "system" {
				return errors.New("custom theme mode must be light, dim, dark, or system")
			}
		}
		if err := validateTokens(ct.Tokens); err != nil {
			return err
		}
	}
	return nil
}

func validateTokens(tokens map[string]string) error {
	if len(tokens) > 128 {
		return errors.New("at most 128 tokens are supported per theme")
	}
	for k, v := range tokens {
		if len(k) > 64 {
			return errors.New("token key cannot exceed 64 characters")
		}
		if len(v) > 512 {
			return errors.New("token value cannot exceed 512 characters")
		}
		for _, illegal := range illegalTokenChars {
			if strings.Contains(k, illegal) || strings.Contains(strings.ToLower(v), illegal) {
				return errors.New("token contains illegal CSS character or expression")
			}
		}
	}
	return nil
}

type ServerSettings struct {
	ID                          string `json:"id"`
	TraefikWildcardIP           string `json:"traefikWildcardIp"`
	RegistrationEnabled         bool   `json:"registrationEnabled"`
	RegistrationDomainAllowlist string `json:"registrationDomainAllowlist,omitempty"`
	CustomDNSResolvers          string `json:"customDnsResolvers"`
	DNSValidationEnabled        bool   `json:"dnsValidationEnabled"`
	IPAllowlist                 string `json:"ipAllowlist"`
	MCPServerEnabled            bool   `json:"mcpServerEnabled"`
	DefaultWildcardDomain       string `json:"defaultWildcardDomain,omitempty"`
	PanelDomain                 string `json:"dashboardDomain,omitempty"`
	SiteName                    string `json:"siteName,omitempty"`
	PublicIPv4                  string `json:"publicIpv4,omitempty"`
	PublicIPv6                  string `json:"publicIpv6,omitempty"`
	ShowSponsorshipPopup        bool   `json:"showSponsorshipPopup"`
	DisableTwoStepConfirmation  bool   `json:"disableTwoStepConfirmation"`
	CloudflareAPIToken          string `json:"cloudflareApiToken,omitempty"`
	NamecheapAPIUser            string `json:"namecheapApiUser,omitempty"`
	NamecheapAPIKey             string `json:"namecheapApiKey,omitempty"`
	NamecheapClientIP           string `json:"namecheapClientIp,omitempty"`
	SpaceshipAPIKey             string `json:"spaceshipApiKey,omitempty"`
	SpaceshipAPISecret          string `json:"spaceshipApiSecret,omitempty"`
	UpdateCheckCron             string `json:"updateCheckCron"`
	AutoUpdateEnabled           bool   `json:"autoUpdateEnabled"`
	TelemetryEnabled            bool   `json:"telemetryEnabled"`
	ConcurrentBuilds            int    `json:"concurrentBuilds"`
	DeploymentTimeout           int    `json:"deploymentTimeout"`
	ServerTimezone              string `json:"serverTimezone"`
	DockerCleanupCron           string `json:"dockerCleanupCron"`
	DiskUsageThreshold          int    `json:"diskUsageThreshold"`
	DiskUsageCron               string `json:"diskUsageCron"`
	CurrentVersion              string `json:"currentVersion"`
	LatestVersion               string `json:"latestVersion"`
	LastUpdateCheck             string `json:"lastUpdateCheck"`
	UpdatedAt                   string `json:"updatedAt"`
}
