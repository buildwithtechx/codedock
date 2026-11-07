import {
  AlertCircle,
  CheckCircle2,
  Cloud,
  Loader2,
  Plus,
  RefreshCw,
  Settings,
  Trash2,
} from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Dialog, DialogContent } from '@codedock/ui/dialog';
import type { IngressConnection } from '#/interfaces/ingress';
import { useDisconnectIngress, useIngressConnections, useRetryIngress } from '#/hooks/use-ingress';
import { IngressConnectForm } from './ingress-connect-form';
import { IngressPreferences } from './ingress-preferences';

interface ServerIngressCardProps {
  serverId: string;
}

export function ServerIngressCard({ serverId }: ServerIngressCardProps) {
  const [connectOpen, setConnectOpen] = useState(false);
  const [preferencesOpen, setPreferencesOpen] = useState(false);
  const [selectedConnection, setSelectedConnection] = useState<IngressConnection | null>(null);

  const {
    data: connections = [],
    isLoading,
    isRefetching,
    refetch,
    isError,
  } = useIngressConnections(serverId);

  const retryMutation = useRetryIngress();
  const disconnectMutation = useDisconnectIngress();

  const retryConnection = (id: string) => {
    retryMutation.mutate(id, {
      onSuccess: () => {
        toast.success('Tunnel reconciliation triggered');
      },
      onError: (err) => {
        toast.error(err instanceof Error ? err.message : 'Failed to retry tunnel');
      },
    });
  };

  const disconnectConnection = (id: string) => {
    disconnectMutation.mutate(id, {
      onSuccess: () => {
        toast.success('Tunnel disconnected');
      },
      onError: (err) => {
        toast.error(err instanceof Error ? err.message : 'Failed to disconnect tunnel');
      },
    });
  };

  const isRetrying = retryMutation.isPending;
  const isDisconnecting = disconnectMutation.isPending;

  const current = connections[0] ?? null;

  return (
    <div
      className="rounded-2xl border border-border/50 bg-card shadow-xs"
      data-testid="server-ingress-card"
    >
      <div className="flex flex-wrap items-center justify-between gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex min-w-0 items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10">
            <Cloud className="size-[18px] text-primary" />
          </div>
          <div>
            <h2 className="font-semibold text-[15px] text-foreground">Managed Ingress</h2>
            <p className="text-muted-foreground text-xs">
              Cloudflare Tunnel per-server ingress and origin routing
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {current && (
            <Button
              variant="ghost"
              size="icon"
              disabled={isRefetching}
              aria-label="Refresh ingress status"
              onClick={() => void refetch()}
            >
              <RefreshCw className={`size-4 ${isRefetching ? 'animate-spin' : ''}`} />
            </Button>
          )}
          {!current && !isLoading && (
            <Button size="sm" onClick={() => setConnectOpen(true)} className="gap-1.5">
              <Plus className="size-4" />
              Connect tunnel
            </Button>
          )}
        </div>
      </div>

      <div className="p-5">
        {isLoading ? (
          <div className="flex h-20 items-center justify-center">
            <Loader2 className="size-5 animate-spin text-muted-foreground" />
          </div>
        ) : isError ? (
          <div className="rounded-xl bg-destructive/10 p-4 text-destructive text-sm">
            Failed to load ingress tunnel status.
          </div>
        ) : current ? (
          <div className="space-y-4">
            <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-muted/30 p-4">
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  <h3 className="font-semibold text-foreground text-sm">{current.name}</h3>
                  <span
                    className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 font-medium text-xs ${
                      current.status === 'ready'
                        ? 'bg-emerald-500/10 text-emerald-500'
                        : current.status === 'error'
                          ? 'bg-destructive/10 text-destructive'
                          : 'bg-amber-500/10 text-amber-500'
                    }`}
                  >
                    {current.status === 'ready' ? (
                      <CheckCircle2 className="size-3" />
                    ) : (
                      <AlertCircle className="size-3" />
                    )}
                    {current.status}
                  </span>
                </div>
                <p className="mt-1 text-muted-foreground text-xs">
                  {current.connector === 'managed'
                    ? 'Managed cloudflared connector'
                    : 'Existing external tunnel'}
                  {' • '}
                  Account: {current.accountId}
                </p>
              </div>

              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setSelectedConnection(current);
                    setPreferencesOpen(true);
                  }}
                  className="gap-1.5"
                >
                  <Settings className="size-3.5" />
                  Preferences
                </Button>

                {current.status !== 'ready' && (
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={isRetrying}
                    onClick={() => retryConnection(current.id)}
                    className="gap-1.5"
                  >
                    <RefreshCw className={`size-3.5 ${isRetrying ? 'animate-spin' : ''}`} />
                    Retry
                  </Button>
                )}

                <Button
                  variant="ghost"
                  size="sm"
                  disabled={isDisconnecting}
                  onClick={() => disconnectConnection(current.id)}
                  className="text-destructive hover:bg-destructive/10 hover:text-destructive"
                >
                  <Trash2 className="size-3.5" />
                </Button>
              </div>
            </div>

            <p className="text-muted-foreground text-xs">
              {current.defaultForDomains
                ? 'Automatic enrollment active: newly added service domains are routed via this tunnel.'
                : 'Manual enrollment: domain hostnames must be explicitly attached.'}
            </p>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center rounded-xl border border-border/80 border-dashed py-8 text-center">
            <Cloud className="mb-2 size-8 text-muted-foreground/60" />
            <p className="font-medium text-foreground text-sm">No ingress tunnel connected</p>
            <p className="mt-1 max-w-sm text-muted-foreground text-xs">
              Connect Cloudflare Tunnel to expose services on this host without opening public
              inbound firewall ports.
            </p>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setConnectOpen(true)}
              className="mt-4 gap-1.5"
            >
              <Plus className="size-3.5" />
              Connect tunnel
            </Button>
          </div>
        )}
      </div>

      <Dialog open={connectOpen} onOpenChange={setConnectOpen}>
        <DialogContent className="max-w-md p-0">
          <IngressConnectForm
            serverId={serverId}
            onClose={() => setConnectOpen(false)}
            onSaved={() => {
              setConnectOpen(false);
              refetch();
            }}
          />
        </DialogContent>
      </Dialog>

      <Dialog open={preferencesOpen} onOpenChange={setPreferencesOpen}>
        <DialogContent className="max-w-md p-0">
          {selectedConnection && (
            <IngressPreferences
              connection={selectedConnection}
              onClose={() => setPreferencesOpen(false)}
              onSaved={() => {
                setPreferencesOpen(false);
                refetch();
              }}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
