package system

import (
	"codedock/internal/modules/ai"
)

type AISettingsHandler = ai.Handler

func NewAISettingsHandler(s *ai.Service) *AISettingsHandler {
	return ai.NewHandler(s)
}
