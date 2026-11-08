import { Button } from '@codedock/ui/button';
import { Link } from '@tanstack/react-router';
import { ArrowRight, Bot, Brain, ExternalLink, MessageSquare, Radio } from 'lucide-react';
import React, { useState } from 'react';
import { toast } from 'sonner';
import { AI_PROVIDER_CATALOG } from '#/constants/ai';
import { useGetAISettings, useUpdateAISettings } from '#/hooks/use-settings';
import { SettingsSection } from '../settings-section';
import { AiProviderCard, type AiProviderItem } from './ai-provider-card';

const PROVIDERS: AiProviderItem[] = [
  {
    id: 'openai',
    name: 'OpenAI',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'openai')?.models ||
      []) as AiProviderItem['models'],
    icon: '/ai-providers/openai.svg',
    keyField: 'openAIKey',
    modelField: 'openAIModel',
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'anthropic')?.models ||
      []) as AiProviderItem['models'],
    icon: '/ai-providers/anthropic.svg',
    keyField: 'anthropicKey',
    modelField: 'anthropicModel',
  },
  {
    id: 'google',
    name: 'Google',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'google')?.models ||
      []) as AiProviderItem['models'],
    icon: '/ai-providers/google.svg',
    keyField: 'googleKey',
    modelField: 'googleModel',
  },
  {
    id: 'mistral',
    name: 'Mistral',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'mistral')?.models ||
      []) as AiProviderItem['models'],
    icon: '/ai-providers/mistral.svg',
    keyField: 'mistralKey',
    modelField: 'mistralModel',
  },
  {
    id: 'groq',
    name: 'Groq',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'groq')?.models ||
      []) as AiProviderItem['models'],
    icon: '/ai-providers/groq.svg',
    keyField: 'groqKey',
    modelField: 'groqModel',
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'deepseek')?.models ||
      []) as AiProviderItem['models'],
    icon: '/ai-providers/deepseek.svg',
    keyField: 'deepSeekKey',
    modelField: 'deepSeekModel',
  },
  {
    id: 'xai',
    name: 'xAI',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'xai')?.models ||
      []) as AiProviderItem['models'],
    icon: '/ai-providers/xai.svg',
    keyField: 'xaiKey',
    modelField: 'xaiModel',
  },
  {
    id: 'moonshot',
    name: 'Moonshot',
    models: (AI_PROVIDER_CATALOG.find((c) => c.id === 'moonshot')?.models ||
      []) as AiProviderItem['models'],
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

  const hasAnyKey = PROVIDERS.some((p) => Boolean(localKeys[p.keyField]));
  const isAssistantEnabled = defaultProvider !== 'none' || hasAnyKey;

  return (
    <div className="space-y-6">
      {/* Your Assistant Card */}
      <section className="rounded-2xl border border-border/50 bg-card p-5 shadow-xs sm:p-6">
        <div className="flex items-start justify-between gap-4">
          <div className="flex items-start gap-4">
            <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-muted text-primary">
              <Bot className="h-6 w-6" />
            </span>
            <div>
              <h2 className="font-medium text-foreground text-lg">Your assistant</h2>
              <p className="mt-1 max-w-xl text-muted-foreground text-sm leading-6">
                Your connected agent, your apps, your permissions. Manage its access here and build
                event-driven flows in Automations.
              </p>
            </div>
          </div>
          <span
            className={`rounded-full px-3 py-1 font-medium text-xs ${
              isAssistantEnabled
                ? 'bg-emerald-500/10 text-emerald-500'
                : 'bg-muted text-muted-foreground'
            }`}
          >
            {isAssistantEnabled ? 'Enabled' : 'Off'}
          </span>
        </div>

        <div className="mt-5">
          <Button asChild size="default" className="gap-2 rounded-xl">
            <Link to="/assistant">
              {isAssistantEnabled ? 'Open assistant' : 'Enable assistant'}
              <ArrowRight className="h-4 w-4" />
            </Link>
          </Button>
        </div>
      </section>

      {/* Chat apps & notifications Card */}
      <section className="rounded-2xl border border-border/50 bg-card p-5 shadow-xs sm:p-6">
        <h3 className="font-medium text-base text-foreground">Chat apps & notifications</h3>
        <p className="mt-2 text-muted-foreground text-sm leading-6">
          Talk to your assistant from chat apps, control deployments with commands, and choose where
          automation results arrive.
        </p>
        <div className="mt-5 flex flex-wrap gap-3">
          <Link
            to="/settings"
            search={{ tab: 'chat-apps' } as never}
            className="inline-flex items-center gap-1.5 rounded-xl border border-border/50 bg-muted/60 px-4 py-2.5 font-medium text-foreground text-sm transition-colors hover:bg-muted"
          >
            <MessageSquare className="h-4 w-4 text-muted-foreground" />
            Chat apps
            <ExternalLink className="h-3.5 w-3.5 opacity-60" />
          </Link>
          <Link
            to="/settings"
            search={{ tab: 'notifications' } as never}
            className="inline-flex items-center gap-1.5 rounded-xl border border-border/50 bg-muted/60 px-4 py-2.5 font-medium text-foreground text-sm transition-colors hover:bg-muted"
          >
            <Radio className="h-4 w-4 text-muted-foreground" />
            Notification channels
            <ExternalLink className="h-3.5 w-3.5 opacity-60" />
          </Link>
        </div>
      </section>

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
                onKeyChange={(val) =>
                  setLocalKeys((prev) => ({ ...prev, [provider.keyField]: val }))
                }
                onSetDefault={() => handleSetDefault(provider.id)}
                onSave={() => handleSaveProvider(provider)}
              />
            );
          })}
        </div>
      </SettingsSection>
    </div>
  );
}
