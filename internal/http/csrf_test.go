package http

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock.run/codedock/internal/config"
	"github.com/labstack/echo/v4"
)

func TestCSRFCookieUsesTLSOrConfiguredHTTPS(t *testing.T) {
	cfg := config.Get()
	previous := cfg.Server.APIHost
	t.Cleanup(func() { cfg.Server.APIHost = previous })
	for _, check := range []struct {
		name, apiHost string
		tls, secure   bool
	}{
		{"spoofed proxy header", "http://localhost:8080", false, false},
		{"direct TLS", "", true, true},
		{"configured HTTPS proxy", "https://api.example.com", false, true},
	} {
		t.Run(check.name, func(t *testing.T) {
			cfg.Server.APIHost = check.apiHost
			request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
			request.Header.Set("X-Forwarded-Proto", "https")
			request.AddCookie(&http.Cookie{Name: "csrf_token", Value: "_echo_csrf_using_sec_fetch_site_"})
			if check.tls {
				request.TLS = &tls.ConnectionState{}
			}
			response := httptest.NewRecorder()
			if err := bootstrapCSRF(echo.New().NewContext(request, response)); err != nil {
				t.Fatal(err)
			}
			cookies := response.Result().Cookies()
			if len(cookies) != 1 || cookies[0].Secure != check.secure || cookies[0].Value == "_echo_csrf_using_sec_fetch_site_" {
				t.Fatalf("invalid CSRF cookie: %+v", cookies)
			}
			expected := http.SameSiteLaxMode
			if check.secure {
				expected = http.SameSiteNoneMode
			}
			if cookies[0].SameSite != expected {
				t.Fatalf("invalid SameSite: %v", cookies[0].SameSite)
			}
		})
	}
}
