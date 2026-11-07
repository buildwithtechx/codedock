import { type AiProviderCatalogItem, aiProviderIdValues } from '#/interfaces/ai';

export const AI_PROVIDER_CATALOG = [
  {
    id: 'openai',
    name: 'OpenAI',
    apiKeyPlaceholder: 'sk-...',
    models: [
      { id: 'gpt-6-astra', name: 'GPT-6 Astra' },
      { id: 'gpt-6.1-sol', name: 'GPT-6.1 Sol' },
      { id: 'gpt-6-luna', name: 'GPT-6 Luna' },
      { id: 'gpt-5.4', name: 'GPT-5.4' },
      { id: 'gpt-5.4-mini', name: 'GPT-5.4 mini' },
      { id: 'gpt-5.2', name: 'GPT-5.2' },
    ],
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    apiKeyPlaceholder: 'sk-ant-...',
    models: [
      { id: 'claude-opus-5-5', name: 'Claude Opus 5.5' },
      { id: 'claude-sonnet-5-5', name: 'Claude Sonnet 5.5' },
      { id: 'claude-haiku-5-5', name: 'Claude Haiku 5.5' },
      { id: 'claude-fable-5-1', name: 'Claude Fable 5.1' },
      { id: 'claude-3-7-sonnet-latest', name: 'Claude 3.7 Sonnet' },
    ],
  },
  {
    id: 'google',
    name: 'Google',
    apiKeyPlaceholder: 'Google AI API key',
    models: [
      { id: 'gemini-4-argon', name: 'Gemini 4 Argon' },
      { id: 'gemini-3.8-flash', name: 'Gemini 3.8 Flash' },
      { id: 'gemini-3.5-flash', name: 'Gemini 3.5 Flash' },
      { id: 'gemini-3.1-pro-preview', name: 'Gemini 3.1 Pro' },
      { id: 'gemini-3.1-flash-lite', name: 'Gemini 3.1 Flash-Lite' },
    ],
  },
  {
    id: 'mistral',
    name: 'Mistral',
    apiKeyPlaceholder: 'Mistral API key',
    models: [
      { id: 'mistral-large-latest', name: 'Mistral Large 4' },
      { id: 'mistral-medium-latest', name: 'Mistral Medium 3.5' },
      { id: 'codestral-latest', name: 'Codestral' },
      { id: 'mistral-small-latest', name: 'Mistral Small 4' },
      { id: 'ministral-8b-latest', name: 'Ministral 8B' },
    ],
  },
  {
    id: 'groq',
    name: 'Groq',
    apiKeyPlaceholder: 'gsk_...',
    models: [
      { id: 'llama-3.3-70b-versatile', name: 'Llama 3.3 70B Versatile' },
      { id: 'deepseek-r1-distill-llama-70b', name: 'DeepSeek R1 Distill Llama 70B' },
      { id: 'llama-3.1-8b-instant', name: 'Llama 3.1 8B Instant' },
      { id: 'qwen-2.5-32b', name: 'Qwen 2.5 32B' },
    ],
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    apiKeyPlaceholder: 'DeepSeek API key',
    models: [
      { id: 'deepseek-chat', name: 'DeepSeek Chat (V4)' },
      { id: 'deepseek-reasoner', name: 'DeepSeek Reasoner' },
      { id: 'deepseek-v4-pro', name: 'DeepSeek V4 Pro' },
      { id: 'deepseek-v4-flash', name: 'DeepSeek V4 Flash' },
    ],
  },
  {
    id: 'xai',
    name: 'xAI',
    apiKeyPlaceholder: 'xAI API key',
    models: [
      { id: 'grok-4.7', name: 'Grok 4.7' },
      { id: 'grok-4.3', name: 'Grok 4.3' },
      { id: 'grok-4.3-latest', name: 'Grok 4.3 Latest' },
      { id: 'grok-4', name: 'Grok 4' },
    ],
  },
  {
    id: 'moonshot',
    name: 'Moonshot',
    apiKeyPlaceholder: 'Moonshot API key',
    models: [
      { id: 'kimi-k3', name: 'Kimi K3' },
      { id: 'kimi-k2.7-code', name: 'Kimi K2.7 Code' },
      { id: 'kimi-k2.7-code-highspeed', name: 'Kimi K2.7 Code Highspeed' },
      { id: 'kimi-k2.6', name: 'Kimi K2.6' },
    ],
  },
] as const satisfies readonly AiProviderCatalogItem[];

export const AI_PROVIDER_IDS = [...aiProviderIdValues];
