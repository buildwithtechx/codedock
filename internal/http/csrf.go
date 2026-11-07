package http

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"codedock/internal/config"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func bootstrapCSRF(c echo.Context) error {
	token := ""
	if cookie, err := c.Cookie("csrf_token"); err == nil {
		token = cookie.Value
	}
	if token == echomiddleware.CSRFUsingSecFetchSite {
		token = ""
	}
	if token == "" {
		if generated, ok := c.Get("csrf").(string); ok && generated != echomiddleware.CSRFUsingSecFetchSite {
			token = generated
		}
	}
	if token == "" {
		random := make([]byte, 16)
		if _, err := rand.Read(random); err != nil {
			return fmt.Errorf("generate CSRF token: %w", err)
		}
		token = hex.EncodeToString(random)
	}
	secure := c.IsTLS() || strings.HasPrefix(config.Get().Server.APIHost, "https://")
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}
	c.SetCookie(&http.Cookie{Name: "csrf_token", Value: token, Path: "/", MaxAge: 86400, Secure: secure, SameSite: sameSite})
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
