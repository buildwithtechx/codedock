package projects

import (
	handlerutils "codedock/internal/handlers/utils"
	"codedock/internal/models"
	"codedock/internal/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var runtimeStreamUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		return handlerutils.IsAllowedWebSocketOrigin(r, origin)
	},
}

func (h *AppHandler) InstanceGraph(c echo.Context) error {
	result, err := h.Runtime.InstanceGraph(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Instance graph", result)
}

func (h *AppHandler) ResourceSummary(c echo.Context) error {
	result, err := h.Runtime.ResourceSummary(c.Request().Context(), c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Resource controls", result)
}

type wsStreamWriter struct {
	ws      *websocket.Conn
	timeout time.Duration
}

func (w wsStreamWriter) Write(data []byte) (int, error) {
	if err := w.ws.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
		return 0, err
	}
	if err := w.ws.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (h *AppHandler) RuntimeLogStream(c echo.Context) error {
	ws, err := runtimeStreamUpgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
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
	pod := c.QueryParam("pod")
	if err := h.Runtime.StreamLogs(ctx, c.Param("id"), pod, wsStreamWriter{ws: ws, timeout: 15 * time.Second}); err != nil && ctx.Err() == nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("Log stream failed: "+err.Error()))
	}
	return nil
}

func (h *AppHandler) RuntimeExecTerminal(c echo.Context) error {
	ws, err := runtimeStreamUpgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	_, open, err := ws.ReadMessage()
	if err != nil {
		return nil
	}
	var request models.RuntimeExecRequest
	if err := json.Unmarshal(open, &request); err != nil || len(request.Command) == 0 {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("First message must be a JSON command"))
		return nil
	}
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()
	stdin, sender := io.Pipe()
	defer sender.Close()
	go func() {
		defer cancel()
		for {
			kind, message, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if kind != websocket.BinaryMessage && kind != websocket.TextMessage {
				continue
			}
			if _, err := sender.Write(message); err != nil {
				return
			}
		}
	}()
	if err := h.Runtime.StreamExec(ctx, c.Param("id"), request.Pod, request.Command, stdin, wsStreamWriter{ws: ws, timeout: 15 * time.Second}); err != nil && ctx.Err() == nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("Terminal failed: "+err.Error()))
	}
	return nil
}
