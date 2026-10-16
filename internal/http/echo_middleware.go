package http

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"codedock/internal/config"
	"codedock/internal/engine/networking"
)

func configureEchoMiddleware(e *echo.Echo) {
	configureEchoIPExtractor(e)
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			slog.Info("request", "method", v.Method, "uri", v.URI, "status", v.Status)
			return nil
		},
	}))
	e.Use(echomiddleware.Rewrite(map[string]string{
		"/api/v1/*": "/api/$1",
	}))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.SecureWithConfig(echomiddleware.SecureConfig{
		XFrameOptions:         "DENY",
		ContentSecurityPolicy: "frame-ancestors 'none'",
	}))
	e.Use(echomiddleware.GzipWithConfig(echomiddleware.GzipConfig{
		Level: 5,
	}))

	allowOrigins := []string{"http://localhost:3000", "http://localhost:8080"}
	if dashboardURL := config.Get().Server.DashboardURL; dashboardURL != "" && !slices.Contains(allowOrigins, dashboardURL) {
		allowOrigins = append(allowOrigins, dashboardURL)
	}
	for _, origin := range config.Get().Server.ExtraCORSOrigins {
		if origin != "" && !slices.Contains(allowOrigins, origin) {
			allowOrigins = append(allowOrigins, origin)
		}
	}

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	e.Use(echomiddleware.CSRFWithConfig(echomiddleware.CSRFConfig{
		TokenLength:  32,
		TokenLookup:  "header:X-CSRF-Token",
		CookieName:   "csrf_token",
		CookiePath:   "/",
		CookieMaxAge: 86400,
		ErrorHandler: func(err error, c echo.Context) error {
			return echo.NewHTTPError(http.StatusForbidden, "invalid or missing CSRF token").SetInternal(fmt.Errorf("validate CSRF request: %w", err))
		},
		Skipper: func(c echo.Context) bool {
			path := strings.TrimPrefix(c.Request().URL.Path, "/api/v1/")
			path = strings.TrimPrefix(path, "/api/")
			if path == "billing/webhook" && c.Request().Method == http.MethodPost && config.Get().Cloud.Enabled {
				return true
			}
			return !strings.HasPrefix(path, "auth/") && strings.HasPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
		},
	}))
}

func configureEchoIPExtractor(e *echo.Echo) {
	trustedNets := networking.ParseCIDRs(networking.GetEdgeTrustedProxies())
	header := networking.ValidateRealIPHeader(config.Get().Traefik.RealIPHeader)

	e.IPExtractor = func(req *http.Request) string {
		peerIPStr := req.RemoteAddr
		if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			peerIPStr = host
		}
		peerIP := net.ParseIP(peerIPStr)
		if peerIP == nil || !networking.IsIPInCIDRs(peerIP, trustedNets) {
			return peerIPStr
		}

		if strings.EqualFold(header, "X-Forwarded-For") {
			rawXFF := req.Header.Get("X-Forwarded-For")
			if rawXFF != "" {
				hops := strings.Split(rawXFF, ",")
				for i := len(hops) - 1; i >= 0; i-- {
					hop := strings.TrimSpace(hops[i])
					hopIP := net.ParseIP(hop)
					if hopIP != nil && !networking.IsIPInCIDRs(hopIP, trustedNets) {
						return hop
					}
				}
				if len(hops) > 0 {
					return strings.TrimSpace(hops[0])
				}
			}
			return peerIPStr
		}

		headerVal := strings.TrimSpace(req.Header.Get(header))
		if headerVal != "" {
			if ip := net.ParseIP(headerVal); ip != nil {
				return headerVal
			}
		}

		return peerIPStr
	}
}

