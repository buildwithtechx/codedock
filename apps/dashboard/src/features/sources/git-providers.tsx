import { Check, CheckCircle2, Link, Trash } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '#/components/ui/dialog';
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

  const [isConnecting, setIsConnecting] = useState<string | null>(null);
  const [isDisconnecting, setIsDisconnecting] = useState<string | null>(null);

  const [accessToken, setAccessToken] = useState('');
  const [accountName, setAccountName] = useState('');

  const getStatus = (providerId: string) => {
    return statuses.find((s) => s.provider === providerId && s.connected);
  };

  const handleConnect = (e: React.FormEvent) => {
    e.preventDefault();
    if (!isConnecting) return;

    connectMutation.mutate(
      { provider: isConnecting, accessToken, accountName: accountName || 'Personal' },
      {
        onSuccess: () => {
          setIsConnecting(null);
          setAccessToken('');
          setAccountName('');
          toast.success(`Successfully connected ${isConnecting}`);
        },
        onError: (err: any) => {
          toast.error(err.message || 'Failed to connect provider');
        },
      }
    );
  };

  const confirmDisconnect = () => {
    if (!isDisconnecting) return;
    disconnectMutation.mutate(isDisconnecting, {
      onSuccess: () => {
        toast.success(`Successfully disconnected ${isDisconnecting}`);
        setIsDisconnecting(null);
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
    <>
      <SettingsSection
        icon={<Link className="size-4 text-primary" />}
        title="Personal Git Providers"
        description="Connect your personal accounts using Access Tokens."
      >
        <div className="space-y-4">
          {PROVIDERS.map((provider) => {
            const status = getStatus(provider.id);
            const isConnected = !!status;

            return (
              <div key={provider.id} className="rounded-xl border border-border/50 p-4">
                <div className="flex items-start gap-3.5">
                  <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50">
                    <img
                      src={provider.icon}
                      alt={provider.name}
                      className="h-5 w-5 object-contain"
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
                      onClick={() => setIsDisconnecting(provider.id)}
                    >
                      Disconnect
                    </Button>
                  ) : (
                    <Button
                      variant="outline"
                      size="sm"
                      className="shrink-0 gap-1 text-xs"
                      onClick={() => {
                        setAccessToken('');
                        setAccountName('');
                        setIsConnecting(provider.id);
                      }}
                    >
                      Connect
                    </Button>
                  )}
                </div>

                <div className="mt-3">
                  {isConnected ? (
                    <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
                      <CheckCircle2 className="size-3" /> Connected as {status.accountName}
                    </span>
                  ) : (
                    <p className="py-0.5 text-muted-foreground text-xs">None added yet.</p>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </SettingsSection>

      <Dialog open={!!isConnecting} onOpenChange={(open) => !open && setIsConnecting(null)}>
        <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="font-bold text-foreground text-xl tracking-tight">
                  Connect {PROVIDERS.find((p) => p.id === isConnecting)?.name}
                </DialogTitle>
                <DialogDescription>Enter your Personal Access Token</DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="h-px w-full bg-border/50" />

          <form onSubmit={handleConnect} className="space-y-6 px-5 pt-4 pb-5">
            <div className="space-y-2">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                Account Name (Optional)
              </Label>
              <Input
                placeholder="e.g. My Personal Account"
                value={accountName}
                onChange={(e) => setAccountName(e.target.value)}
                className="h-11 bg-background/50 font-mono"
              />
            </div>

            <div className="space-y-2">
              <Label className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
                Personal Access Token
              </Label>
              <Input
                type="password"
                placeholder="glpat-xxxxxxxxxxxxxxxxxxxx"
                value={accessToken}
                onChange={(e) => setAccessToken(e.target.value)}
                required
                className="h-11 bg-background/50 font-mono"
              />
            </div>

            <div className="flex justify-end gap-3 pt-4">
              <Button type="button" variant="ghost" onClick={() => setIsConnecting(null)}>
                Cancel
              </Button>
              <Button type="submit" disabled={connectMutation.isPending} className="gap-2">
                <Check className="h-3.5 w-3.5" />
                {connectMutation.isPending ? 'Connecting...' : 'Connect'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={!!isDisconnecting} onOpenChange={(open) => !open && setIsDisconnecting(null)}>
        <DialogContent className="border/50 gap-0 bg-card/95 p-0 backdrop-blur-xl sm:max-w-md [&>button]:hidden">
          <div className="p-5">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-destructive text-xl tracking-tight">
                  <Trash className="h-5 w-5" />
                  Disconnect {PROVIDERS.find((p) => p.id === isDisconnecting)?.name}
                </DialogTitle>
                <DialogDescription>
                  This will prevent deploying new code from this provider.
                </DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button variant="ghost" onClick={() => setIsDisconnecting(null)}>
              Cancel
            </Button>
            <Button
              onClick={(e) => {
                e.preventDefault();
                confirmDisconnect();
              }}
              disabled={disconnectMutation.isPending}
              variant="destructive"
              className="gap-2"
            >
              <Trash className="h-3.5 w-3.5" />
              {disconnectMutation.isPending ? 'Disconnecting...' : 'Disconnect'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
