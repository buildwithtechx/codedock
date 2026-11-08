package ai

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetRawSettings(ctx context.Context) (*AISettings, error) {
	if h.service == nil {
		return nil, ErrFailedLoadSettings
	}
	return h.service.GetAISettings(ctx)
}

func (h *Handler) GetAISettings(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "AI service not initialized")
	}
	s, err := h.service.GetAISettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", RedactAISettings(s))
}

func (h *Handler) UpdateAISettings(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "AI service not initialized")
	}
	var req UpdateAISettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	existing, err := h.service.GetAISettings(c.Request().Context())
	if err == nil && existing != nil {
		if req.OpenAIKey == "********" {
			req.OpenAIKey = existing.OpenAIKey
		}
		if req.AnthropicKey == "********" {
			req.AnthropicKey = existing.AnthropicKey
		}
		if req.GoogleKey == "********" {
			req.GoogleKey = existing.GoogleKey
		}
		if req.MistralKey == "********" {
			req.MistralKey = existing.MistralKey
		}
		if req.GroqKey == "********" {
			req.GroqKey = existing.GroqKey
		}
		if req.DeepSeekKey == "********" {
			req.DeepSeekKey = existing.DeepSeekKey
		}
		if req.XAIKey == "********" {
			req.XAIKey = existing.XAIKey
		}
		if req.MoonshotKey == "********" {
			req.MoonshotKey = existing.MoonshotKey
		}
	}

	if err := h.service.UpdateAISettings(c.Request().Context(), &req.AISettings); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	updated, err := h.service.GetAISettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "AI settings updated successfully", RedactAISettings(updated))
}

func (h *Handler) DiagnoseLogs(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusInternalServerError, "AI service not initialized")
	}
	var req AIDiagnoseRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidRequest.Error())
	}

	if req.Prompt == "" {
		return utils.Error(c, http.StatusBadRequest, ErrPromptRequired.Error())
	}

	diagnosis, err := h.service.DiagnoseLogs(c.Request().Context(), req.Prompt)
	if err != nil {
		if errors.Is(err, ErrFailedLoadSettings) {
			return utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, diagnosis)
}
