package handlers

import (
	"log/slog"
	"net/http"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ServerMetricsWSHandler struct {
	upgrader websocket.Upgrader
}

func NewServerMetricsWSHandler() *ServerMetricsWSHandler {
	return &ServerMetricsWSHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *ServerMetricsWSHandler) Handle(c echo.Context) error {
	serverID := c.Param("serverId")
	if serverID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serverId parameter")
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
