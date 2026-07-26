package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type AIDiagnoseRequest struct {
	Prompt string `json:"prompt"`
}

type AISettingsHandler struct {
	aiSettingsService *services.AISettingsService
}

func NewAISettingsHandler(s *services.AISettingsService) *AISettingsHandler {
	return &AISettingsHandler{aiSettingsService: s}
}

func redactAISettings(s *models.AISettings) *models.AISettings {
	if s == nil {
		return nil
	}
	if s.OpenAIKey != "" {
		s.OpenAIKey = "********"
	}
	if s.AnthropicKey != "" {
		s.AnthropicKey = "********"
	}
	if s.GoogleKey != "" {
		s.GoogleKey = "********"
	}
	if s.MistralKey != "" {
		s.MistralKey = "********"
	}
	if s.GroqKey != "" {
		s.GroqKey = "********"
	}
	if s.DeepSeekKey != "" {
		s.DeepSeekKey = "********"
	}
	if s.XAIKey != "" {
		s.XAIKey = "********"
	}
	if s.MoonshotKey != "" {
		s.MoonshotKey = "********"
	}
	return s
}

func (h *AISettingsHandler) GetAISettings(c echo.Context) error {
	s, err := h.aiSettingsService.GetAISettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", redactAISettings(s))
}

func (h *AISettingsHandler) UpdateAISettings(c echo.Context) error {
	var req models.UpdateAISettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	existing, err := h.aiSettingsService.GetAISettings(c.Request().Context())
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

	if err := h.aiSettingsService.UpdateAISettings(c.Request().Context(), &req.AISettings); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	updated, err := h.aiSettingsService.GetAISettings(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "AI settings updated successfully", redactAISettings(updated))
}

func (h *AISettingsHandler) DiagnoseLogs(c echo.Context) error {
	var req AIDiagnoseRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if req.Prompt == "" {
		return c.String(http.StatusBadRequest, "Prompt is required")
	}

	diagnosis, err := h.aiSettingsService.DiagnoseLogs(c.Request().Context(), req.Prompt)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, diagnosis)
}
