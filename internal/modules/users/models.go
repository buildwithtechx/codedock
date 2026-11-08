package users

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
