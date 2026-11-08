package projects

import (
	"codedock/internal/modules/ai"
	"codedock/internal/repositories"
)

type AISettingsService = ai.Service

func NewAISettingsService(repo repositories.AISettingsRepository) *AISettingsService {
	return ai.NewService(repo)
}
