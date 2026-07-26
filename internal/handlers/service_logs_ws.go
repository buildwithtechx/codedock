package handlers

import (
	"log/slog"
	"net/http"
	"net/url"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ServiceLogsWSHandler struct {
	upgrader     websocket.Upgrader
	tokenService *services.TokenService
	appService   *services.AppService
}

func NewServiceLogsWSHandler(ts *services.TokenService, as *services.AppService) *ServiceLogsWSHandler {
	return &ServiceLogsWSHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				u, err := url.Parse(origin)
				if err != nil {
					return false
				}
				return u.Host == r.Host
			},
		},
		tokenService: ts,
		appService:   as,
	}
}

func (h *ServiceLogsWSHandler) Handle(c echo.Context) error {
	if err := ValidateWebSocketCSWSH(c); err != nil {
		return err
	}
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	if h.tokenService != nil {
		tokenStr := middleware.ExtractTokenFromRequest(c)
		if tokenStr == "" {
			return utils.Error(c, http.StatusUnauthorized, "missing authentication token")
		}

		_, err := h.tokenService.ValidateToken(tokenStr)
		if err != nil {
			return utils.Error(c, http.StatusUnauthorized, "invalid authentication token")
		}
	}

	ws, err := h.upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		slog.Error("failed to upgrade service logs ws", "err", err)
		return err
	}

	engine.GlobalUILogStreamHub.AddClient(serviceID, ws)
	defer engine.GlobalUILogStreamHub.RemoveClient(serviceID, ws)
	defer ws.Close()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}
