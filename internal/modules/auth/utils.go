package auth

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/users"
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

func extractClaims(c echo.Context) (userID, email string) {
	if u := c.Get("user"); u != nil {
		if claims, ok := u.(*users.UserClaims); ok && claims != nil {
			return claims.UserID, claims.Email
		}
		v := reflect.ValueOf(u)
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			if f := v.FieldByName("UserID"); f.IsValid() && f.Kind() == reflect.String {
				userID = f.String()
			}
			if f := v.FieldByName("Email"); f.IsValid() && f.Kind() == reflect.String {
				email = f.String()
			}
		}
	}
	return userID, email
}

func extractUserID(c echo.Context) string {
	id, _ := extractClaims(c)
	return id
}
