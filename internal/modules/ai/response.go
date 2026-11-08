package ai

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
