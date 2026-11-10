package chatapps

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

// ListIdentities returns all external identities paired to this connection
func (h *Handler) ListIdentities(c echo.Context) error {
	id := c.Param("id")
	identities, err := h.service.ListIdentities(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Paired identities retrieved", identities)
}

// DeleteIdentity unlinks an external identity from this Codedock account
func (h *Handler) DeleteIdentity(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteIdentity(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Paired identity removed", nil)
}

// GetChallenge returns metadata for a pairing challenge token
func (h *Handler) GetChallenge(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		token = c.QueryParam("token")
	}
	ch, err := h.service.GetChallenge(c.Request().Context(), token)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Pairing challenge found", ch)
}

// VerifyPairing completes pairing between an external chat identity and the authenticated user
func (h *Handler) VerifyPairing(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req VerifyPairingRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	identity, err := h.service.VerifyPairing(c.Request().Context(), userID, req.Token, req.Code)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Account pairing verified successfully", identity)
}
