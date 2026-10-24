export type AiProviderModel = {
  id: string;
  name: string;
};

export const aiProviderIdValues = [
  'openai',
  'anthropic',
  'google',
  'mistral',
  'groq',
  'deepseek',
  'xai',
  'moonshot',
] as const;

export type AiProviderId = (typeof aiProviderIdValues)[number];

export type AiProviderCatalogItem = {
  id: AiProviderId;
  name: string;
  apiKeyPlaceholder: string;
  models: readonly AiProviderModel[];
};

export interface AiResourceTag {
  id: string;
  type: string;
  name: string;
  subtitle?: string;
  status?: string;
  href?: string;
}

export interface AiToolStep {
  id: string;
  tool: string;
  label: string;
  status: 'pending' | 'done' | 'error';
  detail?: string;
}

export interface AiPendingAction {
  actionId: string;
  tool: string;
  title: string;
  description: string;
  method: string;
  path: string;
  arguments: Record<string, unknown>;
  destructive: boolean;
  resolvedStatus?: 'approved' | 'cancelled' | 'failed';
}

export interface AiTurnMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  createdAt: string;
  mentions?: AiResourceTag[];
  steps?: AiToolStep[];
  pendingActions?: AiPendingAction[];
}

export interface AiConversationThread {
  id: string;
  title: string;
  updatedAt: string;
  messages: AiTurnMessage[];
}

export interface AiChatRequestPayload {
  messages: Array<{ role: 'user' | 'assistant'; content: string }>;
  mode: 'ask' | 'agent';
  provider?: string;
  model?: string;
  currentPath?: string;
  mentions?: AiResourceTag[];
}

export interface AiChatResponsePayload {
  reply?: string;
  provider?: string;
  model?: string;
  configured?: boolean;
  steps?: AiToolStep[];
  pendingActions?: AiPendingAction[];
  executedActions?: string[];
}

export interface AiProviderStatusSummary {
  state: 'ready' | 'no_default' | 'unconfigured';
  defaultProvider: string;
  activeProviderId: AiProviderId | null;
  activeProviderName: string;
  activeModel: string;
  configuredProviders: string[];
}
