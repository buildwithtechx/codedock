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

export interface AiQuestionOption {
  id: string;
  label: string;
  description?: string;
}

export interface AiQuestionItem {
  id: string;
  title: string;
  allowOther?: boolean;
  optional?: boolean;
  multiSelect?: boolean;
  isSecret?: boolean;
  options?: AiQuestionOption[];
}

export interface AiQuestionGroup {
  id: string;
  title: string;
  questions: AiQuestionItem[];
}

export interface AiQuestionAnswer {
  options?: string[];
  text?: string;
}

export interface AiTurnMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  createdAt: string;
  mentions?: AiResourceTag[];
  steps?: AiToolStep[];
  pendingActions?: AiPendingAction[];
  questions?: AiQuestionGroup[];
}

export type AiPersonaId = 'general' | 'logs' | 'security' | 'devops';

export interface AiPersonaDefinition {
  id: AiPersonaId;
  name: string;
  description: string;
  icon: 'bot' | 'terminal' | 'shield-check' | 'rocket';
  readOnly: boolean;
  systemHint: string;
}

export const AI_PERSONAS: readonly AiPersonaDefinition[] = [
  {
    id: 'general',
    name: 'General Assistant',
    description: 'General workspace questions, navigation, and guidance',
    icon: 'bot',
    readOnly: false,
    systemHint: 'General DevOps and platform co-pilot',
  },
  {
    id: 'logs',
    name: 'Logs Investigator',
    description: 'Diagnose container crashes, exit codes, stderr/stdout, and restart loops',
    icon: 'terminal',
    readOnly: true,
    systemHint: 'Container runtime & build log diagnostic specialist',
  },
  {
    id: 'security',
    name: 'Security Reviewer',
    description: 'Audit runtime CVEs, signed security advisories, and semver patches',
    icon: 'shield-check',
    readOnly: true,
    systemHint: 'Vulnerability and container image security specialist',
  },
  {
    id: 'devops',
    name: 'DevOps Co-pilot',
    description: 'Active operations: inspect infrastructure, configure domains, deploy and scale',
    icon: 'rocket',
    readOnly: false,
    systemHint: 'Active infrastructure lifecycle and deployment coordinator',
  },
] as const;

export interface AiConversationThread {
  id: string;
  title: string;
  updatedAt: string;
  messages: AiTurnMessage[];
  persona?: AiPersonaId;
}

export interface AiChatRequestPayload {
  messages: Array<{ role: 'user' | 'assistant'; content: string }>;
  mode: 'ask' | 'agent';
  persona?: AiPersonaId;
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
  questions?: AiQuestionGroup[];
}

export interface AiProviderStatusSummary {
  state: 'ready' | 'no_default' | 'unconfigured';
  defaultProvider: string;
  activeProviderId: AiProviderId | null;
  activeProviderName: string;
  activeModel: string;
  configuredProviders: string[];
}
