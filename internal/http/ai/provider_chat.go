package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"codedock/internal/models"
)

type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type resolvedEndpoint struct {
	Provider string
	BaseURL  string
	APIKey   string
	Model    string
}

func resolveProviderEndpoint(settings *models.AISettings, overrideProvider, overrideModel string) (resolvedEndpoint, error) {
	if settings == nil {
		return resolvedEndpoint{}, errors.New("AI settings are not configured")
	}
	provider := strings.TrimSpace(overrideProvider)
	if provider == "" {
		provider = strings.TrimSpace(settings.DefaultProvider)
	}
	if provider == "" || provider == "none" {
		provider = "openai"
	}

	var ep resolvedEndpoint
	ep.Provider = provider

	switch provider {
	case "anthropic":
		ep.BaseURL = "https://api.anthropic.com/v1"
		ep.APIKey = settings.AnthropicKey
		ep.Model = firstNonEmpty(overrideModel, settings.AnthropicModel, "claude-sonnet-4-6")
	case "google":
		ep.BaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
		ep.APIKey = settings.GoogleKey
		ep.Model = firstNonEmpty(overrideModel, settings.GoogleModel, "gemini-2.5-pro")
	case "groq":
		ep.BaseURL = "https://api.groq.com/openai/v1"
		ep.APIKey = settings.GroqKey
		ep.Model = firstNonEmpty(overrideModel, settings.GroqModel, "llama-3.3-70b-versatile")
	case "mistral":
		ep.BaseURL = "https://api.mistral.ai/v1"
		ep.APIKey = settings.MistralKey
		ep.Model = firstNonEmpty(overrideModel, settings.MistralModel, "mistral-large-latest")
	case "deepseek":
		ep.BaseURL = "https://api.deepseek.com/v1"
		ep.APIKey = settings.DeepSeekKey
		ep.Model = firstNonEmpty(overrideModel, settings.DeepSeekModel, "deepseek-chat")
	case "xai":
		ep.BaseURL = "https://api.x.ai/v1"
		ep.APIKey = settings.XAIKey
		ep.Model = firstNonEmpty(overrideModel, settings.XAIModel, "grok-3")
	case "moonshot":
		ep.BaseURL = "https://api.moonshot.cn/v1"
		ep.APIKey = settings.MoonshotKey
		ep.Model = firstNonEmpty(overrideModel, settings.MoonshotModel, "kimi-k2")
	default:
		ep.Provider = "openai"
		ep.BaseURL = "https://api.openai.com/v1"
		ep.APIKey = settings.OpenAIKey
		ep.Model = firstNonEmpty(overrideModel, settings.OpenAIModel, "gpt-4o")
	}

	if strings.TrimSpace(ep.APIKey) == "" {
		for _, alt := range []string{"openai", "anthropic", "google", "groq", "deepseek", "mistral", "xai", "moonshot"} {
			if alt == provider {
				continue
			}
			if k := keyForProvider(settings, alt); k != "" {
				return resolveProviderEndpoint(settings, alt, "")
			}
		}
		return ep, fmt.Errorf("no API key configured for %s (configure a key in Settings → AI Models)", ep.Provider)
	}
	return ep, nil
}

func keyForProvider(s *models.AISettings, p string) string {
	switch p {
	case "openai":
		return s.OpenAIKey
	case "anthropic":
		return s.AnthropicKey
	case "google":
		return s.GoogleKey
	case "groq":
		return s.GroqKey
	case "mistral":
		return s.MistralKey
	case "deepseek":
		return s.DeepSeekKey
	case "xai":
		return s.XAIKey
	case "moonshot":
		return s.MoonshotKey
	default:
		return ""
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func callCompletion(ctx context.Context, ep resolvedEndpoint, systemPrompt string, messages []ChatMessage, tools []map[string]any) (ChatMessage, error) {
	if ep.Provider == "anthropic" {
		return callAnthropicMessages(ctx, ep, systemPrompt, messages)
	}
	return callOpenAICompatible(ctx, ep, systemPrompt, messages, tools)
}

func callOpenAICompatible(ctx context.Context, ep resolvedEndpoint, systemPrompt string, messages []ChatMessage, tools []map[string]any) (ChatMessage, error) {
	wireMsgs := make([]map[string]any, 0, len(messages)+1)
	wireMsgs = append(wireMsgs, map[string]any{"role": "system", "content": systemPrompt})
	for _, m := range messages {
		entry := map[string]any{"role": m.Role, "content": m.Content}
		if m.ToolCallID != "" {
			entry["tool_call_id"] = m.ToolCallID
		}
		if len(m.ToolCalls) > 0 {
			entry["tool_calls"] = m.ToolCalls
		}
		wireMsgs = append(wireMsgs, entry)
	}

	reqBody := map[string]any{
		"model":       ep.Model,
		"messages":    wireMsgs,
		"temperature": 0.2,
	}
	if len(tools) > 0 {
		reqBody["tools"] = tools
		reqBody["tool_choice"] = "auto"
	}

	raw, _ := json.Marshal(reqBody)
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "POST", ep.BaseURL+"/chat/completions", bytes.NewReader(raw))
	if err != nil {
		return ChatMessage{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ep.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("AI provider request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return ChatMessage{}, fmt.Errorf("%s error (%d): %s", ep.Provider, resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Choices []struct {
			Message ChatMessage `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil || len(parsed.Choices) == 0 {
		return ChatMessage{}, errors.New("failed to parse AI completion response")
	}
	return parsed.Choices[0].Message, nil
}

func callAnthropicMessages(ctx context.Context, ep resolvedEndpoint, systemPrompt string, messages []ChatMessage) (ChatMessage, error) {
	var wireMsgs []map[string]any
	for _, m := range messages {
		role := m.Role
		if role == "tool" {
			role = "user"
			wireMsgs = append(wireMsgs, map[string]any{
				"role":    role,
				"content": fmt.Sprintf("[Tool Result %s]: %s", m.ToolCallID, m.Content),
			})
			continue
		}
		if role != "assistant" {
			role = "user"
		}
		if strings.TrimSpace(m.Content) != "" {
			wireMsgs = append(wireMsgs, map[string]any{"role": role, "content": m.Content})
		}
	}

	reqBody := map[string]any{
		"model":      ep.Model,
		"max_tokens": 2048,
		"system":     systemPrompt,
		"messages":   wireMsgs,
	}

	raw, _ := json.Marshal(reqBody)
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "POST", ep.BaseURL+"/messages", bytes.NewReader(raw))
	if err != nil {
		return ChatMessage{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", ep.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := client.Do(req)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("Anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return ChatMessage{}, fmt.Errorf("Anthropic error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return ChatMessage{}, err
	}
	var out strings.Builder
	for _, block := range parsed.Content {
		if block.Type == "text" {
			out.WriteString(block.Text)
		}
	}
	return ChatMessage{Role: "assistant", Content: out.String()}, nil
}
