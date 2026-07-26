package handlers

import (
	"log/slog"
	"net/http"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ServerMetricsWSHandler struct {
	upgrader      websocket.Upgrader
	tokenService  *services.TokenService
	serverService services.ServerService
}

func NewServerMetricsWSHandler(ts *services.TokenService, ss services.ServerService) *ServerMetricsWSHandler {
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
				return isAllowedWebSocketOrigin(r, origin)
			},
		},
		tokenService:  ts,
		serverService: ss,
	}
}

func (h *ServerMetricsWSHandler) Handle(c echo.Context) error {
	if err := ValidateWebSocketCSWSH(c); err != nil {
		return err
	}
	serverID := c.Param("serverId")
	if serverID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serverId parameter")
	}

	if h.tokenService != nil {
		tokenStr := middleware.ExtractTokenFromRequest(c)
		if tokenStr == "" {
			return utils.Error(c, http.StatusUnauthorized, "missing authentication token")
		}

		claimsMap, err := h.tokenService.ValidateToken(tokenStr)
		if err != nil {
			return utils.Error(c, http.StatusUnauthorized, "invalid authentication token")
		}

		userID, _ := claimsMap["sub"].(string)
		role, _ := claimsMap["role"].(string)

		if role != string(models.UserRoleAdmin) && role != string(models.UserRoleOwner) {
			if h.serverService == nil {
				return utils.Error(c, http.StatusInternalServerError, "server service not configured")
			}
			server, err := h.serverService.GetServer(c.Request().Context(), serverID)
			if err != nil || server == nil {
				return utils.Error(c, http.StatusNotFound, "server not found")
			}
			if server.UserID != userID {
				return utils.Error(c, http.StatusForbidden, "insufficient permissions")
			}
		}
	}

	ws, err := h.upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		slog.Error("failed to upgrade server metrics ws", "err", err)
		return err
	}

	engine.GlobalUIMetricsHub.AddClient(serverID, ws)
	defer engine.GlobalUIMetricsHub.RemoveClient(serverID, ws)
	defer ws.Close()

	// Keep alive / Wait for close
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}
