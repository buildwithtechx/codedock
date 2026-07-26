package system

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WorkerWSHandler struct {
	hub        *engine.WorkerHub
	serverRepo repositories.ServerRepository
	upgrader   websocket.Upgrader
}

func NewWorkerWSHandler(hub *engine.WorkerHub, serverRepo repositories.ServerRepository) *WorkerWSHandler {
	return &WorkerWSHandler{
		hub:        hub,
		serverRepo: serverRepo,
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
	}
}

func (h *WorkerWSHandler) Connect(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}

	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
	}

	server, err := h.serverRepo.GetByToken(c.Request().Context(), token)
	if err != nil || server == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	conn := h.hub.Register(server.ID, ws)

	authAck := models.WorkerAuthResultPayload{Success: true}
	authAckBytes, _ := json.Marshal(authAck)

	_ = conn.Send(&models.WorkerMessage{
		ID:        uuid.New().String(),
		Type:      models.WorkerMessageTypeAuthResult,
		Timestamp: time.Now(),
		Payload:   authAckBytes,
	})

	conn.Listen()

	return nil
}
