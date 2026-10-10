import { Check, CheckCircle2, Link, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { SettingsSection } from '#/features/settings/settings-section';
import { useConnect, useDisconnect, useGetStatus } from '#/hooks/use-git';

const PROVIDERS = [
  { id: 'github', name: 'GitHub', icon: '/git-providers/github-icon.svg' },
  { id: 'gitlab', name: 'GitLab', icon: '/git-providers/gitlab-icon.svg' },
  { id: 'bitbucket', name: 'Bitbucket', icon: '/git-providers/bitbucket-icon.svg' },
  { id: 'gitea', name: 'Gitea', icon: '/git-providers/gitea-icon.svg' },
];

export function GitProviders() {
  const { data, isLoading } = useGetStatus();
  const connectMutation = useConnect();
  const disconnectMutation = useDisconnect();

  const statuses = (data?.data as any[]) || [];

  const [activeProvider, setActiveProvider] = useState<string | null>(null);
  const [accessToken, setAccessToken] = useState('');
  const [accountName, setAccountName] = useState('');

  const getStatus = (providerId: string) => {
    return statuses.find((s) => s.provider === providerId && s.connected);
  };

  const handleConnect = (e: React.FormEvent, providerId: string) => {
    e.preventDefault();
    if (!accessToken.trim()) return;

    connectMutation.mutate(
      {
        provider: providerId,
        accessToken: accessToken.trim(),
        accountName: accountName.trim() || 'Personal',
      },
      {
        onSuccess: () => {
          setActiveProvider(null);
          setAccessToken('');
          setAccountName('');
          toast.success(`Successfully connected ${providerId}`);
        },
        onError: (err: any) => {
          toast.error(err.message || 'Failed to connect provider');
        },
      }
    );
  };

  const handleDisconnect = (providerId: string) => {
    disconnectMutation.mutate(providerId, {
      onSuccess: () => {
        toast.success(`Successfully disconnected ${providerId}`);
      },
      onError: (err: any) => {
        toast.error(err.message || 'Failed to disconnect provider');
      },
    });
  };

  if (isLoading) {
    return <div className="h-64 animate-pulse rounded-xl bg-card" />;
  }

  return (
    <SettingsSection
      icon={<Link className="size-4 text-primary" />}
      title="Personal Git Providers"
      description="Connect your personal accounts using Access Tokens."
    >
      <div className="space-y-4">
        {PROVIDERS.map((provider) => {
          const status = getStatus(provider.id);
          const isConnected = !!status;
          const isExpanded = activeProvider === provider.id;

          return (
            <div key={provider.id} className="rounded-xl border border-border/50 p-4">
              <div className="flex items-start gap-3.5">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50">
                  <img
                    src={provider.icon}
                    alt={provider.name}
                    className={`h-5 w-5 object-contain ${
                      provider.id === 'github'
                        ? 'dim:brightness-0 dim:invert dark:brightness-0 dark:invert'
                        : ''
                    }`}
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="font-medium text-foreground text-sm">{provider.name}</p>
                  <p className="mt-0.5 text-muted-foreground text-xs">Personal Access Token</p>
                </div>
                {isConnected ? (
                  <Button
                    variant="outline"
                    size="sm"
                    className="shrink-0 gap-1 text-destructive text-xs hover:bg-destructive/10 hover:text-destructive"
                    onClick={() => handleDisconnect(provider.id)}
                    disabled={disconnectMutation.isPending}
                  >
                    <Trash2 className="size-3.5" />
                    Disconnect
                  </Button>
                ) : (
                  <Button
                    variant="outline"
                    size="sm"
                    className="shrink-0 gap-1 text-xs"
                    onClick={() => {
                      if (isExpanded) {
                        setActiveProvider(null);
                      } else {
                        setAccessToken('');
                        setAccountName('');
                        setActiveProvider(provider.id);
                      }
                    }}
                  >
                    <Plus className="size-3.5" />
                    {isExpanded ? 'Close' : 'Connect'}
                  </Button>
                )}
              </div>

              <div className="mt-3">
                {isConnected ? (
                  <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
                    <CheckCircle2 className="size-3" /> Connected as {status.accountName}
                  </span>
                ) : (
                  <span className="inline-flex items-center gap-1 rounded-full bg-muted px-2 py-0.5 font-medium text-[10px] text-muted-foreground">
                    Not connected
                  </span>
                )}
              </div>

              {isExpanded && !isConnected && (
                <form
                  onSubmit={(e) => handleConnect(e, provider.id)}
                  className="mt-4 space-y-4 border-border/40 border-t pt-4"
                >
                  <div className="space-y-1.5">
                    <Label className="text-xs">Account Name (Optional)</Label>
                    <Input
                      placeholder="e.g. My Personal Account"
                      value={accountName}
                      onChange={(e) => setAccountName(e.target.value)}
                      className="bg-muted/30 font-mono text-xs"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label className="text-xs">Personal Access Token</Label>
                    <Input
                      type="password"
                      placeholder="Token with repository read access"
                      value={accessToken}
                      onChange={(e) => setAccessToken(e.target.value)}
                      required
                      className="bg-muted/30 font-mono text-xs"
                    />
                  </div>

                  <div className="flex items-center justify-end gap-2 border-border/40 border-t pt-3">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => setActiveProvider(null)}
                      className="text-xs"
                    >
                      Cancel
                    </Button>
                    <Button
                      type="submit"
                      size="sm"
                      disabled={connectMutation.isPending || !accessToken.trim()}
                      className="gap-1.5 text-xs"
                    >
                      <Check className="size-3.5" />
                      {connectMutation.isPending ? 'Connecting...' : `Connect ${provider.name}`}
                    </Button>
                  </div>
                </form>
              )}
            </div>
          );
        })}
      </div>
    </SettingsSection>
  );
}
