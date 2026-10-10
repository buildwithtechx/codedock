import { Brain, Check, ChevronDown, Star } from 'lucide-react';
import React, { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Input } from '#/components/ui/input';
import { useGetAISettings, useUpdateAISettings } from '#/features/settings';
import { aiProviderCatalog } from '#/lib/ai-providers';
import { SettingsSection } from './settings-section';

const PROVIDERS = [
  {
    id: 'openai',
    name: 'OpenAI',
    models: aiProviderCatalog.find((c) => c.id === 'openai')?.models || [],
    icon: '/ai-providers/openai.svg',
    keyField: 'openAIKey' as const,
    modelField: 'openAIModel' as const,
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    models: aiProviderCatalog.find((c) => c.id === 'anthropic')?.models || [],
    icon: '/ai-providers/anthropic.svg',
    keyField: 'anthropicKey' as const,
    modelField: 'anthropicModel' as const,
  },
  {
    id: 'google',
    name: 'Google',
    models: aiProviderCatalog.find((c) => c.id === 'google')?.models || [],
    icon: '/ai-providers/google.svg',
    keyField: 'googleKey' as const,
    modelField: 'googleModel' as const,
  },
  {
    id: 'mistral',
    name: 'Mistral',
    models: aiProviderCatalog.find((c) => c.id === 'mistral')?.models || [],
    icon: '/ai-providers/mistral.svg',
    keyField: 'mistralKey' as const,
    modelField: 'mistralModel' as const,
  },
  {
    id: 'groq',
    name: 'Groq',
    models: aiProviderCatalog.find((c) => c.id === 'groq')?.models || [],
    icon: '/ai-providers/groq.svg',
    keyField: 'groqKey' as const,
    modelField: 'groqModel' as const,
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    models: aiProviderCatalog.find((c) => c.id === 'deepseek')?.models || [],
    icon: '/ai-providers/deepseek.svg',
    keyField: 'deepSeekKey' as const,
    modelField: 'deepSeekModel' as const,
  },
  {
    id: 'xai',
    name: 'xAI',
    models: aiProviderCatalog.find((c) => c.id === 'xai')?.models || [],
    icon: '/ai-providers/xai.svg',
    keyField: 'xaiKey' as const,
    modelField: 'xaiModel' as const,
  },
  {
    id: 'moonshot',
    name: 'Moonshot',
    models: aiProviderCatalog.find((c) => c.id === 'moonshot')?.models || [],
    icon: '/ai-providers/moonshot.svg',
    keyField: 'moonshotKey' as const,
    modelField: 'moonshotModel' as const,
  },
] as const;

export function AISettings() {
  const { data: settings } = useGetAISettings();
  const updateSettings = useUpdateAISettings();
  const [localKeys, setLocalKeys] = useState<Record<string, string>>({});

  const pendingSettings = React.useRef<Record<string, unknown> | null>(null);

  React.useEffect(() => {
    if (settings?.data && !updateSettings.isPending) {
      pendingSettings.current = settings.data;
      const initialKeys: Record<string, string> = {};
      for (const p of PROVIDERS) {
        initialKeys[p.keyField] = (settings.data[p.keyField] as string) || '';
      }
      setLocalKeys(initialKeys);
    }
  }, [settings?.data, updateSettings.isPending]);

  const defaultProvider = (settings?.data?.defaultProvider as string) || 'none';

  const handleSetDefault = (id: string) => {
    const payload = { ...(pendingSettings.current || settings?.data || {}), defaultProvider: id };
    pendingSettings.current = payload;
    updateSettings.mutate(payload, {
      onSuccess: () => {
        toast.success(`Default AI provider set to ${id}`);
      },
    });
  };

  const handleUpdateModel = (modelField: string, value: string) => {
    const payload = { ...(pendingSettings.current || settings?.data || {}), [modelField]: value };
    pendingSettings.current = payload;
    updateSettings.mutate(payload, {
      onSuccess: () => {
        toast.success('Model updated');
      },
    });
  };

  const handleSaveProvider = (provider: (typeof PROVIDERS)[number]) => {
    const keyVal = localKeys[provider.keyField] || '';
    const payload = {
      ...(pendingSettings.current || settings?.data || {}),
      [provider.keyField]: keyVal,
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
        <span className="font-mono text-muted-foreground text-xs uppercase">
          Default: <strong className="text-foreground">{defaultProvider}</strong>
        </span>
      }
    >
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        {PROVIDERS.map((provider) => {
          const isDefault = defaultProvider === provider.id;
          const currentModel =
            (settings?.data?.[provider.modelField] as string) || provider.models[0]?.id || '';
          const currentKeyVal = localKeys[provider.keyField] ?? '';

          return (
            <div
              key={provider.id}
              className="flex flex-col justify-between space-y-4 rounded-xl border border-border/50 bg-card p-4 transition-colors hover:bg-muted/30"
            >
              <div className="flex w-full items-start justify-between">
                <div className="flex items-start gap-3">
                  <div className="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-black">
                    <img
                      src={provider.icon}
                      alt={provider.name}
                      className="h-5 w-5 object-contain"
                    />
                  </div>
                  <div className="space-y-1">
                    <h3 className="font-semibold text-sm leading-none">{provider.name}</h3>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <button
                          type="button"
                          className="flex max-w-48 cursor-pointer appearance-none items-center gap-1 truncate bg-transparent font-medium text-muted-foreground text-xs outline-none hover:text-foreground"
                        >
                          <span className="truncate">
                            {provider.models.find((m) => m.id === currentModel)?.name ||
                              currentModel}
                          </span>
                          <ChevronDown className="h-3 w-3 shrink-0" />
                        </button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="start" className="w-52">
                        {provider.models.map((m) => (
                          <DropdownMenuItem
                            key={m.id}
                            onSelect={() => handleUpdateModel(provider.modelField, m.id)}
                            className="flex cursor-pointer items-center justify-between py-1.5"
                          >
                            <div className="flex flex-col">
                              <span className="font-medium text-xs">{m.name}</span>
                              <span className="font-mono text-[10px] text-muted-foreground">
                                {m.id}
                              </span>
                            </div>
                            {currentModel === m.id && (
                              <Check className="h-3.5 w-3.5 text-emerald-500" />
                            )}
                          </DropdownMenuItem>
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => handleSetDefault(isDefault ? 'none' : provider.id)}
                  title={isDefault ? 'Active default provider' : 'Set as default'}
                  className="flex h-8 w-8 items-center justify-center rounded-lg border border-border bg-background text-muted-foreground hover:text-foreground"
                >
                  <Star
                    className={`h-4 w-4 ${isDefault ? 'fill-foreground text-foreground' : ''}`}
                  />
                </button>
              </div>

              <div className="space-y-3">
                <Input
                  type="password"
                  placeholder={`API key for ${provider.name}`}
                  value={currentKeyVal}
                  onChange={(e) =>
                    setLocalKeys((prev) => ({ ...prev, [provider.keyField]: e.target.value }))
                  }
                  className="bg-muted/30 font-mono text-xs"
                />
                <div className="flex justify-end">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => handleSaveProvider(provider)}
                    disabled={updateSettings.isPending}
                    className="gap-1.5 text-xs"
                  >
                    <Check className="size-3.5" />
                    Save {provider.name}
                  </Button>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </SettingsSection>
  );
}
