package models

import (
	"testing"
)

func TestValidateAppearance(t *testing.T) {
	t.Run("nil is valid", func(t *testing.T) {
		if err := ValidateAppearance(nil); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("valid default", func(t *testing.T) {
		pref := &AppearancePreferences{
			ActiveTheme: "nord",
			Mode:        "dark",
			FontFamily:  "Inter, sans-serif",
			Radius:      "8px",
			Tokens: map[string]string{
				"--primary": "#3b82f6",
				"--bg":      "#1e1e2e",
			},
		}
		if err := ValidateAppearance(pref); err != nil {
			t.Fatalf("expected valid pref, got %v", err)
		}
	})

	t.Run("invalid mode", func(t *testing.T) {
		pref := &AppearancePreferences{Mode: "neon"}
		if err := ValidateAppearance(pref); err == nil {
			t.Fatal("expected error for invalid mode")
		}
	})

	t.Run("illegal token character in value", func(t *testing.T) {
		tokens := map[string]string{
			"--bg": "red; background: black",
		}
		pref := &AppearancePreferences{Tokens: tokens}
		if err := ValidateAppearance(pref); err == nil {
			t.Fatal("expected error for semicolon in token")
		}
	})

	t.Run("illegal token character url() in value", func(t *testing.T) {
		tokens := map[string]string{
			"--bg": "url(https://malicious.site/leak)",
		}
		pref := &AppearancePreferences{Tokens: tokens}
		if err := ValidateAppearance(pref); err == nil {
			t.Fatal("expected error for url() in token")
		}
	})

	t.Run("illegal token character in key", func(t *testing.T) {
		tokens := map[string]string{
			"--bg{color}": "#fff",
		}
		pref := &AppearancePreferences{Tokens: tokens}
		if err := ValidateAppearance(pref); err == nil {
			t.Fatal("expected error for curly braces in token key")
		}
	})

	t.Run("custom theme validation", func(t *testing.T) {
		pref := &AppearancePreferences{
			CustomThemes: []CustomTheme{
				{
					ID:   "my-theme",
					Name: "My Theme",
					Mode: "light",
					Tokens: map[string]string{
						"--primary": "#000",
					},
				},
			},
		}
		if err := ValidateAppearance(pref); err != nil {
			t.Fatalf("expected valid custom theme, got %v", err)
		}

		pref.CustomThemes[0].Tokens["--leak"] = "calc(100% - 10px); color: red"
		if err := ValidateAppearance(pref); err == nil {
			t.Fatal("expected error for invalid token in custom theme")
		}
	})
}
