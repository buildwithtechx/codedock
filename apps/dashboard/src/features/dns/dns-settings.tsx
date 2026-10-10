import { CheckCircle2, KeyRound, Plus } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { useGetSettings, useUpdateSettings } from '#/features/settings';
import { SettingsSection } from '#/features/settings/settings-section';
import { DnsProviderForm } from './components/dns-provider-form';

const PROVIDERS = [
  {
    id: 'cloudflare',
    name: 'Cloudflare',
    sub: 'API KEY / TOKEN + ZONE',
    check: (d: Record<string, string>) => Boolean(d.cloudflareApiToken),
  },
  {
    id: 'namecheap',
    name: 'Namecheap',
    sub: 'API USER + KEY',
    check: (d: Record<string, string>) => Boolean(d.namecheapApiUser && d.namecheapApiKey),
  },
  {
    id: 'spaceship',
    name: 'Spaceship',
    sub: 'API KEY + SECRET',
    check: (d: Record<string, string>) => Boolean(d.spaceshipApiKey && d.spaceshipApiSecret),
  },
] as const;

export const DnsSettings = () => {
  const { data, isLoading } = useGetSettings();
  const { mutateAsync: updateSettings, isPending } = useUpdateSettings();

  const [activeProvider, setActiveProvider] = useState<string | null>(null);
  const [formData, setFormData] = useState<Record<string, string>>({
    cloudflareApiToken: '',
    cloudflareEmail: '',
    cloudflareZoneId: '',
    namecheapApiUser: '',
    namecheapApiKey: '',
    namecheapClientIp: '',
    spaceshipApiKey: '',
    spaceshipApiSecret: '',
  });

  useEffect(() => {
    if (data?.data) {
      setFormData((prev) => ({
        ...prev,
        cloudflareApiToken: data.data.cloudflareApiToken || '',
        namecheapApiUser: data.data.namecheapApiUser || '',
        namecheapApiKey: data.data.namecheapApiKey || '',
        namecheapClientIp: data.data.namecheapClientIp || '',
        spaceshipApiKey: data.data.spaceshipApiKey || '',
        spaceshipApiSecret: data.data.spaceshipApiSecret || '',
      }));
    }
  }, [data?.data]);

  const handleSaveProvider = async (provider: string) => {
    try {
      const payload: Record<string, string> = {};
      if (provider === 'cloudflare') {
        payload.cloudflareApiToken = formData.cloudflareApiToken;
        payload.cloudflareEmail = formData.cloudflareEmail;
        payload.cloudflareZoneId = formData.cloudflareZoneId;
      } else if (provider === 'namecheap') {
        payload.namecheapApiUser = formData.namecheapApiUser;
        payload.namecheapApiKey = formData.namecheapApiKey;
        payload.namecheapClientIp = formData.namecheapClientIp;
      } else if (provider === 'spaceship') {
        payload.spaceshipApiKey = formData.spaceshipApiKey;
        payload.spaceshipApiSecret = formData.spaceshipApiSecret;
      }

      await updateSettings({ payload });
      toast.success(`${provider.charAt(0).toUpperCase() + provider.slice(1)} credentials saved`);
      setActiveProvider(null);
    } catch {
      toast.error('Failed to save provider credentials');
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-40 w-full rounded-2xl" />
        <Skeleton className="h-40 w-full rounded-2xl" />
      </div>
    );
  }

  return (
    <SettingsSection
      icon={<KeyRound className="size-4 text-primary" />}
      title="Credentials"
      description="Connect DNS providers for managed domain records and automated SSL verification."
    >
      <div className="space-y-4">
        {PROVIDERS.map((provider) => {
          const isConfigured = provider.check(formData);
          const isExpanded = activeProvider === provider.id;

          return (
            <div key={provider.id} className="rounded-xl border border-border/50 p-4">
              <div className="flex items-start gap-3.5">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50">
                  <img
                    src={`/dns-providers/${provider.id}.svg`}
                    alt={provider.name}
                    className="h-5 w-auto"
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="font-medium text-foreground text-sm">{provider.name}</p>
                  <p className="mt-0.5 text-muted-foreground text-xs">{provider.sub}</p>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setActiveProvider(isExpanded ? null : provider.id)}
                  className="shrink-0 gap-1 text-xs"
                >
                  <Plus className="size-3.5" />
                  {isExpanded ? 'Close' : isConfigured ? 'Edit' : 'Add'}
                </Button>
              </div>

              <div className="mt-3">
                {isConfigured ? (
                  <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
                    <CheckCircle2 className="size-3" /> Active
                  </span>
                ) : (
                  <p className="py-0.5 text-muted-foreground text-xs">None added yet.</p>
                )}
              </div>

              {isExpanded && (
                <div className="mt-4 border-border/40 border-t pt-4">
                  <DnsProviderForm
                    activeProvider={provider.id}
                    formData={formData}
                    setFormData={setFormData}
                    isPending={isPending}
                    handleSaveProvider={handleSaveProvider}
                    onCancel={() => setActiveProvider(null)}
                  />
                </div>
              )}
            </div>
          );
        })}
      </div>
    </SettingsSection>
  );
};
