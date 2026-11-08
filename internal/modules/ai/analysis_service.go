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

type DeploymentContext struct {
	ServiceName    string
	RuntimeMode    string
	RepositoryURL  string
	Branch         string
	RootDirectory  string
	InstallCommand string
	BuildCommand   string
	StartCommand   string
	Status         string
	BuildLogs      string
}

type AnalysisService struct {
	aiRepo     Repository
	httpClient *http.Client
}

func NewAnalysisService(aiRepo Repository) *AnalysisService {
	return &AnalysisService{
		aiRepo:     aiRepo,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *AnalysisService) ExplainDeploymentFailure(ctx context.Context, d *DeploymentContext) (*AIExplanationResponse, error) {
	if d == nil {
		return nil, errors.New("deployment context is required")
	}
	if d.Status != "failed" {
		return nil, ErrDeploymentNotFailed
	}

	aiSettings, err := s.aiRepo.GetAISettings(ctx)
	if err != nil || aiSettings == nil {
		return nil, ErrFailedLoadSettings
	}

	var baseURL, apiKey, modelName string
	provider := aiSettings.DefaultProvider
	if provider == "" {
		provider = ProviderOpenAI
	}

	switch provider {
	case ProviderAnthropic:
		return nil, errors.New("Anthropic is not fully supported yet in this AI pipeline without a custom proxy")
	case ProviderGoogle:
		return nil, errors.New("Google is not fully supported yet in this AI pipeline without a custom proxy")
	case ProviderGroq:
		baseURL = "https://api.groq.com/openai/v1"
		apiKey = aiSettings.GroqKey
		modelName = aiSettings.GroqModel
		if modelName == "" {
			modelName = "llama-3.3-70b-versatile"
		}
	case ProviderMistral:
		baseURL = "https://api.mistral.ai/v1"
		apiKey = aiSettings.MistralKey
		modelName = aiSettings.MistralModel
		if modelName == "" {
			modelName = "mistral-large-latest"
		}
	case ProviderDeepSeek:
		baseURL = "https://api.deepseek.com/v1"
		apiKey = aiSettings.DeepSeekKey
		modelName = aiSettings.DeepSeekModel
		if modelName == "" {
			modelName = "deepseek-chat"
		}
	case ProviderXAI:
		baseURL = "https://api.x.ai/v1"
		apiKey = aiSettings.XAIKey
		modelName = aiSettings.XAIModel
		if modelName == "" {
			modelName = "grok-4.7"
		}
	case ProviderMoonshot:
		baseURL = "https://api.moonshot.cn/v1"
		apiKey = aiSettings.MoonshotKey
		modelName = aiSettings.MoonshotModel
		if modelName == "" {
			modelName = "kimi-k3"
		}
	case ProviderOpenAI:
		fallthrough
	default:
		baseURL = "https://api.openai.com/v1"
		apiKey = aiSettings.OpenAIKey
		modelName = aiSettings.OpenAIModel
		if modelName == "" {
			modelName = "gpt-6-astra"
		}
	}

	if apiKey == "" {
		return nil, fmt.Errorf("%s API key is missing in AI settings", provider)
	}

	prompt := buildPrompt(d)

	return callOpenAI(ctx, s.httpClient, baseURL, apiKey, modelName, prompt)
}

func buildPrompt(d *DeploymentContext) string {
	maxLogChars := 30000
	logs := d.BuildLogs
	if len(logs) > maxLogChars {
		logs = logs[len(logs)-maxLogChars:]
	}

	return fmt.Sprintf(`You are helping diagnose a failed Codedock deployment.

Explain the failure in plain, specific language for the app owner. Focus on what happened and the smallest likely fix.

Rules:
- Use only the deployment metadata and logs below.
- Do not invent files, commands, ports, packages, or environment variables that are not supported by the logs.
- Prefer concrete fixes over generic advice.
- If evidence is weak, say so and set confidence to low.
- Keep summary, cause, and suggestedFix each under 80 words.
- relatedLogLines should quote the shortest relevant log lines exactly as they appear below.
- commands should contain only commands or setting changes the user can plausibly run or make.

Deployment metadata:
Service: %s
Runtime mode: %s
Repository: %s
Branch: %s
Root directory: %s
Install command: %s
Build command: %s
Start command: %s
Deployment status: %s

Deployment logs:
%s`, d.ServiceName, d.RuntimeMode, d.RepositoryURL, d.Branch, d.RootDirectory, d.InstallCommand, d.BuildCommand, d.StartCommand, d.Status, logs)
}

func callOpenAI(ctx context.Context, client *http.Client, baseURL, apiKey, model, prompt string) (*AIExplanationResponse, error) {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/chat/completions"
	} else {
		baseURL = baseURL + "/chat/completions"
	}

	reqBody := OpenAIChatCompletionRequest{
		Model: model,
		Messages: []OpenAIChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.2,
		ResponseFormat: &OpenAIResponseFormat{
			Type: "json_schema",
			JSONSchema: &OpenAIJSONSchemaSpec{
				Name: "deployment_failure_explanation",
				Schema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"summary":         map[string]any{"type": "string"},
						"cause":           map[string]any{"type": "string"},
						"suggestedFix":    map[string]any{"type": "string"},
						"confidence":      map[string]any{"type": "string", "enum": []string{"low", "medium", "high"}},
						"commands":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
						"relatedLogLines": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					},
					"required":             []string{"summary", "cause", "suggestedFix", "confidence", "commands", "relatedLogLines"},
					"additionalProperties": false,
				},
				Strict: true,
			},
		},
	}

	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result OpenAIChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, errors.New("OpenAI returned no choices")
	}

	var explanation AIExplanationResponse
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &explanation); err != nil {
		return nil, fmt.Errorf("failed to parse explanation JSON: %w", err)
	}

	return &explanation, nil
}
