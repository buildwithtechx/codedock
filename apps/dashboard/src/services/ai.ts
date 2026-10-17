import { AI_PROVIDER_CATALOG, AI_PROVIDER_IDS } from '#/constants/ai';
import type { AiProviderId, AiProviderModel } from '#/interfaces/ai';

export function isAiProviderId(value: unknown): value is AiProviderId {
  return typeof value === 'string' && AI_PROVIDER_IDS.includes(value as AiProviderId);
}

export function aiProviderName(providerId: AiProviderId): string {
  return AI_PROVIDER_CATALOG.find((provider) => provider.id === providerId)?.name ?? providerId;
}

export function aiProviderModels(providerId: AiProviderId): readonly AiProviderModel[] {
  return AI_PROVIDER_CATALOG.find((provider) => provider.id === providerId)?.models ?? [];
}

export function defaultAiModel(providerId: AiProviderId): string {
  return aiProviderModels(providerId)[0]?.id ?? '';
}

export function isAiProviderModel(providerId: AiProviderId, modelId: unknown): modelId is string {
  return (
    typeof modelId === 'string' &&
    aiProviderModels(providerId).some((model) => model.id === modelId)
  );
}
