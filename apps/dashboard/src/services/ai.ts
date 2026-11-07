import { AI_PROVIDER_CATALOG, AI_PROVIDER_IDS } from '#/constants/ai';
import type {
  AiChatRequestPayload,
  AiChatResponsePayload,
  AiProviderId,
  AiProviderModel,
  AiProviderStatusSummary,
  AiResourceTag,
} from '#/interfaces/ai';
import { apiClient } from '#/lib/api-client';

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

const PROVIDER_FIELD_MAP: Record<AiProviderId, { keyField: string; modelField: string }> = {
  openai: { keyField: 'openAIKey', modelField: 'openAIModel' },
  anthropic: { keyField: 'anthropicKey', modelField: 'anthropicModel' },
  google: { keyField: 'googleKey', modelField: 'googleModel' },
  mistral: { keyField: 'mistralKey', modelField: 'mistralModel' },
  groq: { keyField: 'groqKey', modelField: 'groqModel' },
  deepseek: { keyField: 'deepSeekKey', modelField: 'deepSeekModel' },
  xai: { keyField: 'xaiKey', modelField: 'xaiModel' },
  moonshot: { keyField: 'moonshotKey', modelField: 'moonshotModel' },
};

export function evaluateAiProviderStatus(
  rawSettings?: Record<string, unknown> | null
): AiProviderStatusSummary {
  const configuredProviders: AiProviderId[] = [];
  for (const id of AI_PROVIDER_IDS) {
    const { keyField } = PROVIDER_FIELD_MAP[id];
    const keyVal = rawSettings?.[keyField];
    if (typeof keyVal === 'string' && keyVal.trim().length > 0) {
      configuredProviders.push(id);
    }
  }

  const rawDefault =
    typeof rawSettings?.defaultProvider === 'string' ? rawSettings.defaultProvider.trim() : 'none';
  const hasExplicitDefault = isAiProviderId(rawDefault);
  const defaultHasKey = hasExplicitDefault && configuredProviders.includes(rawDefault);

  if (configuredProviders.length === 0) {
    return {
      state: 'unconfigured',
      defaultProvider: rawDefault || 'none',
      activeProviderId: null,
      activeProviderName: 'Local Instance Copilot',
      activeModel: 'instance-copilot',
      configuredProviders: [],
    };
  }

  const activeProviderId = defaultHasKey ? rawDefault : configuredProviders[0];
  const { modelField } = PROVIDER_FIELD_MAP[activeProviderId];
  const savedModel = rawSettings?.[modelField];
  const activeModel =
    typeof savedModel === 'string' && savedModel.trim().length > 0
      ? savedModel.trim()
      : defaultAiModel(activeProviderId);

  return {
    state: defaultHasKey ? 'ready' : 'no_default',
    defaultProvider: rawDefault || 'none',
    activeProviderId,
    activeProviderName: aiProviderName(activeProviderId),
    activeModel,
    configuredProviders: configuredProviders.map((id) => aiProviderName(id)),
  };
}

interface AiResourcesEnvelope {
  data?: AiResourceTag[];
}

interface AiChatEnvelope {
  data?: AiChatResponsePayload;
  reply?: string;
  steps?: AiChatResponsePayload['steps'];
  pendingActions?: AiChatResponsePayload['pendingActions'];
}

export const aiAssistantService = {
  async listResources(): Promise<AiResourceTag[]> {
    const res = await apiClient.get<AiResourcesEnvelope>('/ai/resources');
    return res.data ?? [];
  },

  async sendChat(payload: AiChatRequestPayload): Promise<AiChatResponsePayload> {
    const res = await apiClient.post<AiChatEnvelope>('/ai/chat', payload);
    return res.data ?? res;
  },

  async confirmAction(actionId: string): Promise<Record<string, unknown>> {
    return apiClient.post<Record<string, unknown>>('/ai/actions/confirm', { actionId });
  },
};
