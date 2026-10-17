import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { AlertCircle, CheckCircle2, Cloud } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Dialog, DialogContent } from '@codedock/ui/dialog';
import { Label } from '@codedock/ui/label';
import { Switch } from '@codedock/ui/switch';
import { ingressService } from '#/services/ingress';

interface DomainIngressSettingsProps {
  domainId: string;
  serverId?: string;
  hostname: string;
  onChanged?: () => void;
}

export function DomainIngressSettings({
  domainId,
  serverId,
  hostname,
  onChanged,
}: DomainIngressSettingsProps) {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);

  const { data: host, isLoading } = useQuery({
    queryKey: ['domain-ingress', domainId],
    queryFn: () => ingressService.getDomain(domainId),
  });

  const { data: connections = [] } = useQuery({
    queryKey: ['server-ingress', serverId],
    queryFn: () => ingressService.list(serverId),
    enabled: Boolean(serverId),
  });

  const defaultConnection = connections[0] ?? null;

  const { mutate: updateDomain, isPending: isUpdating } = useMutation({
    mutationFn: (connectionId: string | null) =>
      ingressService.setDomain({
        domainId,
        connectionId,
      }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['domain-ingress', domainId] });
      toast.success(data ? 'Domain mapped to ingress tunnel' : 'Ingress routing disabled');
      onChanged?.();
      setModalOpen(false);
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : 'Failed to update domain ingress');
    },
  });

  const isEnabled = Boolean(host && host.status !== 'removing');

  const handleToggle = (checked: boolean) => {
    if (checked) {
      if (!defaultConnection) {
        toast.error('No ingress tunnel available on this server to route through');
        return;
      }
      updateDomain(defaultConnection.id);
    } else {
      updateDomain(null);
    }
  };

  return (
    <div className="flex items-center gap-2" data-testid="domain-ingress-settings">
      {host && (
        <span
          className={`inline-flex items-center gap-1 font-medium text-xs ${
            host.status === 'ready'
              ? 'text-emerald-500'
              : host.status === 'error'
                ? 'text-destructive'
                : 'text-amber-500'
          }`}
        >
          {host.status === 'ready' ? (
            <CheckCircle2 className="size-3" />
          ) : (
            <AlertCircle className="size-3" />
          )}
          Tunnel: {host.status}
        </span>
      )}

      <Button
        variant="ghost"
        size="sm"
        onClick={() => setModalOpen(true)}
        className="h-8 gap-1.5 text-xs"
      >
        <Cloud className="size-3.5" />
        Tunnel Routing
      </Button>

      <Dialog open={modalOpen} onOpenChange={setModalOpen}>
        <DialogContent className="max-w-md p-5">
          <div className="flex items-center gap-3 border-border/50 border-b pb-4">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10">
              <Cloud className="size-5 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold text-base text-foreground">Tunnel Ingress</h3>
              <p className="text-muted-foreground text-xs">{hostname}</p>
            </div>
          </div>

          <div className="space-y-4 py-4">
            <div className="flex items-center justify-between rounded-xl bg-muted/40 p-4">
              <div>
                <Label className="font-medium text-foreground text-sm">
                  Enable Ingress Routing
                </Label>
                <p className="text-muted-foreground text-xs">
                  Proxy traffic through Cloudflare Tunnel without public ports
                </p>
              </div>
              <Switch
                checked={isEnabled}
                onCheckedChange={handleToggle}
                disabled={isUpdating || isLoading}
              />
            </div>

            {defaultConnection && (
              <div className="rounded-xl border border-border/60 p-3 text-xs">
                <span className="text-muted-foreground">Target Tunnel: </span>
                <span className="font-medium text-foreground">{defaultConnection.name}</span>
                <span className="text-muted-foreground"> ({defaultConnection.connector})</span>
              </div>
            )}
          </div>

          <div className="flex justify-end pt-2">
            <Button variant="outline" size="sm" onClick={() => setModalOpen(false)}>
              Done
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
