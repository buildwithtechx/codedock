package deployments

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

var terminalUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"bearer"},
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
}

type TerminalHandler struct {
	Targets        TerminalDockerTargets
	dockerClient   *client.Client
	tokenValidator TokenValidator
	appLookup      AppServiceLookup
	projectLookup  ProjectLookup
	userLookup     UserStatusLookup
	normalizeName  func(id string) string
}

func NewTerminalHandler(
	dockerClient *client.Client,
	tokenValidator TokenValidator,
	appLookup AppServiceLookup,
	projectLookup ProjectLookup,
	userLookup UserStatusLookup,
) *TerminalHandler {
	return &TerminalHandler{
		dockerClient:   dockerClient,
		tokenValidator: tokenValidator,
		appLookup:      appLookup,
		projectLookup:  projectLookup,
		userLookup:     userLookup,
		normalizeName:  utils.NormalizeContainerName,
	}
}

func (h *TerminalHandler) HandleWebSocket(c echo.Context) error {
	if err := utils.ValidateWebSocketCSWSH(c); err != nil {
		return err
	}

	var userID, role string
	if h.tokenValidator != nil {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			return utils.Error(c, http.StatusUnauthorized, "missing authentication token for terminal access")
		}
		claims, err := h.tokenValidator.ValidateToken(tokenStr)
		if err != nil {
			return utils.Error(c, http.StatusUnauthorized, "invalid authentication token for terminal access")
		}
		sub, ok := claims["sub"].(string)
		if !ok || strings.TrimSpace(sub) == "" {
			sub, ok = claims["Subject"].(string)
		}
		if !ok || strings.TrimSpace(sub) == "" {
			return utils.Error(c, http.StatusUnauthorized, "invalid or missing subject claim in token")
		}
		userID = strings.TrimSpace(sub)
		role, _ = claims["role"].(string)

		if h.userLookup != nil {
			active, err := h.userLookup.IsUserActive(c.Request().Context(), userID)
			if err != nil || !active {
				return utils.Error(c, http.StatusUnauthorized, ErrUserDeactivated.Error())
			}
		}
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}

	containerName := h.normalizeName(id)
	var targetSvc *AppServiceSummary
	if h.appLookup != nil {
		if svc, err := h.appLookup.GetAppService(c.Request().Context(), id); err == nil && svc != nil {
			targetSvc = svc
			if h.projectLookup != nil && userID != "" && role != "admin" {
				if !h.projectLookup.HasPermission(c.Request().Context(), svc.ProjectID, userID, role, "") {
					return utils.Error(c, http.StatusForbidden, ErrTerminalUnauthorized.Error())
				}
			}
			if svc.ContainerID != "" && svc.ContainerID != "-" {
				containerName = svc.ContainerID
			}
		}
	}

	target := h.dockerClient
	if h.Targets != nil {
		resolved, release, err := h.Targets.DockerForService(c.Request().Context(), id)
		if err != nil {
			return utils.Error(c, http.StatusConflict, err.Error())
		}
		defer release()
		target = resolved
	}
	if target == nil {
		return utils.Error(c, http.StatusInternalServerError, "docker client unavailable")
	}

	owned, inspectErr := target.ContainerInspect(c.Request().Context(), containerName)
	if inspectErr != nil || owned.Config == nil || owned.Config.Labels["codedock.service_id"] != id {
		return utils.Error(c, http.StatusConflict, "terminal container ownership changed")
	}

	execConfig := container.ExecOptions{
		Cmd:          []string{"/bin/sh"},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	resp, err := target.ContainerExecCreate(c.Request().Context(), containerName, execConfig)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create exec instance: "+err.Error())
	}
	hijackedResp, err := target.ContainerExecAttach(c.Request().Context(), resp.ID, container.ExecAttachOptions{Tty: true})
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

	if targetSvc != nil && userID != "" {
		stopWatcher := StartTerminalAuthorityWatcher(c.Request().Context(), 5*time.Second, func() error {
			return VerifyTerminalAuthority(c.Request().Context(), h.projectLookup, h.userLookup, targetSvc.ProjectID, userID, role)
		}, func() {
			_ = ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "terminal authority revoked"), time.Now().Add(time.Second))
			_ = ws.Close()
			hijackedResp.Close()
		})
		defer stopWatcher()
	}

	errChan := make(chan error, 2)
	go func() {
		wsReader := wsToReader(ws)
		_, err := io.Copy(hijackedResp.Conn, wsReader)
		errChan <- err
	}()
	go func() {
		wsWriter := wsToWriter(ws)
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

func extractToken(c echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	if auth != "" && strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	if cookie, err := c.Cookie("token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return c.QueryParam("token")
}

func wsToReader(ws *websocket.Conn) io.Reader {
	r, w := io.Pipe()
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				w.CloseWithError(err)
				return
			}
			_, _ = w.Write(message)
		}
	}()
	return r
}

func wsToWriter(ws *websocket.Conn) io.Writer {
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
