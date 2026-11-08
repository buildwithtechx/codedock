export const CHAT_APP_PROVIDERS = ['slack', 'discord', 'telegram', 'teams'] as const;
export type ChatAppProvider = (typeof CHAT_APP_PROVIDERS)[number];

export const CHAT_APP_PURPOSES = ['commands', 'assistant', 'both'] as const;
export type ChatAppPurpose = (typeof CHAT_APP_PURPOSES)[number];

export const CHAT_APP_TRANSPORTS = ['webhook', 'outbound'] as const;
export type ChatAppTransport = (typeof CHAT_APP_TRANSPORTS)[number];

export type ChatAppStatus = 'connected' | 'disconnected' | 'error' | 'pairing_required';

export interface ChatAppConnection {
  id: string;
  userId: string;
  provider: ChatAppProvider;
  transport: ChatAppTransport;
  purpose: ChatAppPurpose;
  name: string;
  botUsername: string;
  config: Record<string, unknown>;
  enabled: boolean;
  status: ChatAppStatus;
  lastError?: string | null;
  lastEventAt?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface ChatAppIdentity {
  id: string;
  connectionId: string;
  userId: string;
  provider: ChatAppProvider;
  externalUserId: string;
  externalUsername: string;
  role: string;
  permissions: string[];
  active: boolean;
  pairedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface ChatAppChallenge {
  id: string;
  connectionId: string;
  provider: ChatAppProvider;
  externalUserId: string;
  externalChannelId?: string;
  token: string;
  code: string;
  expiresAt: string;
  verifiedAt?: string | null;
  createdAt: string;
}

export interface CreateChatAppInput {
  provider: ChatAppProvider;
  transport?: ChatAppTransport;
  purpose?: ChatAppPurpose;
  name?: string;
  botUsername?: string;
  credentials: Record<string, string>;
  config?: Record<string, unknown>;
}

export interface UpdateChatAppInput {
  name?: string;
  transport?: ChatAppTransport;
  purpose?: ChatAppPurpose;
  botUsername?: string;
  enabled?: boolean;
  credentials?: Record<string, string>;
  config?: Record<string, unknown>;
}

export interface VerifyPairingInput {
  token?: string;
  code: string;
}
