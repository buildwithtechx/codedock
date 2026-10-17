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
