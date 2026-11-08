package servers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"codedock/internal/http/middleware"
	"codedock/internal/platform/observability"
	"codedock/internal/utils"
)

type TokenValidator interface {
	ValidateToken(tokenStr string) (map[string]any, error)
}

type UserStatusChecker interface {
	GetUserByID(ctx context.Context, id string) (*UserInfo, error)
}

type ServerMetricsWSHandler struct {
	upgrader websocket.Upgrader
	tokens   TokenValidator
	service  Service
	users    UserStatusChecker
}

func NewServerMetricsWSHandler(tokens TokenValidator, service Service, users UserStatusChecker) *ServerMetricsWSHandler {
	return &ServerMetricsWSHandler{
		upgrader: websocket.Upgrader{
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
			Subprotocols: []string{"auth", "bearer"},
		},
		tokens:  tokens,
		service: service,
		users:   users,
	}
}

func (h *ServerMetricsWSHandler) Handle(c echo.Context) error {
	if err := utils.ValidateWebSocketCSWSH(c); err != nil {
		return err
	}
	serverID := c.Param("serverId")
	if serverID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serverId parameter")
	}

	if h.tokens != nil {
		tokenStr := middleware.ExtractTokenFromRequest(c)
		if tokenStr == "" {
			return utils.Error(c, http.StatusUnauthorized, "missing authentication token")
		}

		claimsMap, err := h.tokens.ValidateToken(tokenStr)
		if err != nil {
			return utils.Error(c, http.StatusUnauthorized, "invalid authentication token")
		}

		userID, _ := claimsMap["sub"].(string)
		if userID != "" && h.users != nil {
			u, err := h.users.GetUserByID(c.Request().Context(), userID)
			if err != nil || u == nil {
				return utils.Error(c, http.StatusUnauthorized, "user account not found or deactivated")
			}
		}

		role, _ := claimsMap["role"].(string)
		if role != "admin" && role != "owner" {
			if h.service == nil {
				return utils.Error(c, http.StatusInternalServerError, "server service not configured")
			}
			server, err := h.service.GetServer(c.Request().Context(), serverID, userID)
			if err != nil || server == nil {
				return utils.Error(c, http.StatusNotFound, "server not found")
			}
			if server.UserID != userID && !(server.IsLocal && server.UserID == "system") {
				return utils.Error(c, http.StatusForbidden, "insufficient permissions")
			}
		}
	}

	responseHeader := http.Header{}
	if c.Request().Header.Get("Sec-WebSocket-Protocol") != "" {
		responseHeader.Set("Sec-WebSocket-Protocol", "auth")
	}

	ws, err := h.upgrader.Upgrade(c.Response().Writer, c.Request(), responseHeader)
	if err != nil {
		slog.Error("failed to upgrade server metrics ws", "err", err)
		return err
	}

	observability.GlobalUIMetricsHub.AddClient(serverID, ws)
	defer observability.GlobalUIMetricsHub.RemoveClient(serverID, ws)
	defer ws.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				return
			}
		}
	}()

	_ = ws.WriteMessage(websocket.TextMessage, observability.GetHostMetricsPayload())

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			if err := ws.WriteMessage(websocket.TextMessage, observability.GetHostMetricsPayload()); err != nil {
				return nil
			}
		}
	}
}
