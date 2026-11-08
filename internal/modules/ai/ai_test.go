package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAIRepo struct {
	settings *AISettings
	getErr   error
	updErr   error
}

func (m *mockAIRepo) GetAISettings(ctx context.Context) (*AISettings, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.settings == nil {
		return &AISettings{DefaultProvider: ProviderOpenAI}, nil
	}
	cp := *m.settings
	return &cp, nil
}

func (m *mockAIRepo) UpdateAISettings(ctx context.Context, cfg *AISettings) error {
	if m.updErr != nil {
		return m.updErr
	}
	m.settings = cfg
	return nil
}

func TestRedactAISettings(t *testing.T) {
	assert.Nil(t, RedactAISettings(nil))

	s := &AISettings{
		OpenAIKey:    "sk-openai-12345",
		AnthropicKey: "sk-ant-12345",
		GoogleKey:    "ai-google-12345",
		MistralKey:   "mis-12345",
		GroqKey:      "gsk-12345",
		DeepSeekKey:  "dsk-12345",
		XAIKey:       "xai-12345",
		MoonshotKey:  "moon-12345",
	}

	redacted := RedactAISettings(s)
	require.NotNil(t, redacted)
	assert.Equal(t, "********", redacted.OpenAIKey)
	assert.Equal(t, "********", redacted.AnthropicKey)
	assert.Equal(t, "********", redacted.GoogleKey)
	assert.Equal(t, "********", redacted.MistralKey)
	assert.Equal(t, "********", redacted.GroqKey)
	assert.Equal(t, "********", redacted.DeepSeekKey)
	assert.Equal(t, "********", redacted.XAIKey)
	assert.Equal(t, "********", redacted.MoonshotKey)

	// Original struct should not be mutated
	assert.Equal(t, "sk-openai-12345", s.OpenAIKey)
}

func TestAIService_GetAndUpdate(t *testing.T) {
	repo := &mockAIRepo{
		settings: &AISettings{
			DefaultProvider: ProviderOpenAI,
			OpenAIKey:       "sk-secret",
			OpenAIModel:     "gpt-4o",
		},
	}
	svc := NewService(repo)

	cfg, err := svc.GetAISettings(context.Background())
	require.NoError(t, err)
	assert.Equal(t, ProviderOpenAI, cfg.DefaultProvider)
	assert.Equal(t, "sk-secret", cfg.OpenAIKey)

	cfg.OpenAIModel = "gpt-4.5"
	err = svc.UpdateAISettings(context.Background(), cfg)
	require.NoError(t, err)

	updated, err := svc.GetAISettings(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "gpt-4.5", updated.OpenAIModel)
}

