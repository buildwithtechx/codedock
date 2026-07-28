package system

import (
	"log/slog"
	"net/http"

	handlerutils "codedock.run/codedock/internal/handlers/utils"

	"codedock.run/codedock/internal/engine/observability"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ServerMetricsWSHandler struct {
	upgrader      websocket.Upgrader
	tokenService  *authservices.TokenService
	serverService systemservices.ServerService
	userRepo      userStatusProvider
}

func NewServerMetricsWSHandler(ts *authservices.TokenService, ss systemservices.ServerService, ur userStatusProvider) *ServerMetricsWSHandler {
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
				return handlerutils.IsAllowedWebSocketOrigin(r, origin)
			},
		},
		tokenService:  ts,
		serverService: ss,
		userRepo:      ur,
	}
}

func (h *ServerMetricsWSHandler) Handle(c echo.Context) error {
	if err := handlerutils.ValidateWebSocketCSWSH(c); err != nil {
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
		if userID != "" && h.userRepo != nil {
			u, err := h.userRepo.GetUserByID(c.Request().Context(), userID)
			if err != nil || u == nil || !u.IsActive {
				return utils.Error(c, http.StatusUnauthorized, "user account not found or deactivated")
			}
		}

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

	observability.GlobalUIMetricsHub.AddClient(serverID, ws)
	defer observability.GlobalUIMetricsHub.RemoveClient(serverID, ws)
	defer ws.Close()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}
