package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	repo       Repository
	httpClient *http.Client
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:       repo,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *Service) GetAISettings(ctx context.Context) (*AISettings, error) {
	return s.repo.GetAISettings(ctx)
}

func (s *Service) UpdateAISettings(ctx context.Context, cfg *AISettings) error {
	return s.repo.UpdateAISettings(ctx, cfg)
}

func (s *Service) DiagnoseLogs(ctx context.Context, logs string) (string, error) {
	settings, err := s.repo.GetAISettings(ctx)
	if err != nil || settings == nil {
		return "", ErrFailedLoadSettings
	}

	var baseURL, apiKey, modelName string
	provider := settings.DefaultProvider
	if provider == "" {
		provider = ProviderOpenAI
	}

	switch provider {
	case ProviderAnthropic:
		baseURL = "https://api.anthropic.com/v1"
		apiKey = settings.AnthropicKey
		modelName = settings.AnthropicModel
		if modelName == "" {
			modelName = "claude-sonnet-4-6"
		}
	case ProviderGoogle:
		baseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
		apiKey = settings.GoogleKey
		modelName = settings.GoogleModel
		if modelName == "" {
			modelName = "gemini-2.5-pro"
		}
	case ProviderGroq:
		baseURL = "https://api.groq.com/openai/v1"
		apiKey = settings.GroqKey
		modelName = settings.GroqModel
		if modelName == "" {
			modelName = "llama-3.3-70b-versatile"
		}
	case ProviderMistral:
		baseURL = "https://api.mistral.ai/v1"
		apiKey = settings.MistralKey
		modelName = settings.MistralModel
		if modelName == "" {
			modelName = "mistral-large-latest"
		}
	case ProviderDeepSeek:
		baseURL = "https://api.deepseek.com/v1"
		apiKey = settings.DeepSeekKey
		modelName = settings.DeepSeekModel
		if modelName == "" {
			modelName = "deepseek-chat"
		}
	case ProviderXAI:
		baseURL = "https://api.x.ai/v1"
		apiKey = settings.XAIKey
		modelName = settings.XAIModel
		if modelName == "" {
			modelName = "grok-4.7"
		}
	case ProviderMoonshot:
		baseURL = "https://api.moonshot.cn/v1"
		apiKey = settings.MoonshotKey
		modelName = settings.MoonshotModel
		if modelName == "" {
			modelName = "kimi-k3"
		}
	case ProviderOpenAI:
		fallthrough
	default:
		baseURL = "https://api.openai.com/v1"
		apiKey = settings.OpenAIKey
		modelName = settings.OpenAIModel
		if modelName == "" {
			modelName = "gpt-6-astra"
		}
	}

	if apiKey == "" {
		return "", fmt.Errorf("%s API key is missing in AI settings", provider)
	}

	prompt := fmt.Sprintf("You are a DevOps AI assistant. Diagnose these logs and explain the error in plain English. Provide a fix if possible:\n\n%s", logs)

	url := baseURL + "/chat/completions"
	reqBody := map[string]any{
		"model": modelName,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.2,
	}

	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", errors.New("AI returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}

