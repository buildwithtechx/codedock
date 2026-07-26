package middleware

import (
	"net"
	"strings"

	"github.com/labstack/echo/v4"
)

func IsIPAllowed(clientIPStr string, allowlistStr string) bool {
	clientIP := net.ParseIP(clientIPStr)
	if clientIP == nil {
		return false
	}
	entries := strings.SplitSeq(allowlistStr, ",")
	for entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			_, cidrNet, err := net.ParseCIDR(entry)
			if err == nil && cidrNet.Contains(clientIP) {
				return true
			}
		} else {
			if clientIPStr == entry {
				return true
			}
		}
	}
	return false
}

func ExtractClientIP(c echo.Context) string {
	return c.RealIP()
}

func ExtractTokenFromRequest(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}
	if wsProtocols := c.Request().Header.Get("Sec-WebSocket-Protocol"); wsProtocols != "" {
		parts := strings.Split(wsProtocols, ",")
		for i := 0; i < len(parts)-1; i++ {
			key := strings.TrimSpace(strings.ToLower(parts[i]))
			val := strings.TrimSpace(parts[i+1])
			if (key == "auth" || key == "bearer" || key == "token") && val != "" {
				return val
			}
		}
	}
	cookie, err := c.Cookie("codedock_token")
	if err == nil && cookie.Value != "" {
		return strings.TrimSpace(cookie.Value)
	}
	return ""
}
