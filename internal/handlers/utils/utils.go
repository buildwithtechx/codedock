package handlerutils

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func SetAuthCookie(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   72 * 3600,
	})
}

func ClearAuthCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

func SetRefreshCookie(c echo.Context, refreshToken string) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   30 * 24 * 3600,
	})
}

func ClearRefreshCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

func ExtractClaims(c echo.Context) *models.UserClaims {
	if claims, ok := c.Get("user").(*models.UserClaims); ok {
		return claims
	}
	return nil
}

func ExtractUserID(c echo.Context) string {
	if claims := ExtractClaims(c); claims != nil {
		return claims.UserID
	}
	return ""
}

func IsAllowedWebSocketOrigin(r *http.Request, origin string) bool {
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if strings.EqualFold(u.Host, r.Host) {
		return true
	}
	allowed := []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:8080",
	}
	if dashURL := config.Get().Server.DashboardURL; dashURL != "" {
		allowed = append(allowed, dashURL)
	}
	for _, a := range allowed {
		if strings.EqualFold(strings.TrimRight(origin, "/"), strings.TrimRight(a, "/")) {
			return true
		}
	}
	return false
}

func ValidateWebSocketCSWSH(c echo.Context) error {
	r := c.Request()

	if auth := r.Header.Get("Authorization"); auth != "" && strings.HasPrefix(auth, "Bearer ") {
		return nil
	}

	if csrfHeader := r.Header.Get("X-CSRF-Token"); csrfHeader != "" {
		if cookie, err := c.Cookie("csrf_token"); err == nil && cookie.Value != "" && cookie.Value == csrfHeader {
			return nil
		}
	}

	if wsProtocols := r.Header.Get("Sec-WebSocket-Protocol"); wsProtocols != "" {
		parts := strings.Split(wsProtocols, ",")
		for i := 0; i < len(parts)-1; i++ {
			key := strings.TrimSpace(strings.ToLower(parts[i]))
			val := strings.TrimSpace(parts[i+1])
			if key == "auth" || key == "bearer" || key == "token" {
				if val != "" {
					return nil
				}
			}
			if key == "csrf" {
				if cookie, err := c.Cookie("csrf_token"); err == nil && cookie.Value != "" && cookie.Value == val {
					return nil
				}
			}
		}
	}

	origin := r.Header.Get("Origin")
	if !IsAllowedWebSocketOrigin(r, origin) {
		return utils.Error(c, http.StatusForbidden, "cross-origin WebSocket connection denied (CSWSH protection)")
	}

	return nil
}
