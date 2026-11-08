package ai

type UpdateAISettingsRequest struct {
	AISettings
}

type AIDiagnoseRequest struct {
	Prompt string `json:"prompt"`
}

type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIJSONSchemaSpec struct {
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema,omitempty"`
	Strict bool           `json:"strict"`
}

type OpenAIResponseFormat struct {
	Type       string                `json:"type"`
	JSONSchema *OpenAIJSONSchemaSpec `json:"json_schema,omitempty"`
}

type OpenAIChatCompletionRequest struct {
	Model          string                `json:"model"`
	Messages       []OpenAIChatMessage   `json:"messages"`
	Temperature    float64               `json:"temperature"`
	ResponseFormat *OpenAIResponseFormat `json:"response_format,omitempty"`
}
