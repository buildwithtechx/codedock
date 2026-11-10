package ai

import "context"

// Repository defines database persistence for AI provider credentials and model selections.
type Repository interface {
	GetAISettings(ctx context.Context) (*AISettings, error)
	UpdateAISettings(ctx context.Context, cfg *AISettings) error
}

// AIServiceOperations provides high-level LLM interactions and settings management.
type AIServiceOperations interface {
	GetAISettings(ctx context.Context) (*AISettings, error)
	UpdateAISettings(ctx context.Context, cfg *AISettings) error
	DiagnoseLogs(ctx context.Context, logs string) (string, error)
}
