package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func (h *ServerHandler) GetRateLimit(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	res := models.ServerRateLimitConfig{
		RPS:       50,
		Burst:     20,
		Whitelist: []string{},
	}
	return utils.Success(c, "Rate limit retrieved", res)
}

func (h *ServerHandler) UpdateRateLimit(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	var req models.ServerRateLimitConfig
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.RPS < 0 {
		req.RPS = 50
	}
	if req.Burst < 0 {
		req.Burst = 20
	}
	if req.Whitelist == nil {
		req.Whitelist = []string{}
	}

	return utils.Success(c, "Rate limit updated", req)
}
