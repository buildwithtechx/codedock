package projects

import (
	"codedock/internal/modules/ai"
	"codedock/internal/repositories"
)

type AIAnalysisService = ai.AnalysisService
type AIExplanationResponse = ai.AIExplanationResponse

func NewAIAnalysisService(d repositories.DeploymentRepository, a repositories.AppServiceRepository, aiRepo repositories.AISettingsRepository) *AIAnalysisService {
	return ai.NewAnalysisService(d, a, aiRepo)
}
