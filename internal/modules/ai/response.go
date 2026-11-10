package ai

type OpenAIChatChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChatChoice struct {
	Index        int                     `json:"index"`
	Message      OpenAIChatChoiceMessage `json:"message"`
	FinishReason string                  `json:"finish_reason,omitempty"`
}

type OpenAIChatCompletionResponse struct {
	ID      string             `json:"id,omitempty"`
	Object  string             `json:"object,omitempty"`
	Created int64              `json:"created,omitempty"`
	Model   string             `json:"model,omitempty"`
	Choices []OpenAIChatChoice `json:"choices"`
}

type AIExplanationResponse struct {
	Summary         string   `json:"summary"`
	Cause           string   `json:"cause"`
	SuggestedFix    string   `json:"suggestedFix"`
	Confidence      string   `json:"confidence"`
	Commands        []string `json:"commands"`
	RelatedLogLines []string `json:"relatedLogLines"`
}

type ChatResponse struct {
	Reply           string          `json:"reply"`
	Provider        string          `json:"provider"`
	Model           string          `json:"model"`
	Configured      bool            `json:"configured"`
	Steps           []ToolStep      `json:"steps,omitempty"`
	PendingActions  []PendingAction `json:"pendingActions,omitempty"`
	ExecutedActions []string        `json:"executedActions,omitempty"`
	Questions       []QuestionGroup `json:"questions,omitempty"`
}

// RedactAISettings masks confidential provider API keys before returning settings over the network
func RedactAISettings(s *AISettings) *AISettings {
	if s == nil {
		return nil
	}
	cp := *s
	if cp.OpenAIKey != "" {
		cp.OpenAIKey = "********"
	}
	if cp.AnthropicKey != "" {
		cp.AnthropicKey = "********"
	}
	if cp.GoogleKey != "" {
		cp.GoogleKey = "********"
	}
	if cp.MistralKey != "" {
		cp.MistralKey = "********"
	}
	if cp.GroqKey != "" {
		cp.GroqKey = "********"
	}
	if cp.DeepSeekKey != "" {
		cp.DeepSeekKey = "********"
	}
	if cp.XAIKey != "" {
		cp.XAIKey = "********"
	}
	if cp.MoonshotKey != "" {
		cp.MoonshotKey = "********"
	}
	return &cp
}
