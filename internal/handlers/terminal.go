package handlers

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

var terminalUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
}

type TerminalHandler struct {
	dockerClient   *client.Client
	tokenService   *services.TokenService
	appService     *services.AppService
	projectService *services.ProjectService
	normalizeName  func(id string) string
}

func NewTerminalHandler(
	dockerClient *client.Client,
	tokenService *services.TokenService,
	appService *services.AppService,
	projectService *services.ProjectService,
) *TerminalHandler {
	return &TerminalHandler{
		dockerClient:   dockerClient,
		tokenService:   tokenService,
		appService:     appService,
		projectService: projectService,
		normalizeName:  utils.NormalizeContainerName,
	}
}

func (h *TerminalHandler) HandleWebSocket(c echo.Context) error {
	if err := ValidateWebSocketCSWSH(c); err != nil {
		return err
	}
	var claimsMap map[string]interface{}
	if h.tokenService != nil {
		tokenStr := middleware.ExtractTokenFromRequest(c)
		if tokenStr == "" {
			return utils.Error(c, http.StatusUnauthorized, "missing authentication token for terminal access")
		}
		cm, err := h.tokenService.ValidateToken(tokenStr)
		if err != nil {
			return utils.Error(c, http.StatusUnauthorized, "invalid authentication token for terminal access")
		}
		claimsMap = cm
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	containerName := h.normalizeName(id)
	if h.appService != nil {
		if svc, err := h.appService.GetAppService(c.Request().Context(), id); err == nil && svc != nil {
			if h.projectService != nil && claimsMap != nil {
				userID, _ := claimsMap["sub"].(string)
				role, _ := claimsMap["role"].(string)
				if role != "admin" {
					if !h.projectService.HasPermission(c.Request().Context(), svc.ProjectID, userID, models.UserRole(role), "") {
						return utils.Error(c, http.StatusForbidden, "insufficient permissions to access this terminal")
					}
				}
			}

			if svc.ContainerID != "" && svc.ContainerID != "-" {
				containerName = svc.ContainerID
			} else {
				containerName = h.normalizeName(svc.ID)
			}
		}
	}
	execConfig := container.ExecOptions{
		Cmd:          []string{"/bin/sh"},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}
	if h.dockerClient == nil {
		return utils.Error(c, http.StatusInternalServerError, "docker client unavailable")
	}
	resp, err := h.dockerClient.ContainerExecCreate(context.Background(), containerName, execConfig)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create exec instance: "+err.Error())
	}
	hijackedResp, err := h.dockerClient.ContainerExecAttach(context.Background(), resp.ID, container.ExecAttachOptions{Tty: true})
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to attach to exec instance: "+err.Error())
	}
	defer hijackedResp.Close()
	responseHeader := http.Header{}
	if reqProto := c.Request().Header.Get("Sec-WebSocket-Protocol"); reqProto != "" {
		parts := strings.Split(reqProto, ",")
		if len(parts) > 0 {
			responseHeader.Set("Sec-WebSocket-Protocol", strings.TrimSpace(parts[0]))
		}
	}
	ws, err := terminalUpgrader.Upgrade(c.Response().Writer, c.Request(), responseHeader)
	if err != nil {
		return err
	}
	defer ws.Close()
	errChan := make(chan error, 2)
	go func() {
		wsReader := h.wsToReader(ws)
		_, err := io.Copy(hijackedResp.Conn, wsReader)
		errChan <- err
	}()
	go func() {
		wsWriter := h.wsToWriter(ws)
		_, err := io.Copy(wsWriter, hijackedResp.Reader)
		errChan <- err
	}()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		}
	}()
	<-errChan
	return nil
}

func (h *TerminalHandler) wsToReader(ws *websocket.Conn) io.Reader {
	r, w := io.Pipe()
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				w.CloseWithError(err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}
		}
	}()
	return r
}

func (h *TerminalHandler) wsToWriter(ws *websocket.Conn) io.Writer {
	r, w := io.Pipe()
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	return w
}
