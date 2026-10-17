import { describe, expect, it } from 'vitest';
import {
  aiProviderModels,
  aiProviderName,
  defaultAiModel,
  isAiProviderId,
  isAiProviderModel,
} from './ai';

describe('services/ai', () => {
  it('validates known provider ids', () => {
    expect(isAiProviderId('openai')).toBe(true);
    expect(isAiProviderId('anthropic')).toBe(true);
    expect(isAiProviderId('unknown')).toBe(false);
  });

  it('resolves provider names and models for all supported providers', () => {
    const providers = [
      { id: 'openai', name: 'OpenAI', defaultModel: 'gpt-6-astra' },
      { id: 'anthropic', name: 'Anthropic', defaultModel: 'claude-opus-5-5' },
      { id: 'google', name: 'Google', defaultModel: 'gemini-4-argon' },
      { id: 'mistral', name: 'Mistral', defaultModel: 'mistral-large-latest' },
      { id: 'groq', name: 'Groq', defaultModel: 'llama-3.3-70b-versatile' },
      { id: 'deepseek', name: 'DeepSeek', defaultModel: 'deepseek-chat' },
      { id: 'xai', name: 'xAI', defaultModel: 'grok-4.7' },
      { id: 'moonshot', name: 'Moonshot', defaultModel: 'kimi-k3' },
    ] as const;

    for (const provider of providers) {
      expect(isAiProviderId(provider.id)).toBe(true);
      expect(aiProviderName(provider.id)).toBe(provider.name);
      const models = aiProviderModels(provider.id);
      expect(models.length).toBeGreaterThan(0);
      expect(defaultAiModel(provider.id)).toBe(provider.defaultModel);
      expect(isAiProviderModel(provider.id, provider.defaultModel)).toBe(true);
      expect(isAiProviderModel(provider.id, 'non-existent')).toBe(false);
    }
  });
});
