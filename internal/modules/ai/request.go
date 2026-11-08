package ai

type UpdateAISettingsRequest struct {
	AISettings
}

type AIDiagnoseRequest struct {
	Prompt string `json:"prompt"`
}

