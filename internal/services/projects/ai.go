package projects

import (
	"context"

	"codedock/internal/models"
	"codedock/internal/modules/ai"
	"codedock/internal/repositories"
)

type AISettingsService = ai.Service

type aiSettingsRepoAdapter struct {
	repo repositories.AISettingsRepository
}

func (a *aiSettingsRepoAdapter) GetAISettings(ctx context.Context) (*ai.AISettings, error) {
	if a.repo == nil {
		return nil, nil
	}
	s, err := a.repo.GetAISettings(ctx)
	if err != nil || s == nil {
		return nil, err
	}
	return &ai.AISettings{
		ID:              s.ID,
		DefaultProvider: s.DefaultProvider,
		OpenAIKey:       s.OpenAIKey,
		OpenAIModel:     s.OpenAIModel,
		AnthropicKey:    s.AnthropicKey,
		AnthropicModel:  s.AnthropicModel,
		GoogleKey:       s.GoogleKey,
		GoogleModel:     s.GoogleModel,
		MistralKey:      s.MistralKey,
		MistralModel:    s.MistralModel,
		GroqKey:         s.GroqKey,
		GroqModel:       s.GroqModel,
		DeepSeekKey:     s.DeepSeekKey,
		DeepSeekModel:   s.DeepSeekModel,
		XAIKey:          s.XAIKey,
		XAIModel:        s.XAIModel,
		MoonshotKey:     s.MoonshotKey,
		MoonshotModel:   s.MoonshotModel,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}, nil
}

func (a *aiSettingsRepoAdapter) UpdateAISettings(ctx context.Context, cfg *ai.AISettings) error {
	if a.repo == nil {
		return nil
	}
	return a.repo.UpdateAISettings(ctx, &models.AISettings{
		ID:              cfg.ID,
		DefaultProvider: cfg.DefaultProvider,
		OpenAIKey:       cfg.OpenAIKey,
		OpenAIModel:     cfg.OpenAIModel,
		AnthropicKey:    cfg.AnthropicKey,
		AnthropicModel:  cfg.AnthropicModel,
		GoogleKey:       cfg.GoogleKey,
		GoogleModel:     cfg.GoogleModel,
		MistralKey:      cfg.MistralKey,
		MistralModel:    cfg.MistralModel,
		GroqKey:         cfg.GroqKey,
		GroqModel:       cfg.GroqModel,
		DeepSeekKey:     cfg.DeepSeekKey,
		DeepSeekModel:   cfg.DeepSeekModel,
		XAIKey:          cfg.XAIKey,
		XAIModel:        cfg.XAIModel,
		MoonshotKey:     cfg.MoonshotKey,
		MoonshotModel:   cfg.MoonshotModel,
		CreatedAt:       cfg.CreatedAt,
		UpdatedAt:       cfg.UpdatedAt,
	})
}

func NewAISettingsService(repo repositories.AISettingsRepository) *AISettingsService {
	return ai.NewService(&aiSettingsRepoAdapter{repo: repo})
}
