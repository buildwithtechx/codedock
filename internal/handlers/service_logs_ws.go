package handlers

import (
	"log/slog"
	"net/http"
	"net/url"

	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ServiceLogsWSHandler struct {
	upgrader       websocket.Upgrader
	tokenService   *services.TokenService
	appService     *services.AppService
	projectService *services.ProjectService
}

func NewServiceLogsWSHandler(ts *services.TokenService, as *services.AppService, ps *services.ProjectService) *ServiceLogsWSHandler {
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
		tokenService:   ts,
		appService:     as,
		projectService: ps,
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
