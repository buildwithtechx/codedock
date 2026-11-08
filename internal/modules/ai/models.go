package ai

import "codedock/internal/models"

type AISettings = models.AISettings

type AIExplanationResponse struct {
	Summary         string   `json:"summary"`
	Cause           string   `json:"cause"`
	SuggestedFix    string   `json:"suggestedFix"`
	Confidence      string   `json:"confidence"`
	Commands        []string `json:"commands"`
	RelatedLogLines []string `json:"relatedLogLines"`
}
