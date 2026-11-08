package monitoring

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type TokenValidator interface {
	ValidateToken(token string) (jwt.MapClaims, error)
}

type UserValidator interface {
	IsUserActive(ctx context.Context, userID string) (bool, error)
}

type ServiceAccessChecker interface {
	HasAccess(ctx context.Context, serviceID, userID, role string) (bool, error)
}

type ServiceLogsWSHandler struct {
	Streams       ServiceLogStreamer
	upgrader      websocket.Upgrader
	tokenService  TokenValidator
	accessChecker ServiceAccessChecker
	userValidator UserValidator
}

func NewServiceLogsWSHandler(
	ts TokenValidator,
	ac ServiceAccessChecker,
	uv UserValidator,
) *ServiceLogsWSHandler {
	return &ServiceLogsWSHandler{
		upgrader: websocket.Upgrader{
			Subprotocols: []string{"bearer"},
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					if r.Header.Get("Authorization") != "" || r.Header.Get("Sec-WebSocket-Protocol") != "" {
						return true
					}
					return false
				}
				return utils.IsAllowedWebSocketOrigin(r, origin)
			},
		},
		tokenService:  ts,
		accessChecker: ac,
		userValidator: uv,
	}
}

func (h *ServiceLogsWSHandler) Handle(c echo.Context) error {
	if err := utils.ValidateWebSocketCSWSH(c); err != nil {
		return err
	}
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	var claimsMap map[string]interface{}
	if h.tokenService != nil {
		tokenStr := extractWSToken(c)
		if tokenStr == "" {
			return utils.Error(c, http.StatusUnauthorized, "missing authentication token")
		}

		cm, err := h.tokenService.ValidateToken(tokenStr)
		if err != nil {
			return utils.Error(c, http.StatusUnauthorized, "invalid authentication token")
		}
		claimsMap = cm

		userID, _ := claimsMap["sub"].(string)
		if userID != "" && h.userValidator != nil {
			active, err := h.userValidator.IsUserActive(c.Request().Context(), userID)
			if err != nil || !active {
				return utils.Error(c, http.StatusUnauthorized, "user account not found or deactivated")
			}
		}
	}

	if h.accessChecker != nil && claimsMap != nil {
		userID, _ := claimsMap["sub"].(string)
		role, _ := claimsMap["role"].(string)
		hasAccess, err := h.accessChecker.HasAccess(c.Request().Context(), serviceID, userID, role)
		if err != nil || !hasAccess {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions to access service logs")
		}
	}

	responseHeader := http.Header{}
	if reqProto := c.Request().Header.Get("Sec-WebSocket-Protocol"); reqProto != "" {
		parts := strings.Split(reqProto, ",")
		if len(parts) > 0 {
			responseHeader.Set("Sec-WebSocket-Protocol", strings.TrimSpace(parts[0]))
		}
	}

	ws, err := h.upgrader.Upgrade(c.Response().Writer, c.Request(), responseHeader)
	if err != nil {
		slog.Error("failed to upgrade service logs ws", "err", err)
		return err
	}
	defer ws.Close()

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()
	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	if h.Streams == nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("Log stream unavailable"))
		return nil
	}
	if err := h.Streams.StreamServiceLogs(ctx, serviceID, serviceLogWSWriter{ws}); err != nil && ctx.Err() == nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("Log stream failed: "+err.Error()))
	}
	return nil
}

func extractWSToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}
	if q := c.QueryParam("token"); q != "" {
		return strings.TrimSpace(q)
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
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if strings.HasPrefix(strings.ToLower(trimmed), "bearer ") {
				return strings.TrimSpace(trimmed[7:])
			}
		}
	}
	return ""
}

type serviceLogWSWriter struct{ ws *websocket.Conn }

func (w serviceLogWSWriter) Write(data []byte) (int, error) {
	if err := w.ws.SetWriteDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return 0, err
	}
	if err := w.ws.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}
