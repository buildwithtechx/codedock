package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetAISettings(t *testing.T) {
	repo := &mockAIRepo{
		settings: &AISettings{
			DefaultProvider: ProviderOpenAI,
			OpenAIKey:       "sk-secret-1234",
			OpenAIModel:     "gpt-4o",
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ai", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAISettings(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data AISettings `json:"data"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "********", resp.Data.OpenAIKey)
	assert.Equal(t, "gpt-4o", resp.Data.OpenAIModel)
}

func TestHandler_UpdateAISettings(t *testing.T) {
	repo := &mockAIRepo{
		settings: &AISettings{
			DefaultProvider: ProviderOpenAI,
			OpenAIKey:       "original-key",
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	e := echo.New()
	updateBody := UpdateAISettingsRequest{
		AISettings: AISettings{
			DefaultProvider: ProviderAnthropic,
			OpenAIKey:       "********", // should retain original key
			AnthropicKey:    "new-anthropic-key",
		},
	}
	payload, _ := json.Marshal(updateBody)
	req := httptest.NewRequest(http.MethodPut, "/ai", bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.UpdateAISettings(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	raw, err := repo.GetAISettings(c.Request().Context())
	require.NoError(t, err)
	assert.Equal(t, "original-key", raw.OpenAIKey)
	assert.Equal(t, "new-anthropic-key", raw.AnthropicKey)
}

func TestHandler_DiagnoseLogs_PromptRequired(t *testing.T) {
	repo := &mockAIRepo{}
	svc := NewService(repo)
	handler := NewHandler(svc)

	e := echo.New()
	payload := []byte(`{"prompt":""}`)
	req := httptest.NewRequest(http.MethodPost, "/ai/diagnose", bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.DiagnoseLogs(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterRoutes_NilSafety(t *testing.T) {
	e := echo.New()
	authGroup := e.Group("/api")
	apiGroup := e.Group("/api/v1")

	// Must not panic with nil handler
	assert.NotPanics(t, func() {
		RegisterRoutes(authGroup, apiGroup, nil, nil)
	})
}

