package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func setAuthCookie(c echo.Context, token string) {
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

func clearAuthCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

func setRefreshCookie(c echo.Context, refreshToken string) {
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

func clearRefreshCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "codedock_refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

func HashAPIToken(token string) string {
	return token
}
