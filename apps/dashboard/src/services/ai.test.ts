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

  it('resolves provider names and models', () => {
    expect(aiProviderName('openai')).toBe('OpenAI');
    const models = aiProviderModels('openai');
    expect(models.length).toBeGreaterThan(0);
    expect(defaultAiModel('openai')).toBe(models[0].id);
    expect(isAiProviderModel('openai', models[0].id)).toBe(true);
    expect(isAiProviderModel('openai', 'non-existent')).toBe(false);
  });
});
