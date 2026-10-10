import { Lock } from 'lucide-react';
import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Skeleton } from '#/components/ui/skeleton';
import { SettingsSection } from '#/features/settings/settings-section';
import { useListOAuthProviders, useSaveOAuthProvider } from '#/hooks/use-oauth';
import type { SaveOAuthProviderRequest } from '#/interfaces/oauth';
import { OAuthProviderCard, type OAuthProviderDef } from './oauth-provider-card';

const PROVIDERS: OAuthProviderDef[] = [
  {
    id: 'github',
    name: 'GitHub',
    icon: '/git-providers/github-icon.svg',
    fields: [
      { key: 'clientId', label: 'Client ID', placeholder: 'Iv1.abc123...' },
      { key: 'clientSecret', label: 'Client Secret', placeholder: '••••••••', type: 'password' },
    ],
  },
  {
    id: 'gitlab',
    name: 'GitLab',
    icon: '/git-providers/gitlab-icon.svg',
    fields: [
      { key: 'clientId', label: 'Application ID', placeholder: 'abc123...' },
      { key: 'clientSecret', label: 'Secret', placeholder: '••••••••', type: 'password' },
      {
        key: 'baseUrl',
        label: 'Self-hosted URL (optional)',
        placeholder: 'https://gitlab.example.com',
      },
    ],
  },
  {
    id: 'google',
    name: 'Google',
    icon: '/oauth-providers/google.svg',
    fields: [
      { key: 'clientId', label: 'Client ID', placeholder: '123456789.apps.googleusercontent.com' },
      { key: 'clientSecret', label: 'Client Secret', placeholder: 'GOCSPX-...', type: 'password' },
    ],
  },
  {
    id: 'microsoft',
    name: 'Microsoft',
    icon: '/oauth-providers/microsoft.svg',
    fields: [
      {
        key: 'clientId',
        label: 'Application (Client) ID',
        placeholder: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx',
      },
      { key: 'clientSecret', label: 'Client Secret', placeholder: '••••••••', type: 'password' },
      { key: 'tenant', label: 'Tenant ID', placeholder: 'common' },
    ],
  },
];

type ProviderState = Record<string, Partial<SaveOAuthProviderRequest>>;

export const OAuthProvidersList = () => {
  const { data, isLoading } = useListOAuthProviders();
  const { mutateAsync: save, isPending } = useSaveOAuthProvider();

  const [form, setForm] = useState<ProviderState>({});
  const [saving, setSaving] = useState<string | null>(null);

  const isInitialized = React.useRef(false);

  useEffect(() => {
    if (data?.data) {
      if (!isInitialized.current) {
        isInitialized.current = true;
        const initial: ProviderState = {};
        for (const p of PROVIDERS) {
          const existing = data.data.find((d) => d.providerName === p.id);
          initial[p.id] = {
            id: existing?.id,
            providerName: p.id,
            clientId: existing?.clientId ?? '',
            clientSecret: existing?.clientSecret ? '********' : '',
            baseUrl: existing?.baseUrl ?? '',
            tenant: existing?.tenant ?? '',
            enabled: existing?.enabled ?? false,
          };
        }
        setForm(initial);
      } else {
        setForm((f) => {
          const next = { ...f };
          let changed = false;
          for (const p of PROVIDERS) {
            const existing = data.data.find((d) => d.providerName === p.id);
            if (existing?.id && next[p.id]?.id !== existing.id) {
              next[p.id] = { ...next[p.id], id: existing.id };
              changed = true;
            }
          }
          return changed ? next : f;
        });
      }
    }
  }, [data]);

  const set = (providerId: string, key: keyof SaveOAuthProviderRequest, value: unknown) => {
    setForm((f) => ({ ...f, [providerId]: { ...f[providerId], [key]: value } }));
  };

  const handleSaveProvider = async (provider: OAuthProviderDef) => {
    setSaving(provider.id);
    try {
      const state = form[provider.id] as SaveOAuthProviderRequest;
      await save({
        payload: {
          ...state,
          redirectUri: `${window.location.origin}/api/auth/oauth/${provider.id}/callback`,
          clientSecret: state.clientSecret === '********' ? undefined : state.clientSecret,
          enabled:
            !state.clientId || (!state.clientSecret && state.clientSecret !== '********')
              ? false
              : state.enabled,
        },
      });
      toast.success(`${provider.name} OAuth settings saved`);
    } catch {
      toast.error(`Failed to save ${provider.name} OAuth settings`);
    } finally {
      setSaving(null);
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[...Array(4)].map((_, i) => (
          <Skeleton key={i} className="h-32 w-full rounded-2xl" />
        ))}
      </div>
    );
  }

  return (
    <SettingsSection
      icon={<Lock className="size-4 text-primary" />}
      title="OAuth Authentication"
      description="Configure single sign-on providers for your workspace users."
    >
      <div className="space-y-4">
        {PROVIDERS.map((provider) => (
          <OAuthProviderCard
            key={provider.id}
            provider={provider}
            state={form[provider.id] ?? {}}
            isPending={isPending}
            isSaving={saving === provider.id}
            onFieldChange={(key, val) => set(provider.id, key, val)}
            onSave={() => handleSaveProvider(provider)}
          />
        ))}
      </div>
    </SettingsSection>
  );
};
