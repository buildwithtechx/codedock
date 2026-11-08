package ai

type AISettings struct {
	ID              string `json:"id"`
	DefaultProvider string `json:"defaultProvider"`
	OpenAIKey       string `json:"openAIKey,omitempty"`
	OpenAIModel     string `json:"openAIModel,omitempty"`
	AnthropicKey    string `json:"anthropicKey,omitempty"`
	AnthropicModel  string `json:"anthropicModel,omitempty"`
	GoogleKey       string `json:"googleKey,omitempty"`
	GoogleModel     string `json:"googleModel,omitempty"`
	MistralKey      string `json:"mistralKey,omitempty"`
	MistralModel    string `json:"mistralModel,omitempty"`
	GroqKey         string `json:"groqKey,omitempty"`
	GroqModel       string `json:"groqModel,omitempty"`
	DeepSeekKey     string `json:"deepSeekKey,omitempty"`
	DeepSeekModel   string `json:"deepSeekModel,omitempty"`
	XAIKey          string `json:"xaiKey,omitempty"`
	XAIModel        string `json:"xaiModel,omitempty"`
	MoonshotKey     string `json:"moonshotKey,omitempty"`
	MoonshotModel   string `json:"moonshotModel,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type AIExplanationResponse struct {
	Summary         string   `json:"summary"`
	Cause           string   `json:"cause"`
	SuggestedFix    string   `json:"suggestedFix"`
	Confidence      string   `json:"confidence"`
	Commands        []string `json:"commands"`
	RelatedLogLines []string `json:"relatedLogLines"`
}
