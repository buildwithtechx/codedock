package system

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	handlerutils "codedock.run/codedock/internal/handlers/utils"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"io"
)

type tokenValidator interface {
	ValidateToken(token string) (jwt.MapClaims, error)
}

type userStatusProvider interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type ServiceLogStreams interface {
	StreamServiceLogs(context.Context, string, io.Writer) error
}
type ServiceLogsWSHandler struct {
	Streams        ServiceLogStreams
	upgrader       websocket.Upgrader
	tokenService   tokenValidator
	appService     *projectservices.AppService
	projectService *projectservices.ProjectService
	userRepo       userStatusProvider
}

func NewServiceLogsWSHandler(ts tokenValidator, as *projectservices.AppService, ps *projectservices.ProjectService, ur userStatusProvider) *ServiceLogsWSHandler {
	return &ServiceLogsWSHandler{
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
		tokenService:   ts,
		appService:     as,
		projectService: ps,
		userRepo:       ur,
	}
}

func (h *ServiceLogsWSHandler) Handle(c echo.Context) error {
	if err := handlerutils.ValidateWebSocketCSWSH(c); err != nil {
		return err
	}
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing serviceId parameter")
	}

	var claimsMap map[string]interface{}
	if h.tokenService != nil {
		tokenStr := middleware.ExtractTokenFromRequest(c)
		if tokenStr == "" {
			return utils.Error(c, http.StatusUnauthorized, "missing authentication token")
		}

		cm, err := h.tokenService.ValidateToken(tokenStr)
		if err != nil {
			return utils.Error(c, http.StatusUnauthorized, "invalid authentication token")
		}
		claimsMap = cm

		userID, _ := claimsMap["sub"].(string)
		if userID != "" && h.userRepo != nil {
			u, err := h.userRepo.GetUserByID(c.Request().Context(), userID)
			if err != nil || u == nil || !u.IsActive {
				return utils.Error(c, http.StatusUnauthorized, "user account not found or deactivated")
			}
		}
	}

	if h.appService != nil {
		svc, err := h.appService.GetAppService(c.Request().Context(), serviceID)
		if err != nil || svc == nil {
			return utils.Error(c, http.StatusNotFound, "service not found")
		}
		if h.projectService != nil && claimsMap != nil {
			userID, _ := claimsMap["sub"].(string)
			role, _ := claimsMap["role"].(string)
			if role != "admin" {
				if !h.projectService.HasPermission(c.Request().Context(), svc.ProjectID, userID, models.UserRole(role), "") {
					return utils.Error(c, http.StatusForbidden, "insufficient permissions to access service logs")
				}
			}
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
	if err := h.Streams.StreamServiceLogs(ctx, serviceID, serviceLogWriter{ws}); err != nil && ctx.Err() == nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("Log stream failed: "+err.Error()))
	}
	return nil
}

type serviceLogWriter struct{ ws *websocket.Conn }

func (w serviceLogWriter) Write(data []byte) (int, error) {
	if err := w.ws.SetWriteDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return 0, err
	}
	if err := w.ws.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}
	return len(data), nil
}
