import { Brain } from 'lucide-react';
import React, { useState } from 'react';
import { toast } from 'sonner';
import { useGetAISettings, useUpdateAISettings } from '#/features/settings';
import { aiProviderCatalog } from '#/lib/ai-providers';
import { AiProviderCard, type AiProviderItem } from './ai-provider-card';
import { SettingsSection } from './settings-section';

const PROVIDERS: AiProviderItem[] = [
  {
    id: 'openai',
    name: 'OpenAI',
    models: aiProviderCatalog.find((c) => c.id === 'openai')?.models || [],
    icon: '/ai-providers/openai.svg',
    keyField: 'openAIKey',
    modelField: 'openAIModel',
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    models: aiProviderCatalog.find((c) => c.id === 'anthropic')?.models || [],
    icon: '/ai-providers/anthropic.svg',
    keyField: 'anthropicKey',
    modelField: 'anthropicModel',
  },
  {
    id: 'google',
    name: 'Google',
    models: aiProviderCatalog.find((c) => c.id === 'google')?.models || [],
    icon: '/ai-providers/google.svg',
    keyField: 'googleKey',
    modelField: 'googleModel',
  },
  {
    id: 'mistral',
    name: 'Mistral',
    models: aiProviderCatalog.find((c) => c.id === 'mistral')?.models || [],
    icon: '/ai-providers/mistral.svg',
    keyField: 'mistralKey',
    modelField: 'mistralModel',
  },
  {
    id: 'groq',
    name: 'Groq',
    models: aiProviderCatalog.find((c) => c.id === 'groq')?.models || [],
    icon: '/ai-providers/groq.svg',
    keyField: 'groqKey',
    modelField: 'groqModel',
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    models: aiProviderCatalog.find((c) => c.id === 'deepseek')?.models || [],
    icon: '/ai-providers/deepseek.svg',
    keyField: 'deepSeekKey',
    modelField: 'deepSeekModel',
  },
  {
    id: 'xai',
    name: 'xAI',
    models: aiProviderCatalog.find((c) => c.id === 'xai')?.models || [],
    icon: '/ai-providers/xai.svg',
    keyField: 'xaiKey',
    modelField: 'xaiModel',
  },
  {
    id: 'moonshot',
    name: 'Moonshot',
    models: aiProviderCatalog.find((c) => c.id === 'moonshot')?.models || [],
    icon: '/ai-providers/moonshot.svg',
    keyField: 'moonshotKey',
    modelField: 'moonshotModel',
  },
];

export function AISettings() {
  const { data: settings } = useGetAISettings();
  const updateSettings = useUpdateAISettings();
  const [localKeys, setLocalKeys] = useState<Record<string, string>>({});
  const [localModels, setLocalModels] = useState<Record<string, string>>({});

  const pendingSettings = React.useRef<Record<string, unknown> | null>(null);

  React.useEffect(() => {
    if (settings?.data && !updateSettings.isPending) {
      pendingSettings.current = settings.data;
      const initialKeys: Record<string, string> = {};
      const initialModels: Record<string, string> = {};
      for (const p of PROVIDERS) {
        initialKeys[p.keyField] = (settings.data[p.keyField] as string) || '';
        initialModels[p.modelField] = (settings.data[p.modelField] as string) || '';
      }
      setLocalKeys(initialKeys);
      setLocalModels(initialModels);
    }
  }, [settings?.data, updateSettings.isPending]);

  const defaultProvider = (settings?.data?.defaultProvider as string) || 'none';

  const handleSetDefault = (id: string) => {
    const nextDefault = defaultProvider === id ? 'none' : id;
    const payload = {
      ...(pendingSettings.current || settings?.data || {}),
      defaultProvider: nextDefault,
    };
    pendingSettings.current = payload;
    updateSettings.mutate(payload, {
      onSuccess: () => {
        toast.success(
          nextDefault === 'none'
            ? 'Default AI provider cleared'
            : `Default AI provider set to ${id}`
        );
      },
    });
  };

  const handleUpdateModel = (modelField: string, value: string) => {
    setLocalModels((prev) => ({ ...prev, [modelField]: value }));
    const payload = { ...(pendingSettings.current || settings?.data || {}), [modelField]: value };
    pendingSettings.current = payload;
    updateSettings.mutate(payload, {
      onSuccess: () => {
        toast.success('Model updated');
      },
    });
  };

  const handleSaveProvider = (provider: AiProviderItem) => {
    const keyVal = localKeys[provider.keyField] || '';
    const modelVal = localModels[provider.modelField] || '';
    const payload = {
      ...(pendingSettings.current || settings?.data || {}),
      [provider.keyField]: keyVal,
      ...(modelVal ? { [provider.modelField]: modelVal } : {}),
    };
    pendingSettings.current = payload;
    updateSettings.mutate(payload, {
      onSuccess: () => {
        toast.success(`${provider.name} credentials saved`);
      },
    });
  };

  return (
    <SettingsSection
      icon={<Brain className="size-4 text-primary" />}
      title="AI Providers"
      description="Configure API keys and model parameters for platform AI services."
      action={
        <div className="flex items-center gap-2">
          <span className="text-muted-foreground text-xs uppercase tracking-wider">Default:</span>
          <span className="rounded-md border border-border/60 bg-muted/40 px-2 py-0.5 font-semibold text-foreground text-xs">
            {defaultProvider === 'none'
              ? 'None'
              : PROVIDERS.find((p) => p.id === defaultProvider)?.name || defaultProvider}
          </span>
        </div>
      }
    >
      <div className="space-y-4">
        {PROVIDERS.map((provider) => {
          const currentModel =
            localModels[provider.modelField] ||
            (settings?.data?.[provider.modelField] as string) ||
            provider.models[0]?.id ||
            '';
          const currentKeyVal = localKeys[provider.keyField] ?? '';

          return (
            <AiProviderCard
              key={provider.id}
              provider={provider}
              currentModel={currentModel}
              currentKey={currentKeyVal}
              isDefault={defaultProvider === provider.id}
              isPending={updateSettings.isPending}
              onModelChange={(modelId) => handleUpdateModel(provider.modelField, modelId)}
              onKeyChange={(val) => setLocalKeys((prev) => ({ ...prev, [provider.keyField]: val }))}
              onSetDefault={() => handleSetDefault(provider.id)}
              onSave={() => handleSaveProvider(provider)}
            />
          );
        })}
      </div>
    </SettingsSection>
  );
}
