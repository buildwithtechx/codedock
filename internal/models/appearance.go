package models

import (
	"errors"
	"strings"
)

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
